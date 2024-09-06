/*
 * This file was last modified at 2024-08-03 12:03 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * postgres.go
 * $Id$
 */
//!+

// Package repo TODO.
package repo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/victor-skurikhin/etcd-client/v1/internal/domain"
	"github.com/victor-skurikhin/etcd-client/v1/internal/domain/entity"
	"github.com/victor-skurikhin/etcd-client/v1/internal/env"
	"log/slog"
	"sync"
	"time"
)

const (
	increase = 1
	tries    = 3
)

var _ domain.Repo[domain.Actioner[*domain.Entity, domain.Entity], *domain.Entity, domain.Entity] = (*Postgres[domain.Actioner[*domain.Entity, domain.Entity], *domain.Entity, domain.Entity])(nil)

var (
	ErrBadPool       = fmt.Errorf("bad Database pool")
	onceKeyValueRepo = new(sync.Once)
	repoKeyValueInst *Postgres[domain.Actioner[*entity.KeyValue, entity.KeyValue], *entity.KeyValue, entity.KeyValue]
)

type Postgres[A domain.Actioner[T, U], T domain.Ptr[U], U domain.Entity] struct {
	pool *pgxpool.Pool
	sLog *slog.Logger
}

type PostgresError struct {
	err  error
	info interface{}
}

func GetKeyValuePostgresRepo(
	cfg env.Config,
) domain.Repo[domain.Actioner[*entity.KeyValue, entity.KeyValue], *entity.KeyValue, entity.KeyValue] {
	onceKeyValueRepo.Do(func() {
		repoKeyValueInst = new(Postgres[domain.Actioner[*entity.KeyValue, entity.KeyValue], *entity.KeyValue, entity.KeyValue])
		repoKeyValueInst.pool = cfg.DBPool()
		repoKeyValueInst.sLog = cfg.Logger()
	})
	return repoKeyValueInst
}

func (p Postgres[A, T, U]) Do(ctx context.Context, action A, unit U, scan func(domain.Scanner) U) (U, error) {

	row, err := rowPostgreSQL(ctx, p.sLog, p.pool, action.SQL(), action.Args(unit)...)

	if err != nil {
		return unit, PostgresError{err: err}
	}
	return scan(row), nil
}

func (p Postgres[A, T, U]) Get(ctx context.Context, action A, unit U, scan func(domain.Scanner) U) ([]U, error) {

	result := make([]U, 0)
	rows, err := rowsPostgreSQL(ctx, p.sLog, p.pool, action.SQL(), action.Args(unit)...)

	if err != nil {
		return nil, PostgresError{err: err}
	}
	for rows.Next() {
		e := scan(rows)
		result = append(result, e)
	}
	return result, nil
}

func (s PostgresError) Error() string {
	return s.err.Error()
}

func (s PostgresError) Err() error {
	return s.err
}

func (s PostgresError) Info() interface{} {
	return s.info
}

func rowPostgreSQL(
	ctx context.Context,
	log *slog.Logger,
	pool *pgxpool.Pool,
	sql string,
	args ...any,
) (pgx.Row, error) {

	if pool == nil {
		return nil, ErrBadPool
	}
	conn, err := pool.Acquire(ctx)

	for i := 1; err != nil && i < tries*increase; i += increase {
		time.Sleep(time.Duration(i) * time.Second)
		log.WarnContext(ctx, env.MSG+"Postgres.rowPostgreSQL", "msg", "retry pool acquire row", "err", err)
		conn, err = pool.Acquire(ctx)
	}
	defer func() {
		if conn != nil {
			conn.Release()
		}
	}()
	if conn == nil || err != nil {
		return nil, PostgresError{err: fmt.Errorf("while connecting %v", err), info: conn}
	}
	return conn.QueryRow(ctx, sql, args...), nil
}

func rowsPostgreSQL(
	ctx context.Context,
	log *slog.Logger,
	pool *pgxpool.Pool,
	sql string,
	args ...any,
) (pgx.Rows, error) {

	if pool == nil {
		return nil, ErrBadPool
	}
	conn, err := pool.Acquire(ctx)

	for i := 1; err != nil && i < tries*increase; i += increase {
		time.Sleep(time.Duration(i) * time.Second)
		log.WarnContext(ctx, env.MSG+"Postgres.rowsPostgreSQL", "msg", "retry pool acquire rows", "err", err)
		conn, err = pool.Acquire(ctx)
	}
	defer func() {
		if conn != nil {
			conn.Release()
		}
	}()
	if conn == nil || err != nil {
		return nil, PostgresError{err: fmt.Errorf("while connecting %v", err), info: conn}
	}
	return conn.Query(ctx, sql, args...)
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
