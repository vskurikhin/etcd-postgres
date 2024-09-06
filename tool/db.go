/*
 * This file was last modified at 2024-07-31 15:59 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * db.go
 * $Id$
 */
//!+

// Package tool TODO.
package tool

import (
	"context"
	"time"

	"github.com/victor-skurikhin/etcd-client/v1/internal/alog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DBConnect(dsn string) *pgxpool.Pool {

	config, err := pgxpool.ParseConfig(dsn)
	IfErrorThenPanic(err)
	sLog.Debug(MSG+"DBConnect", "config", "parsed")

	go func() {
		for alog.GetLogger() == nil {
			time.Sleep(time.Second)
		}
		sLog = alog.GetLogger()
	}()

	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		sLog.Debug(MSG+"DBConnect", "acquire", "connect ping...")
		if err = conn.Ping(ctx); err != nil {
			panic(err)
		}
		sLog.Debug(MSG+"DBConnect", "acquire", "connect Ok")
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.TODO(), config)
	IfErrorThenPanic(err)
	_, err = pool.Acquire(context.TODO())
	IfErrorThenPanic(err)
	sLog.Debug(MSG+"DBConnect", "acquire", "pool Ok")

	return pool
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
