// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"

	"github.com/aliaxy/zero-link/services/link-api/internal/apierror"
	"github.com/aliaxy/zero-link/services/link-api/internal/convert"
	"github.com/aliaxy/zero-link/services/link-api/internal/middleware"
	"github.com/aliaxy/zero-link/services/link-api/internal/svc"
	"github.com/aliaxy/zero-link/services/link-api/internal/types"
	"github.com/aliaxy/zero-link/services/link-rpc/client/linkservice"

	"github.com/zeromicro/go-zero/core/logx"
)

// CreateShortLinkLogic coordinates management short-link creation.
type CreateShortLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateShortLinkLogic creates short-link creation logic.
func NewCreateShortLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateShortLinkLogic {
	return &CreateShortLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateShortLink creates a short link through link-rpc.
func (l *CreateShortLinkLogic) CreateShortLink(
	req *types.CreateShortLinkRequest,
) (resp *types.ShortLinkResponse, err error) {
	subject, ok := middleware.AdminSubjectFromContext(l.ctx)
	if !ok {
		return nil, apierror.ErrUnauthenticated
	}

	rpcResp, err := l.svcCtx.LinkRPC.CreateShortLink(l.ctx, &linkservice.CreateShortLinkRequest{
		OriginUrl:   req.OriginUrl,
		Code:        req.Code,
		Title:       req.Title,
		Description: req.Description,
		ExpireAt:    req.ExpireAt,
		CreatedBy:   subject.ID,
	})
	if err != nil {
		return nil, apierror.FromRPCError(err)
	}

	return convert.OkShortLinkResponse(rpcResp.Link), nil
}
