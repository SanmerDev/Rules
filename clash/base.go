package clash

import (
	"errors"
	"strings"
)

var (
	errPayload = errors.New("payload error")
	noResolve  = "no-resolve"
)

func HasNoResolve(params []string) bool {
	for _, p := range params {
		if p == noResolve {
			return true
		}
	}
	return false
}

func WriteNoResolve(builder *strings.Builder) {
	builder.WriteString(noResolve)
}
