package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ VisitEventModel = (*customVisitEventModel)(nil)

type (
	// VisitEventModel is an interface to be customized, add more methods here,
	// and implement the added methods in customVisitEventModel.
	VisitEventModel interface {
		visitEventModel
		withSession(session sqlx.Session) VisitEventModel
	}

	customVisitEventModel struct {
		*defaultVisitEventModel
	}
)

// NewVisitEventModel returns a model for the database table.
func NewVisitEventModel(conn sqlx.SqlConn) VisitEventModel {
	return &customVisitEventModel{
		defaultVisitEventModel: newVisitEventModel(conn),
	}
}

func (m *customVisitEventModel) withSession(session sqlx.Session) VisitEventModel {
	return NewVisitEventModel(sqlx.NewSqlConnFromSession(session))
}
