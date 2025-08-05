package transaction

import (
	"slices"
	"strings"
)

var ALLOWED_COMMANDS = []string{"MULTI", "EXEC", "DISCARD"}

func (c *controller) IsTransactionCommand(cmd string) bool {
	return slices.Contains(ALLOWED_COMMANDS, strings.ToUpper(cmd))
}
