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

type fakeUpdateModel struct {
	model.ShortLinkModel

	link      *model.ShortLink
	findErr   error
	updateErr error
}

func (f *fakeUpdateModel) FindOneNotDeleted(_ context.Context, _ int64) (*model.ShortLink, error) {
	return f.link, f.findErr
}

func (f *fakeUpdateModel) Update(_ context.Context, _ *model.ShortLink) error {
	return f.updateErr
}

func newUpdateLogic(m model.ShortLinkModel) *UpdateShortLinkLogic {
	return NewUpdateShortLinkLogic(context.Background(), &svc.ServiceContext{
		ShortLinkModel: m,
	})
}

func TestUpdateShortLink_ValidUpdate(t *testing.T) {
	t.Parallel()
	m := &fakeUpdateModel{link: activeShortLink(10, "abc123")}
	resp, err := newUpdateLogic(m).UpdateShortLink(&linkv1.UpdateShortLinkRequest{
		Id:    10,
		Title: "New Title",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Link.Id != 10 {
		t.Fatalf("want id 10, got %d", resp.Link.Id)
	}
}

func TestUpdateShortLink_NotFound(t *testing.T) {
	t.Parallel()
	m := &fakeUpdateModel{findErr: model.ErrNotFound}
	_, err := newUpdateLogic(m).UpdateShortLink(&linkv1.UpdateShortLinkRequest{Id: 99, Title: "X"})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("want NotFound, got %v", err)
	}
}

func TestUpdateShortLink_InvalidID(t *testing.T) {
	t.Parallel()
	_, err := newUpdateLogic(&fakeUpdateModel{}).UpdateShortLink(&linkv1.UpdateShortLinkRequest{Id: 0})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("want InvalidArgument for id=0, got %v", err)
	}
}

func TestUpdateShortLink_InvalidOriginURL(t *testing.T) {
	t.Parallel()
	m := &fakeUpdateModel{link: activeShortLink(1, "abc")}
	_, err := newUpdateLogic(m).UpdateShortLink(&linkv1.UpdateShortLinkRequest{
		Id:        1,
		OriginUrl: "not-a-url",
	})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("want InvalidArgument for bad URL, got %v", err)
	}
}

func TestUpdateShortLink_InvalidStatus(t *testing.T) {
	t.Parallel()
	m := &fakeUpdateModel{link: activeShortLink(1, "abc")}
	_, err := newUpdateLogic(m).UpdateShortLink(&linkv1.UpdateShortLinkRequest{
		Id:     1,
		Status: 99,
	})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("want InvalidArgument for bad status, got %v", err)
	}
}
