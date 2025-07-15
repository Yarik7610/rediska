package rdb

import (
	"encoding/binary"
	"fmt"
)

func (dec *decoder) traverseUInt64() (uint64, error) {
	num, err := dec.traverseUintXBytes(8)
	return num, err
}

func (dec *decoder) traverseUInt32() (uint32, error) {
	num, err := dec.traverseUintXBytes(4)
	return uint32(num), err
}

func (dec *decoder) traverseUInt16() (uint16, error) {
	num, err := dec.traverseUintXBytes(2)
	return uint16(num), err
}

func (dec *decoder) traverseUInt8() (uint8, error) {
	num, err := dec.traverseUintXBytes(1)
	return uint8(num), err
}

func (dec *decoder) traverseUintXBytes(byteCount int) (uint64, error) {
	if dec.len <= dec.pos+byteCount {
		return 0, fmt.Errorf("traverseUInt%d: not enough bytes, length is %d, pos is %d", byteCount*8, dec.len, dec.pos)
	}

	bytes := dec.b[dec.pos : dec.pos+byteCount]
	dec.pos += byteCount

	switch byteCount {
	case 1:
		return uint64(bytes[0]), nil
	case 2:
		return uint64(binary.LittleEndian.Uint16(bytes)), nil
	case 4:
		return uint64(binary.LittleEndian.Uint32(bytes)), nil
	case 8:
		return binary.LittleEndian.Uint64(bytes), nil
	default:
		return 0, fmt.Errorf("traverseUintXBytes: unsupported byte count %d", byteCount)
	}
}

func (dec *decoder) traverseStringLen(offset int) (string, error) {
	if dec.len <= dec.pos+offset {
		return "", fmt.Errorf("traverseStringLen: can't traverse by %d bytes because rdb file length is %d", offset, dec.len)
	}

	str := string(dec.b[dec.pos : dec.pos+offset])
	dec.pos += offset
	return str, nil
}
