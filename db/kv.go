package db

import (
	"fmt"
	"io"
)

const notFoundError = "not found"
const deletionFailedError = "operation failed: delete"
const updateFailedError = "operation failed: update"

type KV struct {
	log Log
	mem map[string][]byte
}

func (kv *KV) Open() error {
	err := kv.log.Open()
	if err != nil {
		return err
	}
	kv.mem = map[string][]byte{}

	eof := false
	for !eof {
		entry := Entry{}
		eof, err = kv.log.Read(&entry)
		// only return if the error needs to be handled explicitly
		if err != nil && err != io.EOF {
			return err
		}
		if entry.deleted {
			deleted, err := kv.Del(entry.key)
			if err != nil {
				return err
			}
			if !deleted {
				return fmt.Errorf(deletionFailedError)
			}
		} else {
			updated, err := kv.Set(entry.key, entry.val)
			if err != nil {
				return err
			}
			if !updated {
				return fmt.Errorf(updateFailedError)
			}
		}

	}
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
	kv.log.Write(&Entry{key: key, val: val})
	return true, nil
}

func (kv *KV) Del(key []byte) (deleted bool, err error) {
	if _, ok := kv.mem[string(key)]; !ok {
		return false, fmt.Errorf(notFoundError)
	}
	delete(kv.mem, string(key))
	kv.log.Write(&Entry{key: key})
	return true, nil
}
