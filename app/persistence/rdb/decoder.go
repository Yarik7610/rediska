package rdb

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/memory"
)

type decoder struct {
	b   []byte
	pos int
	len int
}

func Decode(b []byte) (map[string]memory.Item, error) {
	if len(b) == 0 || b == nil {
		return nil, fmt.Errorf("empty RDB file")
	}

	dec := decoder{b: b, pos: 0, len: len(b)}

	header, err := dec.decodeHeader()
	if err != nil {
		return nil, fmt.Errorf("decode header error: %v", err)
	}
	fmt.Println(header)

	metadata, err := dec.decodeMetadata()
	if err != nil {
		return nil, fmt.Errorf("decode metadata error: %v", err)
	}
	fmt.Println(metadata)

	databases, err := dec.decodeDatabases()
	if err != nil {
		return nil, fmt.Errorf("decode database error: %v", err)
	}
	for _, database := range databases {
		fmt.Println(database)
	}

	end, err := dec.decodeEnd()
	if err != nil {
		return nil, fmt.Errorf("decode end error: %v", err)
	}
	fmt.Println(end)

	var items map[string]memory.Item
	if len(databases) > 0 {
		items = databases[0].items
	} else {
		items = nil
	}
	return items, nil
}

func (dec *decoder) decodeHeader() (*header, error) {
	magicString, err := dec.traverseStringLen(5)
	if err != nil {
		return nil, err
	}
	if magicString != "REDIS" {
		return nil, fmt.Errorf("no magic string detected")
	}

	version, err := dec.traverseStringLen(4)
	if err != nil {
		return nil, err
	}
	versionAtoi, err := strconv.Atoi(version)
	if err != nil {
		return nil, fmt.Errorf("atoi error: %v", err)
	}

	return &header{name: magicString, version: versionAtoi}, nil
}

func (dec *decoder) decodeMetadata() (*metadata, error) {
	data := make(map[string]string)

	for {
		opCode, err := dec.traverseUInt8()
		if err != nil {
			return nil, err
		}
		if opCode == OP_EOF || opCode != OP_AUX {
			dec.pos--
			break
		}

		key, err := dec.decodeString()
		if err != nil {
			return nil, fmt.Errorf("decodeString error: %v", err)
		}
		value, err := dec.decodeString()
		if err != nil {
			return nil, fmt.Errorf("decodeString error: %v", err)
		}
		data[key] = value
	}

	return &metadata{data: data}, nil
}

func (dec *decoder) decodeDatabases() ([]*database, error) {
	var databases []*database

	for {
		opCode, err := dec.traverseUInt8()
		if err != nil {
			return nil, err
		}
		if opCode == OP_EOF || opCode != OP_SELECTDB {
			dec.pos--
			break
		}

		var database database
		database.dbSelector, _, err = dec.decodeLength()
		if err != nil {
			return nil, fmt.Errorf("database decodeLength error: %v", err)
		}

		resizeDbOpCode, err := dec.traverseUInt8()
		if err != nil {
			return nil, fmt.Errorf("database resize db op code error: %v", err)
		}
		if resizeDbOpCode != OP_RESIZEDB {
			return nil, fmt.Errorf("database resize db op code isn't detected, got: %d", resizeDbOpCode)
		}

		database.keysCount, _, err = dec.decodeLength()
		if err != nil {
			return nil, fmt.Errorf("database keys count error: %v", err)
		}

		database.keysWithExpirationCount, _, err = dec.decodeLength()
		if err != nil {
			return nil, fmt.Errorf("database keys with expiration count error: %v", err)
		}

		database.items = make(map[string]memory.Item, database.keysCount)
		err = dec.decodeKeyValuePairs(&database)
		if err != nil {
			if errors.Is(err, ErrorEOF) {
				databases = append(databases, &database)
				break
			}
			return nil, fmt.Errorf("decodeKeyValuePairs error: %v", err)
		}

		databases = append(databases, &database)
	}

	return databases, nil
}

func (dec *decoder) decodeEnd() (*end, error) {
	opCode, err := dec.traverseUInt8()
	if err != nil {
		return nil, err
	}
	if opCode != OP_EOF {
		return nil, fmt.Errorf("decode file end wrong op code: %d", opCode)
	}

	checksum, err := dec.traverseStringLen(8)
	if err != nil {
		return nil, err
	}

	return &end{checksum: checksum}, nil
}

func (dec *decoder) decodeKeyValuePairs(db *database) error {
	for {
		timeStampOpCode, err := dec.traverseUInt8()
		if err != nil {
			return err
		}

		switch timeStampOpCode {
		case OP_EXPIRETIME:
			dec.decodeKeyValueS(db)
		case OP_EXPIRETIMEMS:
			dec.decodeKeyValueMS(db)
		case OP_SELECTDB:
			dec.pos--
			return nil
		case OP_EOF:
			dec.pos--
			return ErrorEOF
		default:
			dec.pos--
			dec.decodeKeyValue(db, time.Time{})
		}
	}
}

func (dec *decoder) decodeKeyValueMS(db *database) error {
	expireTimestampMS, err := dec.traverseUInt64()
	if err != nil {
		return err
	}
	expires := time.UnixMilli(int64(expireTimestampMS))

	err = dec.decodeKeyValue(db, expires)
	if err != nil {
		return fmt.Errorf("decode key value error: %v", err)
	}

	return nil
}

func (dec *decoder) decodeKeyValueS(db *database) error {
	expireTimestampS, err := dec.traverseUInt32()
	if err != nil {
		return err
	}
	expires := time.Unix(int64(expireTimestampS), 0)

	err = dec.decodeKeyValue(db, expires)
	if err != nil {
		return fmt.Errorf("decode key value error: %v", err)
	}

	return nil
}

func (dec *decoder) decodeKeyValue(db *database, expires time.Time) error {
	valueType, err := dec.traverseUInt8()
	if err != nil {
		return fmt.Errorf("value type decode error: %v", err)
	}

	key, err := dec.decodeString()
	if err != nil {
		return fmt.Errorf("key decode string error: %v", err)
	}

	value, err := dec.decodeValue(valueType)
	if err != nil {
		return fmt.Errorf("decode value error: %v", err)
	}

	db.items[key] = memory.Item{Value: value, Expires: expires}
	return nil
}

func (dec *decoder) decodeValue(valueType uint8) (string, error) {
	switch valueType {
	case STRING_ENCODING:
		value, err := dec.decodeString()
		if err != nil {
			return "", fmt.Errorf("decode string error: %v", err)
		}
		return value, nil
	default:
		return "", fmt.Errorf("unsupported value type: %d", valueType)
	}
}
