package replication

import "github.com/codecrafters-io/redis-starter-go/app/resp"

type Replica interface {
	Base
	ReadValueFromMaster() (resp.Value, error)
}
