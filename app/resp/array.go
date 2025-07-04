package resp

import (
	"bytes"
	"fmt"
)

type Array struct {
	Value []Value
}

func (a Array) Encode() ([]byte, error) {
	if a.Value == nil {
		return []byte(NULL_ARRAY_RESP_2), nil
	}

	var b bytes.Buffer

	l := len(a.Value)
	b.WriteString(fmt.Sprintf("*%d\r\n", l))

	for _, val := range a.Value {
		encodedVal, err := val.Encode()
		if err != nil {
			return nil, fmt.Errorf("array encode error: %v", err)
		}
		b.WriteString(string(encodedVal))
	}

	return b.Bytes(), nil
}

func (Array) Decode(b []byte) ([]byte, Value, error) {
	l := len(b)
	if l == 0 {
		return nil, nil, fmt.Errorf("array decode error: expected not fully empty string")
	}

	if string(b) == NULL_ARRAY_RESP_2 {
		return b[len(NULL_ARRAY_RESP_2):], Array{Value: nil}, nil
	}

	if b[0] != '*' {
		return nil, nil, fmt.Errorf("array decode error: didn't find '*' sign")
	}

	resLen, b, err := traverseExpectedLen(b[1:])
	if err != nil {
		return nil, nil, fmt.Errorf("array decode error: %v", err)
	}

	b, err = traverseCRLF(b)
	if err != nil {
		return nil, nil, fmt.Errorf("Array traverse CRLF error: %v", err)
	}

	res := make([]Value, 0, resLen)
	var curVal Value

	for range resLen {
		switch b[0] {
		case '*':
			b, curVal, err = Array{}.Decode(b)
		case '$':
			b, curVal, err = BulkString{}.Decode(b)
		case ':':
			b, curVal, err = Integer{}.Decode(b)
		case '+':
			b, curVal, err = SimpleString{}.Decode(b)
		case '-':
			b, curVal, err = SimpleError{}.Decode(b)
		default:
			return nil, nil, fmt.Errorf("Array decode error: detected unknown RESP type")
		}

		if err != nil {
			return nil, nil, fmt.Errorf("Array decode error: %v", err)
		}
		res = append(res, curVal)
	}

	return b, Array{Value: res}, nil
}
