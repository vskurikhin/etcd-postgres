/*
 * Copyright text:
 * This file was last modified at 2024-07-10 20:32 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * yaml_load_test.go
 * $Id$
 */
//!+

// Package env работа с настройками и окружением.
package env

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type wantLoadConfig struct {
	yamlConfig YamlConfig
	err        error
}

type testLoadConfig struct {
	name  string
	input string
	fRun  func(string) (YamlConfig, error)
	want  wantLoadConfig
}

func TestLoadConfig(t *testing.T) {
	assert.NotNil(t, t)
	var tests = []testLoadConfig{
		{
			name:  `test #1 positive LoadConfig(".")`,
			input: ".",
			fRun:  LoadConfig,
			want: wantLoadConfig{
				yamlConfig: &yamlConfig{EtcdClient: struct {
					Cache struct {
						Enabled     bool
						cacheConfig `mapstructure:",squash"`
					}
					GRPC struct {
						Enabled    bool
						grpcConfig `mapstructure:",squash"`
						TLS        struct {
							Enabled   bool
							tlsConfig `mapstructure:",squash"`
						}
					}
					HTTP struct {
						Enabled    bool
						httpConfig `mapstructure:",squash"`
						TLS        struct {
							Enabled   bool
							tlsConfig `mapstructure:",squash"`
						}
					}
				}{
					Cache: struct {
						Enabled     bool
						cacheConfig `mapstructure:",squash"`
					}{
						Enabled: true,
						cacheConfig: cacheConfig{
							ExpireMs:      1000,
							GCIntervalSec: 10,
						},
					},
					GRPC: struct {
						Enabled    bool
						grpcConfig `mapstructure:",squash"`
						TLS        struct {
							Enabled   bool
							tlsConfig `mapstructure:",squash"`
						}
					}{
						Enabled: true,
						grpcConfig: grpcConfig{
							Address: "localhost",
							Port:    8442,
							Proto:   "tcp",
						},
						TLS: struct {
							Enabled   bool
							tlsConfig `mapstructure:",squash"`
						}{
							Enabled: true,
							tlsConfig: tlsConfig{
								CAFile:   "cert/grpc-test_ca-cert.pem",
								CertFile: "cert/grpc-test_server-cert.pem",
								KeyFile:  "cert/grpc-test_server-key.pem",
							},
						},
					},
					HTTP: struct {
						Enabled    bool
						httpConfig `mapstructure:",squash"`
						TLS        struct {
							Enabled   bool
							tlsConfig `mapstructure:",squash"`
						}
					}{
						Enabled: true,
						httpConfig: httpConfig{
							Address: "localhost",
							Port:    8443,
						},
						TLS: struct {
							Enabled   bool
							tlsConfig `mapstructure:",squash"`
						}{
							Enabled: true,
							tlsConfig: tlsConfig{
								CAFile:   "cert/http-test_ca-cert.pem",
								CertFile: "cert/http-test_server-cert.pem",
								KeyFile:  "cert/http-test_server-key.pem",
							},
						},
					},
				}},
				err: nil,
			},
		},
		{
			name:  `test #2 positive GO_FAVORITES_SKIP_LOAD_CONFIG=True LoadConfig("")`,
			input: "test",
			fRun: func(s string) (YamlConfig, error) {
				t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
				return LoadConfig(s)
			},
			want: wantLoadConfig{
				yamlConfig: &yamlConfig{},
				err:        nil,
			},
		},
	}
	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.fRun(test.input)
			assert.Equal(t, test.want.yamlConfig, got)
			assert.Equal(t, test.want.err, err)
		})
	}
}

func configFileNotFoundError(path string) error {

	wd := chDir()
	defer func() { _ = os.Chdir(wd) }()

	viper.SetConfigName("etcd-client")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/etcd-client/")
	viper.AddConfigPath(path)

	return viper.ReadInConfig()
}

func chDir() string {
	wd, _ := os.Getwd()
	if err := os.Chdir("test"); err != nil {
		fmt.Println("Chdir error:", err)
	}
	return wd
}
