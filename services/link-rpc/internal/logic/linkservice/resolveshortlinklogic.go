package linkservicelogic

import (
	"context"
	"time"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/domain"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/model"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	linkv1 "github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

// ResolveShortLinkLogic handles short-link resolution RPC.
type ResolveShortLinkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewResolveShortLinkLogic creates a ResolveShortLinkLogic.
func NewResolveShortLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResolveShortLinkLogic {
	return &ResolveShortLinkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ResolveShortLink resolves a short code to its origin URL.
func (l *ResolveShortLinkLogic) ResolveShortLink(
	in *linkv1.ResolveShortLinkRequest,
) (*linkv1.ResolveShortLinkResponse, error) {
	link, err := l.svcCtx.ShortLinkModel.FindOneByCode(l.ctx, in.Code)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, rpcError(domain.ErrNotFound)
		}
		return nil, rpcError(err)
	}

	if link.DeletedAt.Valid {
		return nil, rpcError(domain.ErrNotFound)
	}

	if link.Status == domain.LinkStatusDisabled {
		return nil, rpcError(domain.ErrPermissionDenied)
	}

	if link.ExpireAt.Valid && link.ExpireAt.Time.Before(time.Now()) {
		return nil, rpcError(domain.ErrGone)
	}

	return &linkv1.ResolveShortLinkResponse{OriginUrl: link.OriginUrl}, nil
}
