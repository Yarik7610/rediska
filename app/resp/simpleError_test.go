package resp

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestSimpleErrorDecode(t *testing.T) {
// 	tests := []struct {
// 		Name        string
// 		In          []byte
// 		Expected    string
// 		ShouldError bool
// 	}{
// 		{
// 			Name:        "Simple error message",
// 			In:          []byte("-ERR invalid command\r\n"),
// 			Expected:    "ERR invalid command",
// 			ShouldError: false,
// 		},
// 		{
// 			Name:        "Empty error message",
// 			In:          []byte("-\r\n"),
// 			Expected:    "",
// 			ShouldError: false,
// 		},
// 		{
// 			Name:        "Error with spaces",
// 			In:          []byte("-ERR syntax error in query\r\n"),
// 			Expected:    "ERR syntax error in query",
// 			ShouldError: false,
// 		},
// 		{
// 			Name:        "Error with special characters",
// 			In:          []byte("-ERR invalid input!@#$%\r\n"),
// 			Expected:    "ERR invalid input!@#$%",
// 			ShouldError: false,
// 		},
// 		{
// 			Name:        "Single character error",
// 			In:          []byte("-E\r\n"),
// 			Expected:    "E",
// 			ShouldError: false,
// 		},
// 		{
// 			Name:        "Error with numbers",
// 			In:          []byte("-ERR code 404\r\n"),
// 			Expected:    "ERR code 404",
// 			ShouldError: false,
// 		},
// 		{
// 			Name:        "Error with tab",
// 			In:          []byte("-ERR invalid\tinput\r\n"),
// 			Expected:    "ERR invalid\tinput",
// 			ShouldError: false,
// 		},
// 		{
// 			Name:        "Long error message",
// 			In:          []byte("-ERR this is a very long error message to test decoding with more than a few characters\r\n"),
// 			Expected:    "ERR this is a very long error message to test decoding with more than a few characters",
// 			ShouldError: false,
// 		},
// 		{
// 			Name:        "Missing CRLF",
// 			In:          []byte("-ERR invalid command"),
// 			Expected:    "",
// 			ShouldError: true,
// 		},
// 		{
// 			Name:        "Invalid prefix",
// 			In:          []byte("+ERR invalid command\r\n"),
// 			Expected:    "",
// 			ShouldError: true,
// 		},
// 		{
// 			Name:        "Empty input",
// 			In:          []byte(""),
// 			Expected:    "",
// 			ShouldError: true,
// 		},
// 		{
// 			Name:        "CRLF in error message",
// 			In:          []byte("-ERR invalid\r\ninput\r\n"),
// 			Expected:    "ERR invalid",
// 			ShouldError: false,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.Name, func(t *testing.T) {
// 			resp := RESPController{}
// 			out, err := resp.SimpleError.Decode(test.In)

// 			if test.ShouldError {
// 				assert.NotNil(t, err)
// 				assert.Equal(t, test.Expected, out)
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Equal(t, test.Expected, out)
// 			}
// 		})
// 	}
// }
