/*
 * This file was last modified at 2024-08-05 23:27 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * mongo_pool.go
 * $Id$
 */
//!+

// Package etcd_pool TODO.
package etcd_pool

import (
	"fmt"
	"github.com/victor-skurikhin/etcd-client/v1/internal/env"
	"github.com/victor-skurikhin/etcd-client/v1/pool"
	clientV3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/connectivity"
	"log/slog"
	"runtime"
	"sync"
	"time"
)

var _ pool.EtcdPool = (*etcdPool)(nil)

var (
	onceEtcdPool = new(sync.Once)
	etcdPoolInst *etcdPool
)

type etcdPool struct {
	clientConfig clientV3.Config
	connections  int
	pool         chan clientV3.KV
	poolSize     int
	sLog         *slog.Logger
	timeout      time.Duration
}

func GetEtcdPool(cfg env.Config) pool.EtcdPool {

	onceEtcdPool.Do(func() {
		etcdPoolInst = new(etcdPool)
		etcdPoolInst.clientConfig = *cfg.EtcdClientConfig()
		etcdPoolInst.pool = make(chan clientV3.KV, 50*runtime.NumCPU())
		etcdPoolInst.poolSize = 50 * runtime.NumCPU()
		etcdPoolInst.sLog = cfg.Logger()
		etcdPoolInst.timeout = 500 * time.Millisecond
	})
	return etcdPoolInst
}

func (e *etcdPool) AcquireClient() (clientV3.KV, error) {
	for {
		select {
		case client := <-e.pool:
			switch getStateActiveConn(client) {
			case connectivity.Idle:
				fallthrough
			case connectivity.Ready:
				fallthrough
			case connectivity.Connecting:
				return client, nil
			case connectivity.Shutdown:
				return nil, fmt.Errorf("connectivity state: %s", getStateActiveConn(client).String())
			}
		default:
			if e.connections < e.poolSize {
				e.createClientToChan()
			}
		}
	}
}

func (e *etcdPool) ReleaseClient(client clientV3.KV) error {
	select {
	case e.pool <- client:
		return nil
	default:
		if cli, ok := client.(*clientV3.Client); ok {
			if err := cli.Close(); err != nil {
				e.sLog.Error(env.MSG+"etcdPool: Close the client failed", "err", err)
				return err
			} else {
				e.connections--
			}
		}
		return nil
	}
}

func (e *etcdPool) GracefulClose() (err error) {
	if e.pool == nil {
		return fmt.Errorf("close of nil channel")
	}
	close(e.pool)

	for client := range e.pool {
		if cli, ok := client.(*clientV3.Client); ok {
			err = cli.Close()
		}
	}
	return err
}

func (e *etcdPool) createClientToChan() clientV3.KV {

	client, err := clientV3.New(e.clientConfig)

	if err != nil {
		e.sLog.Error(env.MSG+"etcdPool: Create the client failed", "err", err)
	}
	e.pool <- client
	e.connections++

	return client
}

func getStateActiveConn(client clientV3.KV) connectivity.State {
	if cli, ok := client.(*clientV3.Client); ok {
		if client == nil {
			return -2
		} else if cli.ActiveConnection() == nil {
			return -3
		}
		return cli.ActiveConnection().GetState()
	}
	return -1
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
