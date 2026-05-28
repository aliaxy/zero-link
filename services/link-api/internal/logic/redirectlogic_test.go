package logic

import (
	"context"
	"testing"

	"github.com/aliaxy/zero-link/services/link-api/internal/svc"
	"github.com/aliaxy/zero-link/services/link-api/internal/types"
	"github.com/aliaxy/zero-link/services/link-rpc/linkservice"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newRedirectLogic(svc *svc.ServiceContext) *RedirectLogic {
	return NewRedirectLogic(context.Background(), svc)
}

func TestRedirectLogic_Success(t *testing.T) {
	svcCtx := &svc.ServiceContext{
		LinkRPC: fakeLinkService{
			resolveResp: &linkservice.ResolveShortLinkResponse{OriginUrl: "https://example.com"},
		},
	}
	resp, err := newRedirectLogic(svcCtx).Redirect(&types.RedirectRequest{Code: "abc123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.OriginUrl != "https://example.com" {
		t.Fatalf("want https://example.com, got %q", resp.OriginUrl)
	}
}

func TestRedirectLogic_NotFound(t *testing.T) {
	svcCtx := &svc.ServiceContext{
		LinkRPC: fakeLinkService{err: status.Error(codes.NotFound, "not found")},
	}
	_, err := newRedirectLogic(svcCtx).Redirect(&types.RedirectRequest{Code: "missing"})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("want NotFound, got %v", err)
	}
}

func TestRedirectLogic_Disabled(t *testing.T) {
	svcCtx := &svc.ServiceContext{
		LinkRPC: fakeLinkService{err: status.Error(codes.PermissionDenied, "disabled")},
	}
	_, err := newRedirectLogic(svcCtx).Redirect(&types.RedirectRequest{Code: "abc123"})
	if status.Code(err) != codes.PermissionDenied {
		t.Fatalf("want PermissionDenied, got %v", err)
	}
}

func TestRedirectLogic_Expired(t *testing.T) {
	svcCtx := &svc.ServiceContext{
		LinkRPC: fakeLinkService{err: status.Error(codes.FailedPrecondition, "expired")},
	}
	_, err := newRedirectLogic(svcCtx).Redirect(&types.RedirectRequest{Code: "abc123"})
	if status.Code(err) != codes.FailedPrecondition {
		t.Fatalf("want FailedPrecondition, got %v", err)
	}
}
