// Client state machine and event handling
package client

import (
	"log"

	"github.com/IbrahimShahzad/diameter/message"
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
	c.fsm.AddTransition(StateClosed, StateWaitConnAck, EventStart, c.sendConnRequest)

	// State: Wait-Conn-Ack
	c.fsm.AddTransition(StateWaitConnAck, StateWaitCEA, EventConnAck, c.sendCER)

	c.fsm.AddTransition(StateWaitConnAck, StateClosed, EventConnNack, c.cleanup)
	c.fsm.AddTransition(StateWaitConnAck, StateClosed, EventTimeout, c.cleanup)

	// State: Wait-CEA
	c.fsm.AddTransition(StateWaitCEA, StateIOpen, EventCEAReceived, func() error {
		c.startWatchdog()
		return nil
	})

	c.fsm.AddTransition(StateWaitCEA, StateClosed, EventNonCEAReceived, c.cleanup)
	c.fsm.AddTransition(StateWaitCEA, StateClosed, EventTimeout, c.cleanup)

	// State: I-Open
	c.fsm.AddTransition(StateIOpen, StateIOpen, EventSendMessage, c.sendMessage)
	c.fsm.AddTransition(StateIOpen, StateIOpen, EventReceiveDWR, c.sendDWA)
	c.fsm.AddTransition(StateIOpen, StateClosing, EventDisconnect, func() error {
		c.sendDPR()
		c.cleanup()
		return nil
	})
	c.fsm.AddTransition(StateIOpen, StateClosing, EventReceiveDPR, func() error {
		c.sendDPA()
		c.cleanup()
		return nil
	})

	// State: Closing
	c.fsm.AddTransition(StateClosing, StateClosed, EventReceiveDPA, c.cleanup)
	c.fsm.AddTransition(StateClosing, StateClosed, EventTimeout, c.cleanup)
}

// Helper functions for transitions

func (c *Client) sendConnRequest() error {
	log.Println("Sending connection request to server.")
	return c.Connect()
}

func (c *Client) sendCER() error {
	log.Println("Sending Capabilities-Exchange-Request (CER) to server.")
	message, err := message.NewCER()
	if err != nil {
		log.Printf("Error creating CER message: %v", err)
		return err
	}
	c.SendMessage(message)
	if err := c.fsm.Trigger(EventCEAReceived); err != nil {
		log.Printf("Error triggering CEAReceived event: %v", err)
		return err
	}
	return nil
}

func (c *Client) startWatchdog() {
	log.Println("Starting Watchdog.")
	// TODO: Code to start Watchdog timer and send DWR periodically
}

// sendMessage sends a Diameter message to the server.
// The message is taken from the client's message queue.
func (c *Client) sendMessage() error {
	log.Println("Sending Diameter message.")
	for msg := range c.messageQueue {
		if state := c.fsm.GetState(); state != StateIOpen {
			log.Printf("Client not in open state. Current state: %v", state)
			break
		}
		encodedMsg, err := msg.Encode()
		if err != nil {
			log.Printf("Error encoding message: %v", err)
			return err
		}
		_, err = c.conn.Write(encodedMsg)
		if err != nil {
			log.Printf("Error sending message: %v", err)
			return err
		}
	}
	return nil
}

func (c *Client) sendDWA() error {
	log.Println("Sending Diameter Watchdog Answer (DWA) in response to DWR.")
	// TODO: Code to send a DWA message
	return nil
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
	if c.conn != nil {
		c.conn.Close()
	}
	c.fsm.SetState(StateClosed)
	// reset message queue
	c.messageQueue = make(chan *message.DiameterMessage, messageQueueSize)
	// trigger initialisation of FSM
	c.EventChan <- EventStart
	return nil
}

// Any special handling for errors can be done here.
func (c *Client) handleError() error {
	log.Println("Handling error and resetting to closed state.")
	return c.cleanup()
}
