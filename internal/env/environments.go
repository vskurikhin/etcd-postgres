/*
 * This file was last modified at 2024-08-06 17:50 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * environments.go
 * $Id$
 */
//!+

// Package env работа с настройками и окружением.
package env

import (
	"fmt"
	"time"

	c0env "github.com/caarlos0/env"
)

// Environments статичная конфигурация из переменных окружения.
type environments struct {
	CacheExpireMs      int           `env:"CACHE_EXPIRE_MS"`
	CacheGCIntervalSec int           `env:"CACHE_GC_INTERVAL_SEC"`
	DataBaseDSN        string        `env:"DATABASE_DSN"`
	EtcdAddresses      []string      `env:"ETCD_ADDRESSES" envSeparator:","`
	EtcdDialTimeout    time.Duration `env:"ETCD_DIAL_TIMEOUT"`
	GRPCAddress        []string      `env:"GRPC_ADDRESS" envSeparator:":"`
	GRPCCAFile         string        `env:"GRPC_CA_FILE"`
	GRPCCertFile       string        `env:"GRPC_CERT_FILE"`
	GRPCKeyFile        string        `env:"GRPC_KEY_FILE"`
	HTTPAddress        []string      `env:"HTTP_ADDRESS" envSeparator:":"`
	HTTPCAFile         string        `env:"HTTP_CA_FILE"`
	HTTPCertFile       string        `env:"HTTP_CERT_FILE"`
	HTTPKeyFile        string        `env:"HTTP_KEY_FILE"`
}

func getEnvironments() (env *environments, err error) {

	env = new(environments)
	err = c0env.Parse(env)
	return
}

func (e *environments) String() string {
	return fmt.Sprintf(
		`CACHE_EXPIRE_MS: %d
CACHE_GC_INTERVAL_SEC: %d
GRPC_ADDRESS: %s
GRPC_CA_FILE: %s
GRPC_CERT_FILE: %s
GRPC_KEY_FILE: %s
HTTP_ADDRESS: %s
HTTP_CA_FILE: %s
HTTP_CERT_FILE: %s
HTTP_KEY_FILE: %s`,
		e.CacheExpireMs,
		e.CacheGCIntervalSec,
		e.GRPCAddress,
		e.GRPCCAFile,
		e.GRPCCertFile,
		e.GRPCKeyFile,
		e.HTTPAddress,
		e.HTTPCAFile,
		e.HTTPCertFile,
		e.HTTPKeyFile)
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
