package pool

import (
	clientV3 "go.etcd.io/etcd/client/v3"
)

type EtcdPool interface {
	AcquireClient() (*clientV3.Client, error)
	ReleaseClient(*clientV3.Client) error
	GracefulClose() error
}
