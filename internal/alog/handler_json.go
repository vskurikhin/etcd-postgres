/*
 * This file was last modified at 2024-08-06 17:18 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * handler_json.go
 * $Id$
 */
//!+

// Package alog кастомизация slog логгера.
package alog

import (
	"context"
	"io"
	"log"
	"log/slog"

	"github.com/google/uuid"
)

// HandlerJSON промышленный JSON логгер.
type HandlerJSON struct {
	slog.Handler
	l *log.Logger
}

// Handle обработчик промышленного JSON логгера.
func (h *HandlerJSON) Handle(ctx context.Context, r slog.Record) error {

	id := uuid.Max.String()

	if ri := ctx.Value("request-id"); ri != nil {
		if requestId, ok := ri.(string); ok {
			id = requestId
		}
	}
	r.Add("requestId", slog.StringValue(id))

	return h.Handler.Handle(ctx, r)
}

func NewHandlerJSON(out io.Writer, opts *slog.HandlerOptions) *HandlerJSON {

	h := &HandlerJSON{
		Handler: slog.NewJSONHandler(out, opts),
		l:       log.New(out, "", 0),
	}
	return h
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
