package resp

import (
	"fmt"
	"strconv"
)

type integer struct{}

func (integer) Encode(num int) []byte {
  return []byte(fmt.Sprintf(":%d\r\n", num))
}

func (integer) Decode(b []byte) (int, error) {
  l := len(b)
  if l == 0 {
    return 0, fmt.Errorf("integer decode error: expected not fully empty string")
  }

  if b[0] != ':' {
    return 0, fmt.Errorf("integer decode error: didn't find ':' sign")
  }

  payload, err := traversePayloadTillFirstCRLF(b, l)
  if err != nil {
    return 0, fmt.Errorf("integer decode error: %v", err)
  }

  intVal, err := strconv.Atoi(payload)
  if err != nil {
    return 0, fmt.Errorf("integer decode atoi error: %v", err)
  }

  return intVal, nil
}