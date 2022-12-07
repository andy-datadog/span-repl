package spans

import (
	"fmt"
	"span-repl/cmd/srp/state"
	"strconv"
	"strings"
)

type CdCmd struct {
	Target targetWithHelp `arg:"" optional:""`
}

type targetWithHelp string

func (f *targetWithHelp) Help() string {
	return strings.TrimSpace(`
The ID of the span to change to. Select nothing to open an interactive UI.

Special values:
- /  Select no span.
- .  Select currently selected span (no-op).
- .. Select parent span.
- [some number] Select child span by index
`)
}

func (e *CdCmd) Run(s *state.AppState) error {
	return state.WithAppState(s, e.StatefulRun)
}

func (e *CdCmd) StatefulRun(state *spansState) error {
	targetStr := string(e.Target)
	switch targetStr {
	case "":
		// TODO
		return fmt.Errorf("not implemented")
	case "/":
		state.curSpan = nil
	case ".":
	case "..":
		if state.curSpan != nil {
			state.curSpan = state.curSpan.parent
		}
	default:
		if childIndex, err := strconv.Atoi(targetStr); err == nil {
			var spanList []*spanNode
			if state.curSpan != nil {
				spanList = state.curSpan.children
			} else {
				spanList = state.spanTree
			}
			if childIndex >= 0 && childIndex < len(spanList) {
				state.curSpan = spanList[childIndex]
				return nil
			}
		}

		if n := WalkTree(state.spanTree, func(node *spanNode) *spanNode {
			if formatSpanID(node.spanContext.SpanID()) == targetStr {
				return node
			}
			return nil
		}); n != nil {
			state.curSpan = n
		} else {
			return fmt.Errorf("no span matching '%s'", targetStr)
		}
	}
	return nil
}
