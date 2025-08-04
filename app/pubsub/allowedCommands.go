package pubsub

import (
	"fmt"
	"net"
	"slices"
	"strings"
)

var ALLOWED_COMMANDS = []string{"SUBSCRIBE", "UNSUBSCRIBE", "PSUBSCRIBE", "PUNSUBSCRIBE", "PING", "QUIT"}

func (c *controller) ValidateSubscribeModeCommand(cmd string, conn net.Conn) error {
	if !c.InSubscribeMode(conn) {
		return nil
	}
	if !slices.Contains(ALLOWED_COMMANDS, strings.ToUpper(cmd)) {
		return fmt.Errorf("Can't execute '%s': only (P|S)SUBSCRIBE / (P|S)UNSUBSCRIBE / PING / QUIT / RESET are allowed in this context", strings.ToLower(cmd))
	}
	return nil
}
