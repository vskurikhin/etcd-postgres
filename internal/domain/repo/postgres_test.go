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

//import (
//	"context"
//	"fmt"
//	"github.com/stretchr/testify/assert"
//	"github.com/victor-skurikhin/etcd-client/v1/internal/domain"
//	"github.com/victor-skurikhin/etcd-client/v1/internal/domain/entity"
//	"github.com/victor-skurikhin/etcd-client/v1/internal/env"
//	"github.com/victor-skurikhin/etcd-client/v1/tool"
//	"log/slog"
//	"reflect"
//	"testing"
//)
//
//func TestPostgres(t *testing.T) {
//	for _, test := range []struct {
//		name string
//		fRun func(*testing.T) (interface{}, error)
//		want func(*testing.T, interface{}) bool
//	}{
//		{
//			"test #0 positive for function GetKeyValuePostgresRepo(env.Config)",
//			positiveGetKeyValuePostgresRepo,
//			positiveGetKeyValuePostgresRepoCheck,
//		},
//		{
//			"test #1 positive for struct Postgres method Do(context.Context, A, U, func(domain.Scanner))",
//			positivePostgresDo,
//			positivePostgresDoCheck,
//		},
//		{
//			"test #2 positive for struct Postgres method Get(context.Context, A, U, func(domain.Scanner) U)",
//			positivePostgresGet,
//			positivePostgresGetCheck,
//		},
//		{
//			"test #3 negative #1 for struct Postgres method Do(context.Context, A, U, func(domain.Scanner))",
//			negativePostgresDo1,
//			negativePostgresDo1Check,
//		},
//	} {
//		t.Run(test.name, func(t *testing.T) {
//			got, err := test.fRun(t)
//			assert.Nil(t, err)
//			assert.NotNil(t, got)
//			assert.True(t, test.want(t, got))
//		})
//	}
//}
//
//func positiveGetKeyValuePostgresRepo(t *testing.T) (interface{}, error) {
//
//	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
//	t.Setenv("DATABASE_DSN", "")
//	tool.SetLogger(slog.Default())
//	c := env.GetConfig()
//	cfg := c.(env.TestConfig)
//
//	return GetKeyValuePostgresRepo(cfg.GetTestConfig(env.WithTestDBPool("", nil))), nil
//}
//
//func positiveGetKeyValuePostgresRepoCheck(_ *testing.T, i interface{}) bool {
//	_, ok := i.(*Postgres[domain.Actioner[*entity.KeyValue, entity.KeyValue], *entity.KeyValue, entity.KeyValue])
//	return ok
//}
//
//func positivePostgresDo(t *testing.T) (interface{}, error) {
//
//	ctx := context.Background()
//	postgresContainer, dbURL := createPostgresContainer(t, ctx)
//	// Clean up the container
//	defer func() {
//		if err := postgresContainer.Terminate(ctx); err != nil {
//			t.Fatalf("failed to terminate container: %s", err)
//		}
//	}()
//	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
//	t.Setenv("DATABASE_DSN", dbURL)
//	tool.SetLogger(slog.Default())
//	repoPostgres = new(Postgres[domain.Actioner[*IDValue, IDValue], *IDValue, IDValue])
//	repoPostgres.pool = tool.DBConnect(dbURL)
//	fmt.Printf("dbURL: %s\n", dbURL)
//
//	idValue, err := repoPostgres.Do(ctx, IDValueSelect, IDValue{id: 1}, func(scanner domain.Scanner) IDValue {
//		return IDValue{}
//	})
//	return idValue, err
//}
//
//func positivePostgresDoCheck(_ *testing.T, i interface{}) bool {
//	return i == IDValue{id: 1}
//}
//
//func positivePostgresGet(t *testing.T) (interface{}, error) {
//
//	ctx := context.Background()
//	postgresContainer, dbURL := createPostgresContainer(t, ctx)
//	// Clean up the container
//	defer func() {
//		if err := postgresContainer.Terminate(ctx); err != nil {
//			t.Fatalf("failed to terminate container: %s", err)
//		}
//	}()
//	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
//	t.Setenv("DATABASE_DSN", dbURL)
//	tool.SetLogger(slog.Default())
//	repoPostgres = new(Postgres[domain.Actioner[*IDValue, IDValue], *IDValue, IDValue])
//	repoPostgres.pool = tool.DBConnect(dbURL)
//
//	idValues, err := repoPostgres.Get(context.Background(), IDValueGetAll, IDValue{id: 1}, func(s domain.Scanner) IDValue {
//		var r IDValue
//		err := s.Scan(&r.id, &r.value)
//		if err != nil {
//			t.Fatal(err)
//		}
//		return r
//	})
//	return idValues, err
//}
//
//func positivePostgresGetCheck(_ *testing.T, i interface{}) bool {
//	return reflect.DeepEqual(i, []IDValue{{id: 1, value: "test"}})
//}
//
//func negativePostgresDo1(t *testing.T) (interface{}, error) {
//
//	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
//	t.Setenv("DATABASE_DSN", "")
//	tool.SetLogger(slog.Default())
//	repo := GetKeyValuePostgresRepo(env.GetConfig())
//
//	_, err := repo.Do(context.Background(), entity.KeyValueSelect, entity.KeyValue{}, func(scanner domain.Scanner) entity.KeyValue {
//		return entity.KeyValue{}
//	})
//	return err, nil
//}
//
//func negativePostgresDo1Check(_ *testing.T, i interface{}) bool {
//	return i == ErrBadPool
//}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
