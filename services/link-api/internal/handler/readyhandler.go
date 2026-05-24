// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"github.com/aliaxy/zero-link/services/link-api/internal/logic"
	"github.com/aliaxy/zero-link/services/link-api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// ReadyHandler returns the API readiness handler.
func ReadyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewReadyLogic(r.Context(), svcCtx)
		resp, err := l.Ready()
		if err != nil {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusServiceUnavailable, map[string]string{
				"status":  "unavailable",
				"message": err.Error(),
			})
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
