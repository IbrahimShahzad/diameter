// Client state machine and event handling
package client

import (
	"log"

	fsm "github.com/IbrahimShahzad/diameter/state"
)

const (
	// Client-specific states
	_ fsm.State = iota
	StateClosed
	StateWaitConnAck
	StateWaitCEA
	StateIOpen
	StateClosing
)

const (
	// Client-specific events
	_ fsm.Event = iota
	EventStart
	EventConnAck
	EventConnNack
	EventCEAReceived
	EventNonCEAReceived
	EventTimeout
	EventSendMessage
	EventReceiveDWR
	EventReceiveDWA
	EventDisconnect
	EventReceiveDPR
	EventReceiveDPA
)

// InitializeFSM sets up the client FSM with specific states, events, and actions.
func (c *Client) InitializeFSM() {
	c.fsm = fsm.NewFSM(StateClosed)

	// State: Closed
	c.fsm.AddTransition(StateClosed, StateClosed, EventStart, func() error {
		c.sendConnRequest()
		return nil
	})

	// State: Wait-Conn-Ack
	// TODO: Fix the transitions
	// c.fsm.AddTransition(StateWaitConnAck, StateIOpen, EventConnAck, StateWaitCEA, func() error {
	// 	c.sendCER()
	// 	return nil
	// })
	// c.fsm.AddTransition(StateWaitConnAck, EventConnNack, StateClosed, c.cleanup)
	// c.fsm.AddTransition(StateWaitConnAck, EventTimeout, StateClosed, c.handleError)
	//
	// // State: Wait-CEA
	// c.fsm.AddTransition(StateWaitCEA, EventCEAReceived, StateIOpen, func() {
	// 	c.startWatchdog()
	// })
	// c.fsm.AddTransition(StateWaitCEA, EventNonCEAReceived, StateClosed, c.handleError)
	// c.fsm.AddTransition(StateWaitCEA, EventTimeout, StateClosed, c.handleError)
	//
	// // State: I-Open
	// c.fsm.AddTransition(StateIOpen, EventSendMessage, StateIOpen, c.sendMessage)
	// c.fsm.AddTransition(StateIOpen, EventReceiveDWR, StateIOpen, func() {
	// 	c.sendDWA()
	// })
	// c.fsm.AddTransition(StateIOpen, EventDisconnect, StateClosing, func() {
	// 	c.sendDPR()
	// })
	// c.fsm.AddTransition(StateIOpen, EventReceiveDPR, StateClosing, func() {
	// 	c.sendDPA()
	// 	c.cleanup()
	// })
	//
	// // State: Closing
	// c.fsm.AddTransition(StateClosing, EventReceiveDPA, StateClosed, c.cleanup)
	// c.fsm.AddTransition(StateClosing, EventTimeout, StateClosed, c.cleanup)
}

// Helper functions for transitions

func (c *Client) sendConnRequest() {
	log.Println("Sending connection request to server.")
	// TODO: Code to initiate the connection
}

func (c *Client) sendCER() {
	log.Println("Sending Capabilities-Exchange-Request (CER) to server.")
	// TODO: Code to send a CER message
}

func (c *Client) startWatchdog() {
	log.Println("Starting Watchdog.")
	// TODO: Code to start Watchdog timer and send DWR periodically
}

func (c *Client) sendMessage() {
	log.Println("Sending Diameter message.")
	// TODO: Code to send a Diameter message
}

func (c *Client) sendDWA() {
	log.Println("Sending Diameter Watchdog Answer (DWA) in response to DWR.")
	// TODO: Code to send a DWA message
}

func (c *Client) sendDPR() {
	log.Println("Sending Disconnect-Peer-Request (DPR) to server.")
	// TODO: Code to send a DPR message
}

func (c *Client) sendDPA() {
	log.Println("Sending Disconnect-Peer-Answer (DPA) in response to DPR.")
	// TODO: Code to send a DPA message
}

func (c *Client) cleanup() error {
	log.Println("Cleaning up resources and resetting client state.")
	// TODO: Code to close connection and reset resources
	return nil
}

func (c *Client) handleError() error {
	log.Println("Handling error and resetting to closed state.")
	// TODO: Code to handle errors and return to Closed state
	return nil
}
