/*
 * This file was last modified at 2024-07-11 11:30 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * if_error_then_panic.go
 * $Id$
 */
//!+

// Package tool TODO.
package tool

func IfErrorThenPanic(e error) {
	if e != nil {
		panic(e)
	}
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
