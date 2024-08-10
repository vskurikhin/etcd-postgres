/*
 * This file was last modified at 2024-07-19 15:59 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * no_rows_in_result_set.go
 * $Id$
 */
//!+

// Package tool TODO.
package tool

import "github.com/jackc/pgx/v5"

func NoRowsInResultSet(err error) bool {
	return err == pgx.ErrNoRows
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
