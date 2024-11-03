package server

import (
	"log"

	fsm "github.com/IbrahimShahzad/diameter/state"
)

// TODO: Implement the Server struct and methods

// Server-specific states
const (
	_ fsm.State = iota
	StateClosed
	StateROpen
	StateWaitICEA
	StateClosing
)

// Server-specific events
const (
	_ fsm.Event = iota
	EventConnCERReceived
	EventCEAReceived
	EventTimeout
	EventDisconnect
	EventDPRReceived
	EventDPAReceived
	EventDWRReceived
	EventDWAReceived
)

// InitializeFSM sets up the server FSM with specific states, events, and actions.
func (s *Server) InitializeFSM() {
	s.fsm = fsm.NewFSM(StateClosed)

	// State: Closed
	s.fsm.AddTransition(StateClosed, StateROpen, EventConnCERReceived, func() error {
		s.sendCEA()
		return nil
	})

	// State: R-Open
	s.fsm.AddTransition(StateROpen, StateROpen, EventDWRReceived, func() error {
		s.sendDWA()
		return nil
	})
	s.fsm.AddTransition(StateROpen, StateClosing, EventDisconnect, func() error {
		s.sendDPR()
		return nil
	})
	s.fsm.AddTransition(StateROpen, StateClosing, EventDPRReceived, func() error {
		s.sendDPA()
		s.cleanup()
		return nil
	})

	// State: Closing
	s.fsm.AddTransition(StateClosing, StateClosed, EventDPAReceived, s.cleanup)
	s.fsm.AddTransition(StateClosing, StateClosed, EventTimeout, s.cleanup)
}

// Helper functions for transitions

func (s *Server) sendCEA() {
	log.Println("Sending Capabilities-Exchange-Answer (CEA) in response to CER.")
	// Code to send a CEA message
}

func (s *Server) sendDWA() {
	log.Println("Sending Diameter Watchdog Answer (DWA) in response to DWR.")
	// Code to send a DWA message
}

func (s *Server) sendDPR() {
	log.Println("Sending Disconnect-Peer-Request (DPR) to client.")
	// Code to send a DPR message
}

func (s *Server) sendDPA() {
	log.Println("Sending Disconnect-Peer-Answer (DPA) in response to DPR.")
	// Code to send a DPA message
}

func (s *Server) cleanup() error {
	log.Println("Cleaning up server resources.")
	// Code to close connection and reset resources
	return nil
}
