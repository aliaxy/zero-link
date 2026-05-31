package linkservicelogic

import (
	"context"
	"testing"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/model"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	linkv1 "github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeGetModel struct {
	model.ShortLinkModel
	link *model.ShortLink
	err  error
}

func (f fakeGetModel) FindOneNotDeleted(_ context.Context, _ int64) (*model.ShortLink, error) {
	return f.link, f.err
}

func newGetLogic(m model.ShortLinkModel) *GetShortLinkLogic {
	return NewGetShortLinkLogic(context.Background(), &svc.ServiceContext{
		ShortLinkModel: m,
	})
}

func TestGetShortLink_Found(t *testing.T) {
	t.Parallel()
	m := fakeGetModel{link: activeShortLink(5, "abc123")}
	resp, err := newGetLogic(m).GetShortLink(&linkv1.GetShortLinkRequest{Id: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Link.Id != 5 {
		t.Fatalf("want id 5, got %d", resp.Link.Id)
	}
	if resp.Link.Code != "abc123" {
		t.Fatalf("want code abc123, got %q", resp.Link.Code)
	}
}

func TestGetShortLink_NotFound(t *testing.T) {
	t.Parallel()
	m := fakeGetModel{err: model.ErrNotFound}
	_, err := newGetLogic(m).GetShortLink(&linkv1.GetShortLinkRequest{Id: 99})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("want NotFound, got %v", err)
	}
}

func TestGetShortLink_InvalidID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		id   int64
	}{
		{name: "zero", id: 0},
		{name: "negative", id: -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := newGetLogic(fakeGetModel{}).GetShortLink(&linkv1.GetShortLinkRequest{Id: tt.id})
			if status.Code(err) != codes.InvalidArgument {
				t.Fatalf("want InvalidArgument, got %v", err)
			}
		})
	}
}
