package repo

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/victor-skurikhin/etcd-client/v1/internal/env"
	clientV3 "go.etcd.io/etcd/client/v3"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

func createEtcdContainer(ctx context.Context) (testcontainers.Container, string, nat.Port, error) {

	req := testcontainers.ContainerRequest{
		Cmd: []string{
			"/usr/local/bin/etcd",
			"--name=etcd-13",
			"--initial-advertise-peer-urls=http://etcd-13:2380",
			"--listen-peer-urls=http://0.0.0.0:2380",
			"--listen-client-urls=http://0.0.0.0:2379",
			"--advertise-client-urls=http://etcd-13:2379",
			"--initial-cluster=etcd-13=http://etcd-13:2380",
			"--heartbeat-interval=250",
			"--election-timeout=1250",
			"--initial-cluster-state=new",
			"--initial-cluster-token=localhost",
		},
		Image:        "quay.io/coreos/etcd:v3.5.15-arm64",
		ExposedPorts: []string{"2379/tcp", "2380/tcp"},
		WaitingFor: wait.ForExposedPort().
			WithPollInterval(200 * time.Millisecond).
			WithStartupTimeout(5 * time.Second),
	}

	etcdContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", "", err
	}
	ip, err := etcdContainer.Host(ctx)

	if err != nil {
		return nil, "", "", err
	}
	etcdPort, err := etcdContainer.MappedPort(ctx, "2379/tcp")

	if err != nil {
		return nil, "", "", err
	}
	return etcdContainer, ip, etcdPort, nil
}

func createPostgresContainer(ctx context.Context) (*postgres.PostgresContainer, string, error) {

	postgresContainer, err := creatingPostgresContainer(ctx)

	if err != nil {
		return nil, "", err
	}
	dbURL, err := postgresContainer.ConnectionString(ctx)

	if err != nil {
		return nil, "", err
	}
	return postgresContainer, dbURL, nil
}

func creatingPostgresContainer(ctx context.Context) (*postgres.PostgresContainer, error) {

	dbName := "db"
	dbUser := "dbuser"
	dbPassword := "password"

	return postgres.Run(ctx,
		"docker.io/postgres:13.15",
		postgres.WithInitScripts(filepath.Join("testdata", "init-db-test_id_value_test.sh")),
		//postgres.WithConfigFile(filepath.Join("testdata", "my-postgres.conf")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
}

func getTestConfig() env.Config {

	_ = os.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	_ = os.Setenv("DATABASE_DSN", dbURL)
	c := env.GetConfig()
	tst := c.(env.TestConfig)
	return tst.GetTestConfig(
		env.WithLogger(slog.Default()),
		env.WithEtcdClientConfig(clientV3.Config{
			Endpoints:   []string{fmt.Sprintf("%s:%s", etcdIP, etcdPort.Port())},
			DialTimeout: time.Second,
		}))
}
