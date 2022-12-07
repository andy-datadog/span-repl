package main

import (
	"errors"
	"fmt"
	"io"
	"span-repl/cmd/srp/handlers"
	"span-repl/cmd/srp/state"
	"strings"
)
import "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
import "github.com/peterh/liner"

func main() {
	tracer.Start(
		tracer.WithEnv("test"),
		tracer.WithService("span-repl"),
	)
	defer tracer.Stop()

	line := liner.NewLiner()
	defer func() { _ = line.Close() }()

	fmt.Println("")
	fmt.Println("span-repl ready. Type ? to show help or 'exit' to end the session.")
	var continuation bool
	var promptState promptState
	appState := state.NewAppState()
	for {
		if err := prompt(line, &promptState, &continuation, appState); err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				fmt.Printf("ERROR: %v\n", err)
			}
		}
	}
}

type promptState struct {
	parsedSoFar  []string
	curField     strings.Builder
	curQuoteMode quoteMode
	mustAppend   bool
	escapeNext   bool
}

func (p *promptState) flushField() {
	if p.curField.Len() > 0 || p.mustAppend {
		p.parsedSoFar = append(p.parsedSoFar, p.curField.String())
		p.curField.Reset()
		p.mustAppend = false
	}
}

func prompt(line *liner.State, p *promptState, continuation *bool, appState *state.AppState) error {
	ps := "> "
	if *continuation {
		ps = ". "
	}

	cmd, err := line.Prompt(ps)
	if err != nil {
		return err
	}
	line.AppendHistory(cmd)

	if *continuation {
		cmd = "\n" + cmd
		*continuation = false
	}

	parsed := parseCmd(p, cmd)
	if parsed == nil {
		*continuation = true
		return nil
	} else {
		return handlers.HandleCmd(appState, *parsed)
	}
}

type quoteMode int

const (
	quoteModeNone quoteMode = iota
	quoteModeSingle
	quoteModeDouble
)
