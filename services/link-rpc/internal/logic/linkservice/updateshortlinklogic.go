package linkservicelogic

import (
	"context"
	"time"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/domain"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/rpcerror"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	"github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

// UpdateShortLinkLogic coordinates short-link updates.
type UpdateShortLinkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewUpdateShortLinkLogic creates short-link update logic.
func NewUpdateShortLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateShortLinkLogic {
	return &UpdateShortLinkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateShortLink updates mutable fields on a managed short link.
func (l *UpdateShortLinkLogic) UpdateShortLink(
	in *linkv1.UpdateShortLinkRequest,
) (*linkv1.UpdateShortLinkResponse, error) {
	if in.GetId() <= 0 {
		return nil, rpcerror.ToRPC(domain.ErrInvalidArgument)
	}
	if in.GetOriginUrl() != "" {
		if err := domain.ValidateOriginURL(in.GetOriginUrl()); err != nil {
			return nil, rpcerror.ToRPC(domain.ErrInvalidArgument)
		}
	}
	if in.GetStatus() > 0 {
		if err := domain.ValidateLinkStatus(in.GetStatus()); err != nil {
			return nil, rpcerror.ToRPC(domain.ErrInvalidArgument)
		}
	}
	expireAt, err := nullTimeFromString(in.GetExpireAt(), time.Now().UTC())
	if err != nil {
		return nil, rpcerror.ToRPC(domain.ErrInvalidArgument)
	}

	link, err := l.svcCtx.ShortLinkModel.FindOneNotDeleted(l.ctx, in.GetId())
	if err != nil {
		return nil, rpcerror.ToRPC(modelError(err))
	}
	if in.GetOriginUrl() != "" {
		link.OriginUrl = in.GetOriginUrl()
	}
	if in.GetTitle() != "" {
		link.Title = in.GetTitle()
	}
	if in.GetDescription() != "" {
		link.Description = in.GetDescription()
	}
	if in.GetStatus() > 0 {
		link.Status = in.GetStatus()
	}
	if in.GetExpireAt() != "" {
		link.ExpireAt = expireAt
	}
	if err := l.svcCtx.ShortLinkModel.Update(l.ctx, link); err != nil {
		return nil, rpcerror.ToRPC(err)
	}

	link, err = l.svcCtx.ShortLinkModel.FindOneNotDeleted(l.ctx, in.GetId())
	if err != nil {
		return nil, rpcerror.ToRPC(modelError(err))
	}

	return &linkv1.UpdateShortLinkResponse{
		Link: shortLinkFromModel(link),
	}, nil
}
