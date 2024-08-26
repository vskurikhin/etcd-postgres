/*
 * This file was last modified at 2024-08-29 12:33 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * key_value_data_service.go
 * $Id$
 */

package services

import (
	"context"
	"fmt"
	"github.com/victor-skurikhin/etcd-client/v1/internal/domain"
	"github.com/victor-skurikhin/etcd-client/v1/internal/domain/entity"
	"github.com/victor-skurikhin/etcd-client/v1/internal/domain/memory"
	"github.com/victor-skurikhin/etcd-client/v1/internal/domain/repo"
	"github.com/victor-skurikhin/etcd-client/v1/internal/env"
	"github.com/victor-skurikhin/etcd-client/v1/pool"
	"github.com/victor-skurikhin/etcd-client/v1/pool/etcd_pool"
	clientV3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

const CacheInvalidate = "Y2FjaGUtaW52YWxpZGF0ZQo="

type KeyValueDataService interface {
	Delete(context.Context, string) error
	Get(context.Context, string) (entity.KeyValue, error)
	Put(context.Context, entity.KeyValue) error
}

type keyValueDataService struct {
	cache        *memory.Storage
	cacheExpire  time.Duration
	etcdRepo     domain.Repo[domain.Actioner[*entity.KeyValue, entity.KeyValue], *entity.KeyValue, entity.KeyValue]
	hitCounter   atomic.Uint64
	pool         pool.EtcdPool
	postgresRepo domain.Repo[domain.Actioner[*entity.KeyValue, entity.KeyValue], *entity.KeyValue, entity.KeyValue]
	sLog         *slog.Logger
}

var _ KeyValueDataService = (*keyValueDataService)(nil)

var (
	onceKeyValueDataService = new(sync.Once)
	keyValueDataServiceInst *keyValueDataService
)

func (k *keyValueDataService) Delete(ctx context.Context, key string) error {
	return k.delete(ctx, key)
}

func (k *keyValueDataService) Get(ctx context.Context, key string) (entity.KeyValue, error) {
	return k.get(ctx, key)
}

func (k *keyValueDataService) Put(ctx context.Context, value entity.KeyValue) error {
	return k.put(ctx, value)
}

func (k *keyValueDataService) delete(ctx context.Context, key string) error {

	_, err := k.postgresRepo.Do(
		ctx,
		entity.KeyValueDelete,
		entity.MakeKeyValue(key, "", 0, entity.DefaultTAttributes()),
		func(domain.Scanner) entity.KeyValue {
			return entity.KeyValue{}
		})
	if err != nil {
		return err
	}
	_, err = k.etcdRepo.Do(
		ctx,
		entity.KeyValueDelete,
		entity.MakeKeyValue(key, "", 0, entity.DefaultTAttributes()),
		func(domain.Scanner) entity.KeyValue {
			return entity.KeyValue{}
		})
	if err == nil {
		k.keyInvalidate(ctx, key)
	}
	return err
}

const cntKeyValueDataServiceGetJobs = 2

type msgKeyValue struct {
	err   error
	name  string
	value entity.KeyValue
}

func (k *keyValueDataService) get(ctx context.Context, key string) (entity.KeyValue, error) {

	data, err := k.cache.Get(key)

	if err == nil && data != nil {
		k.incrementHitCounter(ctx)
		var result entity.KeyValue
		if er0 := result.FromJSON(data); er0 == nil {
			return result, nil
		} else {
			k.sLog.DebugContext(ctx, env.MSG+"keyValueDataService.get", "err", er0)
		}
	} else {
		k.sLog.DebugContext(ctx, env.MSG+"keyValueDataService.get", "err", err)
	}
	var wg sync.WaitGroup

	wg.Add(cntKeyValueDataServiceGetJobs)
	quit := make(chan struct{})
	results := make(chan msgKeyValue, cntKeyValueDataServiceGetJobs)

	go func() {
		defer wg.Done()
		results <- k.getEtcd(ctx, key)
	}()
	go func() {
		defer wg.Done()
		results <- k.getPostgres(ctx, key)
	}()
	go func() {
		wg.Wait()
		close(results)
		close(quit)
	}()
	for {
		select {
		case result := <-results:
			if result.err == nil {
				return result.value, nil
			}
		case <-quit:
			return entity.KeyValue{}, fmt.Errorf("ERROR") // TODO
		}
	}
}

func (k *keyValueDataService) getEtcd(ctx context.Context, key string) msgKeyValue {

	result, err := k.etcdRepo.Do(ctx,
		entity.KeyValueDelete,
		makeKetValueWithKeyOnly(key),
		func(domain.Scanner) entity.KeyValue {
			return entity.KeyValue{}
		})
	return msgKeyValue{err: err, name: "Etcd", value: result}
}

func (k *keyValueDataService) getPostgres(ctx context.Context, key string) msgKeyValue {

	result, err := k.postgresRepo.Do(ctx,
		entity.KeyValueDelete,
		makeKetValueWithKeyOnly(key),
		func(domain.Scanner) entity.KeyValue {
			return entity.KeyValue{}
		})
	return msgKeyValue{err: err, name: "Postgres", value: result}
}

func (k *keyValueDataService) put(ctx context.Context, unit entity.KeyValue) error {

	g, c := errgroup.WithContext(ctx)
	g.Go(func() error {
		return k.putEtcd(c, unit)
	})
	g.Go(func() error {
		return k.putPostgres(c, unit)
	})
	k.keyInvalidate(ctx, unit.Key())
	return g.Wait()
}

func (k *keyValueDataService) putEtcd(ctx context.Context, unit entity.KeyValue) error {
	_, err := k.etcdRepo.Do(ctx, entity.KeyValueUpsert, unit, func(domain.Scanner) entity.KeyValue {
		return entity.KeyValue{}
	})
	return err
}

func (k *keyValueDataService) putPostgres(ctx context.Context, unit entity.KeyValue) error {
	_, err := k.postgresRepo.Do(ctx, entity.KeyValueUpsert, unit, func(domain.Scanner) entity.KeyValue {
		return entity.KeyValue{}
	})
	return err
}

func (k *keyValueDataService) incrementHitCounter(ctx context.Context) {
	counter := k.hitCounter.Add(1)
	k.sLog.DebugContext(ctx,
		env.MSG+"keyValueDataService.get", "msg",
		fmt.Sprintf("cache hit count: %d", counter),
	)
}

func (k *keyValueDataService) keyInvalidate(ctx context.Context, key string) {

	client, err := k.pool.AcquireClient()

	if err != nil {
		k.sLog.DebugContext(ctx,
			env.MSG+"keyValueDataService.keyInvalidate",
			"msg", "when acquire client",
			"err", err,
		)
		return
	}
	defer func() { _ = k.pool.ReleaseClient(client) }()

	if resp, err := client.Put(ctx, CacheInvalidate, key); err != nil {
		k.sLog.ErrorContext(ctx,
			env.MSG+"keyValueDataService.keyInvalidate",
			"msg", "when put",
			"err", err,
		)
	} else {
		k.sLog.DebugContext(ctx,
			env.MSG+"keyValueDataService.keyInvalidate",
			"msg", fmt.Sprintf("key %s in cache invalidate is done. Metadata is %q\n", key, resp),
		)
	}
}

func (k *keyValueDataService) watch(ctx context.Context, cfg env.Config) {

	cli, err := clientV3.New(*cfg.EtcdClientConfig())

	if err != nil {
		k.sLog.ErrorContext(ctx,
			env.MSG+"keyValueDataService.watch",
			"msg", "new client", "err", err,
		)
	}
	defer func() { _ = cli.Close() }()
	rch := cli.Watch(ctx, CacheInvalidate)

	for watchResp := range rch {
		for _, ev := range watchResp.Events {
			key := string(ev.Kv.Value)
			if err := k.cache.Delete(key); err != nil {
				k.sLog.ErrorContext(ctx,
					env.MSG+"keyValueDataService.watch",
					"msg", "cache.Delete",
					"err", err)
			} else {
				k.sLog.DebugContext(ctx,
					env.MSG+"keyValueDataService.watch",
					"msg", fmt.Sprintf("Cache invalidate by key: \"%s\" is OK.", key),
				)
				k.hitCounter.Store(0)
			}
		}
	}
}

func GetKeyValueDataService(ctx context.Context, cfg env.Config) KeyValueDataService {

	onceKeyValueDataService.Do(func() {
		keyValueDataServiceInst = new(keyValueDataService)
		keyValueDataServiceInst.cache = memory.New(memory.Config{
			GCInterval: cfg.CacheGCInterval(),
		})
		keyValueDataServiceInst.cacheExpire = cfg.CacheExpire()
		keyValueDataServiceInst.etcdRepo = repo.GetKeyValueEtcdRepo(cfg)
		keyValueDataServiceInst.pool = etcd_pool.GetEtcdPool(cfg)
		keyValueDataServiceInst.postgresRepo = repo.GetKeyValuePostgresRepo(cfg)
		keyValueDataServiceInst.sLog = cfg.Logger()
		go func() {
			keyValueDataServiceInst.watch(ctx, cfg)
		}()
	})
	return keyValueDataServiceInst
}

func makeKetValueWithKeyOnly(key string) entity.KeyValue {
	return entity.MakeKeyValue(key, "", 0, entity.DefaultTAttributes())
}
