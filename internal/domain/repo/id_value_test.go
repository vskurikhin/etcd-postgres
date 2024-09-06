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
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/victor-skurikhin/etcd-client/v1/internal/domain"
	"strconv"
	"strings"
	"testing"
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
	id      int
	value   string
	version sql.NullInt64
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

func (f *IDValue) Version() int64 {

	if f == nil {
		return 0
	}
	if f.version.Valid {
		return f.version.Int64
	}
	return 0
}

func (f *IDValue) String() string {
	if f == nil {
		return "{}"
	}
	return fmt.Sprintf(`{"id": %d, "value": "%s", "version": %d}`, f.id, f.value, f.version.Int64)
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
	return IDValue{id: u.id, value: u.value, version: u.version}
}

func (k *idValueCloner) Copy(t *IDValue) *IDValue {

	if t == nil {
		return nil
	}
	return &IDValue{id: t.id, value: t.value, version: t.version}
}

type idValueDelete struct{}

func (k idValueDelete) Args(e IDValue) []any {
	return []any{e.id}
}

func (k idValueDelete) Name() string {
	return domain.DeleteAction
}

func (k idValueDelete) SQL() string {
	return `DELETE FROM test_id_value_test WHERE id = $1 RETURNING id, value, version`
}

type idValueGetAll struct{}

func (k idValueGetAll) Args(_ IDValue) []any {
	return []any{}
}

func (k idValueGetAll) Name() string {
	return domain.GetAllAction
}

func (k idValueGetAll) SQL() string {
	return `SELECT id, value, version FROM test_id_value_test`
}

type idValueSelect struct{}

func (k idValueSelect) Args(e IDValue) []any {
	return []any{e.id}
}

func (k idValueSelect) Name() string {
	return domain.SelectAction
}

func (k idValueSelect) SQL() string {
	return `SELECT id, value, version
	FROM test_id_value_test
	WHERE id = $1`
}

type idValueUpsert struct{}

func (k idValueUpsert) Args(e IDValue) []any {
	return []any{e.id, e.value, e.version}
}

func (k idValueUpsert) Name() string {
	return domain.UpsertAction
}

func (k idValueUpsert) SQL() string {
	return `INSERT INTO test_id_value_test
	(id, value, version) VALUES ($1, $2, $3)
	ON CONFLICT (id)
	DO UPDATE SET value = $2
	RETURNING id, value, version`
}

func getEtcdScanFunc(tb testing.TB, res *IDValue) func(scanner domain.Scanner) IDValue {
	return func(scanner domain.Scanner) IDValue {
		var (
			key     string
			value   string
			version sql.NullInt64
		)
		er0 := scanner.Scan(&key, &value, &version)
		assert.Nil(tb, er0)
		id, er1 := strconv.Atoi(key)
		assert.Nil(tb, er1)
		res.version = version
		return IDValue{id: id, value: value, version: res.version}
	}
}

func getPostgresScanFunc(tb testing.TB) func(scanner domain.Scanner) IDValue {
	return func(scanner domain.Scanner) IDValue {
		var res IDValue
		er0 := scanner.Scan(&res.id, &res.value, &res.version)
		assert.Nil(tb, er0)
		return res
	}
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
