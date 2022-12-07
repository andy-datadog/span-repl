package main

func parseCmd(state *promptState, cmd string) *[]string {
	for _, c := range cmd {
		skip := false

		// TODO(andy) could factor out dealing with quotes to make this cleaner
		if state.escapeNext {
			state.escapeNext = false
		} else if c == '\'' {
			if state.curQuoteMode == quoteModeNone {
				state.curQuoteMode = quoteModeSingle
				state.mustAppend = true
				skip = true
			} else if state.curQuoteMode == quoteModeSingle {
				state.curQuoteMode = quoteModeNone
				skip = true
			}
		} else if c == '"' {
			if state.curQuoteMode == quoteModeNone {
				state.curQuoteMode = quoteModeDouble
				state.mustAppend = true
				skip = true
			} else if state.curQuoteMode == quoteModeDouble {
				state.curQuoteMode = quoteModeNone
				skip = true
			}
		} else if c == '\\' {
			if state.curQuoteMode != quoteModeSingle {
				state.escapeNext = true
				skip = true
			}
		} else if c == ' ' {
			if state.curQuoteMode == quoteModeNone {
				skip = true
				state.flushField()
			}
		}

		if !skip {
			// Always returns nil error, no need to check error
			state.curField.WriteRune(c)
		}
	}

	if state.curQuoteMode == quoteModeNone && !state.escapeNext {
		state.flushField()
		// Copy array pointer because we are about to clear state
		result := state.parsedSoFar
		// Clear state
		*state = promptState{}
		return &result
	} else {
		return nil
	}
}
