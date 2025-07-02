package resp

import (
	"fmt"
	"strconv"
)

type bulkString struct{}

func (bulkString) Encode(str string) []byte {
	l := len(str)
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", l, str))
}

func (bulkString) Decode(b []byte) (string, error) {
	if b[0] != '$' {
		return "", fmt.Errorf("bulk string decode error: didn't find '$' sign")
	}
	b = b[1:]

	bulkStringLen := make([]byte, 0)
	for i := range b {
		if b[i] == '\r' && i+1 < len(b) && b[i+1] == '\n' {
			break
		}
		bulkStringLen = append(bulkStringLen, b[i])
	}
	b = b[len(bulkStringLen)+2:]

	bulkStringLenInt, err := strconv.Atoi(string(bulkStringLen))
	if err != nil {
		return "", fmt.Errorf("bulk string len decode atoi error: %v", err)
	}

	err = requireEndingCRLF(b)
	if err != nil {
		return "", fmt.Errorf("bulk string decode error: %v", err)
	}

	return string(b[:bulkStringLenInt]), nil

}
