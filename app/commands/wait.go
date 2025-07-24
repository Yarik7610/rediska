package commands

import (
	"fmt"
	"net"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *Controller) wait(args []string, conn net.Conn) resp.Value {
	if len(args) != 2 {
		return resp.SimpleError{Value: "WAIT command error: only 2 more arguments supported"}
	}

	if m, ok := c.replication.(replication.Master); ok {
		numReplicas, err := strconv.Atoi(args[0])
		if err != nil {
			return resp.SimpleError{Value: fmt.Sprintf("WAIT command number of replicas atoi error: %v", err)}
		}
		_, err = strconv.Atoi(args[1])
		if err != nil {
			return resp.SimpleError{Value: fmt.Sprintf("WAIT command timeout (MS) atoi error: %v", err)}
		}

		if numReplicas == 0 {
			return resp.Integer{Value: 0}
		}

		// timer := time.After(time.Millisecond * time.Duration(timeoutMS))
		// select  {
		// case timer<-:

		// }
		// case
		// for _, replica := range m.GetReplicas() {

		// }
		return resp.Integer{Value: len(m.GetReplicas())}
	}

	return resp.SimpleError{Value: "WAIT cannot be used with replica instances"}
}
