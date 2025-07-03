package resp

import (
	"fmt"
	"strings"
)

type simpleString struct{}

func (simpleString) Encode(str string) ([]byte, error) {
	if strings.Contains(str, "\r") || strings.Contains(str, "\n") {
		return nil, fmt.Errorf("simple string encode error: can't have '\\r' or '\\n' char in the payload")
	}

	return []byte(fmt.Sprintf("+%s\r\n", str)), nil
}

func (simpleString) Decode(b []byte) (string, error) {
	l := len(b)
	if l == 0 {
		return "", fmt.Errorf("simple string decode error: expected not fully empty string")
	}

	if b[0] != '+' {
		return "", fmt.Errorf("simple string decode error: didn't find '+' sign")
	}

	res, err := traversePayloadTillFirstCRLF(b, l)
	if err != nil {
		return "", fmt.Errorf("simple string decode error: %v", err)
	}

	return res, nil
}
