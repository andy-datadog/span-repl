package state

import "reflect"

type AppState struct {
	data map[reflect.Type]interface{}
}

func NewAppState() *AppState {
	return &AppState{
		data: make(map[reflect.Type]interface{}),
	}
}

func WithAppState[T any](state *AppState, f func(state *T) error) error {
	newState := new(T)
	typ := reflect.TypeOf(*newState)
	if v := state.data[typ]; v != nil {
		// Use existing state if it exists
		newState = v.(*T)
	} else {
		// No state exists, use empty state and put it in the map
		state.data[typ] = newState
	}
	return f(newState)
}
