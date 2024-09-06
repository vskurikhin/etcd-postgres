package repo

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/victor-skurikhin/etcd-client/v1/internal/domain"
	"github.com/victor-skurikhin/etcd-client/v1/internal/env"
	"github.com/victor-skurikhin/etcd-client/v1/pool"
	"github.com/victor-skurikhin/etcd-client/v1/pool/etcd_pool"
	"github.com/victor-skurikhin/etcd-client/v1/tool"
	"log/slog"
	"math"
	"os"
	"runtime"
	"testing"
)

const BenchmarkRepeat = 10

var (
	cfg               env.Config
	dbURL, etcdIP     string
	etcdContainer     testcontainers.Container
	etcdPool          pool.EtcdPool
	etcdPort          nat.Port
	postgresContainer *postgres.PostgresContainer
)

func TestMain(m *testing.M) {

	var err error
	ctx := context.Background()
	etcdContainer, etcdIP, etcdPort, err = createEtcdContainer(ctx)

	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	// Clean up the container
	defer func() { _ = etcdContainer.Terminate(ctx) }()

	postgresContainer, dbURL, err = createPostgresContainer(ctx)

	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	// Clean up the container
	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
	}()
	tool.SetLogger(slog.Default())
	cfg = getTestConfig()
	etcdPool = etcd_pool.GetEtcdPool(cfg)
	defer func() {
		_ = etcdPool.GracefulClose()
	}()
	fmt.Printf("etcdIP: %s, etcdPort: %s\ndbURL: %s\n", etcdIP, etcdPort, dbURL)
	os.Exit(m.Run())
}

func BenchmarkEtcd(b *testing.B) {

	ctx := context.Background()

	etcdRepo = new(Etcd[domain.Actioner[*IDValue, IDValue], *IDValue, IDValue])
	etcdRepo.pool = etcdPool
	etcdRepo.sLog = slog.Default()

	numberOps := runtime.NumCPU() * BenchmarkRepeat
	from := math.MaxInt / 4
	to := from + numberOps

	var res IDValue
	var etcdScan = getEtcdScanFunc(b, &res)

	b.ResetTimer()

	for k := 0; k < b.N; k++ {

		//var wg sync.WaitGroup
		//wg.Add(numberOps)

		for j := from; j < to; j++ {
			//go func(i int) {
			_, _ = etcdRepo.Do(ctx, IDValueUpsert, IDValue{id: j, value: fmt.Sprintf("value%d", j)}, etcdScan)
			//	wg.Done()
			//}(j)
		}
		//wg.Wait()
		for i := from; i < to; i++ {
			res, err := etcdRepo.Do(ctx, IDValueSelect, IDValue{id: i}, etcdScan)
			if err != nil {
				b.Fatal(err)
			}
			expected := IDValue{id: i, value: fmt.Sprintf("value%d", i), version: res.version}
			if res != expected {
				b.Fatal(err)
			}
		}
		for i := from; i < to; i++ {
			_, err := etcdRepo.Do(ctx, IDValueDelete, IDValue{id: i}, etcdScan)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkPostgres(b *testing.B) {

	ctx := context.Background()

	repoPostgres = new(Postgres[domain.Actioner[*IDValue, IDValue], *IDValue, IDValue])
	repoPostgres.pool = tool.DBConnect(dbURL)
	repoPostgres.sLog = slog.Default()

	numberOps := runtime.NumCPU() * BenchmarkRepeat
	from := math.MaxInt / 4
	to := from + numberOps

	var postgresScan = getPostgresScanFunc(b)

	b.ResetTimer()

	for k := 0; k < b.N; k++ {
		for i := from; i < to; i++ {
			_, err := repoPostgres.Do(ctx, IDValueUpsert, IDValue{id: i, value: fmt.Sprintf("value%d", i)}, postgresScan)
			if err != nil {
				b.Fatal(err)
			}
		}
		for i := from; i < to; i++ {
			res, err := repoPostgres.Do(ctx, IDValueSelect, IDValue{id: i}, postgresScan)
			if err != nil {
				b.Fatal(err)
			}
			expected := IDValue{id: i, value: fmt.Sprintf("value%d", i)}
			if res != expected {
				fmt.Printf("res: %v\n", res)
				b.Fatal(err)
			}
		}
		for i := from; i < to; i++ {
			_, err := repoPostgres.Do(ctx, IDValueDelete, IDValue{id: i},
				func(scanner domain.Scanner) IDValue {
					return IDValue{}
				})
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}
