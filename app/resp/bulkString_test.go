package resp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBulkStringEncode(t *testing.T) {
	tests := []struct {
		Name        string
		In          *string
		Expected    []byte
		ShouldError bool
	}{
		{
			Name:        "Simple string",
			In:          strPtr("hello"),
			Expected:    []byte("$5\r\nhello\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Empty string",
			In:          strPtr(""),
			Expected:    []byte("$0\r\n\r\n"),
			ShouldError: false,
		},
		{
			Name:        "String with spaces",
			In:          strPtr("hello world"),
			Expected:    []byte("$11\r\nhello world\r\n"),
			ShouldError: false,
		},
		{
			Name:        "String with special characters",
			In:          strPtr("hello!@#$%"),
			Expected:    []byte("$10\r\nhello!@#$%\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Single character",
			In:          strPtr("A"),
			Expected:    []byte("$1\r\nA\r\n"),
			ShouldError: false,
		},
		{
			Name:        "String with numbers",
			In:          strPtr("12345"),
			Expected:    []byte("$5\r\n12345\r\n"),
			ShouldError: false,
		},
		{
			Name:        "String with tab",
			In:          strPtr("hello\tworld"),
			Expected:    []byte("$11\r\nhello\tworld\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Long string",
			In:          strPtr("This is a very long string to test encoding with more than a few characters"),
			Expected:    []byte("$75\r\nThis is a very long string to test encoding with more than a few characters\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Nil string",
			In:          nil,
			Expected:    []byte("$-1\r\n"),
			ShouldError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			resp := bulkString{}
			out := resp.Encode(test.In)

			if test.ShouldError {
				assert.Fail(t, "Encode should return error")
			} else {
				assert.Equal(t, test.Expected, out)
			}
		})
	}
}

func TestBulkStringDecode(t *testing.T) {
	tests := []struct {
		Name        string
		In          []byte
		Expected    *string
		ShouldError bool
	}{
		{
			Name:        "Simple string",
			In:          []byte("$5\r\nhello\r\n"),
			Expected:    strPtr("hello"),
			ShouldError: false,
		},
		{
			Name:        "Empty string",
			In:          []byte("$0\r\n\r\n"),
			Expected:    strPtr(""),
			ShouldError: false,
		},
		{
			Name:        "String with spaces",
			In:          []byte("$11\r\nhello world\r\n"),
			Expected:    strPtr("hello world"),
			ShouldError: false,
		},
		{
			Name:        "String with special characters",
			In:          []byte("$10\r\nhello!@#$%\r\n"),
			Expected:    strPtr("hello!@#$%"),
			ShouldError: false,
		},
		{
			Name:        "Single character",
			In:          []byte("$1\r\nA\r\n"),
			Expected:    strPtr("A"),
			ShouldError: false,
		},
		{
			Name:        "String with numbers",
			In:          []byte("$5\r\n12345\r\n"),
			Expected:    strPtr("12345"),
			ShouldError: false,
		},
		{
			Name:        "String with tab",
			In:          []byte("$11\r\nhello\tworld\r\n"),
			Expected:    strPtr("hello\tworld"),
			ShouldError: false,
		},
		{
			Name:        "Long string",
			In:          []byte("$75\r\nThis is a very long string to test decoding with more than a few characters\r\n"),
			Expected:    strPtr("This is a very long string to test decoding with more than a few characters"),
			ShouldError: false,
		},
		{
			Name:        "Nil string",
			In:          []byte("$-1\r\n"),
			Expected:    nil,
			ShouldError: false,
		},
		{
			Name:        "Empty input",
			In:          []byte(""),
			Expected:    nil,
			ShouldError: true,
		},
		{
			Name:        "Invalid prefix",
			In:          []byte("+hello\r\n"),
			Expected:    nil,
			ShouldError: true,
		},
		{
			Name:        "Missing CRLF after length",
			In:          []byte("$5\r\nhello"),
			Expected:    nil,
			ShouldError: true,
		},
		{
			Name:        "Invalid length",
			In:          []byte("$abc\r\nhello\r\n"),
			Expected:    nil,
			ShouldError: true,
		},
		{
			Name:        "Length mismatch",
			In:          []byte("$10\r\nhello\r\n"),
			Expected:    nil,
			ShouldError: true,
		},
		{
			Name:        "Extra CRLF in string",
			In:          []byte("$12\r\nhello\r\nworld\r\n"),
			Expected:    strPtr("hello\r\nworld"),
			ShouldError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			resp := bulkString{}
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

func strPtr(s string) *string {
	return &s
}
