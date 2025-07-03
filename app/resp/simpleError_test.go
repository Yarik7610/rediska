package resp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleErrorEncode(t *testing.T) {
	tests := []struct {
		Name        string
		In          Value
		Expected    []byte
		ShouldError bool
	}{
		{
			Name:        "Simple error",
			In:          simpleError{value: "ERR invalid command"},
			Expected:    []byte("-ERR invalid command\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Empty error",
			In:          simpleError{value: ""},
			Expected:    []byte("-\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Error with spaces",
			In:          simpleError{value: "ERR syntax error"},
			Expected:    []byte("-ERR syntax error\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Error with special characters",
			In:          simpleError{value: "ERR!@#$%"},
			Expected:    []byte("-ERR!@#$%\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Single character error",
			In:          simpleError{value: "E"},
			Expected:    []byte("-E\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Error with numbers",
			In:          simpleError{value: "ERR 404"},
			Expected:    []byte("-ERR 404\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Error with tab",
			In:          simpleError{value: "ERR\tinvalid"},
			Expected:    []byte("-ERR\tinvalid\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Long error",
			In:          simpleError{value: "ERR This is a very long error message to test encoding with more than a few characters"},
			Expected:    []byte("-ERR This is a very long error message to test encoding with more than a few characters\r\n"),
			ShouldError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			out, err := test.In.Encode()

			if test.ShouldError {
				assert.NotNil(t, err)
				assert.Nil(t, out)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.Expected, out)
			}
		})
	}
}

func TestSimpleErrorDecode(t *testing.T) {
	tests := []struct {
		Name        string
		In          []byte
		Expected    Value
		ShouldError bool
	}{
		{
			Name:        "Simple error",
			In:          []byte("-ERR invalid command\r\n"),
			Expected:    simpleError{value: "ERR invalid command"},
			ShouldError: false,
		},
		{
			Name:        "Empty error",
			In:          []byte("-\r\n"),
			Expected:    simpleError{value: ""},
			ShouldError: false,
		},
		{
			Name:        "Error with spaces",
			In:          []byte("-ERR syntax error\r\n"),
			Expected:    simpleError{value: "ERR syntax error"},
			ShouldError: false,
		},
		{
			Name:        "Error with special characters",
			In:          []byte("-ERR!@#$%\r\n"),
			Expected:    simpleError{value: "ERR!@#$%"},
			ShouldError: false,
		},
		{
			Name:        "Single character error",
			In:          []byte("-E\r\n"),
			Expected:    simpleError{value: "E"},
			ShouldError: false,
		},
		{
			Name:        "Error with numbers",
			In:          []byte("-ERR 404\r\n"),
			Expected:    simpleError{value: "ERR 404"},
			ShouldError: false,
		},
		{
			Name:        "Error with tab",
			In:          []byte("-ERR\tinvalid\r\n"),
			Expected:    simpleError{value: "ERR\tinvalid"},
			ShouldError: false,
		},
		{
			Name:        "Long error",
			In:          []byte("-ERR This is a very long error message to test decoding with more than a few characters\r\n"),
			Expected:    simpleError{value: "ERR This is a very long error message to test decoding with more than a few characters"},
			ShouldError: false,
		},
		{
			Name:        "Missing CRLF",
			In:          []byte("-ERR invalid"),
			Expected:    nil,
			ShouldError: true,
		},
		{
			Name:        "Invalid prefix",
			In:          []byte("+ERR invalid\r\n"),
			Expected:    nil,
			ShouldError: true,
		},
		{
			Name:        "Empty input",
			In:          []byte(""),
			Expected:    nil,
			ShouldError: true,
		},
		{
			Name:        "Error with forbidden CRLF",
			In:          []byte("-ERR invalid\r\ncommand\r\n"),
			Expected:    nil,
			ShouldError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			resp := simpleError{}
			out, err := resp.Decode(test.In)

			if test.ShouldError {
				assert.NotNil(t, err)
				assert.Equal(t, test.Expected, out)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.Expected, out)
			}
		})
	}
}
