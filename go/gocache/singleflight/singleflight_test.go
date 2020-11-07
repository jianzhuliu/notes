package singleflight

import (
	"testing"
)

func TestEntry(t *testing.T) {
	var e Entry
	value, err := e.Do("key", func() (interface{}, error) {
		t.Logf("do search")
		return "key", nil
	})

	//t.Logf("%v, %v", value, err)
	if err != nil || string(value.(string)) != "key" {
		t.Fatalf("expect key,but %v, err=%v got", value, err)
	}
}
