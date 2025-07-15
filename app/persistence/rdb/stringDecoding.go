package rdb

import (
	"fmt"
	"strconv"
)

func (dec *decoder) decodeString() (string, error) {
	len, isLengthAnIntegerString, err := dec.decodeLength()
	if err != nil {
		return "", fmt.Errorf("decodeLength error: %v", err)
	}

	var str string
	if !isLengthAnIntegerString {
		str, err = dec.traverseStringLen(len)
		if err != nil {
			return "", nil
		}
	} else {
		str = strconv.Itoa(len)
	}

	return str, nil
}

func (dec *decoder) traverseSpecialString(remainingBits uint8) (int, error) {
	switch remainingBits {
	case 0:
		value, err := dec.traverseUInt8()
		if err != nil {
			return 0, fmt.Errorf("failed to read 8-bit integer as string: %v", err)
		}
		return int(value), nil
	case 1:
		value, err := dec.traverseUInt16()
		if err != nil {
			return 0, fmt.Errorf("failed to read 16-bit integer as string: %v", err)
		}
		return int(value), nil
	case 2:
		value, err := dec.traverseUInt32()
		if err != nil {
			return 0, fmt.Errorf("failed to read 16-bit integer as string: %v", err)
		}
		return int(value), nil
	case 3:
		return 0, fmt.Errorf("unsupported compressed string format")
	}
	return 0, fmt.Errorf("unsupported integer string format")
}
