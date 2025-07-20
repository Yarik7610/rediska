package resp

import (
	"fmt"
	"strconv"
)

type Integer struct {
	Value int
}

func (i Integer) Encode() ([]byte, error) {
	return fmt.Appendf(nil, ":%d\r\n", i.Value), nil
}

func (Integer) Decode(b []byte) ([]byte, Value, error) {
	l := len(b)
	if l == 0 || b == nil {
		return nil, nil, fmt.Errorf("integer decode error: expected non-empty data")
	}

	if b[0] != ':' {
		return nil, nil, fmt.Errorf("integer decode error: didn't find ':' sign")
	}

	b, payload, err := traversePayloadTillFirstCRLF(b, l)
	if err != nil {
		return nil, nil, fmt.Errorf("integer decode error: %v", err)
	}

	intVal, err := strconv.Atoi(payload)
	if err != nil {
		return nil, nil, fmt.Errorf("integer decode atoi error: %v", err)
	}

	return b, Integer{Value: intVal}, nil
}
