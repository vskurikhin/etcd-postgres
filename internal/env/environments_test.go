/*
 * Copyright text:
 * This file was last modified at 2024-07-10 20:19 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * environments_test.go
 * $Id$
 */
//!+

// Package env работа с настройками и окружением.
package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetEnvironments(t *testing.T) {
	type want struct {
		env *environments
		err error
	}
	var tests = []struct {
		name string
		fRun func() (env *environments, err error)
		want want
	}{
		{
			name: `positive test #0 nil environments`,
			fRun: getEnvironments,
			want: want{&environments{}, nil},
		},
	}
	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.fRun()
			assert.Equal(t, test.want.env, got)
			assert.Equal(t, test.want.err, err)
		})
	}
}

func Test_environmentsString(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(env *environments) string
		want string
	}{
		{
			name: `positive test #0 nil environments`,
			fRun: func(e *environments) string { return e.String() },
			want: `CACHE_EXPIRE_MS: 0
CACHE_GC_INTERVAL_SEC: 0
GRPC_ADDRESS: []
GRPC_CA_FILE: 
GRPC_CERT_FILE: 
GRPC_KEY_FILE: 
HTTP_ADDRESS: []
HTTP_CA_FILE: 
HTTP_CERT_FILE: 
HTTP_KEY_FILE: `,
		},
	}
	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			envi, err := getEnvironments()
			assert.Nil(t, err)
			got := test.fRun(envi)
			assert.Equal(t, test.want, got)
		})
	}
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
