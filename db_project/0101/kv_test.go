package kv

import (
	"fmt"
	"testing"
)

func TestOpen(t *testing.T) {
	kv := KV{}
	kv.Open()

	if kv.mem == nil {
		t.Errorf("wanted: kv.mem map[string][]byte; got: nil")
	}
}

func TestClose(t *testing.T) {
	kv := KV{}
	kv.Open()

	err := kv.Close()
	if err != nil {
		t.Errorf("wanted: err %s; got: err %v", notFoundError, err)
	}
}

func TestGet(t *testing.T) {
	kv := KV{}
	kv.Open()

	inputKey := []byte("known")
	inputVal := []byte("123")
	kv.mem[string(inputKey)] = inputVal
	val, ok, err := kv.Get(inputKey)

	if string(val) != string(inputVal) {
		t.Errorf("wanted: val = %v; got: val = %v", inputVal, val)
	}

	if !ok {
		t.Errorf("wanted: ok = true; got: ok = false")
	}

	if err != nil {
		t.Errorf("wanted: err = nil; got: err = %v", err)
	}

	inputKey = []byte("unknown")
	val, ok, err = kv.Get(inputKey)

	if val != nil {
		t.Errorf("wanted: val = nil; got val = %v", val)
	}

	if ok {
		t.Errorf("wanted: ok = false; got: ok = true")
	}

	if err == nil {
		fmt.Println(val == nil)
		t.Errorf("wanted: err = not found; got: err = %v", err)
	}
}

func TestSet(t *testing.T) {
	kv := KV{}
	kv.Open()

	inputKey := []byte("known")
	inputVal := []byte("123")
	updated, err := kv.Set(inputKey, inputVal)
	if !updated {
		t.Errorf("wanted: updated = true; got: updated = false")
	}

	if err != nil {
		t.Errorf("wanted: err = nil; got: err %v", err)
	}
}

func TestDel(t *testing.T) {
	kv := KV{}
	kv.Open()

	inputKey := []byte("known")
	inputVal := []byte("123")
	kv.mem[string(inputKey)] = inputVal
	ok, err := kv.Del(inputKey)
	if err != nil {
		t.Errorf("err not nil %v", err)
	}

	if !ok {
		t.Errorf("wanted: ok = false; got: ok = true")
	}

	if err != nil {
		t.Errorf("wanted: err = nil; got: err = %v", err)
	}

	inputKey = []byte("unknown")
	ok, err = kv.Del(inputKey)
	if ok {
		t.Errorf("wanted: ok = false; got: ok = true")
	}

	if err == nil {
		t.Errorf("wanted: err = nil; got: err = %v", err)
	}
}
