package logic

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/domain"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/model"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	linkv1 "github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

// CreateShortLinkLogic coordinates short-link creation.
type CreateShortLinkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewCreateShortLinkLogic creates short-link creation logic.
func NewCreateShortLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateShortLinkLogic {
	return &CreateShortLinkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateShortLink creates a managed short link.
func (l *CreateShortLinkLogic) CreateShortLink(
	in *linkv1.CreateShortLinkRequest,
) (*linkv1.CreateShortLinkResponse, error) {
	if err := domain.ValidateOriginURL(in.GetOriginUrl()); err != nil {
		return nil, rpcError(domain.ErrInvalidArgument)
	}

	code := in.GetCode()
	if code != "" {
		if err := domain.ValidateCustomCode(code); err != nil {
			return nil, rpcError(domain.ErrInvalidArgument)
		}
		if _, err := l.svcCtx.ShortLinkModel.FindOneByCode(l.ctx, code); err == nil {
			return nil, rpcError(domain.ErrConflict)
		} else if !errors.Is(err, model.ErrNotFound) {
			return nil, rpcError(err)
		}
	}
	if code == "" {
		generated, err := domain.GenerateCode()
		if err != nil {
			return nil, rpcError(err)
		}
		code = generated
	}

	expireAt, err := nullTimeFromString(in.GetExpireAt(), time.Now().UTC())
	if err != nil {
		return nil, rpcError(domain.ErrInvalidArgument)
	}

	data := &model.ShortLink{
		Code:        code,
		OriginUrl:   in.GetOriginUrl(),
		Title:       in.GetTitle(),
		Description: in.GetDescription(),
		Status:      domain.LinkStatusActive,
		ExpireAt:    expireAt,
		CreatedBy:   in.GetCreatedBy(),
		DeletedAt:   sql.NullTime{},
	}
	result, err := l.svcCtx.ShortLinkModel.Insert(l.ctx, data)
	if err != nil {
		return nil, rpcError(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, rpcError(err)
	}
	link, err := l.svcCtx.ShortLinkModel.FindOneNotDeleted(l.ctx, id)
	if err != nil {
		return nil, rpcError(modelError(err))
	}

	return &linkv1.CreateShortLinkResponse{
		Link: shortLinkFromModel(link),
	}, nil
}
