package repo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/victor-skurikhin/etcd-client/v1/internal/domain"
	"github.com/victor-skurikhin/etcd-client/v1/tool"
	"log/slog"
	"sync"
	"testing"
)

const IntegrationRecordCount = 8

func TestIntegration(t *testing.T) {

	etcdRepo = new(Etcd[domain.Actioner[*IDValue, IDValue], *IDValue, IDValue])
	etcdRepo.pool = etcdPool
	etcdRepo.sLog = slog.Default()

	repoPostgres = new(Postgres[domain.Actioner[*IDValue, IDValue], *IDValue, IDValue])
	repoPostgres.pool = tool.DBConnect(dbURL)
	repoPostgres.sLog = slog.Default()

	for _, test := range []struct {
		name string
		fRun func(*testing.T) (interface{}, error)
		want func(*testing.T, interface{}) bool
	}{
		{
			"test #0 integration for repo Etcd",
			integrationEtcdRepo,
			integrationEtcdRepoCheck,
		},
		{
			"test #1 integration for repo Postgres",
			integrationPostgresRepo,
			integrationPostgresRepoCheck,
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

func integrationEtcdRepo(t *testing.T) (interface{}, error) {

	var ok bool
	ctx := context.Background()

	ok = etcdUpsertSelect(t, ctx)
	assert.True(t, ok)
	ok = ok && etcdDeleteSelect(t, ctx)
	assert.True(t, ok)
	ok = ok && etcdUpsertGetAll(t, ctx)
	assert.True(t, ok)

	return ok, nil
}

func integrationEtcdRepoCheck(_ *testing.T, i interface{}) bool {
	return i.(bool)
}

func integrationPostgresRepo(t *testing.T) (interface{}, error) {

	var ok bool
	ctx := context.Background()

	ok = postgresUpsertSelect(t, ctx)
	assert.True(t, ok)
	ok = ok && postgresDeleteSelect(t, ctx)
	assert.True(t, ok)
	ok = ok && postgresUpsertGetAll(t, ctx)
	assert.True(t, ok)

	return ok, nil
}

func integrationPostgresRepoCheck(t *testing.T, i interface{}) bool {
	return i.(bool)
}

func etcdUpsertSelect(t *testing.T, ctx context.Context) bool {

	var res IDValue
	var etcdScan = getEtcdScanFunc(t, &res)
	expected := IDValue{id: 99, value: "value99"}
	result, err := etcdRepo.Do(ctx, IDValueUpsert, expected, etcdScan)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
	result, er2 := etcdRepo.Do(ctx, IDValueSelect, IDValue{id: 99}, etcdScan)
	assert.Nil(t, er2)
	expected.version = res.version
	assert.Equal(t, expected, result)
	return result == expected
}

func etcdDeleteSelect(t *testing.T, ctx context.Context) bool {

	var res IDValue
	var etcdScan = getEtcdScanFunc(t, &res)
	expected := IDValue{id: 99}
	result, err := etcdRepo.Do(ctx, IDValueDelete, expected, etcdScan)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
	result, err = etcdRepo.Do(ctx, IDValueSelect, expected, etcdScan)
	assert.Equal(t, expected, result)
	_, ok := err.(EtcdError)
	assert.Equal(t, EtcdError{err: fmt.Errorf("no Kvs, length: 0")}, err)
	return ok
}

func etcdUpsertGetAll(t *testing.T, ctx context.Context) bool {

	var res IDValue
	var etcdScan = getEtcdScanFunc(t, &res)
	var wg sync.WaitGroup
	wg.Add(IntegrationRecordCount)

	for k := 0; k < IntegrationRecordCount; k++ {
		go func(i int) {
			_, _ = etcdRepo.Do(ctx, IDValueUpsert, IDValue{id: i, value: fmt.Sprintf("value%d", i)}, etcdScan)
			wg.Done()
		}(k)
	}
	wg.Wait()
	result, err := etcdRepo.Get(ctx, IDValueSelect, IDValue{id: 5}, etcdScan)
	assert.Nil(t, err)
	assert.True(t, len(result) > 1)
	result, er0 := etcdRepo.Get(ctx, IDValueGetAll, IDValue{}, etcdScan)
	assert.Nil(t, er0)
	return len(result) == IntegrationRecordCount
}

func postgresUpsertSelect(t *testing.T, ctx context.Context) bool {

	var postgresScan = getPostgresScanFunc(t)
	expected := IDValue{id: 99, value: "value99"}
	result, err := repoPostgres.Do(ctx, IDValueUpsert, expected, postgresScan)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
	result, er1 := repoPostgres.Do(ctx, IDValueSelect, IDValue{id: 99}, postgresScan)
	assert.Equal(t, expected, result)
	assert.Nil(t, er1)
	return result == IDValue{id: 99, value: "value99"}
}

func postgresDeleteSelect(t *testing.T, ctx context.Context) bool {

	expected := IDValue{id: 99, value: "value99"}
	result, err := repoPostgres.Do(ctx, IDValueDelete, expected, func(scanner domain.Scanner) IDValue {
		var res IDValue
		er0 := scanner.Scan(&res.id, &res.value, &res.version)
		assert.Nil(t, er0)
		return res
	})
	assert.Nil(t, err)
	assert.Equal(t, expected, result)

	var er2 error
	expected = IDValue{}
	result, er1 := repoPostgres.Do(ctx, IDValueSelect, expected, func(scanner domain.Scanner) IDValue {
		var res IDValue
		er2 = scanner.Scan(&res.id, &res.value, &res.version)
		assert.NotNil(t, er2)
		assert.Equal(t, pgx.ErrNoRows, er2)
		return res
	})
	assert.Nil(t, er1)
	assert.Equal(t, expected, result)
	return er2 == pgx.ErrNoRows
}

func postgresUpsertGetAll(t *testing.T, ctx context.Context) bool {

	var postgresScan = func(scanner domain.Scanner) IDValue {
		var res IDValue
		er0 := scanner.Scan(&res.id, &res.value, &res.version)
		assert.Nil(t, er0)
		return res
	}
	var wg sync.WaitGroup
	wg.Add(IntegrationRecordCount)

	for k := 0; k < IntegrationRecordCount; k++ {
		go func(i int) {
			_, _ = repoPostgres.Do(ctx, IDValueUpsert, IDValue{id: i, value: fmt.Sprintf("value%d", i)}, postgresScan)
			wg.Done()
		}(k)
	}
	wg.Wait()
	result, err := repoPostgres.Get(ctx, IDValueSelect, IDValue{id: 5}, postgresScan)
	assert.Nil(t, err)
	assert.True(t, len(result) > 0)
	result, er0 := repoPostgres.Get(ctx, IDValueGetAll, IDValue{}, postgresScan)
	assert.Nil(t, er0)
	return len(result) == IntegrationRecordCount
}
