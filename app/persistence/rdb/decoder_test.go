package rdb

import (
	"testing"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		name        string
		buffer      []byte
		expected    map[string]memory.Item
		expectedErr bool
	}{
		{
			name: "Valid RDB file",
			buffer: append(append(append(append([]byte{},
				[]byte{'R', 'E', 'D', 'I', 'S', '0', '0', '0', '9'}...),
				[]byte{OP_AUX, 0x03, 'k', 'e', 'y', 0x05, 'v', 'a', 'l', 'u', 'e'}...),
				[]byte{OP_SELECTDB, 0x00, OP_RESIZEDB, 0x01, 0x00, STRING_ENCODING, 0x03, 'k', 'e', 'y', 0x03, 'v', 'a', 'l'}...),
				[]byte{OP_EOF, 0, 0, 0, 0, 0, 0, 0, 0}...),
			expected: map[string]memory.Item{
				"key": {Value: "val", Expires: time.Time{}},
			},
			expectedErr: false,
		},
		{
			name:        "Empty buffer",
			buffer:      []byte{},
			expected:    nil,
			expectedErr: true,
		},
		{
			name:        "Nil buffer",
			buffer:      nil,
			expected:    nil,
			expectedErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Decode(test.buffer)
			assert.Equal(t, test.expected, result)
			if test.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestDecodeHeader(t *testing.T) {
	tests := []struct {
		name        string
		buffer      []byte
		expected    *header
		expectedErr bool
	}{
		{
			name:   "Valid header",
			buffer: []byte("REDIS0009"),
			expected: &header{
				name:    "REDIS",
				version: 9,
			},
			expectedErr: false,
		},
		{
			name:        "Invalid magic string",
			buffer:      []byte("WRONG0009"),
			expected:    nil,
			expectedErr: true,
		},
		{
			name:        "Short buffer",
			buffer:      []byte("RED"),
			expected:    nil,
			expectedErr: true,
		},
		{
			name:        "Invalid version",
			buffer:      []byte("REDISabcd"),
			expected:    nil,
			expectedErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dec := &decoder{b: test.buffer, pos: 0, len: len(test.buffer)}
			result, err := dec.decodeHeader()
			assert.Equal(t, test.expected, result)
			if test.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestDecodeMetadata(t *testing.T) {
	tests := []struct {
		name        string
		buffer      []byte
		expected    *metadata
		expectedErr bool
	}{
		{
			name:   "Valid metadata with one key-value pair",
			buffer: []byte{OP_AUX, 0x03, 'k', 'e', 'y', 0x03, 'v', 'a', 'l', OP_EOF},
			expected: &metadata{
				data: map[string]string{"key": "val"},
			},
			expectedErr: false,
		},
		{
			name:        "Empty metadata",
			buffer:      []byte{OP_EOF},
			expected:    &metadata{data: map[string]string{}},
			expectedErr: false,
		},
		{
			name:        "Invalid opcode",
			buffer:      []byte{0x00},
			expected:    &metadata{data: map[string]string{}},
			expectedErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dec := &decoder{b: test.buffer, pos: 0, len: len(test.buffer)}
			result, err := dec.decodeMetadata()
			assert.Equal(t, test.expected, result)
			if test.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestDecodeDatabases(t *testing.T) {
	tests := []struct {
		name        string
		buffer      []byte
		expected    []*database
		expectedErr bool
	}{
		{
			name:   "Valid database with one key-value pair",
			buffer: []byte{OP_SELECTDB, 0x00, OP_RESIZEDB, 0x01, 0x00, STRING_ENCODING, 0x03, 'k', 'e', 'y', 0x03, 'v', 'a', 'l', OP_EOF},
			expected: []*database{
				{
					dbSelector:              0,
					keysCount:               1,
					keysWithExpirationCount: 0,
					items: map[string]memory.Item{
						"key": {Value: "val", Expires: time.Time{}},
					},
				},
			},
			expectedErr: false,
		},
		{
			name:        "Empty database list",
			buffer:      []byte{OP_EOF},
			expected:    nil,
			expectedErr: false,
		},
		{
			name:        "Invalid resize opcode",
			buffer:      []byte{OP_SELECTDB, 0x3A, 0x00},
			expected:    nil,
			expectedErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dec := &decoder{b: test.buffer, pos: 0, len: len(test.buffer)}
			result, err := dec.decodeDatabases()
			assert.Equal(t, test.expected, result)
			if test.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestDecodeEnd(t *testing.T) {
	tests := []struct {
		name        string
		buffer      []byte
		expected    *end
		expectedErr bool
	}{
		{
			name:   "Valid end",
			buffer: []byte{OP_EOF, 0, 0, 0, 0, 0, 0, 0, 0},
			expected: &end{
				checksum: string([]byte{0, 0, 0, 0, 0, 0, 0, 0}),
			},
			expectedErr: false,
		},
		{
			name:        "Invalid opcode",
			buffer:      []byte{0x00},
			expected:    nil,
			expectedErr: true,
		},
		{
			name:        "Short checksum",
			buffer:      []byte{OP_EOF, 0, 0, 0},
			expected:    nil,
			expectedErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dec := &decoder{b: test.buffer, pos: 0, len: len(test.buffer)}
			result, err := dec.decodeEnd()
			assert.Equal(t, test.expected, result)
			if test.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestDecodeKeyValueFunctions(t *testing.T) {
	tests := []struct {
		name          string
		buffer        []byte
		expires       time.Time
		runTest       func(*decoder, *database) error
		expectedItems map[string]memory.Item
		expectedErr   bool
	}{
		{
			name: "decodeKeyValuePairs: Valid two key-value pairs",
			buffer: []byte{
				STRING_ENCODING, 0x04, 'k', 'e', 'y', '1', 0x04, 'v', 'a', 'l', '1',
				STRING_ENCODING, 0x04, 'k', 'e', 'y', '2', 0x04, 'v', 'a', 'l', '2',
				OP_EOF,
			},
			runTest: func(dec *decoder, db *database) error {
				return dec.decodeKeyValuePairs(db)
			},
			expectedItems: map[string]memory.Item{
				"key1": {Value: "val1", Expires: time.Time{}},
				"key2": {Value: "val2", Expires: time.Time{}},
			},
			expectedErr: true, // rdbEOF
		},
		{
			name:   "decodeKeyValuePairs: Empty buffer",
			buffer: []byte{},
			runTest: func(dec *decoder, db *database) error {
				return dec.decodeKeyValuePairs(db)
			},
			expectedItems: map[string]memory.Item{},
			expectedErr:   true,
		},
		{
			name:   "decodeKeyValuePairs: Stop at OP_SELECTDB",
			buffer: []byte{OP_SELECTDB},
			runTest: func(dec *decoder, db *database) error {
				return dec.decodeKeyValuePairs(db)
			},
			expectedItems: map[string]memory.Item{},
			expectedErr:   false,
		},
		{
			name:   "decodeKeyValueMS: Valid key-value with millisecond expiration",
			buffer: []byte{0xE8, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, STRING_ENCODING, 0x03, 'k', 'e', 'y', 0x03, 'v', 'a', 'l'},
			runTest: func(dec *decoder, db *database) error {
				return dec.decodeKeyValueMS(db)
			},
			expectedItems: map[string]memory.Item{
				"key": {Value: "val", Expires: time.UnixMilli(1000)},
			},
			expectedErr: false,
		},
		{
			name:   "decodeKeyValueMS: Invalid timestamp",
			buffer: []byte{0x00, 0x00},
			runTest: func(dec *decoder, db *database) error {
				return dec.decodeKeyValueMS(db)
			},
			expectedItems: map[string]memory.Item{},
			expectedErr:   true,
		},
		{
			name:   "decodeKeyValueS: Valid key-value with second expiration",
			buffer: []byte{0xE8, 0x03, 0x00, 0x00, STRING_ENCODING, 0x03, 'k', 'e', 'y', 0x03, 'v', 'a', 'l'},
			runTest: func(dec *decoder, db *database) error {
				return dec.decodeKeyValueS(db)
			},
			expectedItems: map[string]memory.Item{
				"key": {Value: "val", Expires: time.Unix(1000, 0)},
			},
			expectedErr: false,
		},
		{
			name:   "decodeKeyValueS: Invalid timestamp",
			buffer: []byte{0x00},
			runTest: func(dec *decoder, db *database) error {
				return dec.decodeKeyValueS(db)
			},
			expectedItems: map[string]memory.Item{},
			expectedErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := &database{items: make(map[string]memory.Item)}
			dec := &decoder{b: test.buffer, pos: 0, len: len(test.buffer)}
			err := test.runTest(dec, db)
			if test.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if len(test.expectedItems) > 0 {
					for key, expectedItem := range test.expectedItems {
						item, exists := db.items[key]
						assert.True(t, exists)
						assert.Equal(t, expectedItem.Value, item.Value)
						assert.Equal(t, expectedItem.Expires, item.Expires)
					}
				}
			}
		})
	}
}

func TestDecodeValue(t *testing.T) {
	tests := []struct {
		name        string
		buffer      []byte
		valueType   uint8
		expected    string
		expectedErr bool
	}{
		{
			name:        "Valid string value",
			buffer:      []byte{0x03, 'v', 'a', 'l'},
			valueType:   STRING_ENCODING,
			expected:    "val",
			expectedErr: false,
		},
		{
			name:        "Unsupported value type",
			buffer:      []byte{},
			valueType:   0xFF,
			expected:    "",
			expectedErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dec := &decoder{b: test.buffer, pos: 0, len: len(test.buffer)}
			result, err := dec.decodeValue(test.valueType)
			assert.Equal(t, test.expected, result)
			if test.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
