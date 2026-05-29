// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"net/http"

	"github.com/aliaxy/zero-link/services/link-api/internal/logic/admin"
	"github.com/aliaxy/zero-link/services/link-api/internal/svc"
	"github.com/aliaxy/zero-link/services/link-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// ListShortLinksHandler handles GET /admin/links.
func ListShortLinksHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListShortLinksRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := admin.NewListShortLinksLogic(r.Context(), svcCtx)
		resp, err := l.ListShortLinks(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
