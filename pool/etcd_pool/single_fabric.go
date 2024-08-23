package etcd_pool

import (
	"github.com/victor-skurikhin/etcd-client/v1/internal/env"
	"github.com/victor-skurikhin/etcd-client/v1/pool"
	clientV3 "go.etcd.io/etcd/client/v3"
	"sync"
)

var _ pool.EtcdPool = (*singleFabricEtcdClient)(nil)

var (
	onceSingleFabricEtcdClient = new(sync.Once)
	singleFabricEtcdClientInst *singleFabricEtcdClient
)

type singleFabricEtcdClient struct {
	clientConfig clientV3.Config
}

func GetSingleFabricEtcdClient(cfg env.Config) pool.EtcdPool {
	onceSingleFabricEtcdClient.Do(func() {
		singleFabricEtcdClientInst = new(singleFabricEtcdClient)
		singleFabricEtcdClientInst.clientConfig = *cfg.EtcdClientConfig()
	})
	return singleFabricEtcdClientInst
}

func (s *singleFabricEtcdClient) AcquireClient() (*clientV3.Client, error) {
	return clientV3.New(s.clientConfig)
}

func (s *singleFabricEtcdClient) ReleaseClient(client *clientV3.Client) error {
	return client.Close()
}

func (s *singleFabricEtcdClient) GracefulClose() error {
	return nil
}
