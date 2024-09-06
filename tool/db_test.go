/*
 * This file was last modified at 2024-07-11 09:38 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * db_test.go
 * $Id$
 */
//!+

package tool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBConnect(t *testing.T) {
	var tests = []struct {
		name string
	}{{"negative test #0 DBConnect"}}

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				p := recover()
				assert.NotNil(t, p)
			}()
			DBConnect("")
		})
	}
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
