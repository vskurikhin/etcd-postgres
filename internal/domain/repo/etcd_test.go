/*
 * This file was last modified at 2024-08-21 09:43 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * etcd_test.go
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
//	"github.com/victor-skurikhin/etcd-client/v1/pool/etcd_pool"
//	"github.com/victor-skurikhin/etcd-client/v1/tool"
//	clientV3 "go.etcd.io/etcd/client/v3"
//	"log/slog"
//	"runtime"
//	"strconv"
//	"sync"
//	"testing"
//)
//
//func TestEtcd(t *testing.T) {
//	for _, test := range []struct {
//		name string
//		fRun func(*testing.T) (interface{}, error)
//		want func(*testing.T, interface{}) bool
//	}{
//		{
//			"test #0 positive for function GetKeyValueEtcdRepo(env.Config)",
//			positiveGetKeyValueEtcdRepo,
//			positiveGetKeyValueEtcdRepoCheck,
//		},
//		{
//			"test #1 positive for struct Etcd method Do(context.Context, A, U, func(domain.Scanner))",
//			positiveEtcdDo,
//			positiveEtcdDoCheck,
//		},
//		//{
//		//	"test #2 positive for struct Etcd method Get(context.Context, A, U, func(domain.Scanner) U)",
//		//	positiveEtcdGet,
//		//	positiveEtcdGetCheck,
//		//},
//		//{
//		//	"test #3 negative #1 for struct Etcd method Do(context.Context, A, U, func(domain.Scanner))",
//		//	negativeEtcdDo1,
//		//	negativeEtcdDo1Check,
//		//},
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
//func positiveGetKeyValueEtcdRepo(t *testing.T) (interface{}, error) {
//
//	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
//	t.Setenv("ETCD_ADDRESSES", "")
//	tool.SetLogger(slog.Default())
//	c := env.GetConfig()
//	cfg := c.(env.TestConfig)
//
//	return GetKeyValueEtcdRepo(cfg.GetTestConfig(env.WithEtcdClientConfig(clientV3.Config{}))), nil
//}
//
//func positiveGetKeyValueEtcdRepoCheck(_ *testing.T, i interface{}) bool {
//	_, ok := i.(*Etcd[domain.Actioner[*entity.KeyValue, entity.KeyValue], *entity.KeyValue, entity.KeyValue])
//	return ok
//}
//
//func positiveEtcdDo(t *testing.T) (interface{}, error) {
//
//	ctx := context.Background()
//	etcdContainer, ip, etcdPort := createEtcdContainer(t, ctx)
//	// Clean up the container
//	defer func() { _ = etcdContainer.Terminate(ctx) }()
//
//	tool.SetLogger(slog.Default())
//	etcdIDValue := new(Etcd[domain.Actioner[*IDValue, IDValue], *IDValue, IDValue])
//	cfg := getTestConfig(t, ip, etcdPort)
//	etcdIDValue.pool = etcd_pool.GetEtcdPool(cfg)
//	etcdIDValue.sLog = slog.Default()
//
//	_, err := etcdIDValue.Do(ctx, IDValueUpsert, IDValue{id: 1, value: "value1"}, func(scanner domain.Scanner) IDValue {
//		var key string
//		var value string
//		_ = scanner.Scan(&key, &value)
//		id, _ := strconv.Atoi(key)
//		return IDValue{id: id, value: value}
//	})
//	if err != nil {
//		t.Fatal(err)
//	}
//	result, err := etcdIDValue.Do(ctx, IDValueSelect, IDValue{id: 1}, func(scanner domain.Scanner) IDValue {
//		var key string
//		var value string
//		_ = scanner.Scan(&key, &value)
//		id, _ := strconv.Atoi(key)
//		return IDValue{id: id, value: value}
//	})
//	if err != nil {
//		t.Fatal(err)
//	}
//	_, err = etcdIDValue.Do(ctx, IDValueDelete, IDValue{id: 1}, func(scanner domain.Scanner) IDValue {
//		var key string
//		var value string
//		_ = scanner.Scan(&key, &value)
//		id, _ := strconv.Atoi(key)
//		return IDValue{id: id, value: value}
//	})
//	if err != nil {
//		t.Fatal(err)
//	}
//	result1, err1 := etcdIDValue.Do(ctx, IDValueSelect, IDValue{id: 1}, func(scanner domain.Scanner) IDValue {
//		var key string
//		var value string
//		_ = scanner.Scan(&key, &value)
//		id, _ := strconv.Atoi(key)
//		return IDValue{id: id, value: value}
//	})
//	fmt.Printf("result1: %v, err1: %v\n", result1, err1)
//	_ = etcdIDValue.pool.GracefulClose()
//
//	return result, err
//}
//
//func positiveEtcdDoCheck(_ *testing.T, i interface{}) bool {
//	return i == IDValue{id: 1, value: "value1"}
//}
//
//func positiveEtcdGet(t *testing.T) (interface{}, error) {
//
//	ctx := context.Background()
//	etcdContainer, ip, etcdPort := createEtcdContainer(t, ctx)
//	// Clean up the container
//	defer func() { _ = etcdContainer.Terminate(ctx) }()
//
//	tool.SetLogger(slog.Default())
//	etcdIDValue := new(Etcd[domain.Actioner[*IDValue, IDValue], *IDValue, IDValue])
//	cfg := getTestConfig(t, ip, etcdPort)
//	etcdIDValue.pool = etcd_pool.GetSingleFabricEtcdClient(cfg)
//	etcdIDValue.sLog = slog.Default()
//
//	var wg sync.WaitGroup
//	wg.Add(runtime.NumCPU() * BenchmarkRepeat)
//
//	for j := 0; j < runtime.NumCPU()*BenchmarkRepeat; j++ {
//		go func() {
//			_, _ = etcdIDValue.Do(ctx, IDValueUpsert, IDValue{id: j, value: fmt.Sprintf("value%d", j)},
//				func(scanner domain.Scanner) IDValue {
//					var key string
//					var value string
//					_ = scanner.Scan(&key, &value)
//					id, _ := strconv.Atoi(key)
//					return IDValue{id: id, value: value}
//				})
//			wg.Done()
//		}()
//	}
//	wg.Wait()
//	result, err := etcdIDValue.Get(ctx, IDValueSelect, IDValue{id: 5}, func(scanner domain.Scanner) IDValue {
//		var key string
//		var value string
//		_ = scanner.Scan(&key, &value)
//		id, _ := strconv.Atoi(key)
//		return IDValue{id: id, value: value}
//	})
//	if err != nil {
//		t.Fatal(err)
//	}
//	assert.True(t, len(result) > 1)
//	for i := 0; i < runtime.NumCPU()*BenchmarkRepeat; i++ {
//		got, er0 := etcdIDValue.Do(ctx, IDValueSelect, IDValue{id: i},
//			func(scanner domain.Scanner) IDValue {
//				var key string
//				var value string
//				_ = scanner.Scan(&key, &value)
//				id, _ := strconv.Atoi(key)
//				return IDValue{id: id, value: value}
//			})
//		if er0 != nil {
//			t.Fatal(err)
//		}
//		expected := IDValue{id: i, value: fmt.Sprintf("value%d", i)}
//		assert.Equal(t, expected, got)
//	}
//	result, err = etcdIDValue.Get(ctx, IDValueGetAll, IDValue{}, func(scanner domain.Scanner) IDValue {
//		var key string
//		var value string
//		_ = scanner.Scan(&key, &value)
//		id, _ := strconv.Atoi(key)
//		return IDValue{id: id, value: value}
//	})
//	if err != nil {
//		t.Fatal(err)
//	}
//	_ = etcdIDValue.pool.GracefulClose()
//	return result, err
//}
//
//func positiveEtcdGetCheck(_ *testing.T, i interface{}) bool {
//	if slice, ok := i.([]IDValue); ok {
//		return len(slice) == runtime.NumCPU()*BenchmarkRepeat
//	}
//	return false
//}
//
//func negativeEtcdDo1(_ *testing.T) (interface{}, error) {
//	return nil, nil
//}
//
//func negativeEtcdDo1Check(_ *testing.T, i interface{}) bool {
//	return false
//}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
