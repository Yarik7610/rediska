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
