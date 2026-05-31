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

type fakeListModel struct {
	model.ShortLinkModel

	links    []*model.ShortLink
	total    int64
	listErr  error
	captured model.ShortLinkListFilter
}

func (f *fakeListModel) List(_ context.Context, filter model.ShortLinkListFilter) ([]*model.ShortLink, int64, error) {
	f.captured = filter
	return f.links, f.total, f.listErr
}

func newListLogic(m model.ShortLinkModel) *ListShortLinksLogic {
	return NewListShortLinksLogic(context.Background(), &svc.ServiceContext{
		ShortLinkModel: m,
	})
}

func TestListShortLinks_DefaultPagination(t *testing.T) {
	t.Parallel()
	m := &fakeListModel{
		links: []*model.ShortLink{activeShortLink(1, "abc123")},
		total: 1,
	}
	resp, err := newListLogic(m).ListShortLinks(&linkv1.ListShortLinksRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Page != 1 {
		t.Fatalf("want page 1, got %d", resp.Page)
	}
	if resp.PageSize != 20 {
		t.Fatalf("want page_size 20, got %d", resp.PageSize)
	}
	if resp.Total != 1 {
		t.Fatalf("want total 1, got %d", resp.Total)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("want 1 item, got %d", len(resp.Items))
	}
}

func TestListShortLinks_StatusFilter(t *testing.T) {
	t.Parallel()
	m := &fakeListModel{links: nil, total: 0}
	resp, err := newListLogic(m).ListShortLinks(&linkv1.ListShortLinksRequest{Status: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.captured.Status != 1 {
		t.Fatalf("want status filter 1, got %d", m.captured.Status)
	}
	if len(resp.Items) != 0 {
		t.Fatalf("want 0 items, got %d", len(resp.Items))
	}
}

func TestListShortLinks_KeywordFilter(t *testing.T) {
	t.Parallel()
	m := &fakeListModel{
		links: []*model.ShortLink{activeShortLink(2, "xyz")},
		total: 1,
	}
	resp, err := newListLogic(m).ListShortLinks(&linkv1.ListShortLinksRequest{Keyword: "xyz"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.captured.Keyword != "xyz" {
		t.Fatalf("want keyword xyz, got %q", m.captured.Keyword)
	}
	if resp.Total != 1 {
		t.Fatalf("want total 1, got %d", resp.Total)
	}
}

func TestListShortLinks_InvalidPagination(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		page     int64
		pageSize int64
	}{
		{name: "negative page", page: -1, pageSize: 10},
		{name: "zero page size after set", page: 1, pageSize: -1},
		{name: "page size exceeds max", page: 1, pageSize: 101},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := newListLogic(&fakeListModel{}).ListShortLinks(&linkv1.ListShortLinksRequest{
				Page:     tt.page,
				PageSize: tt.pageSize,
			})
			if status.Code(err) != codes.InvalidArgument {
				t.Fatalf("want InvalidArgument, got %v", err)
			}
		})
	}
}

func TestListShortLinks_InvalidStatus(t *testing.T) {
	t.Parallel()
	_, err := newListLogic(&fakeListModel{}).ListShortLinks(&linkv1.ListShortLinksRequest{Status: 99})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("want InvalidArgument for unknown status, got %v", err)
	}
}
