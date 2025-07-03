package resp

import (
	"fmt"
)

type simpleError struct{}

func (simpleError) Encode(msg string) []byte {
	return []byte(fmt.Sprintf("-%s\r\n", msg))
}

func (simpleError) Decode(b []byte) (string, error) {
	l := len(b)
	if l == 0 {
		return "", fmt.Errorf("simple error decode error: expected not fully empty string")
	}

	if b[0] != '-' {
		return "", fmt.Errorf("simple error decode error: didn't find '-' sign")
	}

	res, err := traversePayloadTillFirstCRLF(b, l)
	if err != nil {
		return "", fmt.Errorf("simple error decode error: %v", err)
	}

	return res, nil
}
