/*
 * This file was last modified at 2024-07-16 10:11 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * sql_type_converters.go
 * $Id$
 */
//!+

// Package tool TODO.
package tool

import (
	"database/sql"
	"time"
)

func ConvertNullBoolToBoolPointer(b sql.NullBool) *bool {
	if b.Valid {
		return &b.Bool
	}
	return nil
}

func ConvertBoolPointerToNullBool(b *bool) sql.NullBool {
	if b != nil {
		return sql.NullBool{Bool: *b, Valid: true}
	}
	return sql.NullBool{}
}

func ConvertNullTimeToTimePointer(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func ConvertTimePointerToNullTime(t *time.Time) sql.NullTime {
	if t != nil {
		return sql.NullTime{Time: *t, Valid: true}
	}
	return sql.NullTime{}
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
