package server

import (
	"log"

	"github.com/IbrahimShahzad/diameter/message"
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

var ServerEvents = []fsm.Event{
	EventConnCERReceived,
	EventCEAReceived,
	EventTimeout,
	EventDisconnect,
	EventDPRReceived,
	EventDPAReceived,
	EventDWRReceived,
	EventDWAReceived,
}

// InitializeFSM sets up the server FSM with specific states, events, and actions.
func (s *Peer) InitializeFSM() {
	log.Println("Initializing server FSM.")
	s.fsm = fsm.NewFSM(StateClosed)

	// State: Closed
	s.fsm.AddTransition(StateClosed, StateROpen, EventConnCERReceived, func(msg *message.DiameterMessage) error {
		s.sendCEA(msg)
		return nil
	})

	// State: R-Open
	s.fsm.AddTransition(StateROpen, StateROpen, EventDWRReceived, func(msg *message.DiameterMessage) error {
		s.sendDWA()
		return nil
	})

	s.fsm.AddTransition(StateROpen, StateClosing, EventDisconnect, func(msg *message.DiameterMessage) error {
		s.sendDPR()
		return nil
	})

	s.fsm.AddTransition(StateROpen, StateClosing, EventDPRReceived, func(msg *message.DiameterMessage) error {
		s.sendDPA()
		s.cleanup(nil)
		return nil
	})

	// State: Closing
	s.fsm.AddTransition(StateClosing, StateClosed, EventDPAReceived, s.cleanup)
	s.fsm.AddTransition(StateClosing, StateClosed, EventTimeout, s.cleanup)
	s.fsm.AddTransition(StateClosing, StateClosed, EventDisconnect, s.cleanup)

	log.Println("Server FSM initialized.")
}

// Helper functions for transitions

func (p *Peer) sendCEA(cer *message.DiameterMessage) {
	log.Println("00000000000000000000000000000000000000000000000000000------")
	log.Println("Sending Capabilities-Exchange-Answer (CEA) in response to CER.")
	// Code to send a CEA message
	// log.Printf("The current state is %v\n", p.fsm.GetState())
	// s.ShowRegisteredTransitions()

	//receive CER message
	cea := p.generateCEA(cer)
	// receive on msg queue
	// cea := <-p.messageQueue
	log.Println("CEA message: ", cea)
	resp, err := cea.Encode()
	if err != nil {
		log.Printf("Error encoding CEA message: %v", err)
		return
	}
	if resp == nil {
		log.Printf("No response to send")
		return
	}
	// Send the response
	_, err = p.conn.Write(resp)
	if err != nil {
		log.Printf("Error sending response: %v", err)
		return
	}
	p.fsm.Trigger(EventCEAReceived, cea)
	log.Println("CEA message sent.")
}

func (s *Peer) sendDWA() {
	log.Println("Sending Diameter Watchdog Answer (DWA) in response to DWR.")
	// Code to send a DWA message
}

func (s *Peer) sendDPR() {
	log.Println("Sending Disconnect-Peer-Request (DPR) to Peer.")
	// Code to send a DPR message
}

func (s *Peer) sendDPA() {
	log.Println("Sending Disconnect-Peer-Answer (DPA) in response to DPR.")
	// Code to send a DPA message
}

func (s *Peer) cleanup(m *message.DiameterMessage) error {
	log.Println("Cleaning up server resources.")
	// Code to close connection and reset resources
	return nil
}
