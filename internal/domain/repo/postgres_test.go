/*
 * This file was last modified at 2024-08-03 12:03 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * postgres_test.go
 * $Id$
 */
//!+

// Package repo TODO.
package repo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/victor-skurikhin/etcd-client/v1/internal/domain"
	"github.com/victor-skurikhin/etcd-client/v1/internal/domain/entity"
	"github.com/victor-skurikhin/etcd-client/v1/internal/env"
	"github.com/victor-skurikhin/etcd-client/v1/tool"
	"log/slog"
	"testing"
)

func TestPostgres(t *testing.T) {
	for _, test := range []struct {
		name string
		fRun func(*testing.T) (interface{}, error)
		want func(*testing.T, interface{}) bool
	}{
		{
			"test #0 positive for function GetKeyValueEtcdRepo(env.Config)",
			positiveGetKeyValuePostgresRepo,
			positiveGetKeyValuePostgresRepoCheck,
		},
		{
			"test #1 negative #1 for struct Postgres method Do(context.Context, A, U, func(domain.Scanner))",
			negativePostgresDo1,
			negativePostgresDo1Check,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.fRun(t)
			assert.Nil(t, err)
			assert.NotNil(t, got)
			assert.True(t, test.want(t, got))
		})
	}
}

func positiveGetKeyValuePostgresRepo(t *testing.T) (interface{}, error) {
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	t.Setenv("DATABASE_DSN", "")
	tool.SetLogger(slog.Default())
	c := env.GetConfig()
	cfg := c.(env.TestConfig)

	return GetKeyValuePostgresRepo(cfg.GetTestConfig(env.WithTestDBPool("", nil))), nil
}

func positiveGetKeyValuePostgresRepoCheck(t *testing.T, i interface{}) bool {
	_, ok := i.(*Postgres[domain.Actioner[*entity.KeyValue, entity.KeyValue], *entity.KeyValue, entity.KeyValue])
	return ok
}

func negativePostgresDo1(t *testing.T) (interface{}, error) {

	repoKeyValueInst = newTestKeyValuePostgresRepo(nil, slog.Default())
	search := entity.MakeKeyValue("key", "", 0, entity.DefaultTAttributes())
	ctx := context.Background()
	_, err := repoKeyValueInst.Do(ctx, entity.KeyValueUpsert, search, func(domain.Scanner) entity.KeyValue {
		return search
	})
	fmt.Printf("err : %T = %v\n", err, err)
	assert.NotNil(t, err)

	return err, nil
}

func negativePostgresDo1Check(t *testing.T, i interface{}) bool {
	if postgresError, ok := i.(PostgresError); ok {
		return postgresError.err.Error() == fmt.Errorf("bad Database pool").Error()
	}
	return false
}

func newTestKeyValuePostgresRepo(
	pool *pgxpool.Pool, sLog *slog.Logger,
) *Postgres[domain.Actioner[*entity.KeyValue, entity.KeyValue], *entity.KeyValue, entity.KeyValue] {
	repo := new(Postgres[domain.Actioner[*entity.KeyValue, entity.KeyValue], *entity.KeyValue, entity.KeyValue])
	repo.pool = pool
	repo.sLog = sLog
	return repo
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
