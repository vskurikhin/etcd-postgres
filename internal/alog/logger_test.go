/*
 * This file was last modified at 2024-07-31 14:33 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * logger_test.go
 * $Id$
 */
//!+

// Package alog кастомизация slog логгера.
package alog

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLoggerGetLogger(t *testing.T) {
	NewLogger(NewHandlerJSON(os.Stdout, nil))
	assert.Equal(t, logger, GetLogger())
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
