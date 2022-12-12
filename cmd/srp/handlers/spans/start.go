package spans

import (
	"fmt"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"span-repl/cmd/srp/state"
	"time"
)

type StartCmd struct {
	Operation string `arg:"" help:"Operation name."`

	Service  string `help:"Service name of the span."`
	Resource string `help:"Resource name of the span."`
	Type     string `help:"Span type of the span."`
	SpanID   uint64 `help:"Specify a custom ID for the span."`
}

func (e *StartCmd) Run(s *state.AppState) error {
	return state.WithAppState(s, e.StatefulRun)
}

func (e *StartCmd) StatefulRun(state *spansState) error {
	t := time.Now()
	opts := []tracer.StartSpanOption{
		tracer.StartTime(t),
	}
	if state.curSpan != nil {
		opts = append(opts, tracer.ChildOf(state.curSpan.spanContext))
	}
	if e.Service != "" {
		opts = append(opts, tracer.ServiceName(e.Service))
	}
	if e.Resource != "" {
		opts = append(opts, tracer.ResourceName(e.Resource))
	}
	if e.Type != "" {
		opts = append(opts, tracer.SpanType(e.Type))
	}
	if e.SpanID != 0 {
		opts = append(opts, tracer.WithSpanID(e.SpanID))
	}
	span := tracer.StartSpan(e.Operation, opts...)
	spanCtx := span.Context()
	node := spanNode{
		operationName: e.Operation,
		span:          span,
		spanContext:   spanCtx,
		started:       t,
		parent:        state.curSpan,
	}
	if state.curSpan != nil {
		state.curSpan.children = append(state.curSpan.children, &node)
	} else {
		state.spanTree = append(state.spanTree, &node)
	}
	state.curSpan = &node
	fmt.Println("New span created with ID:", formatSpanID(spanCtx.SpanID()))
	return nil
}
