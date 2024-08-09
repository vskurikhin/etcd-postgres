/*
 * This file was last modified at 2024-08-06 20:17 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * config.go
 * $Id$
 */
//!+

// Package env работа с настройками и окружением.
package env

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/victor-skurikhin/etcd-client/v1/tool"
	"google.golang.org/grpc/credentials"
)

const (
	propertyCacheExpireMs            = "cache-expire"
	propertyCacheGCIntervalSec       = "cache-gc-interval"
	propertyConfig                   = "yamlConfig"
	propertyEnvironments             = "environments"
	propertyFlags                    = "flags"
	propertyGRPCAddress              = "grpc-address"
	propertyGRPCTransportCredentials = "grpc-transport-credentials"
	propertyHTTPAddress              = "http-address"
	propertyHTTPHTTPTLSConfig        = "http-tls-yamlConfig"
	propertyLogger                   = "logger"
	MSG                              = "etcd-client "
)

// Config конфигурация собранная из Yaml-файла, переменных окружения и флагов командной строки.
type Config interface {
	fmt.Stringer
	CacheExpire() time.Duration
	CacheGCInterval() time.Duration
	YamlConfig() YamlConfig
	Environments() environments
	Flags() map[string]interface{}
	GRPCAddress() string
	GRPCTransportCredentials() credentials.TransportCredentials
	HTTPAddress() string
	HTTPTLSConfig() *tls.Config
	Logger() *slog.Logger
	SlogJSON() bool
}

type mapProperties struct {
	mp sync.Map
}

var properties Config = (*mapProperties)(nil)
var once = new(sync.Once)

// GetConfig — свойства преобразованные из конфигурации и окружения.
// потокобезопасное (thread-safe) создание.
func GetConfig() Config {

	once.Do(func() {
		yml, err := LoadConfig(".")
		tool.IfErrorThenPanic(err)
		env, err := getEnvironments()
		tool.IfErrorThenPanic(err)
		flm := makeFlagsParse()

		p := preparer{env: env, flagMap: flm, yml: yml}

		cacheExpire, err := p.getCacheExpire()
		slog.Debug(MSG+"GetConfig", "cacheExpire", cacheExpire, "err", err)
		cacheGCInterval, err := p.getCacheGCInterval()
		slog.Debug(MSG+"GetConfig", "cacheGCInterval", cacheGCInterval, "err", err)

		grpcAddress, err := p.getGRPCAddress()
		slog.Debug(MSG+"GetConfig", "grpcAddress", grpcAddress, "err", err)
		gRPCCredentials, err := p.getGRPCTransportCredentials()
		slog.Debug(MSG+"GetConfig", "grpcTransportCredentials", gRPCCredentials, "err", err)

		httpAddress, err := p.getHTTPAddress()
		slog.Debug(MSG+"GetConfig", "httpAddress", httpAddress, "err", err)
		tHTTPConfig, err := p.getHTTPTLSConfig()
		slog.Debug(MSG+"GetConfig", "tHTTPConfig", tHTTPConfig, "err", err)

		properties = getProperties(
			WithCacheExpire(cacheExpire),
			WithCacheGCInterval(cacheGCInterval),
			WithYamlConfig(yml),
			WithEnvironments(*env),
			WithFlags(flm),
			WithGRPCAddress(grpcAddress),
			WithGRPCTransportCredentials(gRPCCredentials),
			WithHTTPAddress(httpAddress),
			WithHTTPTLSConfig(tHTTPConfig),
			WithLogger(setupLogger(slogJSON(flm))),
		)
	})

	return properties
}

// WithCacheExpire — срок действия записи в кэше.
func WithCacheExpire(cacheExpire time.Duration) func(*mapProperties) {
	return func(p *mapProperties) {
		if cacheExpire > 0 {
			p.mp.Store(propertyCacheExpireMs, cacheExpire)
		}
	}
}

// CacheExpire геттер срока действия записи в кэше.
func (p *mapProperties) CacheExpire() time.Duration {
	if a, ok := p.mp.Load(propertyCacheExpireMs); ok {
		if cacheExpire, ok := a.(time.Duration); ok {
			return cacheExpire
		}
	}
	return 0
}

// WithCacheGCInterval — интервал очистки кэша.
func WithCacheGCInterval(cacheGCInterval time.Duration) func(*mapProperties) {
	return func(p *mapProperties) {
		if cacheGCInterval > 0 {
			p.mp.Store(propertyCacheGCIntervalSec, cacheGCInterval)
		}
	}
}

// CacheGCInterval геттер интервала очистки кэша.
func (p *mapProperties) CacheGCInterval() time.Duration {
	if a, ok := p.mp.Load(propertyCacheGCIntervalSec); ok {
		if cacheGCInterval, ok := a.(time.Duration); ok {
			return cacheGCInterval
		}
	}
	return 0
}

// WithYamlConfig — Конфигурация.
func WithYamlConfig(config YamlConfig) func(*mapProperties) {
	return func(p *mapProperties) {
		if config != nil {
			p.mp.Store(propertyConfig, config)
		}
	}
}

// YamlConfig — текущая yaml конфигурация.
func (p *mapProperties) YamlConfig() YamlConfig {
	if c, ok := p.mp.Load(propertyConfig); ok {
		if cfg, ok := c.(YamlConfig); ok {
			return cfg
		}
	}
	return nil
}

// WithEnvironments — Окружение.
func WithEnvironments(env environments) func(*mapProperties) {
	return func(p *mapProperties) {
		p.mp.Store(propertyEnvironments, env)
	}
}

// Environments геттер Окружения.
func (p *mapProperties) Environments() environments {
	if f, ok := p.mp.Load(propertyEnvironments); ok {
		if env, ok := f.(environments); ok {
			return env
		}
	}
	return environments{}
}

// WithFlags — Флаги.
func WithFlags(flags map[string]interface{}) func(*mapProperties) {
	return func(p *mapProperties) {
		if flags != nil {
			p.mp.Store(propertyFlags, flags)
		}
	}
}

// Flags — флаги командной строки.
func (p *mapProperties) Flags() map[string]interface{} {
	if f, ok := p.mp.Load(propertyFlags); ok {
		if flags, ok := f.(map[string]interface{}); ok {
			return flags
		}
	}
	return nil
}

// WithGRPCAddress — адрес gRPC сервера.
func WithGRPCAddress(address string) func(*mapProperties) {
	return func(p *mapProperties) {
		if address != "" {
			p.mp.Store(propertyGRPCAddress, address)
		}
	}
}

// GRPCAddress геттер адреса gRPC сервера.
func (p *mapProperties) GRPCAddress() string {
	if a, ok := p.mp.Load(propertyGRPCAddress); ok {
		if address, ok := a.(string); ok {
			return address
		}
	}
	return ""
}

// WithGRPCTransportCredentials — TLS реквизиты для gRPC-сервера.
func WithGRPCTransportCredentials(tCredentials credentials.TransportCredentials) func(*mapProperties) {
	return func(p *mapProperties) {
		if tCredentials != nil {
			p.mp.Store(propertyGRPCTransportCredentials, tCredentials)
		}
	}
}

// GRPCTransportCredentials геттер TLS реквизитов для gRPC-сервера.
func (p *mapProperties) GRPCTransportCredentials() credentials.TransportCredentials {
	if c, ok := p.mp.Load(propertyGRPCTransportCredentials); ok {
		if tCredentials, ok := c.(credentials.TransportCredentials); ok {
			return tCredentials
		}
	}
	return nil
}

// WithHTTPAddress — адрес HTTP сервера.
func WithHTTPAddress(address string) func(*mapProperties) {
	return func(p *mapProperties) {
		if address != "" {
			p.mp.Store(propertyHTTPAddress, address)
		}
	}
}

// HTTPAddress геттер адреса HTTP сервера.
func (p *mapProperties) HTTPAddress() string {
	if a, ok := p.mp.Load(propertyHTTPAddress); ok {
		if address, ok := a.(string); ok {
			return address
		}
	}
	return ""
}

// WithHTTPTLSConfig — TLS конфигурация для HTTP-сервера.
func WithHTTPTLSConfig(tCredentials *tls.Config) func(*mapProperties) {
	return func(p *mapProperties) {
		if tCredentials != nil {
			p.mp.Store(propertyHTTPHTTPTLSConfig, tCredentials)
		}
	}
}

// HTTPTLSConfig геттер TLS конфигурации для HTTP-сервера.
func (p *mapProperties) HTTPTLSConfig() *tls.Config {
	if c, ok := p.mp.Load(propertyHTTPHTTPTLSConfig); ok {
		if tCredentials, ok := c.(*tls.Config); ok {
			return tCredentials
		}
	}
	return nil
}

// WithLogger — логгер приложения.
func WithLogger(logger *slog.Logger) func(*mapProperties) {
	return func(p *mapProperties) {
		if logger != nil {
			p.mp.Store(propertyLogger, logger)
		}
	}
}

// Logger получение логгера приложения.
func (p *mapProperties) Logger() *slog.Logger {
	if a, ok := p.mp.Load(propertyLogger); ok {
		if logger, ok := a.(*slog.Logger); ok {
			return logger
		}
	}
	return slog.Default()
}

func (p *mapProperties) SlogJSON() bool {
	return slogJSON(p.Flags())
}

func (p *mapProperties) String() string {
	format := `
%s
CacheExpire: %v
CacheGCInterval: %v
Environments: %s
Flags: %v
GRPCAddress: %s
GRPCTransportCredentials: %v
HTTPAddress: %s
HTTPTransportCredentials: %v
`
	return fmt.Sprintf(format,
		p.YamlConfig(),
		p.CacheExpire(),
		p.CacheGCInterval(),
		p.Environments(),
		p.Flags(),
		p.GRPCAddress(),
		p.GRPCTransportCredentials(),
		p.HTTPAddress(),
		p.HTTPTLSConfig(),
	)
}

func slogJSON(flags map[string]interface{}) bool {
	if sj, ok := flags[flagSlogJson]; ok {
		if slogJSON, ok := sj.(*bool); ok {
			return *slogJSON
		}
	}
	return false
}

func getProperties(opts ...func(*mapProperties)) *mapProperties {

	var property = new(mapProperties)

	// вызываем все указанные функции для установки параметров
	for _, opt := range opts {
		opt(property)
	}

	return property
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
