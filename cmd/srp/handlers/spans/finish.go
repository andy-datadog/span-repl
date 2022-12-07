package spans

import (
	"errors"
	"fmt"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"span-repl/cmd/srp/state"
	"time"
)

type FinishCmd struct {
	Error string `help:"Error string to attach to the span."`
}

func (e *FinishCmd) Run(s *state.AppState) error {
	return state.WithAppState(s, e.StatefulRun)
}

func (e *FinishCmd) StatefulRun(state *spansState) error {
	if s := state.curSpan; s != nil {
		if s.finished == nil {
			t := time.Now()
			s.finished = &t
			opts := []tracer.FinishOption{tracer.FinishTime(t)}
			if e.Error != "" {
				opts = append(opts, tracer.WithError(errors.New(e.Error)))
				s.error = e.Error
			}
			s.span.Finish(opts...)
			if s.parent != nil {
				state.curSpan = s.parent
			}
			tracer.Flush()
		} else {
			return fmt.Errorf("span is already finished")
		}
		return nil
	} else {
		return fmt.Errorf("no span selected")
	}
}
