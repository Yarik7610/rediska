package resp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleStringEncode(t *testing.T) {
	tests := []struct {
		Name        string
		In          string
		Expected    []byte
		ShouldError bool
	}{
		{Name: "Simple usual string",
			In:          "PONG",
			Expected:    []byte("+PONG\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Empty string",
			In:          "",
			Expected:    []byte("+\r\n"),
			ShouldError: false,
		},
		{
			Name:        "String with spaces",
			In:          "HELLO WORLD",
			Expected:    []byte("+HELLO WORLD\r\n"),
			ShouldError: false,
		},
		{
			Name:        "String with special characters",
			In:          "OK!@#$%",
			Expected:    []byte("+OK!@#$%\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Single character",
			In:          "A",
			Expected:    []byte("+A\r\n"),
			ShouldError: false,
		},
		{
			Name:        "String with numbers",
			In:          "12345",
			Expected:    []byte("+12345\r\n"),
			ShouldError: false,
		},
		{
			Name:        "String with forbidden CRLF",
			In:          "HELLO\r\nWORLD",
			Expected:    nil,
			ShouldError: true,
		},
		{
			Name:        "String with tab",
			In:          "HELLO\tWORLD",
			Expected:    []byte("+HELLO\tWORLD\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Long string",
			In:          "This is a very long string to test encoding with more than a few characters",
			Expected:    []byte("+This is a very long string to test encoding with more than a few characters\r\n"),
			ShouldError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			resp := RESPController{}
			out, err := resp.SimpleString.Encode(test.In)

			if test.ShouldError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, out, test.Expected)
			}
		})
	}
}

func TestSimpleStringDecode(t *testing.T) {
	tests := []struct {
		Name        string
		In          []byte
		Expected    string
		ShouldError bool
	}{
		{
			Name:        "Simple usual string",
			In:          []byte("+PONG\r\n"),
			Expected:    "PONG",
			ShouldError: false,
		},
		{
			Name:        "Empty string",
			In:          []byte("+\r\n"),
			Expected:    "",
			ShouldError: false,
		},
		{
			Name:        "String with spaces",
			In:          []byte("+HELLO WORLD\r\n"),
			Expected:    "HELLO WORLD",
			ShouldError: false,
		},
		{
			Name:        "String with special characters",
			In:          []byte("+OK!@#$%\r\n"),
			Expected:    "OK!@#$%",
			ShouldError: false,
		},
		{
			Name:        "Single character",
			In:          []byte("+A\r\n"),
			Expected:    "A",
			ShouldError: false,
		},
		{
			Name:        "String with numbers",
			In:          []byte("+12345\r\n"),
			Expected:    "12345",
			ShouldError: false,
		},
		{
			Name:        "String with tab",
			In:          []byte("+HELLO\tWORLD\r\n"),
			Expected:    "HELLO\tWORLD",
			ShouldError: false,
		},
		{
			Name:        "Long string",
			In:          []byte("+This is a very long string to test decoding with more than a few characters\r\n"),
			Expected:    "This is a very long string to test decoding with more than a few characters",
			ShouldError: false,
		},
		{
			Name:        "Missing CRLF",
			In:          []byte("+PONG"),
			Expected:    "",
			ShouldError: true,
		},
		{
			Name:        "Invalid prefix",
			In:          []byte("-PONG\r\n"),
			Expected:    "",
			ShouldError: true,
		},
		{
			Name:        "Empty input",
			In:          []byte(""),
			Expected:    "",
			ShouldError: true,
		},
		{
			Name:        "CRLF in string",
			In:          []byte("+HELLO\r\nWORLD\r\n"),
			Expected:    "HELLO",
			ShouldError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			resp := RESPController{}
			out, err := resp.SimpleString.Decode(test.In)

			if test.ShouldError {
				assert.NotNil(t, err)
				assert.Equal(t, test.Expected, out)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, test.Expected, out)
			}
		})
	}
}
