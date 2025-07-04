package resp

import (
	"fmt"
	"strings"
)

type SimpleString struct {
	Value string
}

func (ss SimpleString) Encode() ([]byte, error) {
	if strings.Contains(ss.Value, "\r") || strings.Contains(ss.Value, "\n") {
		return nil, fmt.Errorf("simple string encode error: can't have '\\r' or '\\n' char in the payload")
	}

	return []byte(fmt.Sprintf("+%s\r\n", ss.Value)), nil
}

func (SimpleString) Decode(b []byte) ([]byte, Value, error) {
	l := len(b)
	if l == 0 {
		return nil, nil, fmt.Errorf("simple string decode error: expected not fully empty string")
	}

	if b[0] != '+' {
		return nil, nil, fmt.Errorf("simple string decode error: didn't find '+' sign")
	}

	b, res, err := traversePayloadTillFirstCRLF(b, l)
	if err != nil {
		return nil, nil, fmt.Errorf("simple string decode error: %v", err)
	}

	return b, SimpleString{Value: res}, nil
}
