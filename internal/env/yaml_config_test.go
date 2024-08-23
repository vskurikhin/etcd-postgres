/*
 * Copyright text:
 * This file was last modified at 2024-07-10 20:02 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * yaml_config_test.go
 * $Id$
 */
//!+

// Package env работа с настройками и окружением.
package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYamlConfig(t *testing.T) {
	var tests = []struct {
		name  string
		fRun  func() YamlConfig
		isNil bool
		want  string
	}{
		{
			name:  `positive test #0 nil yamlConfig`,
			fRun:  nilYamlConfig,
			isNil: true,
			want: `CacheEnabled: false
CacheExpire: 0
CacheGCInterval: 0
DBEnabled: false
DBHost: 
DBName: 
DBPort: 0
DBRetryIncrease: 0
DBRetryTries: 0
DBUserName: 
DBUserPassword: 
GRPCAddress: 
GRPCEnabled: false
GRPCPort: 0
GRPCProto: 
GRPCTLSCAFile: 
GRPCTLSCertFile: 
GRPCTLSKeyFile: 
GRPCTLSEnabled: false
HTTPAddress: 
HTTPEnabled: false
HTTPPort: 0
HTTPTLSCAFile: 
HTTPTLSCertFile: 
HTTPTLSEnabled: false
HTTPTLSKeyFile: `,
		},
		{
			name: `positive test #1 zero yamlConfig`,
			fRun: zeroYamlConfig,
			want: `CacheEnabled: false
CacheExpire: 0
CacheGCInterval: 0
DBEnabled: false
DBHost: 
DBName: 
DBPort: 0
DBRetryIncrease: 0
DBRetryTries: 0
DBUserName: 
DBUserPassword: 
GRPCAddress: 
GRPCEnabled: false
GRPCPort: 0
GRPCProto: 
GRPCTLSCAFile: 
GRPCTLSCertFile: 
GRPCTLSKeyFile: 
GRPCTLSEnabled: false
HTTPAddress: 
HTTPEnabled: false
HTTPPort: 0
HTTPTLSCAFile: 
HTTPTLSCertFile: 
HTTPTLSEnabled: false
HTTPTLSKeyFile: `,
		},
	}
	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.fRun()
			if !test.isNil {
				assert.Equal(t, test.want, got.String())
			} else {
				assert.Equal(t, test.want, (*yamlConfig)(nil).String())
			}
		})
	}
}

func nilYamlConfig() YamlConfig {
	return nil
}

func zeroYamlConfig() YamlConfig {
	return &yamlConfig{}
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
