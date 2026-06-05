// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

// Package main starts the link-api HTTP service.
package main

import (
	"flag"
	"fmt"

	"github.com/aliaxy/zero-link/services/link-api/internal/apierror"
	"github.com/aliaxy/zero-link/services/link-api/internal/config"
	"github.com/aliaxy/zero-link/services/link-api/internal/handler"
	"github.com/aliaxy/zero-link/services/link-api/internal/middleware"
	"github.com/aliaxy/zero-link/services/link-api/internal/svc"
	"github.com/aliaxy/zero-link/services/link-api/pkg/configvalidator"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/link-api.local.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	configvalidator.MustValidate(c)
	httpx.SetErrorHandlerCtx(apierror.ErrorHandler)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	server.Use(middleware.NewCorsMiddleware(c.Cors.AllowOrigins).Handle)
	server.Use(middleware.SecurityHeadersMiddleware)

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
