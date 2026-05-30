package linkservicelogic

import (
	"context"
	"time"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/domain"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	"github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

// DeleteShortLinkLogic coordinates short-link soft deletion.
type DeleteShortLinkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewDeleteShortLinkLogic creates short-link deletion logic.
func NewDeleteShortLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteShortLinkLogic {
	return &DeleteShortLinkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteShortLink soft deletes a managed short link.
func (l *DeleteShortLinkLogic) DeleteShortLink(
	in *linkv1.DeleteShortLinkRequest,
) (*linkv1.DeleteShortLinkResponse, error) {
	if in.GetId() <= 0 {
		return nil, rpcError(domain.ErrInvalidArgument)
	}
	if err := l.svcCtx.ShortLinkModel.SoftDelete(l.ctx, in.GetId(), time.Now().UTC()); err != nil {
		return nil, rpcError(modelError(err))
	}

	return &linkv1.DeleteShortLinkResponse{
		Id:      in.GetId(),
		Deleted: true,
	}, nil
}
