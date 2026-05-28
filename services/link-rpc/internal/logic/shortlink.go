package logic

import (
	"database/sql"
	"errors"
	"time"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/domain"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/model"
	linkv1 "github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"
)

func shortLinkFromModel(link *model.ShortLink) *linkv1.ShortLink {
	if link == nil {
		return nil
	}
	return &linkv1.ShortLink{
		Id:          link.Id,
		Code:        link.Code,
		OriginUrl:   link.OriginUrl,
		Title:       link.Title,
		Description: link.Description,
		Status:      link.Status,
		ExpireAt:    nullTimeString(link.ExpireAt),
		CreatedBy:   link.CreatedBy,
		CreatedAt:   link.CreatedAt.Format(timeFormat),
		UpdatedAt:   link.UpdatedAt.Format(timeFormat),
	}
}

func shortLinkSummaryFromModel(link *model.ShortLink) *linkv1.ShortLinkSummary {
	if link == nil {
		return nil
	}
	return &linkv1.ShortLinkSummary{
		Id:        link.Id,
		Code:      link.Code,
		OriginUrl: link.OriginUrl,
		Title:     link.Title,
		Status:    link.Status,
		ExpireAt:  nullTimeString(link.ExpireAt),
		CreatedAt: link.CreatedAt.Format(timeFormat),
		UpdatedAt: link.UpdatedAt.Format(timeFormat),
	}
}

func nullTimeString(value sql.NullTime) string {
	if !value.Valid {
		return ""
	}
	return value.Time.Format(timeFormat)
}

func nullTimeFromString(value string, now time.Time) (sql.NullTime, error) {
	expireAt, err := domain.ValidateExpireAt(value, now)
	if err != nil {
		return sql.NullTime{}, err
	}
	if expireAt.IsZero() {
		return sql.NullTime{}, nil
	}
	return sql.NullTime{Time: expireAt, Valid: true}, nil
}

func modelError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, model.ErrNotFound) {
		return domain.ErrNotFound
	}
	return err
}
