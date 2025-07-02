package resp

import "fmt"

func requireEndingCRLF(b []byte) error {
	l := len(b)
	if string(b[l-2:]) != "\r\n" {
		return fmt.Errorf("didn't find '\\r\\n' in the end")
	}
	return nil
}
