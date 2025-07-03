package resp

import (
	"fmt"
	"strconv"
)

type integer struct {
	value int
}

func (i integer) Encode() ([]byte, error) {
	return []byte(fmt.Sprintf(":%d\r\n", i.value)), nil
}

func (integer) Decode(b []byte) (Value, error) {
	l := len(b)
	if l == 0 {
		return nil, fmt.Errorf("integer decode error: expected not fully empty string")
	}

	if b[0] != ':' {
		return nil, fmt.Errorf("integer decode error: didn't find ':' sign")
	}

	payload, err := traversePayloadTillFirstCRLF(b, l)
	if err != nil {
		return nil, fmt.Errorf("integer decode error: %v", err)
	}

	intVal, err := strconv.Atoi(payload)
	if err != nil {
		return nil, fmt.Errorf("integer decode atoi error: %v", err)
	}

	return integer{value: intVal}, nil
}
