package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/hatlonely/go-kit/binding"
	"github.com/hatlonely/go-kit/cli"
	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/flag"
	"github.com/hatlonely/go-kit/logger"
	"github.com/hatlonely/go-kit/rpcx"
	sycmApi "github.com/realwrtoff/rest_grpc/sycm/api/gen/go/api"
	"github.com/realwrtoff/rest_grpc/sycm/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var Version string

type Options struct {
	flag.Options

	Http struct {
		Port int
	}
	Grpc struct {
		Port int
	}

	Redis   cli.RedisOptions
	Mysql   cli.MySQLOptions
	Sycm 	service.SycmServiceOptions
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var options Options
	Must(flag.Struct(&options))
	Must(flag.Parse())
	if options.Help {
		fmt.Println(flag.Usage())
		return
	}
	if options.Version {
		fmt.Println(Version)
		return
	}

	if options.ConfigPath == "" {
		options.ConfigPath = "config/sycm_server.json"
	}
	cfg, err := config.NewSimpleFileConfig(options.ConfigPath)
	Must(err)
	Must(binding.Bind(&options, flag.Instance(), binding.NewEnvGetter(binding.WithEnvPrefix("SYCM")), cfg))

	grpcLog, err := logger.NewLoggerWithConfig(cfg.Sub("logger.grpc"))
	Must(err)
	infoLog, err := logger.NewLoggerWithConfig(cfg.Sub("logger.info"))
	Must(err)

	redisCli, err := cli.NewRedisWithOptions(&options.Redis)
	Must(err)
	mysqlCli, err := cli.NewMysqlWithOptions(&options.Mysql)
	Must(err)

	svc, err := service.NewSycmService(
		mysqlCli, redisCli, infoLog,
		service.WithTokenExpiration(options.Sycm.TokenExpiration),
		service.WithTokenMaxRequest(options.Sycm.TokenMaxRequest),
	)
	Must(err)

	rpcServer := grpc.NewServer(
		rpcx.WithGrpcDecorator(grpcLog),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
			PermitWithoutStream: true,            // Allow pings even when there are no active streams
		}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
			MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
			MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
			Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
			Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
		}),
	)
	sycmApi.RegisterSycmServiceServer(rpcServer, svc)

	go func() {
		address, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", options.Grpc.Port))
		Must(err)
		Must(rpcServer.Serve(address))
	}()

	muxServer := runtime.NewServeMux(
		rpcx.WithMuxMetadata(),
		rpcx.WithMuxIncomingHeaderMatcher(),
		rpcx.WithMuxOutgoingHeaderMatcher(),
		rpcx.WithMuxProtoErrorHandler(),
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Must(sycmApi.RegisterSycmServiceHandlerFromEndpoint(
		ctx, muxServer, fmt.Sprintf("0.0.0.0:%v", options.Grpc.Port), []grpc.DialOption{grpc.WithInsecure()},
	))
	infoLog.Info(options)
	Must(http.ListenAndServe(fmt.Sprintf(":%v", options.Http.Port), handlers.CombinedLoggingHandler(os.Stdout, muxServer)))
}
