package linkservicelogic

import (
	"context"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/domain"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/rpcerror"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	"github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetShortLinkLogic coordinates short-link detail retrieval.
type GetShortLinkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewGetShortLinkLogic creates short-link detail logic.
func NewGetShortLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShortLinkLogic {
	return &GetShortLinkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetShortLink returns one non-deleted short link.
func (l *GetShortLinkLogic) GetShortLink(in *linkv1.GetShortLinkRequest) (*linkv1.GetShortLinkResponse, error) {
	if in.GetId() <= 0 {
		return nil, rpcerror.ToRPC(domain.ErrInvalidArgument)
	}

	link, err := l.svcCtx.ShortLinkModel.FindOneNotDeleted(l.ctx, in.GetId())
	if err != nil {
		return nil, rpcerror.ToRPC(modelError(err))
	}

	return &linkv1.GetShortLinkResponse{
		Link: shortLinkFromModel(link),
	}, nil
}
