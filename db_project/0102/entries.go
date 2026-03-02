package entries

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Entry struct {
	key []byte
	val []byte
}

func (ent *Entry) Encode() []byte {
	data := make([]byte, 4+4+len(ent.key)+len(ent.val))
	binary.LittleEndian.PutUint32(data[0:4], uint32(len(ent.key)))
	binary.LittleEndian.PutUint32(data[4:8], uint32(len(ent.val)))
	copy(data[8:], ent.key)
	copy(data[8+len(ent.key):], ent.val)
	return data
}

func (ent *Entry) Decode(r io.Reader) error {
	keyLenSlice := make([]byte, 4)
	n, err := r.Read(keyLenSlice)
	if n != 4 || err != nil {
		return fmt.Errorf("failed to read key-length. err: %v", err)
	}
	keyLen := binary.LittleEndian.Uint32(keyLenSlice)

	valLenSlice := make([]byte, 4)
	n, err = r.Read(valLenSlice)
	if n != 4 || err != nil {
		return fmt.Errorf("failed to read val-length. err: %v", err)
	}
	valLen := binary.LittleEndian.Uint32(valLenSlice)

	key := make([]byte, keyLen)
	n, err = r.Read(key)
	if n != int(keyLen) || err != nil {
		return fmt.Errorf("failed to read key-value. err: %v", err)
	}

	val := make([]byte, valLen)
	n, err = r.Read(val)
	if n != int(valLen) || err != nil {
		return fmt.Errorf("failed to read val-value. err: %v", err)
	}

	ent.key = key
	ent.val = val

	return nil
}
