/*
 * This file was last modified at 2024-07-11 09:38 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * load_tls_test.go
 * $Id$
 */
//!+

package tool

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/credentials"
)

func TestLoadClientTLSCredentials(t *testing.T) {
	type input struct {
		caCertFile string
	}
	var tests = []struct {
		name  string
		input input
		fRun  func(*testing.T, string)
	}{
		{
			name: "negative test #1 LoadServerTLSCredentials",
			input: input{
				caCertFile: "",
			},
			fRun: negativeLoadClientTLSCredentials,
		},
		{
			name: "negative test #2 LoadServerTLSCredentials",
			input: input{
				caCertFile: "test_server-key.pem",
			},
			fRun: negativeLoadClientTLSCredentials,
		},
		{
			name: "positive test #3 LoadServerTLSCredentials",
			input: input{
				caCertFile: "test_ca-cert.pem",
			},
			fRun: positiveLoadClientTLSCredentials,
		},
	}

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t, test.input.caCertFile)
		})
	}
}

func negativeLoadClientTLSCredentials(t *testing.T, caCertFile string) {
	_, err := LoadClientTLSCredentials(caCertFile)
	assert.NotNil(t, err)
}

func positiveLoadClientTLSCredentials(t *testing.T, caCertFile string) {
	expectedInfo := credentials.ProtocolInfo(
		credentials.ProtocolInfo{
			ProtocolVersion:  "",
			SecurityProtocol: "tls",
			SecurityVersion:  "1.2",
			ServerName:       "",
		})
	got, err := LoadClientTLSCredentials(caCertFile)
	assert.Nil(t, err)
	assert.NotNil(t, expectedInfo)
	assert.Equal(t, expectedInfo, got.Info())
}

func TestLoadServerTLSCredentials(t *testing.T) {
	type input struct {
		certFile string
		keyFile  string
	}
	var tests = []struct {
		name  string
		input input
		fRun  func(*testing.T, string, string)
	}{
		{
			name: "positive test #0 LoadServerTLSCredentials",
			input: input{
				certFile: "test_server-cert.pem",
				keyFile:  "test_server-key.pem",
			},
			fRun: positiveLoadServerTLSCredentials,
		},
		{
			name: "negative test #1 LoadServerTLSCredentials",
			input: input{
				certFile: "",
				keyFile:  "",
			},
			fRun: negativeLoadServerTLSCredentials,
		},
	}

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t, test.input.certFile, test.input.keyFile)
		})
	}
}

func negativeLoadServerTLSCredentials(t *testing.T, certFile, keyFile string) {
	_, err := LoadServerTLSCredentials(certFile, keyFile)
	assert.NotNil(t, err)
}

func positiveLoadServerTLSCredentials(t *testing.T, certFile, keyFile string) {
	expectedInfo := credentials.ProtocolInfo(
		credentials.ProtocolInfo{
			ProtocolVersion:  "",
			SecurityProtocol: "tls",
			SecurityVersion:  "1.2",
			ServerName:       "",
		})

	got, err := LoadServerTLSCredentials(certFile, keyFile)
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, expectedInfo, got.Info())
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
