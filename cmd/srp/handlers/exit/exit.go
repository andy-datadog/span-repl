package exit

import (
	"io"
	"span-repl/cmd/srp/state"
)

type Cmd struct{}

func (e *Cmd) Run(_ *state.AppState) error {
	return io.EOF
}
