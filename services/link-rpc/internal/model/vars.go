package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

// ErrNotFound is returned when a model query does not find a row.
var ErrNotFound = sqlx.ErrNotFound

// Cache key prefixes used by the goctl-generated ShortLink model.
// Exported so that packages outside the model layer (e.g. cleanup) can
// invalidate the same keys without importing the generated file directly.
const (
	ShortLinkCacheIDPrefix   = "cache:shortLink:id:"
	ShortLinkCacheCodePrefix = "cache:shortLink:code:"
)
