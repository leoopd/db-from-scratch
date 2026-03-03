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

	testCases := []struct {
		key         []byte
		val         []byte
		kv          KV
		errExpected bool
	}{
		{
			key: []byte("known"),
			val: []byte("123"),
			kv: func(kv KV) KV {
				kv.mem[string("known")] = []byte("123")
				return kv
			}(kv),
			errExpected: false,
		},
		{
			key: []byte("unknown"),
			kv: func(kv KV) KV {
				return kv
			}(kv),
			errExpected: true,
		},
	}

	for _, tc := range testCases {
		val, ok, err := kv.Get(tc.key)
		if tc.errExpected {
			if err == nil {
				fmt.Println(val == nil)
				t.Errorf("[expected to err] err = %v, key = %s", err, string(tc.key))
			}
			if ok {
				t.Errorf("[expected to err] ok = true, key = %s", string(tc.key))
			}
		} else {
			if string(val) != string(tc.val) {
				t.Errorf("wanted: val = %v; got: val = %v", tc.val, val)
			}
			if !ok {
				t.Errorf("wanted: ok = true; got: ok = false")
			}
			if err != nil {
				t.Errorf("unexpected err = %v", err)
			}
		}
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

	testCases := []struct {
		key         []byte
		kv          KV
		errExpected bool
	}{
		{
			key: []byte("known"),
			kv: func(kv KV) KV {
				kv.mem[string("known")] = []byte("123")
				return kv
			}(kv),
			errExpected: false,
		},
		{
			kv: func(kv KV) KV {
				return kv
			}(kv),
			errExpected: true,
		},
	}

	for _, tc := range testCases {
		deleted, err := tc.kv.Del(tc.key)
		if tc.errExpected {
			if deleted {
				t.Errorf("[expected to err] deleted = true. key: %s", string(tc.key))
			}
			if err == nil {
				t.Errorf("[expected to err] err = nil. key: %s", string(tc.key))
			}
		} else {
			if err != nil {
				t.Errorf("err not nil %v. key: %s", err, string(tc.key))
			}

			if !deleted {
				t.Errorf("wanted: deleted = false; got: deleted = true. key: %s", string(tc.key))
			}
		}
	}
}
