package analyticsservicelogic

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/config"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/domain"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/model"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	linkv1 "github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"
)

// FakeShortLinkModel is a test double for ShortLinkModel.
type FakeShortLinkModel struct {
	model.ShortLinkModel
	Link *model.ShortLink
	Err  error
}

func (f FakeShortLinkModel) FindOneByCode(_ context.Context, _ string) (*model.ShortLink, error) {
	return f.Link, f.Err
}

func activeLink() *model.ShortLink {
	return &model.ShortLink{
		Id:        1,
		Code:      "abc123",
		OriginUrl: "https://example.com",
		Status:    domain.LinkStatusActive,
	}
}

type fakeVisitEventModel struct {
	model.VisitEventModel
	Err error
}

func (f fakeVisitEventModel) Insert(_ context.Context, _ *model.VisitEvent) (sql.Result, error) {
	return nil, f.Err
}

type fakeDailyStatModel struct {
	model.LinkDailyStatModel
	UpsertErr error
}

func (f fakeDailyStatModel) UpsertPV(_ context.Context, _ int64, _ string) error {
	return f.UpsertErr
}

func newRecordVisitLogic(
	shortLink model.ShortLinkModel, visitEvent model.VisitEventModel, dailyStat model.LinkDailyStatModel,
) *RecordVisitLogic {
	return NewRecordVisitLogic(context.Background(), &svc.ServiceContext{
		Config:          config.Config{Analytics: config.AnalyticsConfig{IPSalt: "test-salt"}},
		ShortLinkModel:  shortLink,
		VisitEventModel: visitEvent,
		DailyStatModel:  dailyStat,
	})
}

func TestRecordVisit_Success(t *testing.T) {
	logic := newRecordVisitLogic(
		FakeShortLinkModel{Link: activeLink()},
		fakeVisitEventModel{},
		fakeDailyStatModel{},
	)
	_, err := logic.RecordVisit(&linkv1.RecordVisitRequest{
		Code:      "abc123",
		VisitedAt: time.Now().UTC().Format(time.RFC3339Nano),
		Ip:        "1.2.3.4",
		UserAgent: "Mozilla/5.0 Chrome/120",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRecordVisit_LinkNotFound(t *testing.T) {
	logic := newRecordVisitLogic(
		FakeShortLinkModel{Err: model.ErrNotFound},
		fakeVisitEventModel{},
		fakeDailyStatModel{},
	)
	resp, err := logic.RecordVisit(&linkv1.RecordVisitRequest{Code: "missing"})
	if err != nil {
		t.Fatalf("expected nil error for missing link, got %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
}

func TestRecordVisit_InsertFails(t *testing.T) {
	logic := newRecordVisitLogic(
		FakeShortLinkModel{Link: activeLink()},
		fakeVisitEventModel{Err: errors.New("db error")},
		fakeDailyStatModel{},
	)
	_, err := logic.RecordVisit(&linkv1.RecordVisitRequest{
		Code: "abc123",
		Ip:   "1.2.3.4",
	})
	if err == nil {
		t.Fatal("expected error when Insert fails")
	}
}

func TestRecordVisit_UpsertPVFails_StillSucceeds(t *testing.T) {
	logic := newRecordVisitLogic(
		FakeShortLinkModel{Link: activeLink()},
		fakeVisitEventModel{},
		fakeDailyStatModel{UpsertErr: errors.New("stat error")},
	)
	_, err := logic.RecordVisit(&linkv1.RecordVisitRequest{
		Code: "abc123",
		Ip:   "1.2.3.4",
	})
	if err != nil {
		t.Fatalf("UpsertPV failure should not propagate, got %v", err)
	}
}
