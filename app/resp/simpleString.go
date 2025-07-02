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

	for i := range b {
		if b[i] == '\r' {
			if i+1 < l && b[i+1] == '\n' {
				return string(b[1:i]), nil
			}
			return "", fmt.Errorf("simple string decode error: wrong char: %q after '\\r'", b[i+1])
		}
	}

	return "", fmt.Errorf("simple string decode error: didn't find '\\r\\n' in the end")
}
