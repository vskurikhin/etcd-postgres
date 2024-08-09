/*
 * This file was last modified at 2024-07-31 15:37 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * package.go
 * $Id$
 */
//!+

// Package tool TODO.
package tool

import (
	"log/slog"
	"sync"
)

const MSG = "etcd-client.tool "

var (
	once = new(sync.Once)
	sLog *slog.Logger
)

func SetLogger(sl *slog.Logger) *slog.Logger {
	once.Do(func() { sLog = sl })
	return sLog
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
