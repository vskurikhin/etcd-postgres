/*
 * Copyright text:
 * This file was last modified at 2024-07-10 21:51 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * config_test.go
 * $Id$
 */
//!+

// Package env работа с настройками и окружением.
package env

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFirst(t *testing.T) {
	assert.NotNil(t, t)
	var tests = []testLoadConfig{
		{
			name:  `test #0 negative LoadConfig("test")`,
			input: "test",
			fRun: func(s string) (YamlConfig, error) {
				wd := chDir()
				defer func() { _ = os.Chdir(wd) }()
				return LoadConfig(s)
			},
			want: wantLoadConfig{
				yamlConfig: nil,
				err:        configFileNotFoundError("test"),
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

func TestGetProperties(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{
			name: "positive test #0 GetConfig",
			fRun: tryDefaultGetConfig,
		},
	}
	oldCommandLine := pflag.CommandLine
	defer func() { pflag.CommandLine = oldCommandLine }()
	ResetForTesting(func() {})

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func tryDefaultGetConfig(t *testing.T) {
	once = new(sync.Once)
	got := GetConfig()
	expected := fmt.Sprintf("%p", got)
	assert.NotNil(t, got)
	assert.NotNil(t, got.String())
	assert.NotEqual(t, "", got.String())
	assert.Equal(t, expected, fmt.Sprintf("%p", GetConfig()))
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
