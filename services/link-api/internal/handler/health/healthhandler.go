// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

// Package health contains link-api liveness and readiness handlers.
package health

import (
	"net/http"

	"github.com/aliaxy/zero-link/services/link-api/internal/logic/health"
	"github.com/aliaxy/zero-link/services/link-api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// HealthHandler handles GET /healthz.
//
//nolint:revive // goctl convention: type name matches handler name
func HealthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := health.NewHealthLogic(r.Context(), svcCtx)
		resp, err := l.Health()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
