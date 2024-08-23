/*
 * This file was last modified at 2024-08-16 12:12 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * domain.go
 * $Id$
 */
//!+

// Package domain TODO.
package domain

import (
	"context"
)

const (
	DeleteAction = "delete"
	GetAllAction = "getall"
	SelectAction = "select"
	UpsertAction = "upsert"
)

// Actioner the first type param will match pointer types and infer U
type Actioner[T Ptr[U], U Entity] interface {
	Args(U) []any
	Name() string
	SQL() string
}

// Cloner the first type param will match pointer types and infer U
type Cloner[T Ptr[U], U Entity] interface {
	Clone(U) U
	Copy(T) T
}

type Entity interface {
	Key() string
}

// Ptr constraining a type to its pointer type
type Ptr[T Entity] interface {
	*T
}

type Repo[A Actioner[T, U], T Ptr[U], U Entity] interface {
	Do(ctx context.Context, action A, unit U, scan func(Scanner) U) (U, error)
	Get(ctx context.Context, action A, unit U, scan func(Scanner) U) ([]U, error)
}

type Scanner interface {
	Scan(dest ...any) error
}

type Serializable interface {
	Entity
	FromJSON(data []byte) (err error)
	ToJSON() ([]byte, error)
}

type SQLEntity[T Ptr[U], U Entity] interface {
	Entity
	Serializable
	Action(name string) Actioner[T, U]
}

type TransactionalAction interface {
	DeleteTxArgs(...any) TxArgs
	Name() string
	UpsertTxArgs(...any) TxArgs
}

type TxArgs struct {
	Args [][]any
	SQLs []string
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
