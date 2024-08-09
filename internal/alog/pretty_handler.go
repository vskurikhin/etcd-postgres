/*
 * This file was last modified at 2024-08-06 17:24 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * pretty_handler.go
 * $Id$
 */
//!+

// Package alog кастомизация slog логгера.
package alog

import (
	"bytes"
	"context"
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/google/uuid"

	"github.com/fatih/color"
)

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

// PrettyHandler «симпатичный» логгер.
type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

// Handle обработчик «симпатичный» логгера.
func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {

	var level string

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(r.Level.String())
	case slog.LevelInfo:
		level = color.BlueString(r.Level.String())
	case slog.LevelWarn:
		level = color.YellowString(r.Level.String())
	case slog.LevelError:
		level = color.RedString(r.Level.String())
	}
	buf := bytes.NewBuffer(make([]byte, 1024))
	first := true
	r.Attrs(func(a slog.Attr) bool {
		if first {
			first = false
		} else {
			buf.WriteString(" | ")
		}
		buf.WriteString(a.Key)
		buf.WriteString(": ")
		buf.WriteString(a.Value.String())

		return true
	})

	pid := os.Getpid()
	timeStr := r.Time.Format("15:04:05.999999")
	msg := color.CyanString(r.Message)
	id := uuid.Max.String()
	if ri := ctx.Value("request-id"); ri != nil {
		if requestId, ok := ri.(string); ok {
			id = requestId
		}
	}
	h.l.Println(pid, "|", timeStr, "|", level, "|", id, "|", msg, "|", color.WhiteString(buf.String()))

	return nil
}

func NewPrettyHandlerText(out io.Writer, opts PrettyHandlerOptions) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewTextHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}

	return h
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
