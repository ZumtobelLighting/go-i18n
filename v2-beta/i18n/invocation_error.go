package i18n

import (
	"fmt"
	"strings"
)

type invocationError struct {
	name   string
	args   []interface{}
	reason string
}

func (err *invocationError) Error() string {
	format := "invalid invocation %s"
	argc := len(err.args)
	switch argc {
	case 0:
		format += "()"
	case 1:
		format += "(%#v)"
	default:
		format += "(%#v" + strings.Repeat(", %#v", argc-1) + ")"
	}
	format += "; %s"
	formatArgs := append([]interface{}{err.name}, append(err.args, err.reason)...)
	return fmt.Sprintf(format, formatArgs...)
}
