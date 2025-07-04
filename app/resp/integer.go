package resp

import (
	"fmt"
	"strconv"
)

type Integer struct {
	Value int
}

func (i Integer) Encode() ([]byte, error) {
	return []byte(fmt.Sprintf(":%d\r\n", i.Value)), nil
}

func (Integer) Decode(b []byte) ([]byte, Value, error) {
	l := len(b)
	if l == 0 {
		return nil, nil, fmt.Errorf("Integer decode error: expected not fully empty string")
	}

	if b[0] != ':' {
		return nil, nil, fmt.Errorf("Integer decode error: didn't find ':' sign")
	}

	b, payload, err := traversePayloadTillFirstCRLF(b, l)
	if err != nil {
		return nil, nil, fmt.Errorf("Integer decode error: %v", err)
	}

	intVal, err := strconv.Atoi(payload)
	if err != nil {
		return nil, nil, fmt.Errorf("Integer decode atoi error: %v", err)
	}

	return b, Integer{Value: intVal}, nil
}
