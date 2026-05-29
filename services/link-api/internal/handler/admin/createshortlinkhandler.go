// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

// Package admin contains link-api administrator management handlers.
package admin

import (
	"net/http"

	"github.com/aliaxy/zero-link/services/link-api/internal/logic/admin"
	"github.com/aliaxy/zero-link/services/link-api/internal/svc"
	"github.com/aliaxy/zero-link/services/link-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// CreateShortLinkHandler handles POST /admin/links.
func CreateShortLinkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateShortLinkRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := admin.NewCreateShortLinkLogic(r.Context(), svcCtx)
		resp, err := l.CreateShortLink(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
