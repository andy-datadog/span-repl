package spans

import (
	"encoding/base64"
	"encoding/binary"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"time"
)

type spansState struct {
	spanTree []*spanNode
	curSpan  *spanNode
}

type spanNode struct {
	operationName string      // blankable
	span          tracer.Span // nillable
	spanContext   ddtrace.SpanContext
	children      []*spanNode
	parent        *spanNode // nillable
	started       time.Time
	finished      *time.Time // nillable
	error         string     // blankable
}

// WalkTree stops if f returns a non-nil value
func WalkTree[T any](nodes []*spanNode, f func(node *spanNode) *T) *T {
	for _, n := range nodes {
		if r := f(n); r != nil {
			return r
		}
		if r := WalkTree(n.children, f); r != nil {
			return r
		}
	}
	return nil
}

func formatSpanID(code uint64) string {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, code)
	return base64.RawURLEncoding.EncodeToString(buf)
}
