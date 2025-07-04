package resp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBulkStringEncode(t *testing.T) {
	tests := []struct {
		Name        string
		In          Value
		Expected    []byte
		ShouldError bool
	}{
		{
			Name:        "Simple string",
			In:          BulkString{Value: stringPtr("hello")},
			Expected:    []byte("$5\r\nhello\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Empty string",
			In:          BulkString{Value: stringPtr("")},
			Expected:    []byte("$0\r\n\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Nil string",
			In:          BulkString{Value: nil},
			Expected:    []byte(NULL_BULK_STRING_RESP_2),
			ShouldError: false,
		},
		{
			Name:        "String with spaces",
			In:          BulkString{Value: stringPtr("hello world")},
			Expected:    []byte("$11\r\nhello world\r\n"),
			ShouldError: false,
		},
		{
			Name:        "String with special characters",
			In:          BulkString{Value: stringPtr("hello!@#$%")},
			Expected:    []byte("$10\r\nhello!@#$%\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Single character",
			In:          BulkString{Value: stringPtr("A")},
			Expected:    []byte("$1\r\nA\r\n"),
			ShouldError: false,
		},
		{
			Name:        "String with numbers",
			In:          BulkString{Value: stringPtr("12345")},
			Expected:    []byte("$5\r\n12345\r\n"),
			ShouldError: false,
		},
		{
			Name:        "String with tab",
			In:          BulkString{Value: stringPtr("hello\tworld")},
			Expected:    []byte("$11\r\nhello\tworld\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Long string",
			In:          BulkString{Value: stringPtr("This is a very long string to test encoding with more than a few characters")},
			Expected:    []byte("$75\r\nThis is a very long string to test encoding with more than a few characters\r\n"),
			ShouldError: false,
		},
		{
			Name:        "String with embedded CRLF",
			In:          BulkString{Value: stringPtr("hello\r\nworld")},
			Expected:    []byte("$12\r\nhello\r\nworld\r\n"),
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

func TestBulkStringDecode(t *testing.T) {
	tests := []struct {
		Name        string
		In          []byte
		Expected    Value
		ShouldError bool
	}{
		{
			Name:        "Simple string",
			In:          []byte("$5\r\nhello\r\n"),
			Expected:    BulkString{Value: stringPtr("hello")},
			ShouldError: false,
		},
		{
			Name:        "Empty string",
			In:          []byte("$0\r\n\r\n"),
			Expected:    BulkString{Value: stringPtr("")},
			ShouldError: false,
		},
		{
			Name:        "Nil string",
			In:          []byte(NULL_BULK_STRING_RESP_2),
			Expected:    BulkString{Value: nil},
			ShouldError: false,
		},
		{
			Name:        "String with spaces",
			In:          []byte("$11\r\nhello world\r\n"),
			Expected:    BulkString{Value: stringPtr("hello world")},
			ShouldError: false,
		},
		{
			Name:        "String with special characters",
			In:          []byte("$10\r\nhello!@#$%\r\n"),
			Expected:    BulkString{Value: stringPtr("hello!@#$%")},
			ShouldError: false,
		},
		{
			Name:        "Single character",
			In:          []byte("$1\r\nA\r\n"),
			Expected:    BulkString{Value: stringPtr("A")},
			ShouldError: false,
		},
		{
			Name:        "String with numbers",
			In:          []byte("$5\r\n12345\r\n"),
			Expected:    BulkString{Value: stringPtr("12345")},
			ShouldError: false,
		},
		{
			Name:        "String with tab",
			In:          []byte("$11\r\nhello\tworld\r\n"),
			Expected:    BulkString{Value: stringPtr("hello\tworld")},
			ShouldError: false,
		},
		{
			Name:        "Long string",
			In:          []byte("$75\r\nThis is a very long string to test decoding with more than a few characters\r\n"),
			Expected:    BulkString{Value: stringPtr("This is a very long string to test decoding with more than a few characters")},
			ShouldError: false,
		},
		{
			Name:        "String with embedded CRLF",
			In:          []byte("$12\r\nhello\r\nworld\r\n"),
			Expected:    BulkString{Value: stringPtr("hello\r\nworld")},
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
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			bs := BulkString{}
			_, out, err := bs.Decode(test.In)

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
