package entries

import (
	"bytes"
	"fmt"
	"testing"
)

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
		if tc.deleted {
			if !out.deleted {
				t.Errorf("want: deleted = true; got: deleted = false. err: %v, key: %s", err, string(tc.entry.key))
			}
		}
		if err != nil {
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
