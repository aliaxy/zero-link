package logic

import (
	"context"
	"strings"
	"testing"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/config"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	linkv1 "github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"
)

func TestCheckLogic_Check(t *testing.T) {
	tests := []struct {
		name    string
		mysql   string
		redis   string
		wantOK  bool
		wantMsg string
	}{
		{
			name:    "mysql endpoint empty",
			mysql:   "",
			redis:   "127.0.0.1:6379",
			wantOK:  false,
			wantMsg: "mysql endpoint is empty",
		},
		{
			name:    "redis endpoint empty",
			mysql:   "127.0.0.1:3306",
			redis:   "",
			wantOK:  false,
			wantMsg: "redis endpoint is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logic := NewCheckLogic(context.Background(), svc.NewServiceContext(config.Config{
				Dependencies: config.DependenciesConfig{
					MySQL: config.MySQLConfig{Endpoint: tt.mysql},
					Redis: config.RedisConfig{Endpoint: tt.redis},
				},
			}))

			got, err := logic.Check(&linkv1.CheckRequest{})
			if err != nil {
				t.Fatalf("Check() error = %v", err)
			}
			if got.Ok != tt.wantOK {
				t.Fatalf("Check().Ok = %v, want %v", got.Ok, tt.wantOK)
			}
			if !strings.Contains(got.Message, tt.wantMsg) {
				t.Fatalf("Check().Message = %q, want to contain %q", got.Message, tt.wantMsg)
			}
		})
	}
}
