/*
 * This file was last modified at 2024-08-16 12:50 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * key_value.go
 * $Id$
 */
//!+

package repo

import (
	"context"
	"fmt"
	"github.com/victor-skurikhin/etcd-client/v1/internal/domain"
	"strings"
)

const (
	IDValueDeleteAction = "delete"
	IDValueGetAllAction = "getall"
	IDValueSelectAction = "select"
	IDValueUpsertAction = "upsert"
)

var (
	_ domain.Actioner[*IDValue, IDValue]  = (*idValueDelete)(nil)
	_ domain.Actioner[*IDValue, IDValue]  = (*idValueGetAll)(nil)
	_ domain.Actioner[*IDValue, IDValue]  = (*idValueSelect)(nil)
	_ domain.Actioner[*IDValue, IDValue]  = (*idValueUpsert)(nil)
	_ domain.Cloner[*IDValue, IDValue]    = (*idValueCloner)(nil)
	_ domain.Entity                       = (*IDValue)(nil)
	_ domain.Serializable                 = (*IDValue)(nil)
	_ domain.SQLEntity[*IDValue, IDValue] = (*IDValue)(nil)
	_ fmt.Stringer                        = (*IDValue)(nil)
)

var (
	ErrIDValueNil = fmt.Errorf("bad pointer, IDValue is nil")
	IDValueCloner idValueCloner
	IDValueDelete idValueDelete
	IDValueGetAll idValueGetAll
	IDValueSelect idValueSelect
	IDValueUpsert idValueUpsert
	etcdRepo      *Etcd[domain.Actioner[*IDValue, IDValue], *IDValue, IDValue]
	repoPostgres  *Postgres[domain.Actioner[*IDValue, IDValue], *IDValue, IDValue]
)

type IDValue struct {
	id    int
	value string
}

func (f *IDValue) Action(name string) domain.Actioner[*IDValue, IDValue] {
	switch strings.ToLower(name) {
	case IDValueDeleteAction:
		return IDValueDelete
	case IDValueGetAllAction:
		return IDValueGetAll
	case IDValueSelectAction:
		return IDValueSelect
	case IDValueUpsertAction:
		return IDValueUpsert
	default:
		return nil
	}
}

func (f *IDValue) Delete(
	ctx context.Context,
	repo domain.Repo[domain.Actioner[*IDValue, IDValue], *IDValue, IDValue],
) (err error) {

	if f == nil {
		return ErrIDValueNil
	}
	return f.do(ctx, IDValueDelete, repo)
}

func (f *IDValue) FromJSON(data []byte) (err error) {
	return ErrIDValueNil
}

func (f IDValue) Key() string {
	return fmt.Sprintf("%d", f.id)
}

func (f *IDValue) Value() string {

	if f == nil {
		return ""
	}
	return f.value
}

func (f *IDValue) String() string {
	if f == nil {
		return "{}"
	}
	return fmt.Sprintf(`{"id": %d, "value": "%s"}`, f.id, f.value)
}

func (f *IDValue) ToJSON() ([]byte, error) {
	return nil, ErrIDValueNil
}

func (f *IDValue) Upsert(
	ctx context.Context,
	repo domain.Repo[domain.Actioner[*IDValue, IDValue], *IDValue, IDValue],
) error {

	if f == nil {
		return ErrIDValueNil
	}
	return f.do(ctx, IDValueUpsert, repo)
}

func (f *IDValue) do(
	ctx context.Context,
	action domain.Actioner[*IDValue, IDValue],
	repo domain.Repo[domain.Actioner[*IDValue, IDValue], *IDValue, IDValue],
) (err error) {

	_, er0 := repo.Do(ctx, action, *f, func(s domain.Scanner) IDValue {

		t := *f
		err = s.Scan(&t.id, &t.value)

		if err == nil {
			*f = t
		}
		return t
	})
	if er0 != nil {
		return er0
	}
	return err
}

type idValueCloner struct{}

func (k *idValueCloner) Clone(u IDValue) IDValue {

	if k == nil {
		return IDValue{}
	}
	return IDValue{id: u.id, value: u.value}
}

func (k *idValueCloner) Copy(t *IDValue) *IDValue {

	if t == nil {
		return nil
	}
	return &IDValue{id: t.id, value: t.value}
}

type idValueDelete struct{}

func (k idValueDelete) Args(e IDValue) []any {
	return []any{e.id}
}

func (k idValueDelete) Name() string {
	return domain.DeleteAction
}

func (k idValueDelete) SQL() string {
	return `DELETE FROM test_id_value_test WHERE id = $1 RETURNING id, value`
}

type idValueGetAll struct{}

func (k idValueGetAll) Args(_ IDValue) []any {
	return []any{}
}

func (k idValueGetAll) Name() string {
	return domain.GetAllAction
}

func (k idValueGetAll) SQL() string {
	return `SELECT id, value FROM test_id_value_test`
}

type idValueSelect struct{}

func (k idValueSelect) Args(e IDValue) []any {
	return []any{e.id}
}

func (k idValueSelect) Name() string {
	return domain.SelectAction
}

func (k idValueSelect) SQL() string {
	return `SELECT id, value
	FROM test_id_value_test
	WHERE id = $1`
}

type idValueUpsert struct{}

func (k idValueUpsert) Args(e IDValue) []any {
	return []any{e.id, e.value}
}

func (k idValueUpsert) Name() string {
	return domain.UpsertAction
}

func (k idValueUpsert) SQL() string {
	return `INSERT INTO test_id_value_test
	(id, value) VALUES ($1, $2)
	ON CONFLICT (id)
	DO UPDATE SET value = $2
	RETURNING id, value`
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
