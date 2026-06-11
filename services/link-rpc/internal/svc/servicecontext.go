// Package svc wires link-rpc service dependencies.
package svc

import (
	"context"
	"fmt"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/cleanup"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/config"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/model"
	"github.com/aliaxy/zero-link/services/link-rpc/pkg/filter"

	goredis "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const codeCreatedChannel = "zl:code:created"

// ServiceContext holds dependencies shared by link-rpc logic.
type ServiceContext struct {
	Config            config.Config
	DB                sqlx.SqlConn
	Redis             *redis.Redis
	CodeFilter        *filter.CodeFilter
	AdminUserModel    model.AdminUserModel
	ShortLinkModel    model.ShortLinkModel
	VisitEventModel   model.VisitEventModel
	DailyStatModel    model.LinkDailyStatModel
	ReservedCodeModel model.ReservedCodeModel
}

// NewServiceContext creates a link-rpc service context.
func NewServiceContext(c config.Config) *ServiceContext {
	applyDefaults(&c)

	db := sqlx.NewMysql(c.Dependencies.MySQL.DataSource)
	rdb := redis.MustNewRedis(c.CacheRedis[0].RedisConf)
	cf := filter.NewCodeFilter(c.Cuckoo.Capacity)

	rootCtx, rootCancel := context.WithCancel(context.Background())

	// Subscribe before loading DB to avoid a race window where a newly created
	// code lands between the end of the batch scan and the start of Subscribe.
	pubsubClient := goredis.NewClient(&goredis.Options{
		Addr:     c.CacheRedis[0].Host,
		Password: c.CacheRedis[0].Pass,
	})
	pubsub := pubsubClient.Subscribe(rootCtx, codeCreatedChannel)
	go runCodeSubscription(pubsub, cf)

	proc.AddShutdownListener(func() {
		rootCancel()
		_ = pubsub.Close()
		_ = pubsubClient.Close()
	})

	loadCodesIntoFilter(context.Background(), db, cf)

	// Models are created before the cleanup runner so they can be shared.
	shortLinkModel := model.NewShortLinkModel(db, c.CacheRedis)
	reservedCodeModel := model.NewReservedCodeModel(db)
	archiver := model.NewLinkArchiver(db, shortLinkModel, reservedCodeModel)
	cleanupRunner := cleanup.NewRunner(db, shortLinkModel, archiver, c.Retention)
	cleanupRunner.Start(rootCtx)

	return &ServiceContext{
		Config:            c,
		DB:                db,
		Redis:             rdb,
		CodeFilter:        cf,
		AdminUserModel:    model.NewAdminUserModel(db, c.CacheRedis),
		ShortLinkModel:    shortLinkModel,
		VisitEventModel:   model.NewVisitEventModel(db),
		DailyStatModel:    model.NewLinkDailyStatModel(db),
		ReservedCodeModel: reservedCodeModel,
	}
}

func applyDefaults(c *config.Config) {
	if c.Retention.VisitEventRetentionDays == 0 {
		c.Retention.VisitEventRetentionDays = 90
	}
	if c.Retention.ShortLinkRetentionDays == 0 {
		c.Retention.ShortLinkRetentionDays = 365
	}
	if c.Retention.DailyStatRetentionDays == 0 {
		c.Retention.DailyStatRetentionDays = 730
	}
	if c.Retention.CleanupBatchSize == 0 {
		c.Retention.CleanupBatchSize = 1000
	}
	if c.Cuckoo.Capacity == 0 {
		c.Cuckoo.Capacity = 1_000_000
	}
}

type codeRow struct {
	ID   int64  `db:"id"`
	Code string `db:"code"`
}

func loadCodesIntoFilter(ctx context.Context, db sqlx.SqlConn, cf *filter.CodeFilter) {
	const batch = 10_000
	var lastID int64
	total := 0
	for {
		var rows []codeRow
		if err := db.QueryRowsCtx(ctx, &rows,
			fmt.Sprintf("select `id`, `code` from `short_link` where `id` > %d order by `id` limit %d", lastID, batch),
		); err != nil {
			logx.Errorw("filter: load codes failed",
				logx.Field("last_id", lastID), logx.Field("error", err.Error()))
			return
		}
		for _, r := range rows {
			cf.Insert(r.Code)
			lastID = r.ID
		}
		total += len(rows)
		if len(rows) < batch {
			break
		}
	}
	logx.Infow("filter: loaded codes into cuckoo filter", logx.Field("total", total))
}

func runCodeSubscription(pubsub *goredis.PubSub, cf *filter.CodeFilter) {
	for msg := range pubsub.Channel() {
		if msg.Payload != "" {
			cf.Insert(msg.Payload)
		}
	}
}
