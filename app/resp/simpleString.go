package resp

import (
	"fmt"
	"strings"
)

type simpleString struct {
	value string
}

func (ss simpleString) Encode() ([]byte, error) {
	if strings.Contains(ss.value, "\r") || strings.Contains(ss.value, "\n") {
		return nil, fmt.Errorf("simple string encode error: can't have '\\r' or '\\n' char in the payload")
	}

	return []byte(fmt.Sprintf("+%s\r\n", ss.value)), nil
}

func (simpleString) Decode(b []byte) ([]byte, Value, error) {
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

	return b, simpleString{value: res}, nil
}
