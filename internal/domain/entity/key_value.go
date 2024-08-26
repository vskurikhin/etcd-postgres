/*
 * This file was last modified at 2024-08-16 12:50 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * key_value.go
 * $Id$
 */
//!+

package entity

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/victor-skurikhin/etcd-client/v1/internal/domain"
	"strings"
)

var (
	_ domain.Actioner[*KeyValue, KeyValue]  = (*keyValueDelete)(nil)
	_ domain.Actioner[*KeyValue, KeyValue]  = (*keyValueGetAll)(nil)
	_ domain.Actioner[*KeyValue, KeyValue]  = (*keyValueSelect)(nil)
	_ domain.Actioner[*KeyValue, KeyValue]  = (*keyValueUpsert)(nil)
	_ domain.Cloner[*KeyValue, KeyValue]    = (*keyValueCloner)(nil)
	_ domain.Entity                         = (*KeyValue)(nil)
	_ domain.Serializable                   = (*KeyValue)(nil)
	_ domain.SQLEntity[*KeyValue, KeyValue] = (*KeyValue)(nil)
	_ fmt.Stringer                          = (*KeyValue)(nil)
)

var (
	ErrKeyValueNil = fmt.Errorf("bad pointer, KeyValue is nil")
	KeyValueCloner keyValueCloner
	KeyValueDelete keyValueDelete
	KeyValueGetAll keyValueGetAll
	KeyValueSelect keyValueSelect
	KeyValueUpsert keyValueUpsert
)

type KeyValue struct {
	TAttributes
	key     string
	value   string
	version sql.NullInt64
}

type keyValue struct {
	tAttributes
	Key     string `json:"key"`
	Value   string `json:"value,omitempty"`
	Version int64  `json:"version,omitempty"`
}

func GetAllKeyValue(
	ctx context.Context,
	repo domain.Repo[domain.Actioner[*KeyValue, KeyValue], *KeyValue, KeyValue],
) ([]KeyValue, error) {

	var err error

	result, er0 := repo.Get(ctx, KeyValueGetAll, KeyValue{}, func(s domain.Scanner) KeyValue {
		var r KeyValue
		err = s.Scan(&r.key, &r.value, &r.version, &r.deleted, &r.createdAt, &r.updatedAt)
		return r
	})
	if er0 != nil {
		return result, er0
	}
	return result, err
}

func MakeKeyValue(key, value string, version int64, t TAttributes) KeyValue {
	return KeyValue{
		TAttributes: MakeTAttributes(t.deleted, t.createdAt, t.updatedAt),
		key:         key,
		value:       value,
		version:     VersionToNullInt64(version),
	}
}

func NewKeyValue(key, value string, version int64, t TAttributes) *KeyValue {
	n := MakeKeyValue(key, value, version, t)
	return &n
}

func GetKeyValue(
	ctx context.Context,
	repo domain.Repo[domain.Actioner[*KeyValue, KeyValue], *KeyValue, KeyValue],
	key string,
) (r KeyValue, err error) {
	_, er0 := repo.Do(ctx, KeyValueSelect, KeyValue{key: key}, func(s domain.Scanner) KeyValue {
		err = s.Scan(&r.key, &r.value, &r.version, &r.deleted, &r.createdAt, &r.updatedAt)
		return r
	})
	if er0 != nil {
		return r, er0
	}
	return r, err
}

func (f *KeyValue) Action(name string) domain.Actioner[*KeyValue, KeyValue] {
	switch strings.ToLower(name) {
	case domain.DeleteAction:
		return KeyValueDelete
	case domain.GetAllAction:
		return KeyValueGetAll
	case domain.SelectAction:
		return KeyValueSelect
	case domain.UpsertAction:
		return KeyValueUpsert
	default:
		return nil
	}
}

func (f *KeyValue) FromJSON(data []byte) (err error) {

	if f == nil {
		return ErrKeyValueNil
	}
	var t keyValue
	err = json.UnmarshalNoEscape(data, &t)

	if err != nil {
		return err
	}
	*f = MakeKeyValue(t.Key, t.Value, t.Version, MakeTAttributes(
		t.Deleted.ToNullBool(), t.CreatedAt, t.UpdatedAt.ToNullTime(),
	))
	return nil
}

func (f KeyValue) Key() string {
	return f.key
}

func (f *KeyValue) Value() string {

	if f == nil {
		return ""
	}
	return f.value
}

func (f *KeyValue) Version() int64 {

	if f == nil {
		return 0
	}
	if f.version.Valid {
		return f.version.Int64
	}
	return 0
}

func (f *KeyValue) String() string {
	if f == nil {
		return "{}"
	}
	return fmt.Sprintf(
		`{"key": "%s", "value": "%s", "version": %d, %s}`,
		f.key, f.value, f.version.Int64, f.TAttributes.String())
}

func (f *KeyValue) Delete(
	ctx context.Context,
	repo domain.Repo[domain.Actioner[*KeyValue, KeyValue], *KeyValue, KeyValue],
) (err error) {

	if f == nil {
		return ErrKeyValueNil
	}
	return f.do(ctx, KeyValueDelete, repo)
}

func (f *KeyValue) ToJSON() ([]byte, error) {

	if f == nil {
		return nil, ErrKeyValueNil
	}
	result, err := json.MarshalNoEscape(keyValue{
		tAttributes: makeTAttributes(
			FromNullBool(f.deleted), f.createdAt, FromNullTime(f.updatedAt),
		),
		Key:     f.key,
		Value:   f.value,
		Version: FromNullInt64ToVersion(f.version),
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (f *KeyValue) Upsert(
	ctx context.Context,
	repo domain.Repo[domain.Actioner[*KeyValue, KeyValue], *KeyValue, KeyValue],
) error {

	if f == nil {
		return ErrKeyValueNil
	}
	return f.do(ctx, KeyValueUpsert, repo)
}

func (f *KeyValue) do(
	ctx context.Context,
	action domain.Actioner[*KeyValue, KeyValue],
	repo domain.Repo[domain.Actioner[*KeyValue, KeyValue], *KeyValue, KeyValue],
) (err error) {

	_, er0 := repo.Do(ctx, action, *f, func(s domain.Scanner) KeyValue {

		t := *f
		err = s.Scan(&t.key, &t.value, &t.version, &t.deleted, &t.createdAt, &t.updatedAt)

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

type keyValueCloner struct{}

func (k *keyValueCloner) Clone(u KeyValue) KeyValue {
	return MakeKeyValue(u.key, u.value, FromNullInt64ToVersion(u.version), u.TAttributes)
}

func (k *keyValueCloner) Copy(t *KeyValue) *KeyValue {

	if t == nil {
		return nil
	}
	return NewKeyValue(t.key, t.value, FromNullInt64ToVersion(t.version), t.TAttributes)
}

type keyValueDelete struct{}

func (k keyValueDelete) Args(e KeyValue) []any {
	return []any{e.key, e.updatedAt}
}

func (k keyValueDelete) Name() string {
	return domain.DeleteAction
}

func (k keyValueDelete) SQL() string {
	return `UPDATE key_value
	SET deleted = true, updated_at = $2
	WHERE key = $1
	RETURNING key, value, deleted, created_at, updated_at`
}

type keyValueGetAll struct{}

func (k keyValueGetAll) Args(_ KeyValue) []any {
	return []any{}
}

func (k keyValueGetAll) Name() string {
	return domain.GetAllAction
}

func (k keyValueGetAll) SQL() string {
	return `SELECT key, value, deleted, created_at, updated_at FROM key_value`
}

type keyValueSelect struct{}

func (k keyValueSelect) Args(e KeyValue) []any {
	return []any{e.key}
}

func (k keyValueSelect) Name() string {
	return domain.SelectAction
}

func (k keyValueSelect) SQL() string {
	return `SELECT key, value, deleted, created_at, updated_at
	FROM key_value
	WHERE key = $1`
}

type keyValueUpsert struct{}

func (k keyValueUpsert) Args(e KeyValue) []any {
	return []any{e.key, e.value, e.deleted, e.createdAt, e.updatedAt}
}

func (k keyValueUpsert) Name() string {
	return domain.UpsertAction
}

func (k keyValueUpsert) SQL() string {
	return `INSERT INTO key_value
	(key, value, deleted, created_at)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (key)
	DO UPDATE SET value = $2, deleted = $3, updated_at = $5
	RETURNING key, value, deleted, created_at, updated_at`
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
