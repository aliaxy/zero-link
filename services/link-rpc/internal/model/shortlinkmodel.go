package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ShortLinkModel = (*customShortLinkModel)(nil)

type (
	// ShortLinkModel is an interface to be customized, add more methods here,
	// and implement the added methods in customShortLinkModel.
	ShortLinkModel interface {
		shortLinkModel
		withSession(session sqlx.Session) ShortLinkModel
	}

	customShortLinkModel struct {
		*defaultShortLinkModel
	}
)

// NewShortLinkModel returns a model for the database table.
func NewShortLinkModel(conn sqlx.SqlConn) ShortLinkModel {
	return &customShortLinkModel{
		defaultShortLinkModel: newShortLinkModel(conn),
	}
}

func (m *customShortLinkModel) withSession(session sqlx.Session) ShortLinkModel {
	return NewShortLinkModel(sqlx.NewSqlConnFromSession(session))
}
