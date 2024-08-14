/*
 * This file was last modified at 2024-08-06 18:20 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * yaml_config.go
 * $Id$
 */
//!+

// Package env работа с настройками и окружением.
package env

import (
	"fmt"
	"time"
)

var _ YamlConfig = (*yamlConfig)(nil)

// YamlConfig статичная конфигурация собранная из Yaml-файла.
type YamlConfig interface {
	fmt.Stringer
	CacheEnabled() bool
	CacheExpireMs() int
	CacheGCIntervalSec() int
	EtcdAddresses() []string
	EtcdEnabled() bool
	EtcdDialTimeout() time.Duration
	EtcdTLSCAFile() string
	EtcdTLSCertFile() string
	EtcdTLSEnabled() bool
	EtcdTLSKeyFile() string
	GRPCAddress() string
	GRPCEnabled() bool
	GRPCPort() int
	GRPCProto() string
	GRPCTLSCAFile() string
	GRPCTLSCertFile() string
	GRPCTLSEnabled() bool
	GRPCTLSKeyFile() string
	HTTPAddress() string
	HTTPEnabled() bool
	HTTPPort() int
	HTTPTLSCAFile() string
	HTTPTLSCertFile() string
	HTTPTLSEnabled() bool
	HTTPTLSKeyFile() string
}

type yamlConfig struct {
	EtcdClient struct {
		Cache struct {
			Enabled     bool
			cacheConfig `mapstructure:",squash"`
		}
		Etcd struct {
			Enabled    bool
			etcdConfig `mapstructure:",squash"`
			TLS        struct {
				Enabled   bool
				tlsConfig `mapstructure:",squash"`
			}
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
	}
}

type cacheConfig struct {
	ExpireMs      int `mapstructure:"expire_ms"`
	GCIntervalSec int `mapstructure:"gc_interval_sec"`
}

type etcdConfig struct {
	Addresses   []string
	DialTimeout time.Duration `mapstructure:"dial_timeout"`
}

type grpcConfig struct {
	Address string
	Port    int16
	Proto   string
}

type httpConfig struct {
	Address string
	Port    int16
}

type tlsConfig struct {
	CAFile   string `mapstructure:"ca_file"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

// CacheEnabled тумблер включения локального кэша.
func (y *yamlConfig) CacheEnabled() bool {

	if y != nil {
		return y.EtcdClient.Cache.Enabled
	}
	return false
}

// CacheExpireMs срок действия записи в кэше истекает в миллисекундах.
func (y *yamlConfig) CacheExpireMs() int {

	if y != nil {
		return y.EtcdClient.Cache.ExpireMs
	}
	return 0
}

// CacheGCIntervalSec интервал очистки кэша в секундах.
func (y *yamlConfig) CacheGCIntervalSec() int {

	if y != nil {
		return y.EtcdClient.Cache.GCIntervalSec
	}
	return 0
}

func (y *yamlConfig) EtcdAddresses() []string {

	if y != nil {
		return y.EtcdClient.Etcd.Addresses
	}
	return []string{}
}

func (y *yamlConfig) EtcdEnabled() bool {

	if y != nil {
		return y.EtcdClient.Etcd.Enabled
	}
	return false
}

func (y *yamlConfig) EtcdDialTimeout() time.Duration {

	if y != nil {
		return y.EtcdClient.Etcd.DialTimeout
	}
	return 0
}

func (y *yamlConfig) EtcdTLSCAFile() string {

	if y != nil {
		return y.EtcdClient.Etcd.TLS.CAFile
	}
	return ""
}

func (y *yamlConfig) EtcdTLSCertFile() string {

	if y != nil {
		return y.EtcdClient.Etcd.TLS.CertFile
	}
	return ""
}

func (y *yamlConfig) EtcdTLSEnabled() bool {

	if y != nil {
		return y.EtcdClient.Etcd.TLS.Enabled
	}
	return false
}

func (y *yamlConfig) EtcdTLSKeyFile() string {

	if y != nil {
		return y.EtcdClient.Etcd.TLS.KeyFile
	}
	return ""
}

// GRPCAddress адрес для выставления конечных точек gRPC-сервера.
func (y *yamlConfig) GRPCAddress() string {

	if y != nil {
		return y.EtcdClient.GRPC.Address
	}
	return ""
}

// GRPCEnabled тумблер включения gRPC-сервера.
func (y *yamlConfig) GRPCEnabled() bool {

	if y != nil {
		return y.EtcdClient.GRPC.Enabled
	}
	return false
}

// GRPCPort порт для gRPC-сервера.
func (y *yamlConfig) GRPCPort() int {

	if y != nil {
		return int(y.EtcdClient.GRPC.Port)
	}
	return 0
}

// GRPCProto протокол для gRPC-сервера.
func (y *yamlConfig) GRPCProto() string {

	if y != nil {
		return y.EtcdClient.GRPC.Proto
	}
	return ""
}

// GRPCTLSCAFile корневой сертификат центра сертификации который выдал TLS сертификат для gRPC-сервера.
func (y *yamlConfig) GRPCTLSCAFile() string {

	if y != nil {
		return y.EtcdClient.GRPC.TLS.CAFile
	}
	return ""
}

// GRPCTLSCertFile TLS сертификат для gRPC-сервера.
func (y *yamlConfig) GRPCTLSCertFile() string {

	if y != nil {
		return y.EtcdClient.GRPC.TLS.CertFile
	}
	return ""
}

// GRPCTLSKeyFile TLS ключ для gRPC-сервера.
func (y *yamlConfig) GRPCTLSKeyFile() string {

	if y != nil {
		return y.EtcdClient.GRPC.TLS.KeyFile
	}
	return ""
}

// GRPCTLSEnabled тумблер включения на gRPC-сервере TLS шифрования.
func (y *yamlConfig) GRPCTLSEnabled() bool {

	if y != nil {
		return y.EtcdClient.GRPC.TLS.Enabled
	}
	return false
}

// HTTPAddress адрес для выставления конечных точек HTTP-сервера.
func (y *yamlConfig) HTTPAddress() string {

	if y != nil {
		return y.EtcdClient.HTTP.Address
	}
	return ""
}

// HTTPEnabled тумблер включения HTTP-сервера.
func (y *yamlConfig) HTTPEnabled() bool {

	if y != nil {
		return y.EtcdClient.HTTP.Enabled
	}
	return false
}

// HTTPPort порт для HTTP-сервера.
func (y *yamlConfig) HTTPPort() int {

	if y != nil {
		return int(y.EtcdClient.HTTP.Port)
	}
	return 0
}

// HTTPTLSCAFile корневой сертификат центра сертификации который выдал TLS сертификат для HTTP-сервера.
func (y *yamlConfig) HTTPTLSCAFile() string {

	if y != nil {
		return y.EtcdClient.HTTP.TLS.CAFile
	}
	return ""
}

// HTTPTLSCertFile TLS сертификат для HTTP-сервера.
func (y *yamlConfig) HTTPTLSCertFile() string {

	if y != nil {
		return y.EtcdClient.HTTP.TLS.CertFile
	}
	return ""
}

// HTTPTLSKeyFile TLS ключ для HTTP-сервера.
func (y *yamlConfig) HTTPTLSKeyFile() string {

	if y != nil {
		return y.EtcdClient.HTTP.TLS.KeyFile
	}
	return ""
}

// HTTPTLSEnabled тумблер включения на HTTP-сервере TLS шифрования.
func (y *yamlConfig) HTTPTLSEnabled() bool {

	if y != nil {
		return y.EtcdClient.HTTP.TLS.Enabled
	}
	return false
}

func (y *yamlConfig) String() string {
	return fmt.Sprintf(
		`CacheEnabled: %v
CacheExpire: %d
CacheGCInterval: %d
GRPCAddress: %s
GRPCEnabled: %v
GRPCPort: %d
GRPCProto: %s
GRPCTLSCAFile: %s
GRPCTLSCertFile: %s
GRPCTLSKeyFile: %s
GRPCTLSEnabled: %v
HTTPAddress: %s
HTTPEnabled: %v
HTTPPort: %d
HTTPTLSCAFile: %s
HTTPTLSCertFile: %s
HTTPTLSEnabled: %v
HTTPTLSKeyFile: %s`,
		y.CacheEnabled(),
		y.CacheExpireMs(),
		y.CacheGCIntervalSec(),
		y.GRPCAddress(),
		y.GRPCEnabled(),
		y.GRPCPort(),
		y.GRPCProto(),
		y.GRPCTLSCAFile(),
		y.GRPCTLSCertFile(),
		y.GRPCTLSKeyFile(),
		y.GRPCTLSEnabled(),
		y.HTTPAddress(),
		y.HTTPEnabled(),
		y.HTTPPort(),
		y.HTTPTLSCAFile(),
		y.HTTPTLSCertFile(),
		y.HTTPTLSEnabled(),
		y.HTTPTLSKeyFile(),
	)
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
