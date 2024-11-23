// Generalized finite state machine logic
package state

import (
	"errors"
	"fmt"
	"sync"
)

type State int

type Event int

type ActionFunc func() error

type Transition struct {
	From   State
	To     State
	Event  Event
	Action ActionFunc
}

type FSM struct {
	mu          sync.Mutex
	state       State
	transitions map[State]map[Event]Transition
}

const (
	InitialState State = iota
	StateWaitConAck
	StateOpen
	StateWaitDisAck
	StateClosed
)

func NewFSM(s State) *FSM {
	return &FSM{
		state:       s,
		transitions: make(map[State]map[Event]Transition),
	}
}

// Register a transition from one state to another in response to an event.
func (f *FSM) AddTransition(from State, to State, event Event, action ActionFunc) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.transitions[from] == nil {
		f.transitions[from] = make(map[Event]Transition)
	}

	f.transitions[from][event] = Transition{
		From:   from,
		To:     to,
		Event:  event,
		Action: action,
	}
}

// Trigger attempts to transition the FSM to a new state based on the given event.
// It locks the FSM to ensure thread safety, checks for a valid transition from the current state,
// executes the associated action if any, and updates the FSM's state.
//
// Returns an error if no transition is registered for the current state or event, or if the action fails.
func (f *FSM) Trigger(event Event) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	transitionForState, ok := f.transitions[f.state]
	if !ok {
		return errors.Join(errNoTransitionRegisteredForState, fmt.Errorf(" %d", f.state))
	}

	transition, ok := transitionForState[event]
	if !ok {
		return errors.Join(
			errNoTransitionRegisteredForState,
			fmt.Errorf(" %d with event %d", f.state, event),
		)
	}

	// Execute the action associated with the transition.
	if transition.Action != nil {
		if err := transition.Action(); err != nil {
			return err
		}
	}

	f.state = transition.To
	return nil
}

// GetState returns the current state of the FSM.
// It locks the FSM to ensure thread safety before accessing the state.
func (f *FSM) GetState() State {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.state
}

func (f *FSM) SetState(s State) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.state = s
}
