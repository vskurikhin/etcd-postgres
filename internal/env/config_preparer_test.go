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
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPreparerPositive(t *testing.T) {
	for _, test := range []struct {
		name string
		fRun func(*testing.T) (interface{}, error)
	}{
		{
			"test #0 positive #1 for serverAddressPrepareProperty",
			serverAddressPreparePropertyPositiveTest1,
		},
		{
			"test #1 positive #2 for serverAddressPrepareProperty",
			serverAddressPreparePropertyPositiveTest2,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.fRun(t)
			assert.Nil(t, err)
			assert.NotNil(t, got)
		})
	}
}

func serverAddressPreparePropertyPositiveTest2(t *testing.T) (interface{}, error) {
	got, err := serverAddressPrepareProperty("", make(map[string]interface{}), []string{"l", "1"}, "", 0)
	assert.Equal(t, "l:1", got)
	return got, err
}

func serverAddressPreparePropertyPositiveTest1(t *testing.T) (interface{}, error) {
	s := "l:1"
	got, err := serverAddressPrepareProperty("", map[string]interface{}{"": &s}, []string{}, "", 0)
	assert.Equal(t, "l:1", got)
	return got, err
}

func TestPreparerNegative(t *testing.T) {
	for _, test := range []struct {
		name string
		fRun func(*testing.T) (interface{}, error)
	}{
		{
			"test #0 negative #1 for serverAddressPrepareProperty",
			serverAddressPreparePropertyNegativeTest1,
		},
		{
			"test #0 negative #2 for serverAddressPrepareProperty",
			serverAddressPreparePropertyNegativeTest2,
		},
		{
			"test #1 negative for serverTransportCredentialsPrepareProperty",
			serverTransportCredentialsPreparePropertyNegativeTest,
		},
		{
			"test #2 negative #1 for serverTLSConfigPrepareProperty",
			serverTLSConfigPreparePropertyNegativeTest1,
		},
		{
			"test #3 negative #2 for serverTLSConfigPrepareProperty",
			serverTLSConfigPreparePropertyNegativeTest2,
		},
		{
			"test #4 negative for toTimePrepareProperty",
			toTimePreparePropertyNegativeTest,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.fRun(t)
			assert.NotNil(t, err)
			assert.Nil(t, got)
		})
	}
}

func serverAddressPreparePropertyNegativeTest1(t *testing.T) (interface{}, error) {
	got, err := serverAddressPrepareProperty("", make(map[string]interface{}), []string{}, "", 0)
	assert.Equal(t, ":0", got)
	return nil, err
}

func serverAddressPreparePropertyNegativeTest2(t *testing.T) (interface{}, error) {
	s := ""
	got, err := serverAddressPrepareProperty("", map[string]interface{}{"": &s}, []string{}, "", 0)
	assert.Equal(t, "", got)
	return nil, err
}

func serverTransportCredentialsPreparePropertyNegativeTest(t *testing.T) (interface{}, error) {
	got, err := serverTransportCredentialsPrepareProperty("", "", make(map[string]interface{}), "", "", "", "")
	assert.Nil(t, got)
	return nil, err
}

func serverTLSConfigPreparePropertyNegativeTest1(t *testing.T) (interface{}, error) {
	got, err := serverTLSConfigPrepareProperty("", "", make(map[string]interface{}), "", "", "", "")
	assert.Nil(t, got)
	return nil, err
}

func serverTLSConfigPreparePropertyNegativeTest2(t *testing.T) (interface{}, error) {
	got, err := serverTLSConfigPrepareProperty("test", "test", map[string]interface{}{"test": ""}, "", "", "", "")
	assert.Nil(t, got)
	return nil, err
}

func toTimePreparePropertyNegativeTest(t *testing.T) (interface{}, error) {
	f := false
	got, err := toTimePrepareProperty("", &f, 0, 0, 0)
	assert.Equal(t, time.Duration(0), got)
	return nil, err
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
