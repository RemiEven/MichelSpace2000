package genericsync

import (
	"strings"
	"testing"
)

func TestMap(t *testing.T) {
	m := Map[string, string]{}

	if _, ok := m.Load("key"); ok {
		t.Errorf("map that should be empty said it contained a key")
		return
	}

	m.Store("key1", "value1")
	m.Store("key2", "value2")
	if value, ok := m.Load("key1"); !ok {
		t.Errorf("stored value could not be retrieved")
		return
	} else if value != "value1" {
		t.Errorf("unexpected retrieved value: wanted [%s], got [%s]", "value1", value)
	}

	acc := ""
	m.Range(func(key, value string) {
		acc += key + ":" + value
	})
	if !strings.Contains(acc, "key1:value1") {
		t.Errorf("expected Range to add key1:value1 to accumulator")
		return
	} else if !strings.Contains(acc, "key2:value2") {
		t.Errorf("expected Range to add key2:value2 to accumulator")
		return
	}
}
