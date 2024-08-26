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

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/victor-skurikhin/etcd-client/v1/internal/domain"
	"github.com/victor-skurikhin/etcd-client/v1/internal/domain/entity"
	"github.com/victor-skurikhin/etcd-client/v1/internal/env"
	"github.com/victor-skurikhin/etcd-client/v1/pool"
	"github.com/victor-skurikhin/etcd-client/v1/tool"
	"go.etcd.io/etcd/api/v3/etcdserverpb"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientV3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/mock/gomock"
	"log/slog"
	"math"
	"testing"
)

func TestEtcd(t *testing.T) {
	for _, test := range []struct {
		name string
		fRun func(*testing.T) (interface{}, error)
		want func(*testing.T, interface{}) bool
	}{
		{
			"test #0 positive for function GetKeyValueEtcdRepo(env.Config)",
			positiveGetKeyValueEtcdRepo,
			positiveGetKeyValueEtcdRepoCheck,
		},
		{
			"test #1 positive #1 for struct Etcd method Do(context.Context, A, U, func(domain.Scanner))",
			positiveEtcdDo1,
			positiveEtcdDo1Check,
		},
		{
			"test #2 positive #2 for struct Etcd method Do(context.Context, A, U, func(domain.Scanner))",
			positiveEtcdDo2,
			positiveEtcdDo2Check,
		},
		{
			"test #3 positive #3 for struct Etcd method Do(context.Context, A, U, func(domain.Scanner))",
			positiveEtcdDo3,
			positiveEtcdDo3Check,
		},
		{
			"test #4 negative #1 for struct Etcd method Do(context.Context, A, U, func(domain.Scanner))",
			negativeEtcdDo1,
			negativeEtcdDo1Check,
		},
		{
			"test #5 negative #2 for struct Etcd method Do(context.Context, A, U, func(domain.Scanner))",
			negativeEtcdDo2,
			negativeEtcdDo2Check,
		},
		{
			"test #6 negative #3 for struct Etcd method Do(context.Context, A, U, func(domain.Scanner))",
			negativeEtcdDoScanner3,
			negativeEtcdDoScanner3Check,
		},
		{
			"test #6 negative #4 for struct Etcd method Do(context.Context, A, U, func(domain.Scanner))",
			negativeEtcdDoScanner4,
			negativeEtcdDoScanner4Check,
		},
		{
			"test #7 negative #5 for struct Etcd method Do(context.Context, A, U, func(domain.Scanner))",
			negativeEtcdDoScanner5,
			negativeEtcdDoScanner5Check,
		},
		{
			"test #8 negative #6 for struct Etcd method Do(context.Context, A, U, func(domain.Scanner))",
			negativeEtcdDoScanner6,
			negativeEtcdDoScanner6Check,
		},
		{
			"test #9 positive #1 for struct Etcd method Get(context.Context, A, U, func(domain.Scanner))",
			positiveEtcdGet1,
			positiveEtcdGet1Check,
		},
		{
			"test #10 positive #2 for struct Etcd method Get(context.Context, A, U, func(domain.Scanner))",
			positiveEtcdGet2,
			positiveEtcdGet2Check,
		},
		{
			"test #11 negative #1 for struct Etcd method Get(context.Context, A, U, func(domain.Scanner))",
			negativeEtcdGet1,
			negativeEtcdGet1Check,
		},
		{
			"test #12 negative #2 for struct Etcd method Get(context.Context, A, U, func(domain.Scanner))",
			negativeEtcdGet2,
			negativeEtcdGet2Check,
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

func positiveGetKeyValueEtcdRepo(t *testing.T) (interface{}, error) {

	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	t.Setenv("ETCD_ADDRESSES", "")
	tool.SetLogger(slog.Default())
	c := env.GetConfig()
	cfg := c.(env.TestConfig)

	return GetKeyValueEtcdRepo(cfg.GetTestConfig(env.WithEtcdClientConfig(clientV3.Config{}))), nil
}

func positiveGetKeyValueEtcdRepoCheck(_ *testing.T, i interface{}) bool {
	_, ok := i.(*Etcd[domain.Actioner[*entity.KeyValue, entity.KeyValue], *entity.KeyValue, entity.KeyValue])
	return ok
}

func positiveEtcdDo1(t *testing.T) (interface{}, error) {

	ctrl := gomock.NewController(t)
	etcdPoolMock := NewMockEtcdPool(ctrl)
	etcdKeyValueInst = newTestKeyValueEtcdRepo(etcdPoolMock, slog.Default())

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

	expected := entity.MakeKeyValue("key1", "value1", 1, entity.DefaultTAttributes())
	ctx := context.Background()
	result, err := etcdKeyValueInst.Do(ctx, entity.KeyValueUpsert, expected, func(domain.Scanner) entity.KeyValue {
		return expected
	})
	assert.Nil(t, err)
	assert.Equal(t, expected, result)

	return result, nil
}

func positiveEtcdDo1Check(_ *testing.T, i interface{}) bool {
	return i == entity.MakeKeyValue("key1", "value1", 1, entity.DefaultTAttributes())
}

func positiveEtcdDo2(t *testing.T) (interface{}, error) {

	ctrl := gomock.NewController(t)
	etcdPoolMock := NewMockEtcdPool(ctrl)
	etcdKeyValueInst = newTestKeyValueEtcdRepo(etcdPoolMock, slog.Default())

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

	search := entity.MakeKeyValue("key1", "", 0, entity.DefaultTAttributes())
	expected := entity.MakeKeyValue("key1", "value1", 1, entity.DefaultTAttributes())
	ctx := context.Background()
	result, err := etcdKeyValueInst.Do(ctx, entity.KeyValueSelect, search, func(s domain.Scanner) entity.KeyValue {
		var key, value string
		var version sql.NullInt64
		_ = s.Scan(&key, &value, &version)
		return entity.MakeKeyValue(key, value, version.Int64, entity.DefaultTAttributes())
	})
	assert.Nil(t, err)
	assert.Equal(t, expected, result)

	return result, nil
}

func positiveEtcdDo2Check(t *testing.T, i interface{}) bool {
	return i == entity.MakeKeyValue("key1", "value1", 1, entity.DefaultTAttributes())
}

func positiveEtcdDo3(t *testing.T) (interface{}, error) {

	ctrl := gomock.NewController(t)
	etcdPoolMock := NewMockEtcdPool(ctrl)
	etcdKeyValueInst = newTestKeyValueEtcdRepo(etcdPoolMock, slog.Default())

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

	search := entity.MakeKeyValue("key1", "", 0, entity.DefaultTAttributes())
	ctx := context.Background()
	result, err := etcdKeyValueInst.Do(ctx, entity.KeyValueDelete, search, func(domain.Scanner) entity.KeyValue {
		return search
	})
	fmt.Printf("err: %v\n", err)
	assert.Nil(t, err)

	return result, nil
}

func positiveEtcdDo3Check(t *testing.T, i interface{}) bool {
	return i == entity.MakeKeyValue("key1", "", 0, entity.DefaultTAttributes())
}

func negativeEtcdDo1(t *testing.T) (interface{}, error) {

	ctrl := gomock.NewController(t)
	etcdPoolMock := NewMockEtcdPool(ctrl)
	etcdKeyValueInst = newTestKeyValueEtcdRepo(etcdPoolMock, slog.Default())

	etcdPoolMock.
		EXPECT().
		AcquireClient().
		Return(nil, fmt.Errorf("connectivity state: INVALID_STATE")).
		AnyTimes()

	etcdPoolMock.
		EXPECT().
		ReleaseClient(gomock.Any()).
		Return(nil).
		AnyTimes()

	test := entity.MakeKeyValue("key1", "value1", 1, entity.DefaultTAttributes())
	ctx := context.Background()
	_, err := etcdKeyValueInst.Do(ctx, entity.KeyValueUpsert, test, func(domain.Scanner) entity.KeyValue {
		return test
	})
	fmt.Printf("err: %T = %v\n", err, err)
	assert.NotNil(t, err)

	return err, nil
}

func negativeEtcdDo1Check(t *testing.T, i interface{}) bool {
	if etcdError, ok := i.(EtcdError); ok {
		return etcdError.err.Error() == fmt.Errorf("connectivity state: %s", "INVALID_STATE").Error()
	}
	return false
}

func negativeEtcdDo2(t *testing.T) (interface{}, error) {

	ctrl := gomock.NewController(t)
	etcdPoolMock := NewMockEtcdPool(ctrl)
	etcdKeyValueInst = newTestKeyValueEtcdRepo(etcdPoolMock, slog.Default())

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

	test := entity.MakeKeyValue("key1", "value1", 1, entity.DefaultTAttributes())
	ctx := context.Background()
	_, err := etcdKeyValueInst.Do(ctx, keyValueActionUnknownStub{}, test, func(domain.Scanner) entity.KeyValue {
		return test
	})
	assert.NotNil(t, err)

	return err, nil
}

func negativeEtcdDo2Check(t *testing.T, i interface{}) bool {
	if etcdError, ok := i.(EtcdError); ok {
		return etcdError.err.Error() == fmt.Errorf("unknown action, name: %s", keyValueActionUnknownStub{}.Name()).Error()
	}
	return false
}

func negativeEtcdDoScanner3(t *testing.T) (interface{}, error) {

	ctrl := gomock.NewController(t)
	etcdPoolMock := NewMockEtcdPool(ctrl)
	etcdKeyValueInst = newTestKeyValueEtcdRepo(etcdPoolMock, slog.Default())

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

	var erResult error
	search := entity.MakeKeyValue("key1", "", 0, entity.DefaultTAttributes())
	ctx := context.Background()
	_, err := etcdKeyValueInst.Do(ctx, entity.KeyValueSelect, search, func(s domain.Scanner) entity.KeyValue {
		var key, value string
		var version sql.NullInt64
		erResult = s.Scan()
		return entity.MakeKeyValue(key, value, version.Int64, entity.DefaultTAttributes())
	})
	assert.Nil(t, err)
	assert.NotNil(t, erResult)

	return erResult, nil
}

func negativeEtcdDoScanner3Check(t *testing.T, i interface{}) bool {
	if scannerError, ok := i.(ScannerError); ok {
		return scannerError.err.Error() == fmt.Errorf("no required parameters, length: %d", 0).Error()
	}
	return false
}

func negativeEtcdDoScanner4(t *testing.T) (interface{}, error) {

	ctrl := gomock.NewController(t)
	etcdPoolMock := NewMockEtcdPool(ctrl)
	etcdKeyValueInst = newTestKeyValueEtcdRepo(etcdPoolMock, slog.Default())

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

	var erResult error
	search := entity.MakeKeyValue("key1", "", 0, entity.DefaultTAttributes())
	ctx := context.Background()
	_, err := etcdKeyValueInst.Do(ctx, entity.KeyValueSelect, search, func(s domain.Scanner) entity.KeyValue {
		var badKeyType bool
		var value string
		var version sql.NullInt64
		erResult = s.Scan(&badKeyType, &value, &version)
		return entity.MakeKeyValue("badKeyType", value, version.Int64, entity.DefaultTAttributes())
	})
	assert.Nil(t, err)
	assert.NotNil(t, erResult)

	return erResult, nil
}

func negativeEtcdDoScanner4Check(t *testing.T, i interface{}) bool {
	if scannerError, ok := i.(ScannerError); ok {
		var b bool
		return scannerError.err.Error() == fmt.Errorf("argument #0 is not pointer to string, a type: %T", &b).Error()
	}
	return false
}

func negativeEtcdDoScanner5(t *testing.T) (interface{}, error) {

	ctrl := gomock.NewController(t)
	etcdPoolMock := NewMockEtcdPool(ctrl)
	etcdKeyValueInst = newTestKeyValueEtcdRepo(etcdPoolMock, slog.Default())

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

	var erResult error
	search := entity.MakeKeyValue("key1", "", 0, entity.DefaultTAttributes())
	ctx := context.Background()
	_, err := etcdKeyValueInst.Do(ctx, entity.KeyValueSelect, search, func(s domain.Scanner) entity.KeyValue {
		var key string
		var badValueType bool
		var version sql.NullInt64
		erResult = s.Scan(&key, &badValueType, &version)
		return entity.MakeKeyValue(key, "badValueType", version.Int64, entity.DefaultTAttributes())
	})
	assert.Nil(t, err)
	assert.NotNil(t, erResult)

	return erResult, nil
}

func negativeEtcdDoScanner5Check(t *testing.T, i interface{}) bool {
	if scannerError, ok := i.(ScannerError); ok {
		var b bool
		return scannerError.err.Error() == fmt.Errorf("argument #1 is not pointer to string, a type: %T", &b).Error()
	}
	return false
}

func negativeEtcdDoScanner6(t *testing.T) (interface{}, error) {

	ctrl := gomock.NewController(t)
	etcdPoolMock := NewMockEtcdPool(ctrl)
	etcdKeyValueInst = newTestKeyValueEtcdRepo(etcdPoolMock, slog.Default())

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

	var erResult error
	search := entity.MakeKeyValue("key1", "", 0, entity.DefaultTAttributes())
	ctx := context.Background()
	_, err := etcdKeyValueInst.Do(ctx, entity.KeyValueSelect, search, func(s domain.Scanner) entity.KeyValue {
		var key string
		var value string
		var badVersionType bool
		erResult = s.Scan(&key, &value, &badVersionType)
		return entity.MakeKeyValue(key, value, math.MinInt64, entity.DefaultTAttributes())
	})
	assert.Nil(t, err)
	assert.NotNil(t, erResult)

	return erResult, nil
}

func negativeEtcdDoScanner6Check(t *testing.T, i interface{}) bool {
	if scannerError, ok := i.(ScannerError); ok {
		var b bool
		return scannerError.err.Error() == fmt.Errorf("argument #2 is not pointer to sql.NullInt64, a type: %T", &b).Error()
	}
	return false
}

func positiveEtcdGet1(t *testing.T) (interface{}, error) {

	ctrl := gomock.NewController(t)
	etcdPoolMock := NewMockEtcdPool(ctrl)
	etcdKeyValueInst = newTestKeyValueEtcdRepo(etcdPoolMock, slog.Default())

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

	expected := entity.MakeKeyValue("key1", "value1", 1, entity.DefaultTAttributes())
	ctx := context.Background()
	result, err := etcdKeyValueInst.Get(ctx, entity.KeyValueGetAll, expected, func(domain.Scanner) entity.KeyValue {
		return expected
	})
	assert.Nil(t, err)
	assert.NotNil(t, result)

	return result, nil
}

func positiveEtcdGet1Check(t *testing.T, i interface{}) bool {
	if result, ok := i.([]entity.KeyValue); ok && len(result) > 0 {
		return result[0] == entity.MakeKeyValue("key1", "value1", 1, entity.DefaultTAttributes())
	}
	return false
}

func positiveEtcdGet2(t *testing.T) (interface{}, error) {

	ctrl := gomock.NewController(t)
	etcdPoolMock := NewMockEtcdPool(ctrl)
	etcdKeyValueInst = newTestKeyValueEtcdRepo(etcdPoolMock, slog.Default())

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

	search := entity.MakeKeyValue("key", "", 0, entity.DefaultTAttributes())
	expected := entity.MakeKeyValue("key1", "value1", 1, entity.DefaultTAttributes())
	ctx := context.Background()
	result, err := etcdKeyValueInst.Get(ctx, entity.KeyValueSelect, search, func(domain.Scanner) entity.KeyValue {
		return expected
	})
	assert.Nil(t, err)
	assert.NotNil(t, result)

	return result, nil
}

func positiveEtcdGet2Check(t *testing.T, i interface{}) bool {
	if result, ok := i.([]entity.KeyValue); ok && len(result) > 0 {
		return result[0] == entity.MakeKeyValue("key1", "value1", 1, entity.DefaultTAttributes())
	}
	return false
}

func negativeEtcdGet1(t *testing.T) (interface{}, error) {

	ctrl := gomock.NewController(t)
	etcdPoolMock := NewMockEtcdPool(ctrl)
	etcdKeyValueInst = newTestKeyValueEtcdRepo(etcdPoolMock, slog.Default())

	etcdPoolMock.
		EXPECT().
		AcquireClient().
		Return(nil, fmt.Errorf("connectivity state: INVALID_STATE")).
		AnyTimes()

	etcdPoolMock.
		EXPECT().
		ReleaseClient(gomock.Any()).
		Return(nil).
		AnyTimes()

	test := entity.MakeKeyValue("key1", "value1", 1, entity.DefaultTAttributes())
	ctx := context.Background()
	_, err := etcdKeyValueInst.Get(ctx, entity.KeyValueGetAll, test, func(domain.Scanner) entity.KeyValue {
		return test
	})
	fmt.Printf("err: %T = %v\n", err, err)
	assert.NotNil(t, err)

	return err, nil
}

func negativeEtcdGet1Check(t *testing.T, i interface{}) bool {
	if etcdError, ok := i.(EtcdError); ok {
		return etcdError.err.Error() == fmt.Errorf("connectivity state: %s", "INVALID_STATE").Error()
	}
	return false
}

func negativeEtcdGet2(t *testing.T) (interface{}, error) {

	ctrl := gomock.NewController(t)
	etcdPoolMock := NewMockEtcdPool(ctrl)
	etcdKeyValueInst = newTestKeyValueEtcdRepo(etcdPoolMock, slog.Default())

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

	test := entity.MakeKeyValue("key1", "value1", 1, entity.DefaultTAttributes())
	ctx := context.Background()
	_, err := etcdKeyValueInst.Get(ctx, keyValueActionUnknownStub{}, test, func(domain.Scanner) entity.KeyValue {
		return test
	})
	assert.NotNil(t, err)

	return err, nil
}

func negativeEtcdGet2Check(t *testing.T, i interface{}) bool {
	if etcdError, ok := i.(EtcdError); ok {
		return etcdError.err.Error() == fmt.Errorf("unknown action, name: %s", keyValueActionUnknownStub{}.Name()).Error()
	}
	return false
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

type keyValueActionUnknownStub struct{}

func (k keyValueActionUnknownStub) Args(e entity.KeyValue) []any {
	return []any{}
}

func (k keyValueActionUnknownStub) Name() string {
	return "UnknownStub"
}

func (k keyValueActionUnknownStub) SQL() string {
	return ""
}

func newTestKeyValueEtcdRepo(
	pool pool.EtcdPool, sLog *slog.Logger,
) *Etcd[domain.Actioner[*entity.KeyValue, entity.KeyValue], *entity.KeyValue, entity.KeyValue] {
	repo := new(Etcd[domain.Actioner[*entity.KeyValue, entity.KeyValue], *entity.KeyValue, entity.KeyValue])
	repo.pool = pool
	repo.sLog = sLog
	return repo
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
