/*
 * This file was last modified at 2024-07-31 16:17 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * trace_test.go
 * $Id$
 */
//!+

// Package tool TODO.
package tool

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTraceInOut(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{
			name: "test #1 positive ",
			fRun: testTraceInOutLevelInfo,
		},
		{
			name: "test #2 positive ",
			fRun: testTraceInOutLevelDebug,
		},
	}

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func testTraceInOutLevelInfo(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelInfo)
	defer TraceInOut(context.TODO(), "name", "%v, %d, %s, %T", true, 13, "test", func() {})
}

func testTraceInOutLevelDebug(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	SetLogger(slog.Default())
	defer TraceInOut(context.TODO(), "name", "%v, %d, %s, %T", true, 13, "test", func() {})
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
