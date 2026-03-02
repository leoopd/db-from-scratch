package entries

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Entry struct {
	key     []byte
	val     []byte
	deleted bool
}

func (ent *Entry) Encode() []byte {
	// 4 bytes for key len
	// 4 bytes for value len
	// 1 byte to indicate deleted entries
	// key-len bytes for key
	// val-len bytes for value
	data := make([]byte, 4+4+1+len(ent.key)+len(ent.val))
	binary.LittleEndian.PutUint32(data[0:4], uint32(len(ent.key)))
	binary.LittleEndian.PutUint32(data[4:8], uint32(len(ent.val)))
	data[8] = 0
	if ent.deleted {
		data[8] = 1
	}
	copy(data[9:], ent.key)
	copy(data[9+len(ent.key):], ent.val)
	return data
}

func (ent *Entry) Decode(r io.Reader) error {
	// read first 4 bytes to determine the length of our key
	keyLenSlice := make([]byte, 4)
	n, err := r.Read(keyLenSlice)
	if n != 4 || err != nil {
		return fmt.Errorf("failed to read key-length. err: %v", err)
	}
	keyLen := binary.LittleEndian.Uint32(keyLenSlice)

	// read next 4 bytes to determine the length of our value
	valLenSlice := make([]byte, 4)
	n, err = r.Read(valLenSlice)
	if n != 4 || err != nil {
		return fmt.Errorf("failed to read val-length. err: %v", err)
	}
	valLen := binary.LittleEndian.Uint32(valLenSlice)

	// read next byte to determine if our entry was deleted
	delSlice := make([]byte, 1)
	n, err = r.Read(delSlice)
	if n != 1 || err != nil {
		return fmt.Errorf("failed to read key-value. err: %v", err)
	}
	deleted := false
	if delSlice[0] == 1 {
		deleted = true
	}

	// read next full key, using the length determined earlier
	key := make([]byte, keyLen)
	n, err = r.Read(key)
	if n != int(keyLen) || err != nil {
		return fmt.Errorf("failed to read key-value. err: %v", err)
	}

	// read next full value, using the length determined earlier
	val := make([]byte, valLen)
	n, err = r.Read(val)
	if n != int(valLen) || err != nil {
		return fmt.Errorf("failed to read val-value. err: %v", err)
	}

	// persist values after all operations succeeded
	ent.key = key
	ent.val = val
	ent.deleted = deleted

	return nil
}
