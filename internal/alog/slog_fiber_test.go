/*
 * This file was last modified at 2024-07-31 13:56 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * slog_fiber_test.go
 * $Id$
 */
//!+

// Package alog кастомизация slog логгера.
package alog

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	slogf "github.com/samber/slog-fiber"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/bytebufferpool"
)

func TestNewNewWithConfig(t *testing.T) {
	slogLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	handler0 := New(slogLogger)
	handler1 := NewWithConfig(slogLogger, Config{
		DefaultLevel:     slog.LevelInfo,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,

		WithUserAgent:      false,
		WithRequestID:      true,
		WithRequestBody:    false,
		WithRequestHeader:  false,
		WithResponseBody:   false,
		WithResponseHeader: false,
		WithSpanID:         false,
		WithTraceID:        false,

		Filters: []slogf.Filter{},
	})
	NewLogger(NewHandlerJSON(os.Stdout, nil))
	assert.NotNil(t, handler0)
	assert.NotNil(t, handler1)
}

func TestLogHandler(t *testing.T) {
	for _, test := range []struct {
		name string
		fRun func(*testing.T) (string, error)
		want string
	}{
		{
			"none",
			testMethodGet,
			`","level":"ERROR","msg":"some random error","request":{"time":"`,
		},
		{
			"replace",
			testMethodPost,
			`","level":"ERROR","msg":"some random error","request":{"time":"`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.fRun(t)
			assert.Nil(t, err)
			assert.True(t, strings.Contains(got, test.want))
		})
	}
}

func testMethodGet(t *testing.T) (string, error) {

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	slogLogger := slog.New(slog.NewJSONHandler(buf, nil))
	handler := New(slogLogger)
	app := fiber.New()
	app.Use(handler)
	app.Get("/", func(c *fiber.Ctx) error {
		return errors.New("some random error")
	})
	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
	utils.AssertEqual(t, fiber.StatusInternalServerError, resp.StatusCode)

	return buf.String(), err
}

func testMethodPost(t *testing.T) (string, error) {

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	slogLogger := slog.New(slog.NewJSONHandler(buf, nil))
	handler := NewWithConfig(slogLogger, Config{
		DefaultLevel:     slog.LevelInfo,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,

		WithUserAgent:      true,
		WithRequestID:      true,
		WithRequestBody:    true,
		WithRequestHeader:  true,
		WithResponseBody:   true,
		WithResponseHeader: true,
		WithSpanID:         true,
		WithTraceID:        true,

		Filters: []slogf.Filter{},
	})
	app := fiber.New()
	app.Use(handler)
	app.Post("/", func(c *fiber.Ctx) error {
		return errors.New("some random error")
	})
	req := httptest.NewRequest(fiber.MethodPost, "/", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, err := app.Test(req)
	utils.AssertEqual(t, fiber.StatusInternalServerError, resp.StatusCode)

	return buf.String(), err
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
