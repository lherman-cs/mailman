package peer

import (
	"errors"
	"fmt"
	"strings"
)

// ReinitializeError is an error that happens when the caller calls "Connect" or "Wait"
// more than once after either call is successfull
type ReinitializeError struct{}

func (e *ReinitializeError) Error() string {
	return "Peer has been initialized. You have either called Connect or Wait"
}

// sprintError prints err in top-down order
func sprintError(err error) string {
	var stack []string
	level := 0
	spacesPerLevel := 2

	for err != nil {
		space := strings.Repeat(" ", level*spacesPerLevel)
		msg := fmt.Sprintf("%s%s", space, err)
		stack = append(stack, msg)
		err = errors.Unwrap(err)
		level++
	}

	msg := strings.Join(stack, "\n")
	return msg
}
