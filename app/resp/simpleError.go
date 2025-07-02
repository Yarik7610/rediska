package resp

import (
	"fmt"
)

type simpleError struct{}

func (simpleError) Encode(msg string) []byte {
	return []byte(fmt.Sprintf("-%s\r\n", msg))
}

func (simpleError) Decode(b []byte) (string, error) {
	if b[0] != '-' {
		return "", fmt.Errorf("simple error decode error: didn't find '-' sign")
	}

	l := len(b)
	err := requireEndingCRLF(b)
	if err != nil {
		return "", fmt.Errorf("simple error decode error: %v", err)
	}

	return string(b[1 : l-2]), nil
}
