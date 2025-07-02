package resp

import (
	"fmt"
	"strconv"
)

type integer struct{}

func (integer) Encode(str string) ([]byte, error) {
	intVal, err := strconv.Atoi(str)
	if err != nil {
		return nil, fmt.Errorf("integer encode atoi error: %v", err)
	}

	return []byte(fmt.Sprintf(":%d\r\n", intVal)), nil
}

func (integer) Decode(b []byte) (string, error) {
	if b[0] != ':' {
		return "", fmt.Errorf("integer decode error: didn't find ':' sign")
	}

	l := len(b)
	err := requireEndingCRLF(b)
	if err != nil {
		return "", fmt.Errorf("integer decode error: %v", err)
	}

	intVal, err := strconv.Atoi(string(b[1 : l-2]))
	if err != nil {
		return "", fmt.Errorf("integer decode atoi error: %v", err)
	}

	return strconv.FormatInt(int64(intVal), 10), nil
}
