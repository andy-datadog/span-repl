package spans

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hako/durafmt"
	"span-repl/cmd/srp/state"
	"span-repl/cmd/srp/tree"
	"time"
)

type ListCmd struct {
	All bool `short:"a" help:"Include finished spans."`
}

func (e *ListCmd) Run(s *state.AppState) error {
	return state.WithAppState(s, e.StatefulRun)
}

func (e *ListCmd) StatefulRun(state *spansState) error {
	for _, span := range state.spanTree {
		t := tree.New(e.formatSpan(state, span))
		keep := span.finished == nil || span == state.curSpan
		for _, child := range span.children {
			ct, childKeep := e.buildTree(state, child)
			if childKeep {
				keep = true
			}
			if childKeep || e.All {
				t.AddTree(ct)
			}
		}
		if keep || e.All {
			fmt.Println(t.Print())
		}
	}

	return nil
}

// buildTree returns true if all the children are finished
func (e *ListCmd) buildTree(gs *spansState, s *spanNode) (t tree.Tree, keep bool) {
	t = tree.New(e.formatSpan(gs, s))
	keep = s.finished == nil || s == gs.curSpan
	for _, child := range s.children {
		ct, childKeep := e.buildTree(gs, child)
		if childKeep {
			keep = true
		}
		if childKeep || e.All {
			t.AddTree(ct)
		}
	}
	return
}

func (e *ListCmd) formatSpan(gs *spansState, s *spanNode) string {
	operationName := s.operationName
	if operationName == "" {
		operationName = "?"
	}
	st := s.spanContext.TraceID()
	sc := s.spanContext.SpanID()
	if st != sc {
		operationName += fmt.Sprintf(
			" (%s/%s)",
			formatSpanID(st),
			formatSpanID(sc),
		)
	} else {
		operationName += fmt.Sprintf(
			" (%s)",
			formatSpanID(st),
		)
	}
	if s.finished != nil {
		duration := s.finished.Sub(s.started)
		operationName += fmt.Sprintf(
			" [finished: %s]",
			formatDuration(duration),
		)
		operationName = color.New(color.Faint).Sprint(operationName)
		if s.error != "" {
			operationName += color.New(color.Bold, color.FgHiRed).Sprintf(" {error: %s}", s.error)
		}
	} else {
		duration := time.Now().Sub(s.started)
		operationName += fmt.Sprintf(
			" <started %s ago>",
			formatDuration(duration),
		)
	}
	if gs.curSpan == s {
		operationName += color.New(color.Bold, color.FgHiBlue).Sprint(" <--")
	}
	return operationName
}

var units, _ = durafmt.DefaultUnitsCoder.Decode("y:y,w:w,d:d,h:h,m:m,s:s,ms:ms,μs:μs")

func formatDuration(d time.Duration) string {
	return durafmt.Parse(d).LimitFirstN(1).Format(units)
}
