package resp

import "fmt"

type array struct{}

func (array) Encode(arr []Value) []byte {
	l := len(arr)
	b := []byte(fmt.Sprintf("*%d\r\n", l))

	for _, val := range arr {
		val.Encode()
	}

}

func (array) Decode(str string) ([]byte, error) {

}
