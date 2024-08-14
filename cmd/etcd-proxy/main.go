/*
 * Copyright text:
 * This file was last modified at 2024-07-10 21:53 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * main.go
 * $Id$
 */
//!+

package main

import (
	"context"
	"crypto/tls"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/victor-skurikhin/etcd-client/v1/internal/alog"
	"github.com/victor-skurikhin/etcd-client/v1/internal/controllers"
	"github.com/victor-skurikhin/etcd-client/v1/internal/controllers/dto"
	"github.com/victor-skurikhin/etcd-client/v1/internal/env"
	"github.com/victor-skurikhin/etcd-client/v1/internal/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	pb "github.com/victor-skurikhin/etcd-client/v1/proto"
	_ "google.golang.org/grpc/encoding/gzip"
)

const MSG = "etcd-proxy"

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
	sLog         *slog.Logger
)

func main() {
	run(context.Background())
}

func run(ctx context.Context) {
	slog.Info(MSG+"meta info",
		"build_version", buildVersion,
		"build_date", buildDate,
		"build_commit", buildCommit,
	)
	serve(ctx, env.GetConfig())
}

func serve(ctx context.Context, cfg env.Config) {

	sLog = cfg.Logger()
	listen, err := net.Listen("tcp", cfg.GRPCAddress())

	if err != nil {
		log.Fatal(err)
	}
	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	httpServer := makeHTTP(ctx, cfg)
	grpcServer := makeGRPC(ctx, cfg)

	go func() {
		<-sigint
		grpcServer.GracefulStop()
		sLog.Info(MSG+"graceful stop", "msg", "Выключение сервера gRPC")
		if err := httpServer.Shutdown(); err != nil {
			sLog.Error(MSG+"graceful stop", "msg", "Ошибка при выключение сервера HTTP", "err", err)
		}
		sLog.Info(MSG+"graceful stop", "msg", "Выключение сервера HTTP")
		close(idleConnsClosed)
	}()
	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
		sLog.Info(MSG+"graceful stop", "msg", "Выключение сервера gRPC")
		if err := httpServer.Shutdown(); err != nil {
			sLog.Error(MSG+"graceful stop", "msg", "Ошибка при выключение сервера HTTP", "err", err)
		}
		sLog.Info(MSG+"graceful stop", "msg", "Выключение сервера HTTP")
		close(idleConnsClosed)
	}()
	go func() {
		sLog.Info(MSG+"start app", "msg", "Сервер gRPC начал работу")
		if err := grpcServer.Serve(listen); err != nil {
			log.Fatal(err)
		}
	}()
	sLog.Info(MSG+"start app", "msg", "Сервер HTTP начал работу")
	if cfg.YamlConfig().HTTPTLSEnabled() {

		ln, err := tls.Listen("tcp", cfg.HTTPAddress(), cfg.HTTPTLSConfig())
		if err != nil {
			panic(err)
		}
		if err := httpServer.Listener(ln); err != nil {
			sLog.Error(MSG+"start app", "msg", "Ошибка при выключение сервера HTTP", "err", err)
		}
	} else if err := httpServer.Listen(cfg.HTTPAddress()); err != nil {
		sLog.Error(MSG+"start app", "msg", "Ошибка при выключение сервера HTTP", "err", err)
	}
	<-idleConnsClosed
	sLog.Info(MSG+"shutdown app", "msg", "Корректное завершение работы сервера")
}

func makeHTTP(ctx context.Context, cfg env.Config) *fiber.App {

	logHandler := logger.New(logger.Config{
		Format:       "${pid} | ${time} | ${status} | ${locals:requestid} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
		TimeFormat:   "15:04:05.000000",
		TimeZone:     "Local",
		TimeInterval: 500 * time.Millisecond,
		Output:       os.Stdout,
	})
	slogLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	app := fiber.New(fiber.Config{DisableHeaderNormalizing: true})
	micro := fiber.New()
	app.Mount("/api", micro)
	app.Use(requestid.New())
	micro.Use(requestid.New())

	if cfg.SlogJSON() {
		app.Use(alog.New(slogLogger))
		micro.Use(alog.New(slogLogger))
	} else {
		app.Use(logHandler)
		micro.Use(logHandler)
	}
	ctrl := controllers.GetEtcdProxyController(ctx, cfg)
	micro.Delete("/delete/:name", ctrl.Delete)
	micro.Get("/get/:name", ctrl.Get)
	micro.Put("/put/:name", ctrl.Put)
	micro.All("*", func(c *fiber.Ctx) error {
		path := c.Path()
		return c.
			Status(fiber.StatusNotFound).
			JSON(dto.StatusMessagePathDoesNotExists(path))
	})
	return app
}

func makeGRPC(ctx context.Context, cfg env.Config) *grpc.Server {

	var opts []grpc.ServerOption

	if cfg.YamlConfig().GRPCTLSEnabled() {
		opts = []grpc.ServerOption{
			grpc.Creds(cfg.GRPCTransportCredentials()),
		}
	} else {
		opts = []grpc.ServerOption{
			grpc.Creds(insecure.NewCredentials()),
		}
	}
	srv := services.GetEtcdProxyService(ctx, cfg)
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterEtcdClientServiceServer(grpcServer, srv)
	reflection.Register(grpcServer)

	return grpcServer
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
