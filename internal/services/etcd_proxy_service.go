/*
 * This file was last modified at 2024-08-16 11:35 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * etcd_proxy_service.go
 * $Id$
 */
//!+

package services

import (
	"context"
	"fmt"
	"github.com/victor-skurikhin/etcd-client/v1/internal/controllers/dto"
	"github.com/victor-skurikhin/etcd-client/v1/internal/domain/memory"
	"github.com/victor-skurikhin/etcd-client/v1/internal/env"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	pb "github.com/victor-skurikhin/etcd-client/v1/proto"
	clientV3 "go.etcd.io/etcd/client/v3"
)

const BadOneOfUnionValue = "bad oneOf Union value"
const CacheInvalidate = "Y2FjaGUtaW52YWxpZGF0ZQo="

type EtcdProxyService interface {
	pb.EtcdClientServiceServer
	ApiDelete(ctx context.Context, key string) error
	ApiGet(ctx context.Context, key string) (dto.Result, error)
	ApiPut(ctx context.Context, data dto.KeyValue) error
}

type etcdProxyService struct {
	pb.UnimplementedEtcdClientServiceServer
	cache        *memory.Storage
	cacheExpire  time.Duration
	clientConfig clientV3.Config
	hitCounter   atomic.Uint64
	sLog         *slog.Logger
}

var _ EtcdProxyService = (*etcdProxyService)(nil)
var (
	ErrNotFound   = fmt.Errorf("not found")
	onceEtcdProxy = new(sync.Once)
	etcdProxyServ *etcdProxyService
)

// GetEtcdProxyService — потокобезопасное (thread-safe) создание
// REST веб-сервиса etcd proxy.
func GetEtcdProxyService(ctx context.Context, cfg env.Config) EtcdProxyService {

	onceEtcdProxy.Do(func() {
		etcdProxyServ = new(etcdProxyService)
		etcdProxyServ.cache = memory.New(memory.Config{
			GCInterval: cfg.CacheGCInterval(),
		})
		etcdProxyServ.cacheExpire = cfg.CacheExpire()
		etcdProxyServ.clientConfig = *cfg.EtcdClientConfig()
		etcdProxyServ.hitCounter = atomic.Uint64{}
		etcdProxyServ.sLog = cfg.Logger()
		go func() {
			etcdProxyServ.watch(ctx, cfg)
		}()
	})
	return etcdProxyServ
}

func (f *etcdProxyService) ApiDelete(ctx context.Context, key string) error {
	return f.delete(ctx, key)
}

func (f *etcdProxyService) ApiGet(ctx context.Context, key string) (dto.Result, error) {
	return f.get(ctx, key)
}

func (f *etcdProxyService) ApiPut(ctx context.Context, data dto.KeyValue) error {
	return f.put(ctx, data)
}

func (f *etcdProxyService) Delete(ctx context.Context, request *pb.EtcdClientRequest) (*pb.EtcdClientResponse, error) {

	f.sLog.InfoContext(ctx, env.MSG+"EtcdProxyService.Delete", "msg", "gRPC", "request", *request)

	var err error
	var response = pb.EtcdClientResponse{Status: pb.Status_UNKNOWN}

	switch u := request.Union.(type) {
	case *pb.EtcdClientRequest_Key: // u.Number contains the number.

		key := u.Key.GetKey()

		if err = f.delete(ctx, key); err != nil {
			response.Error = err.Error()
			response.Status = pb.Status_FAIL
		} else {
			response.Status = pb.Status_OK
		}
	case *pb.EtcdClientRequest_KeyValue: // u.Name contains the string.
		response.Error = BadOneOfUnionValue
		response.Status = pb.Status_FAIL
	default:
		response.Error = BadOneOfUnionValue
		response.Status = pb.Status_FAIL
	}
	return &response, err
}

func (f *etcdProxyService) Get(ctx context.Context, request *pb.EtcdClientRequest) (*pb.EtcdClientResponse, error) {

	f.sLog.InfoContext(ctx, env.MSG+"EtcdProxyService.Get", "msg", "gRPC", "request", *request)

	var err error
	var response = pb.EtcdClientResponse{Status: pb.Status_UNKNOWN}

	switch u := request.Union.(type) {
	case *pb.EtcdClientRequest_Key:

		key := u.Key.GetKey()

		if got, err := f.get(ctx, key); err != nil {
			response.Error = err.Error()
			response.Status = pb.Status_FAIL
		} else {
			response.KeyValue = &pb.KeyValue{Key: key, Value: got.Value}
			response.Status = pb.Status_OK
		}
	case *pb.EtcdClientRequest_KeyValue:
		response.Error = BadOneOfUnionValue
		response.Status = pb.Status_FAIL
	default:
		response.Error = BadOneOfUnionValue
		response.Status = pb.Status_FAIL
	}
	return &response, err
}

func (f *etcdProxyService) Put(ctx context.Context, request *pb.EtcdClientRequest) (*pb.EtcdClientResponse, error) {

	f.sLog.InfoContext(ctx, env.MSG+"EtcdProxyService.Put", "msg", "gRPC", "request", *request)

	var err error
	var response = pb.EtcdClientResponse{Status: pb.Status_UNKNOWN}

	switch u := request.Union.(type) {
	case *pb.EtcdClientRequest_KeyValue:

		key := u.KeyValue.GetKey()
		value := u.KeyValue.GetValue()

		if err := f.put(ctx, dto.KeyValue{Key: key, Value: value}); err != nil {
			response.Error = err.Error()
			response.Status = pb.Status_FAIL
		} else {
			response.Status = pb.Status_OK
		}
	case *pb.EtcdClientRequest_Key:
		response.Error = BadOneOfUnionValue
		response.Status = pb.Status_FAIL
	default:
		response.Error = BadOneOfUnionValue
		response.Status = pb.Status_FAIL
	}
	return &response, err
}

func (f *etcdProxyService) delete(ctx context.Context, key string) error {

	cli, err := clientV3.New(f.clientConfig)

	if err != nil {
		f.sLog.ErrorContext(ctx, env.MSG+"EtcdProxyService.delete", "err", err)
	}
	defer func() { _ = cli.Close() }()

	if resp, err := cli.Delete(ctx, key); err != nil {
		f.sLog.ErrorContext(ctx,
			env.MSG+"EtcdProxyService.delete",
			"msg", "cli.ApiDelete", "err", err, "resp", resp,
		)
		return err
	} else {
		f.sLog.DebugContext(ctx,
			env.MSG+"EtcdProxyService.delete",
			"msg", fmt.Sprintf("Delete is done. Metadata is %q\n", resp),
		)
		f.keyInvalidate(ctx, cli, key)
	}
	return nil
}

func (f *etcdProxyService) get(ctx context.Context, key string) (dto.Result, error) {

	data, err := f.cache.Get(key)

	if err == nil && data != nil {

		counter := f.hitCounter.Add(1)
		f.sLog.DebugContext(ctx, env.MSG+"EtcdProxyService.get", "msg", fmt.Sprintf("cache hit count: %d", counter))

		return dto.Result{Value: string(data)}, nil
	} else {
		f.sLog.ErrorContext(ctx, env.MSG+"EtcdProxyService.get", "msg", "cache.Get", "err", err)
	}
	if result, err := f.cliGet(ctx, key); err != nil {
		return dto.Result{}, err
	} else {
		f.cacheSet(ctx, dto.KeyValue{Key: key, Value: result.Value})
		return result, nil
	}
}

func (f *etcdProxyService) cliGet(ctx context.Context, key string) (result dto.Result, err error) {

	cli, err := clientV3.New(f.clientConfig)

	if err != nil {
		f.sLog.ErrorContext(ctx, env.MSG+"EtcdProxyService.cliGet", "err", err)
		return dto.Result{}, err
	}
	defer func() { _ = cli.Close() }()
	got, err := cli.Get(ctx, key)

	if err != nil {
		f.sLog.ErrorContext(ctx, env.MSG+"EtcdProxyService.cliGet", "err", err)
		return dto.Result{}, err
	} else {
		f.sLog.DebugContext(ctx,
			env.MSG+"EtcdProxyService.cliGet",
			"msg", fmt.Sprintf("cli.Get is done. Metadata is %v\n", got),
		)
		if len(got.Kvs) < 1 {
			return dto.Result{}, ErrNotFound
		}
		result = dto.Result{Value: string(got.Kvs[0].Value)}
		f.sLog.DebugContext(ctx,
			env.MSG+"EtcdProxyService.cliGet",
			"msg", fmt.Sprintf("the value: %s", string(got.Kvs[0].Value)),
		)
	}
	return result, nil
}

func (f *etcdProxyService) put(ctx context.Context, data dto.KeyValue) error {

	cli, err := clientV3.New(f.clientConfig)

	if err != nil {
		f.sLog.ErrorContext(ctx, env.MSG+"EtcdProxyService.put", "err", err)
		return err
	}
	defer func() { _ = cli.Close() }()

	if resp, err := cli.Put(ctx, data.Key, data.Value); err != nil {
		f.sLog.ErrorContext(ctx,
			env.MSG+"EtcdProxyService.put",
			"msg", "cli.Put", "err", err, "resp", resp,
		)
		return err
	} else {
		f.sLog.DebugContext(ctx,
			env.MSG+"EtcdProxyService.put",
			"msg", fmt.Sprintf("cli.Put is done. Metadata is %q\n", resp),
		)
		f.keyInvalidate(ctx, cli, data.Key)
		f.cacheSet(ctx, data)
	}
	return nil
}

func (f *etcdProxyService) cacheSet(ctx context.Context, data dto.KeyValue) {

	if err := f.cache.Set(data.Key, []byte(data.Value), f.cacheExpire); err != nil {
		f.sLog.ErrorContext(ctx, env.MSG+"EtcdProxyService.cacheSet", "err", err)
	}
}

func (f *etcdProxyService) keyInvalidate(ctx context.Context, cli *clientV3.Client, key string) {

	if resp, err := cli.Put(ctx, CacheInvalidate, key); err != nil {
		f.sLog.ErrorContext(ctx, env.MSG+"EtcdProxyService.keyInvalidate", "err", err)
	} else {
		f.sLog.DebugContext(ctx,
			env.MSG+"EtcdProxyService.keyInvalidate",
			"msg", fmt.Sprintf("key %s in cache invalidate is done. Metadata is %q\n", key, resp),
		)
	}
}

func (f *etcdProxyService) watch(ctx context.Context, cfg env.Config) {

	cli, err := clientV3.New(*cfg.EtcdClientConfig())

	if err != nil {
		f.sLog.ErrorContext(ctx, "err", err)
	}
	defer func() { _ = cli.Close() }()
	rch := cli.Watch(ctx, CacheInvalidate)

	for watchResp := range rch {
		for _, ev := range watchResp.Events {
			key := string(ev.Kv.Value)
			if err := f.cache.Delete(key); err != nil {
				f.sLog.ErrorContext(ctx,
					env.MSG+"EtcdProxyService.watch",
					"msg", "cache.Delete",
					"err", err)
			} else {
				f.sLog.DebugContext(ctx,
					env.MSG+"EtcdProxyService.watch",
					"msg", fmt.Sprintf("Cache invalidate by key: \"%s\" is OK.", key),
				)
				f.hitCounter.Store(0)
			}
		}
	}
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
