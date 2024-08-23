package entity

import (
	"github.com/victor-skurikhin/etcd-client/v1/internal/domain"
	"regexp"
	"testing"
)

func actionCheckSQLArgs[T domain.Ptr[U], U domain.Entity](t *testing.T, action domain.Actioner[T, U], entity U) bool {
	re1, _ := regexp.Compile(`(\$\d+)\b`)
	rs1 := re1.FindAllStringSubmatch(action.SQL(), -1)

	if len(rs1) < 1 {
		return 0 == len(action.Args(entity))
	}
	set := make(map[string]struct{})
	for _, r := range rs1 {
		if len(r) < 1 {
			t.Fatalf("bad len(r)")
		}
		set[r[0]] = struct{}{}
	}
	return len(set) == len(action.Args(entity))
}

func checkTrue(i interface{}) bool {
	if u, ok := i.(bool); ok {
		return u == true
	}
	return false
}
