/*
 * This file was last modified at 2024-07-16 10:17 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * randoms_test.go
 * $Id$
 */
//!+

// Package tool TODO.
package tool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandStringBytes(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{name: "positive test #0 RandStringBytes", fRun: testRandStringBytes},
	}

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func testRandStringBytes(t *testing.T) {
	assert.Equal(t, 32, len(RandStringBytes(32)))
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
