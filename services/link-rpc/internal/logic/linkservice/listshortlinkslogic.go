package linkservicelogic

import (
	"context"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/domain"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/model"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/rpcerror"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	"github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

// ListShortLinksLogic coordinates short-link listing.
type ListShortLinksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewListShortLinksLogic creates short-link listing logic.
func NewListShortLinksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListShortLinksLogic {
	return &ListShortLinksLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListShortLinks returns paginated non-deleted short links.
func (l *ListShortLinksLogic) ListShortLinks(in *linkv1.ListShortLinksRequest) (*linkv1.ListShortLinksResponse, error) {
	page, pageSize, err := domain.NormalizePagination(in.GetPage(), in.GetPageSize())
	if err != nil {
		return nil, rpcerror.ToRPC(domain.ErrInvalidArgument)
	}
	if in.GetStatus() > 0 {
		if err := domain.ValidateLinkStatus(in.GetStatus()); err != nil {
			return nil, rpcerror.ToRPC(domain.ErrInvalidArgument)
		}
	}

	links, total, err := l.svcCtx.ShortLinkModel.List(l.ctx, model.ShortLinkListFilter{
		Page:     page,
		PageSize: pageSize,
		Status:   in.GetStatus(),
		Keyword:  in.GetKeyword(),
	})
	if err != nil {
		return nil, rpcerror.ToRPC(err)
	}

	items := make([]*linkv1.ShortLinkSummary, 0, len(links))
	for _, link := range links {
		items = append(items, shortLinkSummaryFromModel(link))
	}

	return &linkv1.ListShortLinksResponse{
		Items:    items,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}
