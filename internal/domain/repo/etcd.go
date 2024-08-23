/*
 * This file was last modified at 2024-08-21 09:43 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * etcd.go
 * $Id$
 */
//!+

// Package repo TODO.
package repo

import (
	"context"
	"fmt"
	"github.com/victor-skurikhin/etcd-client/v1/internal/domain"
	"github.com/victor-skurikhin/etcd-client/v1/internal/domain/entity"
	"github.com/victor-skurikhin/etcd-client/v1/internal/env"
	"github.com/victor-skurikhin/etcd-client/v1/pool"
	"github.com/victor-skurikhin/etcd-client/v1/pool/etcd_pool"
	clientV3 "go.etcd.io/etcd/client/v3"
	"log/slog"
	"sync"
)

var _ domain.Repo[domain.Actioner[*domain.Entity, domain.Entity], *domain.Entity, domain.Entity] = (*Etcd[domain.Actioner[*domain.Entity, domain.Entity], *domain.Entity, domain.Entity])(nil)

var (
	onceKeyValueEtcd = new(sync.Once)
	etcdKeyValueInst *Etcd[domain.Actioner[*entity.KeyValue, entity.KeyValue], *entity.KeyValue, entity.KeyValue]
)

type Etcd[A domain.Actioner[T, U], T domain.Ptr[U], U domain.Entity] struct {
	pool pool.EtcdPool
	sLog *slog.Logger
}

type EtcdError struct {
	err  error
	info interface{}
}

type ScannerError struct {
	err error
}

func GetKeyValueEtcdRepo(
	cfg env.Config,
) domain.Repo[domain.Actioner[*entity.KeyValue, entity.KeyValue], *entity.KeyValue, entity.KeyValue] {
	onceKeyValueEtcd.Do(func() {
		etcdKeyValueInst = new(Etcd[domain.Actioner[*entity.KeyValue, entity.KeyValue], *entity.KeyValue, entity.KeyValue])
		etcdKeyValueInst.pool = etcd_pool.GetSingleFabricEtcdClient(cfg)
		etcdKeyValueInst.sLog = cfg.Logger()
	})
	return etcdKeyValueInst
}

func (e Etcd[A, T, U]) Do(ctx context.Context, action A, unit U, scan func(domain.Scanner) U) (U, error) {

	client, err := e.pool.AcquireClient()

	if err != nil {
		return unit, EtcdError{err: err}
	}
	defer func() { _ = e.pool.ReleaseClient(client) }()

	switch action.Name() {
	case domain.DeleteAction:
		return e.delete(ctx, client, unit)
	case domain.SelectAction:
		return e.get(ctx, client, unit, scan)
	case domain.UpsertAction:
		return e.put(ctx, client, unit, action.Args(unit)...)
	}
	return unit, EtcdError{err: fmt.Errorf("unknown action, name: %s", action.Name())}
}

func (e Etcd[A, T, U]) Get(ctx context.Context, action A, unit U, scan func(domain.Scanner) U) ([]U, error) {

	client, err := e.pool.AcquireClient()

	if err != nil {
		return nil, EtcdError{err: err}
	}
	defer func() { _ = e.pool.ReleaseClient(client) }()

	switch action.Name() {
	case domain.GetAllAction:
		return e.getAll(ctx, client, "\x00", scan)
	case domain.SelectAction:
		return e.getAll(ctx, client, unit.Key(), scan)
	}
	return nil, EtcdError{err: fmt.Errorf("unknown action, name: %s", action.Name())}
}

func (e Etcd[A, T, U]) delete(ctx context.Context, client *clientV3.Client, unit U) (U, error) {

	if resp, err := client.Delete(ctx, unit.Key()); err != nil {
		return unit, EtcdError{err: err, info: resp}
	}
	return unit, nil
}

func (e Etcd[A, T, U]) get(ctx context.Context, client *clientV3.Client, unit U, scan func(domain.Scanner) U) (U, error) {

	got, err := client.Get(ctx, unit.Key())

	if err != nil {
		return unit, EtcdError{err: err, info: got}
	}
	if len(got.Kvs) < 1 {
		return unit, EtcdError{err: fmt.Errorf("no Kvs, length: %d", len(got.Kvs))}
	}

	return scan(keyValueScanner{key: string(got.Kvs[0].Key), value: got.Kvs[0].Value}), nil
}

func (e Etcd[A, T, U]) getAll(ctx context.Context, client *clientV3.Client, key string, scan func(domain.Scanner) U) ([]U, error) {

	got, err := client.Get(ctx, key, clientV3.WithFromKey())

	if err != nil {
		return nil, EtcdError{err: err}
	}
	result := make([]U, 0)

	for _, kv := range got.Kvs {
		u := scan(keyValueScanner{key: string(kv.Key), value: kv.Value})
		result = append(result, u)
	}
	return result, nil
}

func (e Etcd[A, T, U]) put(ctx context.Context, client *clientV3.Client, unit U, args ...any) (U, error) {

	if len(args) < 2 {
		return unit, EtcdError{err: fmt.Errorf("no required parameters, length: %d", len(args))}
	}
	if s, ok := args[1].(string); !ok {
		return unit, EtcdError{err: fmt.Errorf(
			"second argument for scanner is not pointer to string, type: %T", args[1],
		)}
	} else if resp, err := client.Put(ctx, unit.Key(), s); err != nil {
		return unit, EtcdError{err: err, info: resp}
	}
	return unit, nil
}

func (s EtcdError) Error() string {
	return s.err.Error()
}

func (s EtcdError) Err() error {
	return s.err
}

func (s EtcdError) Info() interface{} {
	return s.info
}

func (s ScannerError) Error() string {
	return s.err.Error()
}

type keyValueScanner struct {
	key   string
	value []byte
}

func (v keyValueScanner) Scan(dest ...any) error {

	if len(dest) < 2 {
		return ScannerError{err: fmt.Errorf("no required parameters, length: %d", len(dest))}
	}
	if pKey, ok0 := dest[0].(*string); ok0 {
		*pKey = v.key

		if pValue, ok1 := dest[1].(*string); ok1 {
			*pValue = string(v.value)
			return nil
		}
		return ScannerError{err: fmt.Errorf("argument for scanner is not pointer to string, type:  %T", dest[0])}
	}
	return ScannerError{err: fmt.Errorf("second argument for scanner is not pointer to string, type: %T", dest[1])}
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
