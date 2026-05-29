package logic

import (
	"time"

	"github.com/aliaxy/zero-link/services/link-api/internal/auth"
	"github.com/aliaxy/zero-link/services/link-api/internal/types"
	"github.com/aliaxy/zero-link/services/link-rpc/linkservice"
)

const timeFormat = time.RFC3339

func adminSubjectFromRPC(admin *linkservice.AdminProfile) auth.AdminSubject {
	if admin == nil {
		return auth.AdminSubject{}
	}
	return auth.AdminSubject{
		ID:       admin.Id,
		Username: admin.Username,
	}
}

func adminInfoFromRPC(admin *linkservice.AdminProfile) types.AdminInfo {
	if admin == nil {
		return types.AdminInfo{}
	}
	return types.AdminInfo{
		Id:        admin.Id,
		Username:  admin.Username,
		Status:    admin.Status,
		CreatedAt: admin.CreatedAt,
	}
}
