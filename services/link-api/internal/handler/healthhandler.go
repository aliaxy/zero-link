// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

// Package handler defines link-api HTTP handlers.
package handler

import (
	"net/http"

	"github.com/aliaxy/zero-link/services/link-api/internal/logic"
	"github.com/aliaxy/zero-link/services/link-api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// HealthHandler returns the API process liveness handler.
func HealthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewHealthLogic(r.Context(), svcCtx)
		resp, err := l.Health()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
