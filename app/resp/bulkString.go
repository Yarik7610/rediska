package resp

import (
	"fmt"
)

type bulkString struct {
	value *string
}

func (bs bulkString) Encode() ([]byte, error) {
	if bs.value == nil {
		return []byte(NULL_RESP_2), nil
	}

	l := len(*bs.value)
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", l, *bs.value)), nil
}

func (bulkString) Decode(b []byte) ([]byte, Value, error) {
	l := len(b)
	if l == 0 {
		return nil, nil, fmt.Errorf("bulk string decode error: expected not fully empty string")
	}

	if string(b) == NULL_RESP_2 {
		return b[len(NULL_RESP_2):], bulkString{value: nil}, nil
	}

	if b[0] != '$' {
		return nil, nil, fmt.Errorf("bulk string decode error: didn't find '$' sign")
	}

	expectedLen, b, err := traverseExpectedLen(b[1:])
	if err != nil {
		return nil, nil, fmt.Errorf("bulk string parse expected len error: %v", err)
	}

	b, err = traverseCRLF(b)
	if err != nil {
		return nil, nil, fmt.Errorf("bulk string traverse CRLF error: %v", err)
	}

	err = requireEndingCRLF(b)
	if err != nil {
		return nil, nil, fmt.Errorf("bulk string decode error: %v", err)
	}

	outLen := len(b) - 2
	if expectedLen != outLen {
		return nil, nil, fmt.Errorf("bulk string decode error: expected len (%d) != out len (%d)", expectedLen, outLen)
	}

	res := string(b[:expectedLen])
	return b[expectedLen+2:], bulkString{value: &res}, nil
}
