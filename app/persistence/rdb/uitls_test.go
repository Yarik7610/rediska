package rdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {
	tests := []struct {
		name        string
		buffer      []byte
		runTest     func(*decoder) (any, error)
		expected    any
		expectedErr error
	}{
		{
			name:   "traverseUInt8: Read one byte",
			buffer: []byte{0x42},
			runTest: func(dec *decoder) (any, error) {
				return dec.traverseUInt8()
			},
			expected:    uint8(0x42),
			expectedErr: nil,
		},
		{
			name:   "traverseUInt16: Read two bytes",
			buffer: []byte{0x34, 0x12}, // 0x1234 in little-endian
			runTest: func(dec *decoder) (any, error) {
				return dec.traverseUInt16()
			},
			expected:    uint16(0x1234),
			expectedErr: nil,
		},
		{
			name:   "traverseUInt32: Read four bytes",
			buffer: []byte{0x78, 0x56, 0x34, 0x12},
			runTest: func(dec *decoder) (any, error) {
				return dec.traverseUInt32()
			},
			expected:    uint32(0x12345678),
			expectedErr: nil,
		},
		{
			name:   "traverseUInt64: Read eight bytes",
			buffer: []byte{0xF0, 0xDE, 0xBC, 0x9A, 0x78, 0x56, 0x34, 0x12},
			runTest: func(dec *decoder) (any, error) {
				return dec.traverseUInt64()
			},
			expected:    uint64(0x123456789ABCDEF0),
			expectedErr: nil,
		},
		{
			name:   "traverseUintXBytes: Read four bytes",
			buffer: []byte{0x78, 0x56, 0x34, 0x12},
			runTest: func(dec *decoder) (any, error) {
				return dec.traverseUintXBytes(4)
			},
			expected:    uint64(0x12345678),
			expectedErr: nil,
		},
		{
			name:   "traverseStringLen: Read string",
			buffer: []byte{'h', 'e', 'l', 'l', 'o'},
			runTest: func(dec *decoder) (any, error) {
				return dec.traverseStringLen(5)
			},
			expected:    "hello",
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dec := &decoder{b: test.buffer, pos: 0, len: len(test.buffer)}
			result, err := test.runTest(dec)

			assert.Equal(t, result, test.expected)
			assert.Equal(t, err, test.expectedErr)
		})
	}
}
