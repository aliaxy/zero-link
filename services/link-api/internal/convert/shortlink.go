package convert

import (
	"github.com/aliaxy/zero-link/services/link-api/internal/types"
	"github.com/aliaxy/zero-link/services/link-rpc/linkservice"
)

// ShortLinkInfoFromRPC converts an RPC ShortLink to an API ShortLinkInfo.
func ShortLinkInfoFromRPC(link *linkservice.ShortLink) types.ShortLinkInfo {
	if link == nil {
		return types.ShortLinkInfo{}
	}
	return types.ShortLinkInfo{
		Id:          link.Id,
		Code:        link.Code,
		OriginUrl:   link.OriginUrl,
		Title:       link.Title,
		Description: link.Description,
		Status:      link.Status,
		ExpireAt:    link.ExpireAt,
		CreatedBy:   link.CreatedBy,
		CreatedAt:   link.CreatedAt,
		UpdatedAt:   link.UpdatedAt,
	}
}

// ShortLinkSummaryFromRPC converts an RPC ShortLinkSummary to an API ShortLinkSummary.
func ShortLinkSummaryFromRPC(link *linkservice.ShortLinkSummary) types.ShortLinkSummary {
	if link == nil {
		return types.ShortLinkSummary{}
	}
	return types.ShortLinkSummary{
		Id:        link.Id,
		Code:      link.Code,
		OriginUrl: link.OriginUrl,
		Title:     link.Title,
		Status:    link.Status,
		ExpireAt:  link.ExpireAt,
		CreatedAt: link.CreatedAt,
		UpdatedAt: link.UpdatedAt,
	}
}

// OkShortLinkResponse builds a success ShortLinkResponse from an RPC ShortLink.
func OkShortLinkResponse(link *linkservice.ShortLink) *types.ShortLinkResponse {
	return &types.ShortLinkResponse{
		Code:    "OK",
		Message: "ok",
		Data:    ShortLinkInfoFromRPC(link),
	}
}
