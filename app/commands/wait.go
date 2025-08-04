package commands

import (
	"fmt"
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *controller) wait(args []string) resp.Value {
	if len(args) != 2 {
		return resp.SimpleError{Value: "WAIT command error: only 2 more arguments supported"}
	}

	if m, ok := c.replication.(replication.Master); ok {
		numReplicas, err := strconv.Atoi(args[0])
		if err != nil {
			return resp.SimpleError{Value: fmt.Sprintf("WAIT command number of replicas atoi error: %v", err)}
		}
		timeoutMS, err := strconv.Atoi(args[1])
		if err != nil {
			return resp.SimpleError{Value: fmt.Sprintf("WAIT command timeout (MS) atoi error: %v", err)}
		}

		if numReplicas == 0 {
			return resp.Integer{Value: 0}
		}

		replicas := m.GetReplicas()
		if !m.HasPendingWrites() {
			return resp.Integer{Value: len(replicas)}
		}

		if numReplicas > len(replicas) {
			numReplicas = len(replicas)
		}

		m.Propagate([]string{"REPLCONF", "GETACK", "*"})

		timer := time.After(time.Millisecond * time.Duration(timeoutMS))
		ackedReplicas := make(map[string]bool)

	loop:
		for len(ackedReplicas) < numReplicas {
			select {
			case <-timer:
				break loop
			case ack := <-m.GetAckCh():
				if !ackedReplicas[ack.Addr] {
					ackedReplicas[ack.Addr] = true
				}
			}
		}

		// If write commands are very slow, race condition can appear
		// Thus, replica will send REPLCONF ACK ... faster, then replying to propagated command
		// This will lead to a case where all ACKS are delivered but the replicas still work on propagated command
		// Ideally, we need to wait on replica side and send REPLCONF ACK ... only when there are no more propagated write commands from master
		// Or we need to compare that ack.offset >= master.MasterReplOffset and only than push ack to ackedReplicas
		if len(ackedReplicas) == len(m.GetReplicas()) {
			m.SetHasPendingWrites(false)
		}

		return resp.Integer{Value: len(ackedReplicas)}
	}

	return resp.SimpleError{Value: "WAIT cannot be used with replica instances"}
}
