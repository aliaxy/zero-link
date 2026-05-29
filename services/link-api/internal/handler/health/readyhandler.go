// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package health

import (
	"net/http"

	"github.com/aliaxy/zero-link/services/link-api/internal/logic/health"
	"github.com/aliaxy/zero-link/services/link-api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// ReadyHandler handles GET /readyz.
func ReadyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := health.NewReadyLogic(r.Context(), svcCtx)
		resp, err := l.Ready()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
