package logic

import (
	"context"
	"testing"

	"github.com/aliaxy/zero-link/services/link-api/internal/auth"
	"github.com/aliaxy/zero-link/services/link-api/internal/middleware"
	"github.com/aliaxy/zero-link/services/link-api/internal/svc"
	"github.com/aliaxy/zero-link/services/link-api/internal/types"
	"github.com/aliaxy/zero-link/services/link-rpc/linkservice"

	"google.golang.org/grpc"
)

type fakeLinkService struct {
	authResp    *linkservice.AuthenticateAdminResponse
	profileResp *linkservice.GetAdminProfileResponse
	createResp  *linkservice.CreateShortLinkResponse
	listResp    *linkservice.ListShortLinksResponse
	getResp     *linkservice.GetShortLinkResponse
	updateResp  *linkservice.UpdateShortLinkResponse
	deleteResp  *linkservice.DeleteShortLinkResponse
	err         error
}

func (f fakeLinkService) Check(
	_ context.Context,
	_ *linkservice.CheckRequest,
	_ ...grpc.CallOption,
) (*linkservice.CheckResponse, error) {
	return nil, nil
}

func (f fakeLinkService) AuthenticateAdmin(
	_ context.Context,
	_ *linkservice.AuthenticateAdminRequest,
	_ ...grpc.CallOption,
) (*linkservice.AuthenticateAdminResponse, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.authResp, nil
}

func (f fakeLinkService) GetAdminProfile(
	_ context.Context,
	_ *linkservice.GetAdminProfileRequest,
	_ ...grpc.CallOption,
) (*linkservice.GetAdminProfileResponse, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.profileResp, nil
}

func (f fakeLinkService) CreateShortLink(
	_ context.Context,
	_ *linkservice.CreateShortLinkRequest,
	_ ...grpc.CallOption,
) (*linkservice.CreateShortLinkResponse, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.createResp, nil
}

func (f fakeLinkService) ListShortLinks(
	_ context.Context,
	_ *linkservice.ListShortLinksRequest,
	_ ...grpc.CallOption,
) (*linkservice.ListShortLinksResponse, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.listResp, nil
}

func (f fakeLinkService) GetShortLink(
	_ context.Context,
	_ *linkservice.GetShortLinkRequest,
	_ ...grpc.CallOption,
) (*linkservice.GetShortLinkResponse, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.getResp, nil
}

func (f fakeLinkService) UpdateShortLink(
	_ context.Context,
	_ *linkservice.UpdateShortLinkRequest,
	_ ...grpc.CallOption,
) (*linkservice.UpdateShortLinkResponse, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.updateResp, nil
}

func (f fakeLinkService) DeleteShortLink(
	_ context.Context,
	_ *linkservice.DeleteShortLinkRequest,
	_ ...grpc.CallOption,
) (*linkservice.DeleteShortLinkResponse, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.deleteResp, nil
}

func (f fakeLinkService) ResolveShortLink(
	_ context.Context,
	_ *linkservice.ResolveShortLinkRequest,
	_ ...grpc.CallOption,
) (*linkservice.ResolveShortLinkResponse, error) {
	if f.err != nil {
		return nil, f.err
	}
	return nil, nil
}

func TestLoginLogic_Login(t *testing.T) {
	tokenManager := auth.NewTokenManager(auth.Config{
		Secret:          "local-test-secret",
		TokenTTLSeconds: 3600,
	})
	svcCtx := &svc.ServiceContext{
		TokenManager: tokenManager,
		LinkRPC: fakeLinkService{
			authResp: &linkservice.AuthenticateAdminResponse{
				Admin: &linkservice.AdminProfile{
					Id:       42,
					Username: "admin",
				},
			},
		},
	}

	logic := NewLoginLogic(context.Background(), svcCtx)
	resp, err := logic.Login(&types.LoginRequest{
		Username: "admin",
		Password: "zerolink",
	})
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	if resp.Code != "OK" {
		t.Fatalf("response code = %q, want OK", resp.Code)
	}
	if resp.Data.Token == "" {
		t.Fatal("response token is empty")
	}
	if resp.Data.Admin.Id != 42 {
		t.Fatalf("admin ID = %d, want 42", resp.Data.Admin.Id)
	}
}

func TestCreateShortLinkLogic_CreateShortLink(t *testing.T) {
	logic := NewCreateShortLinkLogic(authenticatedContext(), &svc.ServiceContext{
		LinkRPC: fakeLinkService{
			createResp: &linkservice.CreateShortLinkResponse{
				Link: sampleRPCLink(),
			},
		},
	})

	resp, err := logic.CreateShortLink(&types.CreateShortLinkRequest{
		OriginUrl: "https://example.com/page",
		Code:      "campaign1",
	})
	if err != nil {
		t.Fatalf("CreateShortLink() error = %v", err)
	}
	if resp.Code != "OK" {
		t.Fatalf("response code = %q, want OK", resp.Code)
	}
	if resp.Data.Id != 1001 {
		t.Fatalf("link ID = %d, want 1001", resp.Data.Id)
	}
}

func TestListShortLinksLogic_ListShortLinks(t *testing.T) {
	logic := NewListShortLinksLogic(authenticatedContext(), &svc.ServiceContext{
		LinkRPC: fakeLinkService{
			listResp: &linkservice.ListShortLinksResponse{
				Items:    []*linkservice.ShortLinkSummary{sampleRPCLinkSummary()},
				Page:     1,
				PageSize: 20,
				Total:    1,
			},
		},
	})

	resp, err := logic.ListShortLinks(&types.ListShortLinksRequest{})
	if err != nil {
		t.Fatalf("ListShortLinks() error = %v", err)
	}
	if len(resp.Data.Items) != 1 {
		t.Fatalf("items length = %d, want 1", len(resp.Data.Items))
	}
}

func TestGetUpdateDeleteShortLinkLogic(t *testing.T) {
	svcCtx := &svc.ServiceContext{
		LinkRPC: fakeLinkService{
			getResp: &linkservice.GetShortLinkResponse{
				Link: sampleRPCLink(),
			},
			updateResp: &linkservice.UpdateShortLinkResponse{
				Link: sampleRPCLink(),
			},
			deleteResp: &linkservice.DeleteShortLinkResponse{
				Id:      1001,
				Deleted: true,
			},
		},
	}

	getResp, err := NewGetShortLinkLogic(authenticatedContext(), svcCtx).GetShortLink(&types.LinkIdRequest{Id: 1001})
	if err != nil {
		t.Fatalf("GetShortLink() error = %v", err)
	}
	if getResp.Data.Id != 1001 {
		t.Fatalf("get link ID = %d, want 1001", getResp.Data.Id)
	}

	updateResp, err := NewUpdateShortLinkLogic(authenticatedContext(), svcCtx).UpdateShortLink(
		&types.UpdateShortLinkRequest{Id: 1001, Title: "Updated"},
	)
	if err != nil {
		t.Fatalf("UpdateShortLink() error = %v", err)
	}
	if updateResp.Data.Id != 1001 {
		t.Fatalf("update link ID = %d, want 1001", updateResp.Data.Id)
	}

	deleteResp, err := NewDeleteShortLinkLogic(authenticatedContext(), svcCtx).DeleteShortLink(
		&types.LinkIdRequest{Id: 1001},
	)
	if err != nil {
		t.Fatalf("DeleteShortLink() error = %v", err)
	}
	if !deleteResp.Data.Deleted {
		t.Fatal("delete response deleted = false, want true")
	}
}

func authenticatedContext() context.Context {
	return middleware.ContextWithAdminSubject(context.Background(), auth.AdminSubject{
		ID:       42,
		Username: "admin",
	})
}

func sampleRPCLink() *linkservice.ShortLink {
	return &linkservice.ShortLink{
		Id:          1001,
		Code:        "campaign1",
		OriginUrl:   "https://example.com/page",
		Title:       "Campaign 1",
		Description: "Optional note",
		Status:      1,
		CreatedBy:   42,
		CreatedAt:   "2026-05-27T12:00:00Z",
		UpdatedAt:   "2026-05-27T12:00:00Z",
	}
}

func sampleRPCLinkSummary() *linkservice.ShortLinkSummary {
	return &linkservice.ShortLinkSummary{
		Id:        1001,
		Code:      "campaign1",
		OriginUrl: "https://example.com/page",
		Title:     "Campaign 1",
		Status:    1,
		CreatedAt: "2026-05-27T12:00:00Z",
		UpdatedAt: "2026-05-27T12:00:00Z",
	}
}

func TestProfileLogic_Profile(t *testing.T) {
	svcCtx := &svc.ServiceContext{
		LinkRPC: fakeLinkService{
			profileResp: &linkservice.GetAdminProfileResponse{
				Admin: &linkservice.AdminProfile{
					Id:       42,
					Username: "admin",
					Status:   1,
				},
			},
		},
	}

	ctx := middleware.ContextWithAdminSubject(context.Background(), auth.AdminSubject{
		ID:       42,
		Username: "admin",
	})
	logic := NewProfileLogic(ctx, svcCtx)

	resp, err := logic.Profile()
	if err != nil {
		t.Fatalf("Profile() error = %v", err)
	}
	if resp.Code != "OK" {
		t.Fatalf("response code = %q, want OK", resp.Code)
	}
	if resp.Data.Id != 42 {
		t.Fatalf("admin ID = %d, want 42", resp.Data.Id)
	}
}
