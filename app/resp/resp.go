package resp

type Value interface {
	Encode() ([]byte, error)
	Decode([]byte) ([]byte, Value, error)
}

type RESPController struct {
	SimpleString simpleString
	SimpleError  simpleError
	Integer      integer
	BulkString   bulkString
	Array        array
}

var resp = RESPController{}

// func (rc *RESPController) Decode(b []byte) ([]Value, error) {
// 	l := len(b)
// 	if l == 0 {
// 		return nil, fmt.Errorf("RESP decode error: expected not fully empty string")
// 	}

// 	res := make([]Value, 0)
// 	var curVal Value
// 	var err error

// 	for {
// 		switch b[0] {
// 		case '*':
// 			b, curVal, err = rc.Array.Decode(b)
// 		case '$':
// 			b, curVal, err = rc.BulkString.Decode(b)
// 		case ':':
// 			b, curVal, err = rc.Integer.Decode(b)
// 		case '+':
// 			b, curVal, err = rc.SimpleString.Decode(b)
// 		case '-':
// 			b, curVal, err = rc.SimpleError.Decode(b)
// 		default:
// 			return nil, fmt.Errorf("RESP decode error: detected unknown RESP type")
// 		}

// 		if err != nil {
// 			return nil, fmt.Errorf("RESP decode error: %v", err)
// 		}
// 		res = append(res, curVal)

// 		if len(b) == 0 {
// 			break
// 		}
// 	}

// 	return res, nil
// }
