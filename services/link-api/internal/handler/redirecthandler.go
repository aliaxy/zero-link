// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"github.com/aliaxy/zero-link/services/link-api/internal/logic"
	"github.com/aliaxy/zero-link/services/link-api/internal/response"
	"github.com/aliaxy/zero-link/services/link-api/internal/svc"
	"github.com/aliaxy/zero-link/services/link-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// RedirectHandler handles short-link redirect requests.
func RedirectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RedirectRequest
		if err := httpx.Parse(r, &req); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		l := logic.NewRedirectLogic(r.Context(), svcCtx)
		resp, err := l.Redirect(&req)
		if err != nil {
			response.RedirectError(w, err)
			return
		}
		http.Redirect(w, r, resp.OriginUrl, http.StatusFound)
	}
}
