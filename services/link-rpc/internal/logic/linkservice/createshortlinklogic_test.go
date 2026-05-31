package linkservicelogic

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/model"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	linkv1 "github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// fakeCreateModel supports FindOneByCode, Insert, and FindOneNotDeleted.
type fakeCreateModel struct {
	model.ShortLinkModel

	findByCodeLink *model.ShortLink
	findByCodeErr  error

	insertedID int64
	insertErr  error

	findNotDeletedLink *model.ShortLink
	findNotDeletedErr  error
}

func (f *fakeCreateModel) FindOneByCode(_ context.Context, _ string) (*model.ShortLink, error) {
	return f.findByCodeLink, f.findByCodeErr
}

func (f *fakeCreateModel) Insert(_ context.Context, _ *model.ShortLink) (sql.Result, error) {
	if f.insertErr != nil {
		return nil, f.insertErr
	}
	return fakeResult{id: f.insertedID}, nil
}

func (f *fakeCreateModel) FindOneNotDeleted(_ context.Context, _ int64) (*model.ShortLink, error) {
	return f.findNotDeletedLink, f.findNotDeletedErr
}

type fakeResult struct{ id int64 }

func (r fakeResult) LastInsertId() (int64, error) { return r.id, nil }
func (r fakeResult) RowsAffected() (int64, error) { return 1, nil }

func activeShortLink(id int64, code string) *model.ShortLink {
	return &model.ShortLink{
		Id:        id,
		Code:      code,
		OriginUrl: "https://example.com",
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func newCreateLogic(m model.ShortLinkModel) *CreateShortLinkLogic {
	return NewCreateShortLinkLogic(context.Background(), &svc.ServiceContext{
		ShortLinkModel: m,
	})
}

func TestCreateShortLink_InvalidURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		url  string
	}{
		{name: "empty", url: ""},
		{name: "relative", url: "/page"},
		{name: "ftp scheme", url: "ftp://example.com"},
		{name: "malformed", url: "://bad"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			logic := newCreateLogic(&fakeCreateModel{})
			_, err := logic.CreateShortLink(&linkv1.CreateShortLinkRequest{OriginUrl: tt.url, Code: "mycode"})
			if status.Code(err) != codes.InvalidArgument {
				t.Fatalf("want InvalidArgument, got %v", err)
			}
		})
	}
}

func TestCreateShortLink_InvalidCustomCode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		code string
	}{
		{name: "too short", code: "ab"},
		{name: "too long", code: "aaaaaaaaaabbbbbbbbbbccccccccccddddd"},
		{name: "space", code: "bad code"},
		{name: "reserved admin", code: "admin"},
		{name: "slash", code: "bad/code"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			logic := newCreateLogic(&fakeCreateModel{})
			_, err := logic.CreateShortLink(&linkv1.CreateShortLinkRequest{
				OriginUrl: "https://example.com",
				Code:      tt.code,
			})
			if status.Code(err) != codes.InvalidArgument {
				t.Fatalf("want InvalidArgument, got %v", err)
			}
		})
	}
}

func TestCreateShortLink_DuplicateCustomCode(t *testing.T) {
	t.Parallel()
	m := &fakeCreateModel{
		findByCodeLink: activeShortLink(1, "mycode"),
		findByCodeErr:  nil,
	}
	logic := newCreateLogic(m)
	_, err := logic.CreateShortLink(&linkv1.CreateShortLinkRequest{
		OriginUrl: "https://example.com",
		Code:      "mycode",
	})
	if status.Code(err) != codes.AlreadyExists {
		t.Fatalf("want AlreadyExists (conflict), got %v", err)
	}
}

func TestCreateShortLink_ValidCustomCode(t *testing.T) {
	t.Parallel()
	m := &fakeCreateModel{
		findByCodeErr:      model.ErrNotFound,
		insertedID:         7,
		findNotDeletedLink: activeShortLink(7, "mycode"),
	}
	logic := newCreateLogic(m)
	resp, err := logic.CreateShortLink(&linkv1.CreateShortLinkRequest{
		OriginUrl: "https://example.com",
		Code:      "mycode",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Link.Code != "mycode" {
		t.Fatalf("want code mycode, got %q", resp.Link.Code)
	}
}

func TestCreateShortLink_GeneratedCode(t *testing.T) {
	t.Parallel()
	m := &fakeCreateModel{
		insertedID:         3,
		findNotDeletedLink: activeShortLink(3, "ABC123"),
	}
	logic := newCreateLogic(m)
	resp, err := logic.CreateShortLink(&linkv1.CreateShortLinkRequest{
		OriginUrl: "https://example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Link.Id != 3 {
		t.Fatalf("want id 3, got %d", resp.Link.Id)
	}
}

func TestCreateShortLink_ModelInsertError(t *testing.T) {
	t.Parallel()
	m := &fakeCreateModel{
		findByCodeErr: model.ErrNotFound,
		insertErr:     errors.New("db error"),
	}
	logic := newCreateLogic(m)
	_, err := logic.CreateShortLink(&linkv1.CreateShortLinkRequest{
		OriginUrl: "https://example.com",
		Code:      "mycode",
	})
	if status.Code(err) != codes.Internal {
		t.Fatalf("want Internal, got %v", err)
	}
}
