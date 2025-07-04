package resp

import (
	"bytes"
	"fmt"
)

type array struct {
	value []Value
}

func (a array) Encode() ([]byte, error) {
	var b bytes.Buffer

	l := len(a.value)
	b.WriteString(fmt.Sprintf("*%d\r\n", l))

	for _, val := range a.value {
		encodedVal, err := val.Encode()
		if err != nil {
			return nil, fmt.Errorf("array encode error: %v", err)
		}
		b.WriteString(fmt.Sprintf("%s\r\n", encodedVal))
	}

	return b.Bytes(), nil
}

func (array) Decode(b []byte) ([]byte, Value, error) {
	l := len(b)
	if l == 0 {
		return nil, nil, fmt.Errorf("array decode error: expected not fully empty string")
	}

	if string(b) == NULL_RESP_2 {
		return b[0:0], array{value: nil}, nil
	}

	if b[0] != '*' {
		return nil, nil, fmt.Errorf("array decode error: didn't find '$' sign")
	}

	res := make([]Value, 0)
	var curVal Value
	var err error

	for {
		switch b[0] {
		case '*':
			b, curVal, err = resp.Array.Decode(b)
		case '$':
			b, curVal, err = resp.BulkString.Decode(b)
		case ':':
			b, curVal, err = resp.Integer.Decode(b)
		case '+':
			b, curVal, err = resp.SimpleString.Decode(b)
		case '-':
			b, curVal, err = resp.SimpleError.Decode(b)
		default:
			return nil, nil, fmt.Errorf("array decode error: detected unknown RESP type")
		}

		if err != nil {
			return nil, nil, fmt.Errorf("array decode error: %v", err)
		}
		res = append(res, curVal)

		if len(b) == 0 {
			break
		}
	}

	return b[0:0], array{value: res}, nil
}
