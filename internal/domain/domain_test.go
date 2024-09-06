/*
 * This file was last modified at 2024-08-16 12:12 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * domain_test.go
 * $Id$
 */
//!+

// Package domain TODO.
package domain

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCloner(t *testing.T) {
	for _, test := range []struct {
		name string
		fRun func(*testing.T) (interface{}, error)
		want func(*testing.T, interface{}) bool
	}{
		{
			"test #0 positive for interface Cloner[T Ptr[U], U any] method Clone(U)",
			positiveClonerClone,
			positiveClonerCloneCheck,
		},
		{
			"test #1 positive for interface Factory[T Ptr[U], U any] method Create()",
			positiveClonerCopy,
			positiveClonerCopyCheck,
		},
		{
			"test #2 negative for interface Factory[T Ptr[U], U any]",
			negativeClonerCopy,
			negativeClonerCopyCheck,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.fRun(t)
			assert.Nil(t, err)
			assert.NotNil(t, got)
			assert.True(t, test.want(t, got))
		})
	}
}

type unit struct{}

type unitCloner struct{}

func (k *unitCloner) Clone(_ unit) unit {
	return unit{}
}

func (k *unitCloner) Copy(t *unit) *unit {

	if t == nil {
		return nil
	}
	return &unit{}
}

func positiveClonerClone(_ *testing.T) (interface{}, error) {
	return (&unitCloner{}).Clone(unit{}), nil
}

func positiveClonerCloneCheck(_ *testing.T, i interface{}) bool {
	if u, ok := i.(unit); ok {
		return u == unit{}
	}
	return false
}

func positiveClonerCopy(_ *testing.T) (interface{}, error) {
	return (&unitCloner{}).Copy(&unit{}), nil
}

func positiveClonerCopyCheck(_ *testing.T, i interface{}) bool {
	if u, ok := i.(*unit); ok {
		return *u == unit{}
	}
	return false
}

func negativeClonerCopy(_ *testing.T) (interface{}, error) {
	if n := (&unitCloner{}).Copy(nil); n == nil {
		return true, nil
	}
	return nil, fmt.Errorf("then Copy(nil) when isn't nil")
}

func negativeClonerCopyCheck(_ *testing.T, _ interface{}) bool {
	return true
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
