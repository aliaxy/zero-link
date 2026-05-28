package logic

import "github.com/aliaxy/zero-link/services/link-api/internal/response"

func errUnauthenticated() error {
	return response.ErrUnauthenticated
}
