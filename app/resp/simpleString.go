package resp

import (
	"fmt"
	"log"
	"strings"
)

type SimpleString struct {
	Value string
}

func (ss SimpleString) Encode() ([]byte, error) {
	if strings.Contains(ss.Value, "\r") || strings.Contains(ss.Value, "\n") {
		return nil, fmt.Errorf("simple string encode error: can't have '\\r' or '\\n' char in the payload")
	}

	return fmt.Appendf(nil, "+%s\r\n", ss.Value), nil
}

func (SimpleString) Decode(b []byte) ([]byte, Value, error) {
	l := len(b)
	if l == 0 || b == nil {
		return nil, nil, fmt.Errorf("simple string decode error: expected non-empty data")
	}

	if b[0] != '+' {
		return nil, nil, fmt.Errorf("simple string decode error: didn't find '+' sign")
	}

	b, res, err := traversePayloadTillFirstCRLF(b, l)
	if err != nil {
		return nil, nil, fmt.Errorf("simple string decode error: %v", err)
	}

	return b, SimpleString{Value: res}, nil
}

func AssertEqualSimpleString(value Value, raw string) {
	v, ok := value.(SimpleString)
	if !ok {
		log.Fatalf("assertion failed: expected SimpleString, got: %T", value)
	}
	if v.Value != raw {
		log.Fatalf("assertion failed: expected %s, got: %v", raw, v.Value)
	}
}
