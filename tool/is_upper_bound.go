/*
 * This file was last modified at 2024-07-29 22:46 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * is_upper_bound.go
 * $Id$
 */

package tool

import "time"

func IsUpperBound(index int, duration time.Duration) bool {
	result := 25 * time.Millisecond * time.Duration(index) * time.Duration(index+1)
	return result < duration
}

func IsUpperBoundWithSleep(index, sleep int, duration time.Duration) bool {
	result := time.Duration(sleep/2) * time.Millisecond * time.Duration(index) * time.Duration(index+1)
	return result < duration
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
