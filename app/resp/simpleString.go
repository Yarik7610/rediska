package resp

import (
	"fmt"
	"strings"
)

type simpleString struct{}

func (simpleString) Encode(str string) ([]byte, error) {
	l := len(str)
	if strings.Contains(str[:l-2], "\r") || strings.Contains(str[:l-2], "\n") {
		return nil, fmt.Errorf("simple string encode error: can't have '\\r' or '\\n' char in the payload")
	}

	return []byte(fmt.Sprintf("+%s\r\n", str)), nil
}

func (simpleString) Decode(b []byte) (string, error) {
	if b[0] != '+' {
		return "", fmt.Errorf("simple string decode error: didn't find '+' sign")
	}
	b = b[1:]

	for i := range b {
		if b[i] == '\r' {
			if i+1 < len(b) && b[i+1] == '\n' {
				return fmt.Sprintf("%s", b[1:i]), nil
			}
			return "", fmt.Errorf("simple string decode error: wrong char: %q after '\\r'", b[i+1])
		}
	}

	return "", fmt.Errorf("simple string decode error: didn't find '\\r\\n' in the end")
}
