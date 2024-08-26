/*
 * This file was last modified at 2024-08-16 12:50 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * key_value_test.go
 * $Id$
 */
//!+

package entity

import (
	"database/sql"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/victor-skurikhin/etcd-client/v1/internal/domain"
	"golang.org/x/net/context"
	"testing"
)

func TestKeyValue(t *testing.T) {
	for _, test := range []struct {
		name string
		fRun func(*testing.T) (interface{}, error)
		want func(*testing.T, interface{}) bool
	}{
		{
			"test #0 positive for struct KeyValue method Action(DeleteAction)",
			positiveActionKeyValueDeleteAction,
			positiveActionKeyValueDeleteActionCheck,
		},
		{
			"test #1 positive for struct KeyValue method Action(GetAllAction)",
			positiveActionKeyValueGetAllAction,
			positiveActionKeyValueGetAllActionCheck,
		},
		{
			"test #2 positive for struct KeyValue method Action(SelectAction)",
			positiveActionKeyValueSelectAction,
			positiveActionKeyValueSelectActionCheck,
		},
		{
			"test #3 positive for struct KeyValue method Action(UpsertAction)",
			positiveActionKeyValueUpsertAction,
			positiveActionKeyValueUpsertActionCheck,
		},
		{
			`test #4 negative for struct KeyValue method Action("")`,
			negativeActionKeyValueEmptyAction,
			negativeActionKeyValueEmptyActionCheck,
		},
		{
			"test #5 positive for struct KeyValue method Delete(context.Context, domain.Repo)",
			positiveKeyValueDelete,
			positiveKeyValueDeleteCheck,
		},
		{
			"test #6 negative for struct KeyValue method Delete(context.Context, domain.Repo)",
			negativeKeyValueDelete,
			negativeKeyValueDeleteCheck,
		},
		{
			"test #7 positive case #1 for struct KeyValue methods ToJSON() and FromJSON([]byte)",
			positiveKeyValueJSON1,
			positiveKeyValueJSON1Check,
		},
		{
			"test #8 positive case #2 for struct KeyValue methods ToJSON() and FromJSON([]byte)",
			positiveKeyValueJSON2,
			positiveKeyValueJSON2Check,
		},
		{
			"test #9 positive case #3 for struct KeyValue methods ToJSON() and FromJSON([]byte)",
			positiveKeyValueJSON3,
			positiveKeyValueJSON3Check,
		},
		{
			"test #10 negative #1 for struct KeyValue method FromJSON([]byte)",
			negativeKeyValueJSON1,
			negativeKeyValueJSON1Check,
		},
		{
			"test #11 negative #2 for struct KeyValue method FromJSON([]byte)",
			negativeKeyValueJSON2,
			negativeKeyValueJSON2Check,
		},
		{
			"test #12 negative #3 for struct KeyValue method Delete(context.Context, domain.Repo)",
			negativeKeyValueJSON3,
			negativeKeyValueJSON3Check,
		},
		{
			"test #13 positive for struct KeyValue method Key()",
			positiveKeyValueKey,
			positiveKeyValueKeyCheck,
		},
		{
			"test #15 positive for struct KeyValue method Value()",
			positiveKeyValueValue,
			positiveKeyValueValueCheck,
		},
		{
			"test #16 negative for struct KeyValue method Value()",
			negativeKeyValueValue,
			negativeKeyValueValueCheck,
		},
		{
			"test #17 positive for struct KeyValue method Value()",
			positiveKeyValueString,
			positiveKeyValueStringCheck,
		},
		{
			"test #18 negative for struct KeyValue method Value()",
			negativeKeyValueValue,
			negativeKeyValueValueCheck,
		},

		{
			"test #19 positive for struct KeyValue method Upsert(context.Context, domain.Repo)",
			positiveKeyValueUpsert,
			positiveKeyValueUpsertCheck,
		},
		{
			"test #20 negative #1 for struct KeyValue method Upsert(context.Context, domain.Repo)",
			negativeKeyValueUpsert,
			negativeKeyValueUpsertCheck,
		},
		{
			"test #21 negative #2 for struct KeyValue method Upsert(context.Context, domain.Repo)",
			negativeRepoErrKeyValueUpsert,
			negativeRepoErrKeyValueUpsertCheck,
		},
		{
			"test #22 negative #3 for struct KeyValue method Upsert(context.Context, domain.Repo)",
			negativeScannerErrKeyValueUpsert,
			negativeScannerErrKeyValueUpsertCheck,
		},
		{
			"test #23 positive for function GetKeyValue(context.Context, domain.Repo, string)",
			positiveSelectKeyValue,
			positiveSelectKeyValueCheck,
		},
		{
			"test #24 negative #1 for function GetKeyValue(context.Context, domain.Repo, string)",
			negativeRepoErrSelectKeyValue,
			negativeRepoErrSelectKeyValueCheck,
		},
		{
			"test #25 negative #2 for function GetKeyValue(context.Context, domain.Repo, string)",
			negativeScannerErrSelectKeyValue,
			negativeScannerErrSelectKeyValueCheck,
		},
		{
			"test #26 positive for function GetAllKeyValue(context.Context, domain.Repo)",
			positiveGetAllKeyValue,
			positiveGetAllKeyValueCheck,
		},
		{
			"test #27 negative #1 for function GetAllKeyValue(context.Context, domain.Repo)",
			negativeRepoErrGetAllKeyValue,
			negativeRepoErrGetAllKeyValueCheck,
		},
		{
			"test #28 negative #1 for function GetAllKeyValue(context.Context, domain.Repo)",
			negativeScannerErrGetAllKeyValue,
			negativeScannerErrGetAllKeyValueCheck,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.fRun(t)
			assert.Nil(t, err)
			assert.NotNil(t, got)
			assert.True(t, test.want(t, got))
		})
	}
}

func positiveActionKeyValueDeleteAction(_ *testing.T) (interface{}, error) {
	return (&KeyValue{}).Action(domain.DeleteAction) == KeyValueDelete, nil
}

func positiveActionKeyValueDeleteActionCheck(_ *testing.T, i interface{}) bool {
	return checkTrue(i)
}

func positiveActionKeyValueGetAllAction(_ *testing.T) (interface{}, error) {
	return (&KeyValue{}).Action(domain.GetAllAction) == KeyValueGetAll, nil
}

func positiveActionKeyValueGetAllActionCheck(_ *testing.T, i interface{}) bool {
	return checkTrue(i)
}

func positiveActionKeyValueSelectAction(_ *testing.T) (interface{}, error) {
	return (&KeyValue{}).Action(domain.SelectAction) == KeyValueSelect, nil
}

func positiveActionKeyValueSelectActionCheck(_ *testing.T, i interface{}) bool {
	return checkTrue(i)
}

func positiveActionKeyValueUpsertAction(_ *testing.T) (interface{}, error) {
	return (&KeyValue{}).Action(domain.UpsertAction) == KeyValueUpsert, nil
}

func positiveActionKeyValueUpsertActionCheck(_ *testing.T, i interface{}) bool {
	return checkTrue(i)
}

func negativeActionKeyValueEmptyAction(_ *testing.T) (interface{}, error) {
	return (&KeyValue{}).Action("") == nil, nil
}

func negativeActionKeyValueEmptyActionCheck(_ *testing.T, i interface{}) bool {
	return checkTrue(i)
}

func positiveKeyValueDelete(_ *testing.T) (interface{}, error) {
	err := (&KeyValue{}).Delete(context.Background(), &stubRepoOk[domain.Actioner[*KeyValue, KeyValue], *KeyValue, KeyValue]{})
	return true, err
}

func positiveKeyValueDeleteCheck(_ *testing.T, i interface{}) bool {
	return checkTrue(i)
}

func negativeKeyValueDelete(_ *testing.T) (interface{}, error) {
	var kv *KeyValue
	err := kv.Delete(context.Background(), &stubRepoOk[domain.Actioner[*KeyValue, KeyValue], *KeyValue, KeyValue]{})
	return err == ErrKeyValueNil, nil
}

func negativeKeyValueDeleteCheck(_ *testing.T, i interface{}) bool {
	return checkTrue(i)
}

func positiveKeyValueJSON1(t *testing.T) (interface{}, error) {
	expected := MakeKeyValue("key1", "value1", 0, DefaultTAttributes())
	j, err := expected.ToJSON()
	fmt.Printf("j: %s\n", string(j))
	assert.Nil(t, err)
	assert.NotNil(t, j)
	got := KeyValue{}
	err = (&got).FromJSON(j)
	assert.Nil(t, err)
	assert.Equal(t, expected, got)
	assert.Equal(t, expected.String(), got.String())
	return &expected, nil
}

func positiveKeyValueJSON1Check(_ *testing.T, i interface{}) bool {
	return true
}

func positiveKeyValueJSON2(t *testing.T) (interface{}, error) {
	expected := MakeKeyValue("key1", "", 0, TAttributes{
		deleted:   sql.NullBool{Valid: true},
		updatedAt: sql.NullTime{Valid: true},
	})
	j, err := expected.ToJSON()
	assert.Nil(t, err)
	assert.NotNil(t, j)
	got := KeyValue{}
	err = (&got).FromJSON(j)
	assert.Nil(t, err)
	assert.Equal(t, expected, got)
	assert.Equal(t, expected.String(), got.String())
	return &expected, nil
}

func positiveKeyValueJSON2Check(_ *testing.T, i interface{}) bool {
	return true
}

func positiveKeyValueJSON3(t *testing.T) (interface{}, error) {
	expected := MakeKeyValue("key1", "", 0, TAttributes{deleted: sql.NullBool{true, true}})
	j, err := expected.ToJSON()
	assert.Nil(t, err)
	assert.NotNil(t, j)
	got := KeyValue{}
	err = (&got).FromJSON(j)
	assert.Nil(t, err)
	assert.Equal(t, expected, got)
	assert.Equal(t, expected.String(), got.String())
	return &expected, nil
}

func positiveKeyValueJSON3Check(_ *testing.T, i interface{}) bool {
	return true
}

func negativeKeyValueJSON1(_ *testing.T) (interface{}, error) {
	var kv *KeyValue
	err := kv.FromJSON([]byte("{}"))
	return err, nil
}

func negativeKeyValueJSON1Check(_ *testing.T, i interface{}) bool {
	return i == ErrKeyValueNil
}

func negativeKeyValueJSON2(_ *testing.T) (interface{}, error) {
	var kv KeyValue
	err := kv.FromJSON([]byte(""))
	return err, nil
}

func negativeKeyValueJSON2Check(_ *testing.T, i interface{}) bool {
	_, ok := i.(*json.SyntaxError)
	return ok
}

func negativeKeyValueJSON3(_ *testing.T) (interface{}, error) {
	var kv *KeyValue
	_, err := kv.ToJSON()
	return err, nil
}

func negativeKeyValueJSON3Check(_ *testing.T, i interface{}) bool {
	return i == ErrKeyValueNil
}

func positiveKeyValueKey(_ *testing.T) (interface{}, error) {
	kv := KeyValue{key: "key1"}
	return kv.Key(), nil
}

func positiveKeyValueKeyCheck(_ *testing.T, i interface{}) bool {
	return i == "key1"
}

func positiveKeyValueValue(_ *testing.T) (interface{}, error) {
	kv := KeyValue{value: "value1"}
	return kv.Value(), nil
}

func positiveKeyValueValueCheck(_ *testing.T, i interface{}) bool {
	return i == "value1"
}

func negativeKeyValueValue(_ *testing.T) (interface{}, error) {
	var kv *KeyValue
	return kv.Value(), nil
}

func negativeKeyValueValueCheck(_ *testing.T, i interface{}) bool {
	return i == ""
}

func positiveKeyValueString(_ *testing.T) (interface{}, error) {
	kv := KeyValue{key: "key1", value: "value1"}
	return kv.String(), nil
}

func positiveKeyValueStringCheck(_ *testing.T, i interface{}) bool {
	return i == `{"key": "key1", "value": "value1", "version": 0, "deleted": {false false}, "createdAt": "0001-01-01 00:00:00 +0000 UTC", "updatedAt": "{0001-01-01 00:00:00 +0000 UTC false}"}`
}

func positiveKeyValueUpsert(_ *testing.T) (interface{}, error) {
	err := (&KeyValue{}).Upsert(context.Background(), &stubRepoOk[domain.Actioner[*KeyValue, KeyValue], *KeyValue, KeyValue]{})
	return true, err
}

func positiveKeyValueUpsertCheck(_ *testing.T, i interface{}) bool {
	return checkTrue(i)
}

func negativeKeyValueUpsert(_ *testing.T) (interface{}, error) {
	var kv *KeyValue
	err := kv.Upsert(context.Background(), &stubRepoOk[domain.Actioner[*KeyValue, KeyValue], *KeyValue, KeyValue]{})
	return err == ErrKeyValueNil, nil
}

func negativeKeyValueUpsertCheck(_ *testing.T, i interface{}) bool {
	return checkTrue(i)
}

func negativeRepoErrKeyValueUpsert(_ *testing.T) (interface{}, error) {
	err := (&KeyValue{}).Upsert(context.Background(), &stubRepoErr[domain.Actioner[*KeyValue, KeyValue], *KeyValue, KeyValue]{})
	return err == ErrStub, nil
}

func negativeRepoErrKeyValueUpsertCheck(_ *testing.T, i interface{}) bool {
	return checkTrue(i)
}

func negativeScannerErrKeyValueUpsert(_ *testing.T) (interface{}, error) {
	err := (&KeyValue{}).Upsert(context.Background(), &stubRepoScannerErr[domain.Actioner[*KeyValue, KeyValue], *KeyValue, KeyValue]{})
	return err == ErrStub, nil
}

func negativeScannerErrKeyValueUpsertCheck(_ *testing.T, i interface{}) bool {
	return checkTrue(i)
}

func positiveSelectKeyValue(_ *testing.T) (interface{}, error) {
	return GetKeyValue(context.Background(), &stubRepoOk[domain.Actioner[*KeyValue, KeyValue], *KeyValue, KeyValue]{}, "")
}

func positiveSelectKeyValueCheck(_ *testing.T, i interface{}) bool {
	return i == KeyValue{}
}

func negativeRepoErrSelectKeyValue(_ *testing.T) (interface{}, error) {
	_, err := GetKeyValue(context.Background(), &stubRepoErr[domain.Actioner[*KeyValue, KeyValue], *KeyValue, KeyValue]{}, "")
	return err, nil
}

func negativeRepoErrSelectKeyValueCheck(_ *testing.T, i interface{}) bool {
	return i == ErrStub
}

func negativeScannerErrSelectKeyValue(_ *testing.T) (interface{}, error) {
	_, err := GetKeyValue(context.Background(), &stubRepoScannerErr[domain.Actioner[*KeyValue, KeyValue], *KeyValue, KeyValue]{}, "")
	return err, nil
}

func negativeScannerErrSelectKeyValueCheck(_ *testing.T, i interface{}) bool {
	return i == ErrStub
}

func positiveGetAllKeyValue(_ *testing.T) (interface{}, error) {
	return GetAllKeyValue(context.Background(), &stubRepoOk[domain.Actioner[*KeyValue, KeyValue], *KeyValue, KeyValue]{})
}

func positiveGetAllKeyValueCheck(_ *testing.T, i interface{}) bool {
	if r, ok := i.([]KeyValue); ok {
		return len(r) == 0
	}
	return false
}

func negativeRepoErrGetAllKeyValue(_ *testing.T) (interface{}, error) {
	_, err := GetAllKeyValue(context.Background(), &stubRepoErr[domain.Actioner[*KeyValue, KeyValue], *KeyValue, KeyValue]{})
	return err, nil
}

func negativeRepoErrGetAllKeyValueCheck(_ *testing.T, i interface{}) bool {
	return i == ErrStub
}

func negativeScannerErrGetAllKeyValue(_ *testing.T) (interface{}, error) {
	_, err := GetAllKeyValue(context.Background(), &stubRepoScannerErr[domain.Actioner[*KeyValue, KeyValue], *KeyValue, KeyValue]{})
	return err, nil
}

func negativeScannerErrGetAllKeyValueCheck(_ *testing.T, i interface{}) bool {
	return i == ErrStub
}

func TestKeyValueCloner(t *testing.T) {
	for _, test := range []struct {
		name string
		fRun func(*testing.T) (interface{}, error)
		want func(*testing.T, interface{}) bool
	}{
		{
			"test #0 positive for interface Cloner[T Ptr[U], U any] method Clone()",
			positiveKeyValueClonerClone,
			positiveKeyValueClonerCloneCheck,
		},
		{
			"test #1 positive for interface Cloner[T Ptr[U], U any] method Copy()",
			positiveKeyValueClonerCopy,
			positiveKeyValueClonerCopyCheck,
		},
		{
			"test #2 negative for interface Cloner[T Ptr[U], U any] method Copy()",
			negativeKeyValueClonerCopy,
			negativeKeyValueClonerCopyCheck,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.fRun(t)
			assert.Nil(t, err)
			assert.NotNil(t, got)
			assert.True(t, test.want(t, got))
		})
	}
}

func positiveKeyValueClonerClone(_ *testing.T) (interface{}, error) {
	f := KeyValue{key: "key1", value: "value1"}
	return KeyValueCloner.Clone(f), nil
}

func positiveKeyValueClonerCloneCheck(_ *testing.T, i interface{}) bool {
	return i == KeyValue{key: "key1", value: "value1"}
}

func positiveKeyValueClonerCopy(_ *testing.T) (interface{}, error) {
	f := KeyValue{key: "key1", value: "value1"}
	return KeyValueCloner.Copy(&f), nil
}

func positiveKeyValueClonerCopyCheck(_ *testing.T, i interface{}) bool {
	if u, ok := i.(*KeyValue); ok {
		return *u == KeyValue{key: "key1", value: "value1"}
	}
	return false
}

func negativeKeyValueClonerCopy(_ *testing.T) (interface{}, error) {
	if n := KeyValueCloner.Copy(nil); n == nil {
		return true, nil
	}
	return nil, fmt.Errorf("then Copy(nil) when isn't nil")
}

func negativeKeyValueClonerCopyCheck(_ *testing.T, i interface{}) bool {
	return i == true
}

func TestKeyValueDeleteAction(t *testing.T) {
	for _, test := range []struct {
		name string
		fRun func(*testing.T) (interface{}, error)
		want func(*testing.T, interface{}) bool
	}{
		{
			"test #0 positive for struct keyValueDelete",
			positiveKeyValueDeleteAction,
			positiveKeyValueDeleteActionCheck,
		},
		{
			"test #1 positive for struct KeyValueGetAll",
			positiveKeyValueGetAllAction,
			positiveKeyValueGetAllActionCheck,
		},
		{
			"test #2 positive for struct KeyValueSelect",
			positiveKeyValueSelectAction,
			positiveKeyValueSelectActionCheck,
		},
		{
			"test #3 positive for struct KeyValueUpsert",
			positiveKeyValueUpsertAction,
			positiveKeyValueUpsertActionCheck,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.fRun(t)
			assert.Nil(t, err)
			assert.NotNil(t, got)
			assert.True(t, test.want(t, got))
		})
	}
}

func positiveKeyValueDeleteAction(t *testing.T) (interface{}, error) {
	return actionCheckSQLArgs[*KeyValue, KeyValue](t, KeyValueDelete, KeyValue{}), nil
}

func positiveKeyValueDeleteActionCheck(t *testing.T, i interface{}) bool {
	return checkTrue(i)
}

func positiveKeyValueGetAllAction(t *testing.T) (interface{}, error) {
	return actionCheckSQLArgs[*KeyValue, KeyValue](t, KeyValueGetAll, KeyValue{}), nil
}

func positiveKeyValueGetAllActionCheck(t *testing.T, i interface{}) bool {
	return checkTrue(i)
}

func positiveKeyValueSelectAction(t *testing.T) (interface{}, error) {
	return actionCheckSQLArgs[*KeyValue, KeyValue](t, KeyValueSelect, KeyValue{}), nil
}

func positiveKeyValueSelectActionCheck(t *testing.T, i interface{}) bool {
	return checkTrue(i)
}

func positiveKeyValueUpsertAction(t *testing.T) (interface{}, error) {
	return actionCheckSQLArgs[*KeyValue, KeyValue](t, KeyValueUpsert, KeyValue{}), nil
}

func positiveKeyValueUpsertActionCheck(t *testing.T, i interface{}) bool {
	return checkTrue(i)
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
