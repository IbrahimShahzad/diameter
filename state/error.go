package state

import "errors"

var errNoTransitionRegisteredForState = errors.New("no transition registered for state")
