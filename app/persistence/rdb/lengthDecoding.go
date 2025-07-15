package rdb

import "fmt"

func (dec *decoder) decodeLength() (int, bool, error) {
	lenByte, err := dec.traverseUInt8()
	if err != nil {
		return 0, false, err
	}

	switch lenByte >> 6 {
	case 0:
		return int(lenByte & 0x3F), false, nil
	case 1:
		nextLenByte, err := dec.traverseUInt8()
		if err != nil {
			return 0, false, err
		}
		return int(lenByte&0x3F)<<8 | int(nextLenByte), false, nil
	case 2:
		len, err := dec.traverseUInt32()
		if err != nil {
			return 0, false, err
		}
		return int(len), false, nil
	case 3:
		remainingBits := lenByte & 0x3F
		val, err := dec.traverseSpecialString(remainingBits)
		return val, true, err
	}
	return 0, false, fmt.Errorf("invalid length decode format")
}
