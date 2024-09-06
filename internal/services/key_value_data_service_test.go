package services

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/victor-skurikhin/etcd-client/v1/internal/domain"
	"github.com/victor-skurikhin/etcd-client/v1/internal/domain/entity"
	"github.com/victor-skurikhin/etcd-client/v1/internal/domain/memory"
	"github.com/victor-skurikhin/etcd-client/v1/internal/env"
	"github.com/victor-skurikhin/etcd-client/v1/pool"
	"github.com/victor-skurikhin/etcd-client/v1/tool"
	"go.etcd.io/etcd/api/v3/etcdserverpb"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientV3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/mock/gomock"
	"log/slog"
	"testing"
)

func TestKeyValueDataService(t *testing.T) {
	for _, test := range []struct {
		name string
		fRun func(*testing.T) (interface{}, error)
		want func(*testing.T, interface{}) bool
	}{
		{
			"test #0 positive for function GetKeyValueDataService(env.Config)",
			positiveGetKeyValueDataService,
			positiveGetKeyValueDataServiceCheck,
		},
		{
			"test #1 positive #1 for struct KeyValueDataService method Delete(context.Context, string)",
			positivePostgresDelete1,
			positivePostgresDelete1Check,
		},
		{
			"test #1 positive #1 for struct KeyValueDataService method Delete(context.Context, string)",
			positivePostgresGet1,
			positivePostgresGet1Check,
		},
		{
			"test #1 positive #1 for struct KeyValueDataService method Delete(context.Context, string)",
			positivePostgresPut1,
			positivePostgresPut1Check,
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

func positiveGetKeyValueDataService(t *testing.T) (interface{}, error) {
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	t.Setenv("DATABASE_DSN", "")
	tool.SetLogger(slog.Default())
	c := env.GetConfig()
	cfg := c.(env.TestConfig)

	return GetKeyValueDataService(context.Background(), cfg.GetTestConfig(
		env.WithTestDBPool("", nil),
		env.WithEtcdClientConfig(clientV3.Config{Endpoints: []string{"localhost:0"}}),
	)), nil
}

func positiveGetKeyValueDataServiceCheck(_ *testing.T, i interface{}) bool {
	_, ok := i.(*keyValueDataService)
	return ok
}

func positivePostgresDelete1(t *testing.T) (interface{}, error) {
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	t.Setenv("DATABASE_DSN", "")

	cfg0 := env.GetConfig()
	cfg1 := cfg0.(env.TestConfig)
	cfg2 := cfg1.GetTestConfig(
		env.WithTestDBPool("", nil),
		env.WithEtcdClientConfig(clientV3.Config{Endpoints: []string{"localhost:0"}}),
	)
	ctrl := gomock.NewController(t)
	etcdRepo := NewMockRepo[domain.Actioner[*entity.KeyValue, entity.KeyValue], *entity.KeyValue, entity.KeyValue](ctrl)
	etcdPoolMock := NewMockEtcdPool(ctrl)
	postgresRepo := NewMockRepo[domain.Actioner[*entity.KeyValue, entity.KeyValue], *entity.KeyValue, entity.KeyValue](ctrl)
	etcdRepo.
		EXPECT().
		Do(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(entity.KeyValue{}, nil).
		AnyTimes()
	etcdPoolMock.
		EXPECT().
		AcquireClient().
		Return(&kvTestStub{}, nil).
		AnyTimes()
	etcdPoolMock.
		EXPECT().
		ReleaseClient(gomock.Any()).
		Return(nil).
		AnyTimes()
	postgresRepo.
		EXPECT().
		Do(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(entity.KeyValue{}, nil).
		AnyTimes()

	srv := newTestKeyValueDataService(cfg2, etcdRepo, etcdPoolMock, postgresRepo)

	return srv, srv.Delete(context.Background(), "key1")
}

func positivePostgresDelete1Check(_ *testing.T, i interface{}) bool {
	_, ok := i.(*keyValueDataService)
	return ok
}

func positivePostgresGet1(t *testing.T) (interface{}, error) {
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	t.Setenv("DATABASE_DSN", "")

	cfg0 := env.GetConfig()
	cfg1 := cfg0.(env.TestConfig)
	cfg2 := cfg1.GetTestConfig(
		env.WithTestDBPool("", nil),
		env.WithEtcdClientConfig(clientV3.Config{Endpoints: []string{"localhost:0"}}),
	)
	ctrl := gomock.NewController(t)
	etcdRepo := NewMockRepo[domain.Actioner[*entity.KeyValue, entity.KeyValue], *entity.KeyValue, entity.KeyValue](ctrl)
	etcdPoolMock := NewMockEtcdPool(ctrl)
	postgresRepo := NewMockRepo[domain.Actioner[*entity.KeyValue, entity.KeyValue], *entity.KeyValue, entity.KeyValue](ctrl)
	etcdRepo.
		EXPECT().
		Do(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(entity.KeyValue{}, nil).
		AnyTimes()
	etcdPoolMock.
		EXPECT().
		AcquireClient().
		Return(&kvTestStub{}, nil).
		AnyTimes()
	etcdPoolMock.
		EXPECT().
		ReleaseClient(gomock.Any()).
		Return(nil).
		AnyTimes()
	postgresRepo.
		EXPECT().
		Do(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(entity.KeyValue{}, nil).
		AnyTimes()

	srv := newTestKeyValueDataService(cfg2, etcdRepo, etcdPoolMock, postgresRepo)
	unit, err := srv.Get(context.Background(), "key1")

	return unit, err
}

func positivePostgresGet1Check(t *testing.T, i interface{}) bool {
	_, ok := i.(entity.KeyValue)
	return ok
}

func positivePostgresPut1(t *testing.T) (interface{}, error) {
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	t.Setenv("DATABASE_DSN", "")

	cfg0 := env.GetConfig()
	cfg1 := cfg0.(env.TestConfig)
	cfg2 := cfg1.GetTestConfig(
		env.WithTestDBPool("", nil),
		env.WithEtcdClientConfig(clientV3.Config{Endpoints: []string{"localhost:0"}}),
	)
	ctrl := gomock.NewController(t)
	etcdRepo := NewMockRepo[domain.Actioner[*entity.KeyValue, entity.KeyValue], *entity.KeyValue, entity.KeyValue](ctrl)
	etcdPoolMock := NewMockEtcdPool(ctrl)
	postgresRepo := NewMockRepo[domain.Actioner[*entity.KeyValue, entity.KeyValue], *entity.KeyValue, entity.KeyValue](ctrl)
	etcdRepo.
		EXPECT().
		Do(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(entity.KeyValue{}, nil).
		AnyTimes()
	etcdPoolMock.
		EXPECT().
		AcquireClient().
		Return(&kvTestStub{}, nil).
		AnyTimes()
	etcdPoolMock.
		EXPECT().
		ReleaseClient(gomock.Any()).
		Return(nil).
		AnyTimes()
	postgresRepo.
		EXPECT().
		Do(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(entity.KeyValue{}, nil).
		AnyTimes()

	srv := newTestKeyValueDataService(cfg2, etcdRepo, etcdPoolMock, postgresRepo)

	return srv, srv.Put(context.Background(), entity.KeyValue{})
}

func positivePostgresPut1Check(t *testing.T, i interface{}) bool {
	_, ok := i.(*keyValueDataService)
	return ok
}

func newTestKeyValueDataService(
	cfg env.Config,
	etcdRepo domain.Repo[domain.Actioner[*entity.KeyValue, entity.KeyValue], *entity.KeyValue, entity.KeyValue],
	pool pool.EtcdPool,
	postgresRepo domain.Repo[domain.Actioner[*entity.KeyValue, entity.KeyValue], *entity.KeyValue, entity.KeyValue],
) *keyValueDataService {
	inst := new(keyValueDataService)
	inst.cache = memory.New(memory.Config{
		GCInterval: cfg.CacheGCInterval(),
	})
	inst.cacheExpire = cfg.CacheExpire()
	inst.etcdRepo = etcdRepo
	inst.pool = pool
	inst.postgresRepo = postgresRepo
	inst.sLog = slog.Default()

	return inst
}

var _ clientV3.KV = (*kvTestStub)(nil)

type kvTestStub struct {
}

func (k *kvTestStub) Put(_ context.Context, _, _ string, _ ...clientV3.OpOption) (*clientV3.PutResponse, error) {
	return &clientV3.PutResponse{}, nil
}

func (k *kvTestStub) Get(_ context.Context, _ string, _ ...clientV3.OpOption) (*clientV3.GetResponse, error) {
	kvs := make([]*mvccpb.KeyValue, 0)
	kvs = append(kvs, &mvccpb.KeyValue{Key: []byte("key1"), Value: []byte("value1"), Version: 1})
	return &clientV3.GetResponse{
		Header: &etcdserverpb.ResponseHeader{},
		Kvs:    kvs,
	}, nil
}

func (k *kvTestStub) Delete(_ context.Context, _ string, _ ...clientV3.OpOption) (*clientV3.DeleteResponse, error) {
	return &clientV3.DeleteResponse{}, nil
}

func (k *kvTestStub) Compact(_ context.Context, _ int64, _ ...clientV3.CompactOption) (*clientV3.CompactResponse, error) {
	return &clientV3.CompactResponse{}, nil
}

func (k *kvTestStub) Do(_ context.Context, _ clientV3.Op) (clientV3.OpResponse, error) {
	return clientV3.OpResponse{}, nil
}

func (k *kvTestStub) Txn(_ context.Context) clientV3.Txn {
	return nil
}
