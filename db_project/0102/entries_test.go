package entries

import (
	"bytes"
	"fmt"
	"testing"
)

func TestDecode(t *testing.T) {
	testCases := []struct {
		entry Entry
	}{
		{
			entry: Entry{
				key: []byte("thisIsSomeKey"),
				val: []byte("thisIsSomeValue"),
			},
		},
	}

	for _, tc := range testCases {
		input := tc.entry.Encode()
		b := bytes.NewBuffer(input)

		out := Entry{}
		err := out.Decode(b)
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
	}
}
