package linkservicelogic

import (
	"context"
	"time"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/domain"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/metrics"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/model"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/rpcerror"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	linkv1 "github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

// ResolveShortLinkLogic handles short-link resolution RPC.
type ResolveShortLinkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewResolveShortLinkLogic creates a ResolveShortLinkLogic.
func NewResolveShortLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResolveShortLinkLogic {
	return &ResolveShortLinkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ResolveShortLink resolves a short code to its origin URL.
func (l *ResolveShortLinkLogic) ResolveShortLink(
	in *linkv1.ResolveShortLinkRequest,
) (*linkv1.ResolveShortLinkResponse, error) {
	// Fast path: cuckoo filter definitively says code does not exist — skip Redis and DB.
	// go-zero's internal singleflight already covers cache breakdown; this handles penetration.
	if l.svcCtx.CodeFilter != nil {
		if !l.svcCtx.CodeFilter.Lookup(in.Code) {
			l.Infow("resolve miss", logx.Field("code", in.Code), logx.Field("result", "miss"))
			metrics.FilterRequestsTotal.Inc("miss")
			metrics.RedirectRequestsTotal.Inc("miss")
			return nil, rpcerror.ToRPC(domain.ErrNotFound)
		}
		metrics.FilterRequestsTotal.Inc("hit")
	}

	link, err := l.svcCtx.ShortLinkModel.FindOneByCode(l.ctx, in.Code)
	if err != nil {
		if err == model.ErrNotFound {
			l.Infow("resolve miss", logx.Field("code", in.Code), logx.Field("result", "miss"))
			metrics.RedirectRequestsTotal.Inc("miss")
			return nil, rpcerror.ToRPC(domain.ErrNotFound)
		}
		l.Errorw("resolve failed", logx.Field("code", in.Code), logx.Field("error", err.Error()))
		metrics.RedirectRequestsTotal.Inc("error")
		return nil, rpcerror.ToRPC(err)
	}

	if link.DeletedAt.Valid {
		l.Infow("resolve miss", logx.Field("code", in.Code), logx.Field("result", "miss"))
		metrics.RedirectRequestsTotal.Inc("miss")
		return nil, rpcerror.ToRPC(domain.ErrNotFound)
	}

	if link.Status == domain.LinkStatusDisabled {
		l.Infow("resolve disabled", logx.Field("code", in.Code), logx.Field("result", "disabled"))
		metrics.RedirectRequestsTotal.Inc("disabled")
		return nil, rpcerror.ToRPC(domain.ErrPermissionDenied)
	}

	if link.ExpireAt.Valid && link.ExpireAt.Time.Before(time.Now()) {
		l.Infow("resolve expired", logx.Field("code", in.Code), logx.Field("result", "expired"))
		metrics.RedirectRequestsTotal.Inc("expired")
		return nil, rpcerror.ToRPC(domain.ErrGone)
	}

	l.Infow("resolve hit", logx.Field("code", in.Code), logx.Field("result", "hit"), logx.Field("url", link.OriginUrl))
	metrics.RedirectRequestsTotal.Inc("hit")
	return &linkv1.ResolveShortLinkResponse{OriginUrl: link.OriginUrl}, nil
}
