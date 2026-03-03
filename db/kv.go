package db

import "fmt"

const notFoundError = "not found"

type KV struct {
	mem map[string][]byte
}

func (kv *KV) Open() error {
	kv.mem = map[string][]byte{}
	return nil
}

func (kv *KV) Close() error { return nil }

func (kv *KV) Get(key []byte) (val []byte, ok bool, err error) {
	val, ok = kv.mem[string(key)]
	if !ok {
		return nil, false, fmt.Errorf(notFoundError)
	}
	return val, ok, nil
}

func (kv *KV) Set(key []byte, val []byte) (updated bool, err error) {
	kv.mem[string(key)] = val
	return true, nil
}

func (kv *KV) Del(key []byte) (deleted bool, err error) {
	if _, ok := kv.mem[string(key)]; !ok {
		return false, fmt.Errorf(notFoundError)
	}
	delete(kv.mem, string(key))
	return true, nil
}
