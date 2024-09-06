/*
 * This file was last modified at 2024-07-16 21:11 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * stub_repo_test.go
 * $Id$
 */
//!+

// Package entity TODO.
package entity

import (
	"context"
	"fmt"
	"github.com/victor-skurikhin/etcd-client/v1/internal/domain"
)

var _ domain.Repo[domain.Actioner[*domain.Entity, domain.Entity], *domain.Entity, domain.Entity] = (*stubRepoOk[domain.Actioner[*domain.Entity, domain.Entity], *domain.Entity, domain.Entity])(nil)
var ErrStub = fmt.Errorf("stub error")

type stubRepoOk[A domain.Actioner[T, U], T domain.Ptr[U], U domain.Entity] struct {
}

func (s stubRepoOk[A, T, U]) Do(_ context.Context, action A, entity U, scan func(domain.Scanner) U) (U, error) {
	if action.SQL() == "" {
		panic(`action.SQL() == ""`)
	}
	if len(action.Args(entity)) < 1 {
		panic(`len(action.Args(entity)) < 1`)
	}
	scan(&stubScannerOk{})
	return entity, nil
}

func (s stubRepoOk[A, T, U]) Get(_ context.Context, action A, _ U, scan func(domain.Scanner) U) ([]U, error) {
	if action.SQL() == "" {
		panic(`action.SQL() == ""`)
	}
	scan(&stubScannerOk{})
	return make([]U, 0), nil
}

type stubScannerOk struct {
}

func (s *stubScannerOk) Scan(_ ...any) error {
	return nil
}

type stubRepoErr[A domain.Actioner[T, U], T domain.Ptr[U], U domain.Entity] struct {
}

func (s stubRepoErr[A, T, U]) Do(_ context.Context, _ A, entity U, scan func(domain.Scanner) U) (U, error) {
	var u U
	scan(&stubScannerOk{})
	return u, ErrStub
}

func (s stubRepoErr[A, T, U]) Get(_ context.Context, _ A, _ U, scan func(domain.Scanner) U) ([]U, error) {
	scan(&stubScannerOk{})
	return make([]U, 0), ErrStub
}

type stubRepoScannerErr[A domain.Actioner[T, U], T domain.Ptr[U], U domain.Entity] struct {
}

func (s stubRepoScannerErr[A, T, U]) Do(_ context.Context, _ A, entity U, scan func(domain.Scanner) U) (U, error) {
	scan(&stubScannerErr{})
	return entity, nil
}

func (s stubRepoScannerErr[A, T, U]) Get(_ context.Context, _ A, _ U, scan func(domain.Scanner) U) ([]U, error) {
	scan(&stubScannerErr{})
	return make([]U, 0), nil
}

type stubScannerErr struct {
}

func (s *stubScannerErr) Scan(_ ...any) error {
	return ErrStub
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
