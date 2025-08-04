package pubsub

import (
	"slices"
	"strings"
)

var ALLOWED_COMMANDS = []string{"SUBSCRIBE", "UNSUBSCRIBE", "PSUBSCRIBE", "PUNSUBSCRIBE", "PING", "QUIT"}

func (c *controller) IsSubscribeModeCommand(cmd string) bool {
	return slices.Contains(ALLOWED_COMMANDS, strings.ToUpper(cmd))
}
