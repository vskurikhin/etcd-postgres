/*
 * This file was last modified at 2024-09-05 12:29 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * repo.go
 * $Id$
 */

package domain

import (
	"context"
)

type Repo[A Actioner[T, U], T Ptr[U], U Entity] interface {
	Do(ctx context.Context, action A, unit U, scan func(Scanner) U) (U, error)
	Get(ctx context.Context, action A, unit U, scan func(Scanner) U) ([]U, error)
}
