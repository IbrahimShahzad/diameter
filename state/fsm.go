// Generalized finite state machine logic
// Diameter State Machine as per [RFC 6733]
// [RFC 6733]: (https://tools.ietf.org/html/rfc6733)
package state

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type State string
type Event string

// type ActionFunc func(msg *message.DiameterMessage) error
type ActionFunc[T any] func(ctx context.Context, args *T) (*T, error)

type Transition[T any] struct {
	From   State
	To     State
	Event  Event
	Action []Action[T]
}

type FSM[T any] struct {
	states       []State
	mu           sync.Mutex
	currentState State
	transitions  []Transition[T]
}

func NewFSM[T any](initialState State) *FSM[T] {
	fsm := &FSM[T]{
		currentState: initialState,
	}
	fsm.RegisterState(initialState)
	return fsm
}

func (f *FSM[T]) RegisterState(state State) {
	if f.states == nil {
		f.states = make([]State, 0)
	}
	f.states = append(f.states, state)
}

func (f *FSM[T]) AddTransition(from, to State, event Event, actions []Action[T]) {
	if f.transitions == nil {
		f.transitions = make([]Transition[T], 0)
	}
	f.transitions = append(f.transitions, Transition[T]{
		From:   from,
		To:     to,
		Event:  event,
		Action: actions,
	})
}

// Trigger attempts to transition the FSM to a new state based on the given event.
// checks for a valid transition from the current state,
// executes the associated action if any, and updates the FSM's state.
//
// Returns an error if no transition is registered for the current state or event, or if the action fails.
func (f *FSM[T]) Trigger(ctx context.Context, event Event, args *T) (*T, error) {
	var err error
	var nextState State
	var handlers []Action[T]

	for _, transition := range f.transitions {
		if transition.From == f.currentState && transition.Event == event {
			nextState = transition.To
			handlers = transition.Action
			break
		}
	}

	for _, handler := range handlers {
		if handler.Fn == nil {
			return args, errors.New(fmt.Sprintf("No handler found for event %s in state %s", event, f.currentState))
		}
		if args, err = handler.Fn(ctx, args); err != nil {
			return args, err
		}
	}
	f.currentState = nextState
	return args, nil
}

// GetState returns the current state of the FSM.
// It locks the FSM to ensure thread safety before accessing the state.
func (f *FSM[T]) GetState() State {
	return f.currentState
}

func (f *FSM[T]) SetState(s State) {
	f.currentState = s
}
