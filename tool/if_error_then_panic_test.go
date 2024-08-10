/*
 * This file was last modified at 2024-07-11 09:38 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * if_error_then_panic_test.go
 * $Id$
 */
//!+

package tool

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIfErrorThenPanic(t *testing.T) {
	var tests = []struct {
		name string
	}{{"positive test #0 LoadServerTLSCredentials"}}

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				p := recover()
				assert.NotNil(t, p)
			}()
			IfErrorThenPanic(fmt.Errorf("test"))
		})
	}
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
