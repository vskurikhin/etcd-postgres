/*
 * This file was last modified at 2024-07-31 14:12 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * flag_parse.go
 * $Id$
 */
//!+

// Package env работа с настройками и окружением.
package env

import (
	"github.com/spf13/pflag"
	"time"
)

const (
	flagCacheExpireMs      = "cache-expire-ms"
	flagCacheGCIntervalSec = "cache-gc-interval-sec"
	flagDatabaseDSN        = "database-dsn"
	flagDebug              = "debug"
	flagEtcdAddresses      = "etcd-addresses"
	flagEtcdDialTimeout    = "etcd-dial-timeout"
	flagGRPCAddress        = "grpc-address"
	flagGRPCCAFile         = "grpc-ca-file"
	flagGRPCCertFile       = "grpc-cert-file"
	flagGRPCKeyFile        = "grpc-key-file"
	flagHTTPAddress        = "http-address"
	flagHTTPCAFile         = "http-ca-file"
	flagHTTPCertFile       = "http-cert-file"
	flagHTTPKeyFile        = "http-key-file"
	flagSlogJson           = "slog-json"
)

func makeFlagsParse() map[string]interface{} {

	var flagsMap = make(map[string]interface{})

	if !pflag.Parsed() {
		flagsMap[flagCacheExpireMs] = pflag.Int(
			flagCacheExpireMs,
			1000,
			"time to expire key in millisecond",
		)
		flagsMap[flagCacheGCIntervalSec] = pflag.Int(
			flagCacheGCIntervalSec,
			10,
			"time before deleting expired keys in second",
		)
		flagsMap[flagDatabaseDSN] = pflag.StringP(
			flagDatabaseDSN,
			"b",
			"postgres://dbuser:password@localhost:5432/db?sslmode=disable",
			"database DSN",
		)
		flagsMap[flagDebug] = pflag.BoolP(
			flagDebug,
			"d",
			false,
			"debug logging level",
		)
		flagsMap[flagEtcdAddresses] = pflag.StringP(
			flagEtcdAddresses,
			"e",
			"localhost:1379,localhost:2379,localhost:3379",
			"etcd servers host and port",
		)
		flagsMap[flagEtcdDialTimeout] = pflag.Duration(
			flagEtcdDialTimeout,
			2*time.Second,
			"etcd servers host and port",
		)
		flagsMap[flagGRPCAddress] = pflag.StringP(
			flagGRPCAddress,
			"g",
			"localhost:8443",
			"gRPC server host and port",
		)
		flagsMap[flagGRPCCAFile] = pflag.String(
			flagGRPCCAFile,
			"cert/grpc-test_ca-cert.pem",
			"gRPC CA file",
		)
		flagsMap[flagGRPCCertFile] = pflag.String(
			flagGRPCCertFile,
			"cert/grpc-test_server-cert.pem",
			"gRPC server certificate file",
		)
		flagsMap[flagGRPCKeyFile] = pflag.String(
			flagGRPCKeyFile,
			"cert/grpc-test_server-key.pem",
			"gRPC server key file",
		)
		flagsMap[flagHTTPAddress] = pflag.StringP(
			flagHTTPAddress,
			"a",
			"localhost:443",
			"HTTP server host and port",
		)
		flagsMap[flagHTTPCAFile] = pflag.String(
			flagHTTPCAFile,
			"cert/http-test_ca-cert.pem",
			"HTTP CA file",
		)
		flagsMap[flagHTTPCertFile] = pflag.String(
			flagHTTPCertFile,
			"cert/http-test_server-cert.pem",
			"HTTP server certificate file",
		)
		flagsMap[flagHTTPKeyFile] = pflag.String(
			flagHTTPKeyFile,
			"cert/http-test_server-key.pem",
			"HTTP server key file",
		)
		flagsMap[flagSlogJson] = pflag.BoolP(
			flagSlogJson,
			"s",
			false,
			"slog JSON output",
		)
		pflag.Parse()
	}
	return flagsMap
}

func setIfFlagChanged(name string, set func()) {
	pflag.VisitAll(func(f *pflag.Flag) {
		if f.Changed && f.Name == name {
			set()
		}
	})
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
