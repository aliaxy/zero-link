// Package main starts the link-rpc service.
package main

import (
	"flag"
	"fmt"

	"github.com/aliaxy/zero-link/services/link-api/pkg/configvalidator"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/config"
	adminserver "github.com/aliaxy/zero-link/services/link-rpc/internal/server/adminservice"
	analyticsserver "github.com/aliaxy/zero-link/services/link-rpc/internal/server/analyticsservice"
	healthserver "github.com/aliaxy/zero-link/services/link-rpc/internal/server/healthservice"
	linkserver "github.com/aliaxy/zero-link/services/link-rpc/internal/server/linkservice"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	linkv1 "github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/link-rpc.local.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	configvalidator.MustValidate(c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		linkv1.RegisterHealthServiceServer(grpcServer, healthserver.NewHealthServiceServer(ctx))
		linkv1.RegisterAdminServiceServer(grpcServer, adminserver.NewAdminServiceServer(ctx))
		linkv1.RegisterLinkServiceServer(grpcServer, linkserver.NewLinkServiceServer(ctx))
		linkv1.RegisterAnalyticsServiceServer(grpcServer, analyticsserver.NewAnalyticsServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
