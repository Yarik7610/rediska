package rdb

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/codecrafters-io/redis-starter-go/app/memory"
)

type Header struct {
	name    string
	version int
}

func (h *Header) String() string {
	return fmt.Sprintf("HEADER\nName: %s, version: %d\n", h.name, h.version)
}

type Metadata struct {
	data map[string]string
}

func (m *Metadata) String() string {
	var b bytes.Buffer

	b.WriteString(fmt.Sprintln("METADATA"))
	for key, value := range m.data {
		b.WriteString(fmt.Sprintf("Key: %s, Value: %s\n", key, value))
	}

	return b.String()
}

// No division on 2 maps: expired and unexpired
type Database struct {
	dbSelector              int
	keysCount               int
	keysWithExpirationCount int
	items                   map[string]memory.Item
}

func (d *Database) String() string {
	var b bytes.Buffer

	b.WriteString(fmt.Sprintf("DATABASE #%d\n", d.dbSelector))
	b.WriteString(fmt.Sprintf("Keys count: %d, Keys with expiration count: %d\n", d.keysCount, d.keysWithExpirationCount))
	for key, value := range d.items {
		b.WriteString(fmt.Sprintf("Key: %s, Value: %+v\n", key, value))
	}

	return b.String()
}

var ErrorEOF error = errors.New("EOF")

const (
	OP_EOF          = 0xFF
	OP_SELECTDB     = 0xFE
	OP_EXPIRETIME   = 0xFD
	OP_EXPIRETIMEMS = 0xFC
	OP_RESIZEDB     = 0xFB
	OP_AUX          = 0xFA
)

const (
	STRING_ENCODING             = 0
	LIST_ENCODING               = 1
	SET_ENCODING                = 2
	SORTED_SET_ENCODING         = 3
	HASH_ENCODING               = 4
	ZIPMAP_ENCODING             = 9
	ZIPLIST_ENCODING            = 10
	INTSET_ENCODING             = 11
	SORTED_SET_ZIPLIST_ENCODING = 12
	HASHMAP_ZIPLIST_ENCODING    = 13
	LIST_QUICKLIST_ENCODING     = 14
)

func IsFileExists(dir, filename string) bool {
	path := filepath.Join(dir, filename)
	fmt.Println(path)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

func ReadRDB(dir, filename string) ([]byte, error) {
	path := filepath.Join(dir, filename)
	data, err := os.ReadFile(path)
	return data, err
}
