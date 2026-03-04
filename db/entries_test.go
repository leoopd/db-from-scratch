package db

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestEncode(t *testing.T) {
	// create byte slice res: https://go.dev/play/p/nRZvqliUHtM
	testCases := []struct {
		entry   Entry
		deleted bool
		res     []byte
	}{
		{
			entry: Entry{
				key:     []byte("thisIsSomeKey"),
				val:     []byte("thisIsSomeValue"),
				deleted: false,
			},
			deleted: false,
			res: []byte{245, 148, 245, 50, 13, 0, 0, 0, 15, 0, 0, 0, 0, 116, 104, 105, 115, 73, 115,
				83, 111, 109, 101, 75, 101, 121, 116, 104, 105, 115, 73, 115, 83, 111, 109, 101, 86,
				97, 108, 117, 101},
		},
		{
			entry: Entry{
				key:     []byte("thisIsSomeKey"),
				val:     []byte("thisIsSomeValue"),
				deleted: true,
			},
			res: []byte{209, 77, 242, 81, 13, 0, 0, 0, 15, 0, 0, 0, 1, 116, 104, 105, 115, 73, 115,
				83, 111, 109, 101, 75, 101, 121, 116, 104, 105, 115, 73, 115, 83, 111, 109, 101, 86,
				97, 108, 117, 101},
			deleted: true,
		},
	}

	for _, tc := range testCases {
		out := tc.entry.Encode()
		if !bytes.Equal(tc.res, out) {
			t.Errorf("out doesn't match tc.res")
		}
	}
}

func TestDecode(t *testing.T) {
	testCases := []struct {
		entry   Entry
		deleted bool
	}{
		{
			entry: Entry{
				key:     []byte("thisIsSomeKey"),
				val:     []byte("thisIsSomeValue"),
				deleted: false,
			},
			deleted: false,
		},
		{
			entry: Entry{
				key:     []byte("thisIsSomeKey"),
				val:     []byte("thisIsSomeValue"),
				deleted: true,
			},
			deleted: true,
		},
	}

	for _, tc := range testCases {
		input := tc.entry.Encode()
		b := bytes.NewBuffer(input)

		out := Entry{}
		err := out.Decode(b)

		if tc.deleted && !out.deleted {
			t.Errorf("want: deleted = true; got: deleted = false. err: %v, key: %s", err, string(tc.entry.key))
		}

		if err != io.EOF {
			t.Errorf("encountered err: %v", err)
		}

		fmt.Println("Key:", string(out.key))
		if string(out.key) != string(tc.entry.key) {
			t.Errorf("keys don't match. want: %s, got: %s", string(tc.entry.key), string(out.key))
		}

		fmt.Println("Value:", string(out.val))
		if string(out.val) != string(tc.entry.val) {
			t.Errorf("values don't match. want: %s, got: %s", string(tc.entry.val), string(out.val))
		}

		fmt.Println("Deleted:", out.deleted)
		if string(out.val) != string(tc.entry.val) {
			t.Errorf("values don't match. want: %s, got: %s", string(tc.entry.val), string(out.val))
		}
	}
}
