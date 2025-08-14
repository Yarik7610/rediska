package commands

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

func (c *controller) info(args []string) resp.Value {
	if len(args) != 1 {
		return resp.SimpleError{Value: "INFO command error: only 1 argument supported"}
	}

	section := args[0]
	switch section {
	case "replication":
		replicationInfo := c.replicationController.Info().String()
		return resp.BulkString{Value: &replicationInfo}
	default:
		return resp.SimpleError{Value: fmt.Sprintf("INFO unsupported section: %s", section)}
	}
}

func (c *controller) replconf(args []string, conn net.Conn) resp.Value {
	if len(args) != 2 {
		return resp.SimpleError{Value: "REPLCONF command error: only 2 more arguments supported"}
	}

	addr := utils.GetRemoteAddr(conn)

	secondCommand := args[0]
	arg := args[1]
	switch replicationController := c.replicationController.(type) {
	case replication.MasterController:
		switch strings.ToLower(secondCommand) {
		case "listening-port":
			replicationController.AddReplicaConn(conn)
			return resp.SimpleString{Value: "OK"}
		case "capa":
			if arg != "psync2" {
				return resp.SimpleError{Value: fmt.Sprintf("REPLCONF capa unsupported argument: %s", arg)}
			}
			return resp.SimpleString{Value: "OK"}
		case "ack":
			ackOffset, err := strconv.Atoi(arg)
			if err != nil {
				return resp.SimpleError{Value: fmt.Sprintf("REPLCONF ACK master offset atoi error: %s", secondCommand)}
			}
			replicationController.SendAck(addr, ackOffset)
			return nil
		default:
			return resp.SimpleError{Value: fmt.Sprintf("REPLCONF master unsupported second command: %s", secondCommand)}
		}
	case replication.ReplicaController:
		switch strings.ToUpper(secondCommand) {
		case "GETACK":
			if arg != "*" {
				return resp.SimpleError{Value: fmt.Sprintf("REPLCONF GETACK replica unsupported argument: %s", arg)}
			}
			if replicationController.GetMasterConn() != conn {
				return resp.SimpleError{Value: "REPLCONF GETACK * can be send only by master"}
			}

			// No syncing with propagated write commands from master
			// Ideally, we should wait here, until replica won't have any write commands from master in processing
			// Or we need to compare that Ack.offset >= master.MasterReplOffset and only than push Ack to ackedReplicas (wait.go)
			response := resp.CreateBulkStringArray("REPLCONF", "ACK", strconv.Itoa(replicationController.Info().MasterReplOffset))
			if err := utils.WriteCommand(response, conn); err != nil {
				return resp.SimpleError{Value: fmt.Sprintf("REPLCONF GETACK * write to master error: %v", err)}
			}
			return nil
		}
		return resp.SimpleError{Value: fmt.Sprintf("REPLCONF isn't supported for replica: %s", secondCommand)}
	default:
		return resp.SimpleError{Value: fmt.Sprintf("REPLCONF command detected unknown type assertion: %T", replicationController)}
	}
}

func (c *controller) psync(args []string, conn net.Conn) resp.Value {
	if len(args) != 2 {
		return resp.SimpleError{Value: "PSYNC command error: only 2 argument supported"}
	}

	requestedReplID := args[0]
	requestedReplOffset := args[1]

	switch replicationController := c.replicationController.(type) {
	case replication.MasterController:
		if requestedReplID == "?" && requestedReplOffset == "-1" {
			if !replicationController.IsReplica(conn) {
				return resp.SimpleError{Value: "PSYNC command error: failed to send FULLRESYNC, because no such replica exists"}
			}

			response := "FULLRESYNC" + " " + replicationController.Info().MasterReplID + " " + "0"
			if err := utils.WriteCommand(resp.SimpleString{Value: response}, conn); err != nil {
				return resp.SimpleError{Value: fmt.Sprintf("PSYNC command error: failed to send FULLRESYNC: %v", err)}
			}
			replicationController.SendRDBFile(conn)
			return nil
		}
		return resp.SimpleError{Value: fmt.Sprintf("PSYNC master unsupported replication id: %s and replication offset: %s", requestedReplID, requestedReplOffset)}
	case replication.ReplicaController:
		return resp.SimpleError{Value: "PSYNC isn't supported for replica"}
	default:
		return resp.SimpleError{Value: fmt.Sprintf("PSYNC command detected unknown type assertion: %T", replicationController)}
	}
}

func (c *controller) wait(args []string) resp.Value {
	if len(args) != 2 {
		return resp.SimpleError{Value: "WAIT command error: only 2 more arguments supported"}
	}

	if mc, ok := c.replicationController.(replication.MasterController); ok {
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

		replicas := mc.GetReplicas()
		if !mc.HasPendingWrites() {
			return resp.Integer{Value: len(replicas)}
		}

		if numReplicas > len(replicas) {
			numReplicas = len(replicas)
		}

		mc.Propagate([]string{"REPLCONF", "GETACK", "*"})

		timer := time.After(time.Millisecond * time.Duration(timeoutMS))
		ackedReplicas := make(map[string]bool)

	loop:
		for len(ackedReplicas) < numReplicas {
			select {
			case <-timer:
				break loop
			case ack := <-mc.GetAcksCh():
				if !ackedReplicas[ack.Addr] {
					ackedReplicas[ack.Addr] = true
				}
			}
		}

		// If write commands are very slow, race condition can appear
		// Thus, replica will send REPLCONF ACK ... faster, then replying to propagated command
		// This will lead to a case where all ACKS are delivered but the replicas still work on propagated command
		// Ideally, we need to wait on replica side and send REPLCONF ACK ... only when there are no more propagated write commands from master
		// Or we need to compare that Ack.offset >= master.MasterReplOffset and only than push Ack to ackedReplicas
		if len(ackedReplicas) == len(mc.GetReplicas()) {
			mc.SetHasPendingWrites(false)
		}

		return resp.Integer{Value: len(ackedReplicas)}
	}

	return resp.SimpleError{Value: "WAIT cannot be used with replica instances"}
}
