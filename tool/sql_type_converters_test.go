/*
 * This file was last modified at 2024-07-16 10:11 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * sql_type_converters_test.go
 * $Id$
 */
//!+

// Package tool TODO.
package tool

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSQLTypeConverters(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{name: "positive test #0 ConvertNullBoolToBoolPointer", fRun: testConvertNullBoolToBoolPointer},
		{name: "positive test #1 ConvertBoolPointerToNullBool", fRun: testConvertBoolPointerToNullBool},
		{name: "positive test #2 ConvertNullTimeToTimePointer", fRun: testConvertNullTimeToTimePointer},
		{name: "positive test #3 ConvertBoolPointerToNullBool", fRun: testConvertTimePointerToNullTime},
	}

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func testConvertNullBoolToBoolPointer(t *testing.T) {
	assert.Nil(t, ConvertNullBoolToBoolPointer(sql.NullBool{}))
	assert.True(t, *ConvertNullBoolToBoolPointer(sql.NullBool{Bool: true, Valid: true}))
}

func testConvertBoolPointerToNullBool(t *testing.T) {
	fl, tr := false, true
	assert.Equal(t, sql.NullBool{}, ConvertBoolPointerToNullBool(nil))
	assert.Equal(t, sql.NullBool{Valid: true}, ConvertBoolPointerToNullBool(&fl))
	assert.Equal(t, sql.NullBool{Bool: true, Valid: true}, ConvertBoolPointerToNullBool(&tr))
}

func testConvertNullTimeToTimePointer(t *testing.T) {
	assert.Nil(t, ConvertNullTimeToTimePointer(sql.NullTime{}))
	assert.Equal(t, time.Time{}, *ConvertNullTimeToTimePointer(sql.NullTime{Valid: true}))
}

func testConvertTimePointerToNullTime(t *testing.T) {
	tm := time.Time{}
	assert.Equal(t, sql.NullTime{}, ConvertTimePointerToNullTime(nil))
	assert.Equal(t, sql.NullTime{Time: tm, Valid: true}, ConvertTimePointerToNullTime(&tm))
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
