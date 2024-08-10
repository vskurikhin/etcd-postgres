/*
 * This file was last modified at 2024-07-31 13:56 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * handler_json.go
 * $Id$
 */
//!+

// Package alog кастомизация slog логгера.
package alog

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

var testTime = time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

func TestJSONHandlerContextWithValue(t *testing.T) {
	for _, test := range []struct {
		name string
		opts slog.HandlerOptions
		want string
	}{
		{
			"none",
			slog.HandlerOptions{},
			`{"time":"2000-01-02T03:04:05Z","level":"INFO","msg":"m","a":1,"m":{"b":2},"requestId":"00010203-0405-0607-0809-0a0b0c0d0e0f"}`,
		},
		{
			"replace",
			slog.HandlerOptions{ReplaceAttr: upperCaseKey},
			`{"TIME":"2000-01-02T03:04:05Z","LEVEL":"INFO","MSG":"m","A":1,"M":{"b":2},"REQUESTID":"00010203-0405-0607-0809-0a0b0c0d0e0f"}`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			h := NewHandlerJSON(&buf, &test.opts)
			r := slog.NewRecord(testTime, slog.LevelInfo, "m", 0)
			r.AddAttrs(slog.Int("a", 1), slog.Any("m", map[string]int{"b": 2}))
			id := uuid.UUID{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
			ctx := context.WithValue(context.Background(), "request-id", id.String())
			if err := h.Handle(ctx, r); err != nil {
				t.Fatal(err)
			}
			got := strings.TrimSuffix(buf.String(), "\n")
			if got != test.want {
				t.Errorf("\ngot  %s\nwant %s", got, test.want)
			}
		})
	}
}

func TestJSONHandlerContextBackground(t *testing.T) {
	for _, test := range []struct {
		name string
		opts slog.HandlerOptions
		want string
	}{
		{
			"none",
			slog.HandlerOptions{},
			`{"time":"2000-01-02T03:04:05Z","level":"INFO","msg":"m","a":1,"m":{"b":2},"requestId":"ffffffff-ffff-ffff-ffff-ffffffffffff"}`,
		},
		{
			"replace",
			slog.HandlerOptions{ReplaceAttr: upperCaseKey},
			`{"TIME":"2000-01-02T03:04:05Z","LEVEL":"INFO","MSG":"m","A":1,"M":{"b":2},"REQUESTID":"ffffffff-ffff-ffff-ffff-ffffffffffff"}`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			h := NewHandlerJSON(&buf, &test.opts)
			r := slog.NewRecord(testTime, slog.LevelInfo, "m", 0)
			r.AddAttrs(slog.Int("a", 1), slog.Any("m", map[string]int{"b": 2}))
			if err := h.Handle(context.Background(), r); err != nil {
				t.Fatal(err)
			}
			got := strings.TrimSuffix(buf.String(), "\n")
			if got != test.want {
				t.Errorf("\ngot  %s\nwant %s", got, test.want)
			}
		})
	}
}

func upperCaseKey(_ []string, a slog.Attr) slog.Attr {
	a.Key = strings.ToUpper(a.Key)
	return a
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
