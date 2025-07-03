package resp

type RESPController struct {
	SimpleString simpleString
	SimpleError  simpleError
	Integer      integer
	BulkString   bulkString
	Array        array
}

type Value interface {
	Encode() ([]byte, error)
	Decode([]byte) (Value, error)
}
