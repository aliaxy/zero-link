package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

// ErrNotFound is returned when a model query does not find a row.
var ErrNotFound = sqlx.ErrNotFound
