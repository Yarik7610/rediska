package rdb

import (
	"fmt"
)

func (dec *decoder) traverseUInt32() (uint32, error) {
	var num uint32
	for range 4 {
		b, err := dec.traverseUInt8()
		if err != nil {
			return 0, fmt.Errorf("traverseUInt32 error: %v", err)
		}
		num = (num << 8) | uint32(b)
	}
	return num, nil
}

func (dec *decoder) traverseUInt16() (uint16, error) {
	var num uint16
	for range 2 {
		b, err := dec.traverseUInt8()
		if err != nil {
			return 0, fmt.Errorf("traverseUInt16 error: %v", err)
		}
		num = (num << 8) | uint16(b)
	}
	return num, nil
}

func (dec *decoder) traverseUInt8() (uint8, error) {
	if dec.len <= dec.pos+1 {
		return 0, fmt.Errorf("traverseUInt8: can't traverse by 1 byte because rdb file length is %d", dec.len)
	}

	num := uint8(dec.b[dec.pos])
	dec.pos += 1
	return num, nil
}

func (dec *decoder) traverseStringLen(offset int) (string, error) {
	if dec.len <= dec.pos+offset {
		return "", fmt.Errorf("traverseStringLen: can't traverse by %d bytes because rdb file length is %d", offset, dec.len)
	}

	str := string(dec.b[dec.pos : dec.pos+offset])
	dec.pos += offset
	return str, nil
}
