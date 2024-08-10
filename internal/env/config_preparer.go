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

func setupLogger(slogJSON bool) *slog.Logger {
	if slogJSON {
		alog.NewLogger(alog.NewHandlerJSON(os.Stdout, nil))
	} else {
		opts := alog.PrettyHandlerOptions{
			SlogOpts: slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}
		alog.NewLogger(alog.NewPrettyHandlerText(os.Stdout, opts))
	}
	return tool.SetLogger(alog.GetLogger())
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
