package resp

import (
	"fmt"
)

type BulkString struct {
	Value *string
}

func (bs BulkString) Encode() ([]byte, error) {
	if bs.Value == nil {
		return []byte(NULL_BULK_STRING_RESP_2), nil
	}

	l := len(*bs.Value)
	return fmt.Appendf(nil, "$%d\r\n%s\r\n", l, *bs.Value), nil
}

func (BulkString) Decode(b []byte) ([]byte, Value, error) {
	l := len(b)
	if l == 0 || b == nil {
		return nil, nil, fmt.Errorf("bulk string decode error: expected non-empty data")
	}

	if string(b) == NULL_BULK_STRING_RESP_2 {
		return b[len(NULL_BULK_STRING_RESP_2):], BulkString{Value: nil}, nil
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

	if len(b) < expectedLen+2 {
		return nil, nil, fmt.Errorf("bulk string decode error: not enough bytes for content (%d < %d)", len(b), expectedLen+2)
	}

	if err = requireEndingCRLF(b[expectedLen : expectedLen+2]); err != nil {
		return nil, nil, fmt.Errorf("bulk string decode error: %v", err)
	}

	res := string(b[:expectedLen])
	return b[expectedLen+2:], BulkString{Value: &res}, nil
}
