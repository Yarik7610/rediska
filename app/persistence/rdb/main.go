package rdb

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/codecrafters-io/redis-starter-go/app/memory"
)

type header struct {
	name    string
	version int
}

func (h *header) String() string {
	return fmt.Sprintf("HEADER\nName: %s, version: %d\n", h.name, h.version)
}

type metadata struct {
	data map[string]string
}

func (m *metadata) String() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintln("METADATA"))
	for key, value := range m.data {
		b.WriteString(fmt.Sprintf("Key: %s, Value: %s\n", key, value))
	}
	return b.String()
}

// No division on 2 maps: expired and unexpired
type database struct {
	dbSelector              int
	keysCount               int
	keysWithExpirationCount int
	items                   map[string]memory.Item
}

func (d *database) String() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("DATABASE #%d\n", d.dbSelector))
	b.WriteString(fmt.Sprintf("Keys count: %d, Keys with expiration count: %d\n", d.keysCount, d.keysWithExpirationCount))
	for key, value := range d.items {
		b.WriteString(fmt.Sprintf("Key: %s, Value: %+v\n", key, value))
	}
	return b.String()
}

type end struct {
	checksum string
}

func (e *end) String() string {
	return fmt.Sprintf("END\nChecksum: %x\n", e.checksum)
}

var rdbEOF error = errors.New("EOF")

const (
	OP_EOF          = 0xFF
	OP_SELECTDB     = 0xFE
	OP_EXPIRETIME   = 0xFD
	OP_EXPIRETIMEMS = 0xFC
	OP_RESIZEDB     = 0xFB
	OP_AUX          = 0xFA
)

const STRING_ENCODING = 0

const EMPTY_DB_HEX = "524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2"

func IsFileExists(dir, filename string) bool {
	path := filepath.Join(dir, filename)
	fmt.Println(path)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

func ReadRDBFile(dir, filename string) ([]byte, error) {
	path := filepath.Join(dir, filename)
	data, err := os.ReadFile(path)
	return data, err
}
