package rdb

import (
	"fmt"
	"strconv"
)

type decoder struct {
	b   []byte
	pos int
	len int
}

func Decode(b []byte) error {
	dec := decoder{b: b, pos: 0, len: len(b)}

	header, err := dec.decodeHeader()
	if err != nil {
		return fmt.Errorf("decode header error: %v", err)
	}
	fmt.Println(header)

	metadata, err := dec.decodeMetadata()
	if err != nil {
		return fmt.Errorf("decode metadata error: %v", err)
	}
	fmt.Println(metadata)

	return nil
}

func (dec *decoder) decodeHeader() (*Header, error) {
	magicString, err := dec.traverseStringLen(5)
	if err != nil {
		return nil, err
	}
	if magicString != "REDIS" {
		return nil, fmt.Errorf("no magic string detected")
	}

	version, err := dec.traverseStringLen(4)
	if err != nil {
		return nil, err
	}
	versionAtoi, err := strconv.Atoi(version)
	if err != nil {
		return nil, fmt.Errorf("atoi error: %v", err)
	}

	return &Header{name: magicString, version: versionAtoi}, nil
}

func (dec *decoder) decodeMetadata() (*Metadata, error) {
	data := make(map[string]string)

	for {
		opCode, err := dec.traverseUInt8()
		if err != nil {
			return nil, err
		}
		if opCode == OP_EOF || opCode != OP_AUX {
			break
		}

		key, err := dec.decodeString()
		if err != nil {
			return nil, fmt.Errorf("decodeString error: %v", err)
		}
		value, err := dec.decodeString()
		if err != nil {
			return nil, fmt.Errorf("decodeString error: %v", err)
		}
		data[key] = value
	}

	return &Metadata{data: data}, nil
}

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
		val, err := dec.traverseIntegerString(remainingBits)
		return val, true, err
	}
	return 0, false, fmt.Errorf("invalid length decode format")
}

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
