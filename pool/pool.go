package pool

import (
	clientV3 "go.etcd.io/etcd/client/v3"
)

type EtcdPool interface {
	AcquireClient() (clientV3.KV, error)
	ReleaseClient(clientV3.KV) error
	GracefulClose() error
}
