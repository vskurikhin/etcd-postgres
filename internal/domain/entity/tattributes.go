/*
 * This file was last modified at 2024-07-18 19:05 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * tattributes.go
 * $Id$
 */

package entity

import (
	"database/sql"
	"fmt"
	"time"
)

type TAttributes struct {
	deleted   sql.NullBool
	createdAt time.Time
	updatedAt sql.NullTime
}

func DefaultTAttributes() TAttributes {
	return TAttributes{deleted: sql.NullBool{}, createdAt: time.Time{}, updatedAt: sql.NullTime{}}
}

func MakeTAttributes(deleted sql.NullBool, createdAt time.Time, updatedAt sql.NullTime) TAttributes {
	return TAttributes{deleted: deleted, createdAt: createdAt, updatedAt: updatedAt}
}

func (t *TAttributes) String() string {
	if t == nil {
		return ""
	}
	return fmt.Sprintf(
		`"deleted": %v, "createdAt": "%v", "updatedAt": "%v"`,
		t.deleted, t.createdAt, t.updatedAt,
	)
}

type tAttributes struct {
	Deleted   JsonNullBool `json:"deleted,inline"`
	CreatedAt time.Time    `json:"created_at,inline"`
	UpdatedAt JsonNullTime `json:"updated_at,inline"`
}

func makeTAttributes(deleted JsonNullBool, createdAt time.Time, updatedAt JsonNullTime) tAttributes {
	return tAttributes{Deleted: deleted, CreatedAt: createdAt, UpdatedAt: updatedAt}
}
