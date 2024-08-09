/*
 * This file was last modified at 2024-07-16 10:17 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * randoms.go
 * $Id$
 */
//!+

// Package tool TODO.
package tool

import (
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var rnd *rand.Rand

func RandStringBytes(n int) string {

	b := make([]byte, n)

	for i := range b {
		b[i] = letterBytes[rnd.Intn(len(letterBytes))]
	}
	return string(b)
}

func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
