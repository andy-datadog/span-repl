package help

import (
	"github.com/alecthomas/kong"
	"span-repl/cmd/srp/state"
)

type Cmd struct{}

type State struct {
	K *kong.Kong
}

func (e *Cmd) Run(s *state.AppState) error {
	return state.WithAppState(s, e.StatefulRun)
}

func (e *Cmd) StatefulRun(state *State) error {
	printHelp(state.K)
	return nil
}

func printHelp(k *kong.Kong) {
	ctx, _ := kong.Trace(k, nil)

	ctxCopy := *ctx

	kongCopy := *ctxCopy.Kong
	ctxCopy.Kong = &kongCopy

	modelCopy := *kongCopy.Model
	kongCopy.Model = &modelCopy

	modelCopy.HelpFlag = nil

	kong.DefaultHelpPrinter(kong.HelpOptions{
		NoAppSummary: true,
	}, &ctxCopy)
}
