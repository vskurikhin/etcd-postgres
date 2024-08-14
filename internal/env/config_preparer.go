/*
 * This file was last modified at 2024-08-06 20:17 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * config_preparer.go
 * $Id$
 */
//!+

// Package env работа с настройками и окружением.
package env

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/victor-skurikhin/etcd-client/v1/internal/alog"
	"github.com/victor-skurikhin/etcd-client/v1/tool"
	"google.golang.org/grpc/credentials"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

type preparer struct {
	env     *environments
	flagMap map[string]interface{}
	yml     YamlConfig
}

var ErrEmptyAddress = fmt.Errorf("can't configure epmty address")

func (p *preparer) getCacheExpire() (time.Duration, error) {
	return toTimePrepareProperty(
		flagCacheExpireMs,
		p.flagMap[flagCacheExpireMs],
		p.env.CacheExpireMs,
		p.yml.CacheExpireMs(),
		time.Millisecond,
	)
}

func (p *preparer) getCacheGCInterval() (time.Duration, error) {
	return toTimePrepareProperty(
		flagCacheGCIntervalSec,
		p.flagMap[flagCacheGCIntervalSec],
		p.env.CacheGCIntervalSec,
		p.yml.CacheGCIntervalSec(),
		time.Second,
	)
}

func (p *preparer) getEtcdAddresses() ([]string, error) {
	if p.yml.EtcdEnabled() {
		return serverAddressesPrepareProperty(
			flagEtcdAddresses, p.flagMap,
			p.env.EtcdAddresses,
			p.yml.EtcdAddresses())
	}
	return []string{}, fmt.Errorf("etcd servers disabled")
}

func (p *preparer) getEtcdDialTimeout() (time.Duration, error) {
	if p.yml.EtcdEnabled() {
		result, err := timePrepareProperty(
			flagEtcdDialTimeout, p.flagMap,
			p.env.EtcdDialTimeout,
			p.yml.EtcdDialTimeout())
		if result == 0 {
			return time.Second, err
		}
		return result, err
	}
	return 0, fmt.Errorf("etcd servers disabled")
}

func (p *preparer) getGRPCAddress() (string, error) {
	if p.yml.GRPCEnabled() {
		return serverAddressPrepareProperty(
			flagGRPCAddress, p.flagMap,
			p.env.GRPCAddress,
			p.yml.GRPCAddress(),
			p.yml.GRPCPort())
	}
	return "", fmt.Errorf("gRPC server disabled")
}

func (p *preparer) getGRPCTransportCredentials() (credentials.TransportCredentials, error) {
	if p.yml.GRPCEnabled() {
		return serverTransportCredentialsPrepareProperty(
			flagGRPCCertFile,
			flagGRPCKeyFile, p.flagMap,
			p.env.GRPCCertFile,
			p.env.GRPCKeyFile,
			p.yml.GRPCTLSCertFile(),
			p.yml.GRPCTLSKeyFile(),
		)
	}
	return nil, fmt.Errorf("gRPC server disabled")
}

func (p *preparer) getHTTPAddress() (string, error) {
	if p.yml.GRPCEnabled() {
		return serverAddressPrepareProperty(
			flagHTTPAddress, p.flagMap,
			p.env.HTTPAddress,
			p.yml.HTTPAddress(),
			p.yml.HTTPPort(),
		)
	}
	return "", fmt.Errorf("HTTP server disabled")
}

func (p *preparer) getHTTPTLSConfig() (*tls.Config, error) {
	if p.yml.HTTPTLSEnabled() {
		return serverTLSConfigPrepareProperty(
			flagHTTPCertFile,
			flagHTTPKeyFile, p.flagMap,
			p.env.HTTPCertFile,
			p.env.HTTPKeyFile,
			p.yml.HTTPTLSCertFile(),
			p.yml.HTTPTLSKeyFile(),
		)
	}
	return nil, fmt.Errorf("HTTP server disabled")
}

func parseEnvAddress(address []string) string {

	port, err := strconv.Atoi(address[len(address)-1])
	tool.IfErrorThenPanic(err)
	var bb bytes.Buffer

	if len(address) > 1 {
		for i := 0; i < len(address)-1; i++ {
			bb.WriteString(address[i])
			bb.WriteRune(':')
		}
	} else {
		bb.WriteRune(':')
	}
	return fmt.Sprintf("%s%d", bb.String(), port)
}

func serverAddressesPrepareProperty(
	name string,
	flm map[string]interface{},
	envAddresses []string,
	ymlAddresses []string,
) ([]string, error) {

	addresses := make([]string, 0)
	var err error

	getFlagAddress := func() {
		if ps, ok := flm[name].(*string); !ok {
			err = fmt.Errorf("bad value of %s : %v", name, flm[name])
		} else if ps != nil {
			addresses = make([]string, 0)
			for _, a := range strings.Split(*ps, ",") {
				address := strings.Split(a, ":")
				addresses = append(addresses, parseEnvAddress(address))
			}
		}
	}
	if len(envAddresses) > 0 {
		for _, a := range envAddresses {
			address := strings.Split(a, ":")
			addresses = append(addresses, parseEnvAddress(address))
		}
	} else if len(ymlAddresses) < 1 {
		getFlagAddress()
	} else {
		for _, a := range ymlAddresses {
			address := strings.Split(a, ":")
			addresses = append(addresses, parseEnvAddress(address))
		}
	}
	setIfFlagChanged(name, getFlagAddress)

	if len(addresses) < 1 {
		err = ErrEmptyAddress
	}
	return addresses, err
}

func serverAddressPrepareProperty(
	name string,
	flm map[string]interface{},
	envAddress []string,
	ymlAddress string,
	ymlPort int,
) (string, error) {

	var address string
	var err error

	getFlagAddress := func() {
		if a, ok := flm[name].(*string); !ok {
			err = fmt.Errorf("bad value of %s : %v", name, flm[name])
		} else {
			address = *a
		}
	}
	address = fmt.Sprintf("%s:%d", ymlAddress, ymlPort)

	if len(envAddress) > 0 {
		address = parseEnvAddress(envAddress)
	} else if ymlAddress == "" && ymlPort == 0 {
		getFlagAddress()
	}
	setIfFlagChanged(name, getFlagAddress)

	if address == "" {
		err = ErrEmptyAddress
	}
	return address, err
}

func serverTransportCredentialsPrepareProperty(
	nameCertFile string,
	nameKeyFile string,
	flm map[string]interface{},
	envTLSCertFile string,
	envTLSKeyFile string,
	ymlTLSCertFile string,
	ymlTLSKeyFile string,
) (tCredentials credentials.TransportCredentials, err error) {

	certFile, keyFile := ymlTLSCertFile, ymlTLSKeyFile
	getFlagGRPCCertFile := func() {
		if cf, ok := flm[nameCertFile].(*string); !ok {
			err = fmt.Errorf("bad value of %s : %v", flagGRPCCertFile, flm[flagGRPCCertFile])
		} else {
			certFile = *cf
		}
	}
	getFlagGRPCKeyFile := func() {
		if kf, ok := flm[nameKeyFile].(*string); !ok {
			err = fmt.Errorf("bad value of %s : %v", flagGRPCKeyFile, flm[flagGRPCKeyFile])
		} else {
			keyFile = *kf
		}
	}
	if envTLSCertFile != "" {
		certFile = envTLSCertFile
	}
	if envTLSKeyFile != "" {
		keyFile = envTLSKeyFile
	}
	if certFile == "" {
		getFlagGRPCCertFile()
	}
	if keyFile == "" {
		getFlagGRPCKeyFile()
	}
	setIfFlagChanged(nameCertFile, getFlagGRPCCertFile)
	setIfFlagChanged(nameKeyFile, getFlagGRPCKeyFile)
	if err != nil {
		return nil, err
	}
	return tool.LoadServerTLSCredentials(certFile, keyFile)
}

func serverTLSConfigPrepareProperty(
	nameCertFile string,
	nameKeyFile string,
	flm map[string]interface{},
	envTLSCertFile string,
	envTLSKeyFile string,
	ymlTLSCertFile string,
	ymlTLSKeyFile string,
) (tConfig *tls.Config, err error) {
	certFile, keyFile := ymlTLSCertFile, ymlTLSKeyFile
	getFlagCertFile := func() {
		if cf, ok := flm[nameCertFile].(*string); !ok {
			err = fmt.Errorf("bad value of %s : %v", flagGRPCCertFile, flm[flagGRPCCertFile])
		} else {
			certFile = *cf
		}
	}
	getFlagKeyFile := func() {
		if kf, ok := flm[nameKeyFile].(*string); !ok {
			err = fmt.Errorf("bad value of %s : %v", flagGRPCKeyFile, flm[flagGRPCKeyFile])
		} else {
			keyFile = *kf
		}
	}
	if envTLSCertFile != "" {
		certFile = envTLSCertFile
	}
	if envTLSKeyFile != "" {
		keyFile = envTLSKeyFile
	}
	if certFile == "" {
		getFlagCertFile()
	}
	if keyFile == "" {
		getFlagKeyFile()
	}
	setIfFlagChanged(nameCertFile, getFlagCertFile)
	setIfFlagChanged(nameKeyFile, getFlagKeyFile)
	if err != nil {
		return nil, err
	}
	cer, err := tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		return nil, err
	}
	return &tls.Config{Certificates: []tls.Certificate{cer}}, nil
}

func setupLogger(debug bool, slogJSON bool) *slog.Logger { // *flm[propertyDebug].(*bool)
	var level slog.Level
	if debug {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}
	slog.SetLogLoggerLevel(level)
	if slogJSON {
		alog.NewLogger(alog.NewHandlerJSON(os.Stdout, &slog.HandlerOptions{Level: level}))
	} else {
		opts := alog.PrettyHandlerOptions{
			SlogOpts: slog.HandlerOptions{
				Level: level,
			},
		}
		alog.NewLogger(alog.NewPrettyHandlerText(os.Stdout, opts))
	}
	return tool.SetLogger(alog.GetLogger())
}

func timePrepareProperty(name string, flag interface{}, env time.Duration, yaml time.Duration) (time.Duration, error) {

	var result time.Duration
	var err error

	getFlag := func() {
		if a, ok := flag.(*time.Duration); !ok {
			err = fmt.Errorf("bad value")
		} else {
			result = *a
		}
	}
	if yaml > 0 {
		result = yaml
	}
	if env > 0 {
		result = env
	} else if result == 0 {
		getFlag()
	}
	setIfFlagChanged(name, getFlag)

	return result, err
}

func toTimePrepareProperty(name string, flag interface{}, env int, yaml int, scale time.Duration) (time.Duration, error) {

	var result time.Duration
	var err error

	getFlag := func() {
		if a, ok := flag.(*int); !ok {
			err = fmt.Errorf("bad value")
		} else {
			result = time.Duration(*a) * scale
		}
	}
	if yaml > 0 {
		result = time.Duration(yaml) * scale
	}
	if env > 0 {
		result = time.Duration(env) * scale
	} else if result == 0 {
		getFlag()
	}
	setIfFlagChanged(name, getFlag)

	return result, err
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
