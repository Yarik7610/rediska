package resp

import (
	"fmt"
)

type simpleError struct {
	value string
}

func (se simpleError) Encode() ([]byte, error) {
	return []byte(fmt.Sprintf("-%s\r\n", se.value)), nil
}

func (simpleError) Decode(b []byte) (Value, error) {
	l := len(b)
	if l == 0 {
		return nil, fmt.Errorf("simple error decode error: expected not fully empty string")
	}

	if b[0] != '-' {
		return nil, fmt.Errorf("simple error decode error: didn't find '-' sign")
	}

	res, err := traversePayloadTillFirstCRLF(b, l)
	if err != nil {
		return nil, fmt.Errorf("simple error decode error: %v", err)
	}

	return simpleError{value: res}, nil
}
