package resp

import "fmt"

func parseTillFirstCRLF(b []byte, l int) (string, error) {
	for i := range b {
		if b[i] == '\r' {
			if i+1 < l && b[i+1] == '\n' {
				return string(b[1:i]), nil
			}
			return "", fmt.Errorf("wrong char: %q after '\\r'", b[i+1])
		}
	}

	return "", fmt.Errorf("didn't find '\\r\\n' in the end")
}

func requireEndingCRLF(b []byte) error {
	l := len(b)
	if string(b[l-2:]) != "\r\n" {
		return fmt.Errorf("didn't find '\\r\\n' in the end")
	}
	return nil
}
