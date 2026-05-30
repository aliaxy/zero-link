package linkservicelogic

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/domain"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/model"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	linkv1 "github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func newResolveLogic(m model.ShortLinkModel) *ResolveShortLinkLogic {
	return NewResolveShortLinkLogic(context.Background(), &svc.ServiceContext{
		ShortLinkModel: m,
	})
}

func TestResolveShortLink_NotFound(t *testing.T) {
	logic := newResolveLogic(FakeShortLinkModel{Err: model.ErrNotFound})
	_, err := logic.ResolveShortLink(&linkv1.ResolveShortLinkRequest{Code: "abc123"})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("want NotFound, got %v", err)
	}
}

func TestResolveShortLink_SoftDeleted(t *testing.T) {
	link := activeLink()
	link.DeletedAt = sql.NullTime{Time: time.Now(), Valid: true}
	logic := newResolveLogic(FakeShortLinkModel{Link: link})
	_, err := logic.ResolveShortLink(&linkv1.ResolveShortLinkRequest{Code: "abc123"})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("want NotFound for soft-deleted, got %v", err)
	}
}

func TestResolveShortLink_Disabled(t *testing.T) {
	link := activeLink()
	link.Status = domain.LinkStatusDisabled
	logic := newResolveLogic(FakeShortLinkModel{Link: link})
	_, err := logic.ResolveShortLink(&linkv1.ResolveShortLinkRequest{Code: "abc123"})
	if status.Code(err) != codes.PermissionDenied {
		t.Fatalf("want PermissionDenied, got %v", err)
	}
}

func TestResolveShortLink_Expired(t *testing.T) {
	link := activeLink()
	link.ExpireAt = sql.NullTime{Time: time.Now().Add(-time.Hour), Valid: true}
	logic := newResolveLogic(FakeShortLinkModel{Link: link})
	_, err := logic.ResolveShortLink(&linkv1.ResolveShortLinkRequest{Code: "abc123"})
	if status.Code(err) != codes.FailedPrecondition {
		t.Fatalf("want FailedPrecondition for expired, got %v", err)
	}
}

func TestResolveShortLink_Active(t *testing.T) {
	logic := newResolveLogic(FakeShortLinkModel{Link: activeLink()})
	resp, err := logic.ResolveShortLink(&linkv1.ResolveShortLinkRequest{Code: "abc123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.OriginUrl != "https://example.com" {
		t.Fatalf("want origin_url https://example.com, got %q", resp.OriginUrl)
	}
}

func TestResolveShortLink_NotYetExpired(t *testing.T) {
	link := activeLink()
	link.ExpireAt = sql.NullTime{Time: time.Now().Add(time.Hour), Valid: true}
	logic := newResolveLogic(FakeShortLinkModel{Link: link})
	resp, err := logic.ResolveShortLink(&linkv1.ResolveShortLinkRequest{Code: "abc123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.OriginUrl != "https://example.com" {
		t.Fatalf("want origin_url https://example.com, got %q", resp.OriginUrl)
	}
}
