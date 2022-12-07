package handlers

import (
	"span-repl/cmd/srp/handlers/exit"
	"span-repl/cmd/srp/handlers/help"
	"span-repl/cmd/srp/handlers/spans"
	"span-repl/cmd/srp/state"
)
import "github.com/alecthomas/kong"

func HandleCmd(appState *state.AppState, cmd []string) error {
	if len(cmd) < 1 {
		return nil
	}

	var cli struct {
		Exit   exit.Cmd        `cmd:"" help:"Exit the CLI." aliases:"quit"`
		Help   help.Cmd        `cmd:"" help:"Show usage help." aliases:"?"`
		Start  spans.StartCmd  `cmd:"" help:"Start a span."`
		List   spans.ListCmd   `cmd:"" help:"List all known spans." aliases:"ls"`
		Cd     spans.CdCmd     `cmd:"" help:"Change currently selected span." aliases:"ls"`
		Finish spans.FinishCmd `cmd:"" help:"Finish the currently selected span."`
	}

	skip := false
	parser, err := kong.New(&cli, kong.Exit(func(i int) {
		skip = true
	}), kong.Name(">"))
	if err != nil {
		return err
	}

	ctx, err := parser.Parse(cmd)
	if err != nil {
		return err
	}

	if skip {
		return nil
	}

	_ = state.WithAppState(appState, func(state *help.State) error {
		state.K = parser
		return nil
	})

	return ctx.Run(appState)
}
