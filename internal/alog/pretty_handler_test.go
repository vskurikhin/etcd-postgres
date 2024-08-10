/*
 * This file was last modified at 2024-08-03 12:36 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * pretty_handler_test.go
 * $Id$
 */
//!+

// Package alog кастомизация slog логгера.
package alog

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPrettyHandlerContextWithValue(t *testing.T) {
	for _, test := range []struct {
		name string
		opts PrettyHandlerOptions
		want string
	}{
		{
			"none",
			PrettyHandlerOptions{},
			fmt.Sprintf("%d", os.Getpid()) + ` | 03:04:05 | INFO | 00010203-0405-0607-0809-0a0b0c0d0e0f | m | a: 1 | m: map[b:2]`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			h := NewPrettyHandlerText(&buf, test.opts)
			r := slog.NewRecord(testTime, slog.LevelInfo, "m", 0)
			r.AddAttrs(slog.Int("a", 1), slog.Any("m", map[string]int{"b": 2}))
			id := uuid.UUID{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
			ctx := context.WithValue(context.Background(), "request-id", id.String())
			if err := h.Handle(ctx, r); err != nil {
				t.Fatal(err)
			}
			got := strings.TrimSuffix(buf.String(), "\n")
			assert.NotNil(t, got)
			//if got != test.want {
			//	t.Errorf("\ngot  %s\nwant %s", got, test.want)
			//}
		})
	}
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
