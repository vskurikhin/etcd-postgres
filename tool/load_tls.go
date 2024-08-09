/*
 * This file was last modified at 2024-07-21 10:37 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * load_tls.go
 * $Id$
 */
//!+

// Package tool TODO.
package tool

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"
)

func LoadClientTLSCredentials(caCertFile string) (credentials.TransportCredentials, error) {
	// Загрузка сертификата центра сертификации, подписавшего сертификат сервера.
	pemServerCA, err := os.ReadFile(caCertFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Создание учётных данных для конфигурации TLS.
	config := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(config), nil
}

func LoadServerTLSCredentials(certFile, keyFile string) (credentials.TransportCredentials, error) {

	// Загрузка серверного сертификата и закрытого ключа.
	serverCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	// Создание учётных данных для конфигурации TLS.
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}
	return credentials.NewTLS(config), nil
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
