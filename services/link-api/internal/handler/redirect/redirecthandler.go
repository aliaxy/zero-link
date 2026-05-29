// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

// Package redirect contains link-api short-link redirect handlers.
package redirect

import (
	"net/http"

	"github.com/aliaxy/zero-link/services/link-api/internal/logic/redirect"
	"github.com/aliaxy/zero-link/services/link-api/internal/svc"
	"github.com/aliaxy/zero-link/services/link-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// RedirectHandler handles GET /:code.
//
//nolint:revive // goctl convention: type name matches handler name
func RedirectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RedirectRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := redirect.NewRedirectLogic(r.Context(), svcCtx)
		resp, err := l.Redirect(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
