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
	//FOR DEBUG
	fmt.Println(header)

	return nil
}

func (dec *decoder) decodeHeader() (*Header, error) {
	magicString, err := dec.traverseNBytes(5)
	if err != nil {
		return nil, fmt.Errorf("traverseNBytes error: %v", err)
	}
	if magicString != "REDIS" {
		return nil, fmt.Errorf("no magic string detected")
	}

	version, err := dec.traverseNBytes(4)
	if err != nil {
		return nil, fmt.Errorf("traverseNBytes error: %v", err)
	}
	versionInt, err := strconv.Atoi(version)
	if err != nil {
		return nil, fmt.Errorf("atoi error: %v", err)
	}

	return &Header{name: magicString, version: versionInt}, nil
}

func (dec *decoder) traverseNBytes(offset int) (string, error) {
	if dec.len <= dec.pos+offset {
		return "", fmt.Errorf("can't traverse by %d bytes because rdb file length is %d", offset, dec.len)
	}

	str := string(dec.b[dec.pos : dec.pos+offset])
	dec.pos += offset
	return str, nil
}
