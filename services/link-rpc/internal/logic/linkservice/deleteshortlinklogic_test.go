package linkservicelogic

import (
	"context"
	"testing"
	"time"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/model"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	linkv1 "github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeDeleteModel struct {
	model.ShortLinkModel
	softDeleteErr error
}

func (f fakeDeleteModel) SoftDelete(_ context.Context, _ int64, _ time.Time) error {
	return f.softDeleteErr
}

func newDeleteLogic(m model.ShortLinkModel) *DeleteShortLinkLogic {
	return NewDeleteShortLinkLogic(context.Background(), &svc.ServiceContext{
		ShortLinkModel: m,
	})
}

func TestDeleteShortLink_Success(t *testing.T) {
	t.Parallel()
	resp, err := newDeleteLogic(fakeDeleteModel{}).DeleteShortLink(&linkv1.DeleteShortLinkRequest{Id: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Deleted {
		t.Fatal("want Deleted=true")
	}
	if resp.Id != 5 {
		t.Fatalf("want id 5, got %d", resp.Id)
	}
}

func TestDeleteShortLink_NotFound(t *testing.T) {
	t.Parallel()
	m := fakeDeleteModel{softDeleteErr: model.ErrNotFound}
	_, err := newDeleteLogic(m).DeleteShortLink(&linkv1.DeleteShortLinkRequest{Id: 99})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("want NotFound, got %v", err)
	}
}

func TestDeleteShortLink_InvalidID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		id   int64
	}{
		{name: "zero", id: 0},
		{name: "negative", id: -5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := newDeleteLogic(fakeDeleteModel{}).DeleteShortLink(&linkv1.DeleteShortLinkRequest{Id: tt.id})
			if status.Code(err) != codes.InvalidArgument {
				t.Fatalf("want InvalidArgument, got %v", err)
			}
		})
	}
}
