package resp

import (
	"fmt"
	"strconv"
)

func traversePayloadTillFirstCRLF(b []byte, l int) (string, error) {
	for i := range b {
		if b[i] == '\r' {
			if i+1 < l && b[i+1] == '\n' {
				return string(b[1:i]), nil
			}
			return "", fmt.Errorf("wrong char: %q after '\\r'", b[i+1])
		}
	}

	return "", fmt.Errorf("didn't find '\\r\\n' in the end")
}

func traverseExpectedLen(b []byte) (int, []byte, error) {
	expectedLen := make([]byte, 0)

	for i := range b {
		if b[i] == '\r' && i+1 < len(b) && b[i+1] == '\n' {
			break
		}
		expectedLen = append(expectedLen, b[i])
	}

	expectedLenInt, err := strconv.Atoi(string(expectedLen))
	if err != nil {
		return 0, nil, fmt.Errorf("len atoi error: %v", err)
	}

	return expectedLenInt, b[len(expectedLen):], nil
}

func traverseCRLF(b []byte) ([]byte, error) {
	if len(b) >= 2 && b[0] == '\r' && b[1] == '\n' {
		return b[2:], nil
	}
	return nil, fmt.Errorf("didn't detect '\\r\\n'")
}

func requireEndingCRLF(b []byte) error {
	l := len(b)
	if string(b[l-2:]) != "\r\n" {
		return fmt.Errorf("didn't find '\\r\\n' in the end")
	}
	return nil
}
