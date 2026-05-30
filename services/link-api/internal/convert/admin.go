// Package convert translates link-rpc protobuf types to link-api types.
package convert

import (
	"github.com/aliaxy/zero-link/services/link-api/internal/auth"
	"github.com/aliaxy/zero-link/services/link-api/internal/types"
	"github.com/aliaxy/zero-link/services/link-rpc/client/adminservice"
)

// AdminSubjectFromRPC converts an RPC AdminProfile to an auth subject.
func AdminSubjectFromRPC(admin *adminservice.AdminProfile) auth.AdminSubject {
	if admin == nil {
		return auth.AdminSubject{}
	}
	return auth.AdminSubject{
		ID:       admin.Id,
		Username: admin.Username,
	}
}

// AdminInfoFromRPC converts an RPC AdminProfile to an API AdminInfo.
func AdminInfoFromRPC(admin *adminservice.AdminProfile) types.AdminInfo {
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
