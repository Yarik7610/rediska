package resp

import (
	"fmt"
)

type SimpleError struct {
	Value string
}

func (se SimpleError) Encode() ([]byte, error) {
	return []byte(fmt.Sprintf("-%s\r\n", se.Value)), nil
}

func (SimpleError) Decode(b []byte) ([]byte, Value, error) {
	l := len(b)
	if l == 0 {
		return nil, nil, fmt.Errorf("simple error decode error: expected not fully empty string")
	}

	if b[0] != '-' {
		return nil, nil, fmt.Errorf("simple error decode error: didn't find '-' sign")
	}

	b, res, err := traversePayloadTillFirstCRLF(b, l)
	if err != nil {
		return nil, nil, fmt.Errorf("simple error decode error: %v", err)
	}

	return b, SimpleError{Value: res}, nil
}
