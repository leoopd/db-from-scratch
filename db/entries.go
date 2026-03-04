package db

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
)

var ErrBadSum = errors.New("bad checksum")

const (
	standardLen = 4
	singleLen   = 1
)

type Entry struct {
	key     []byte
	val     []byte
	deleted bool
}

func (ent *Entry) Encode() []byte {
	// 4 bytes for checksum
	// 4 bytes for key len
	// 4 bytes for value len
	// 1 byte to indicate deleted entries
	// key-len bytes for key
	// val-len bytes for value
	data := make([]byte, 4+4+4+1+len(ent.key)+len(ent.val))
	binary.LittleEndian.PutUint32(data[4:8], uint32(len(ent.key)))
	binary.LittleEndian.PutUint32(data[8:12], uint32(len(ent.val)))
	data[12] = 0
	if ent.deleted {
		data[12] = 1
	}
	copy(data[13:], ent.key)
	copy(data[13+len(ent.key):], ent.val)

	checksum := crc32.ChecksumIEEE(data[4:])
	binary.LittleEndian.PutUint32(data[:4], checksum)
	return data
}

func (ent *Entry) Decode(r io.Reader) error {
	// read first 4 bytes to determine the checksum of our data
	checksum := make([]byte, standardLen)
	_, err := io.ReadFull(r, checksum)
	if err != nil {
		return err
	}

	// read next 4 bytes to determine the length of our key
	keyLenSlice := make([]byte, standardLen)
	_, err = io.ReadFull(r, keyLenSlice)
	if err != nil {
		return err
	}
	keyLen := binary.LittleEndian.Uint32(keyLenSlice)

	// read next 4 bytes to determine the length of our value
	valLenSlice := make([]byte, standardLen)
	_, err = io.ReadFull(r, valLenSlice)
	if err != nil {
		return err
	}
	valLen := binary.LittleEndian.Uint32(valLenSlice)

	// read next byte to determine if our entry was deleted
	delSlice := make([]byte, singleLen)
	_, err = io.ReadFull(r, delSlice)
	if err != nil {
		return err
	}
	deleted := false
	if delSlice[0] == 1 {
		deleted = true
	}

	// read next full key, using the length determined earlier
	key := make([]byte, keyLen)
	_, err = io.ReadFull(r, key)
	if err != nil {
		return err
	}

	// read next full value, using the length determined earlier
	val := make([]byte, valLen)
	_, err = io.ReadFull(r, val)
	if err != nil {
		return err
	}

	record := bytes.Join([][]byte{keyLenSlice, valLenSlice, delSlice, key, val}, nil)
	if csCalculated := crc32.ChecksumIEEE(record); csCalculated != binary.LittleEndian.Uint32(checksum) {
		return ErrBadSum
	}

	// persist values after all operations succeeded
	ent.key = key
	ent.val = val
	ent.deleted = deleted

	return io.EOF
}
