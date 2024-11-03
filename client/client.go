// Client struct and methods for starting/stopping
package client

import (
	"fmt"
	"time"

	"github.com/IbrahimShahzad/diameter/message"
	fsm "github.com/IbrahimShahzad/diameter/state"
	"github.com/IbrahimShahzad/diameter/transport"
)

const eventBufferSize = 10
const watchdogTTL = 10

type Client struct {
	conn              *transport.DiameterConnection
	fsm               *fsm.FSM
	EventChan         chan fsm.Event
	serverAddr        string
	protocol          transport.ProtocolType
	connectionTimeout time.Duration
	watchdogTTL       time.Duration
}

func NewClient(
	serverAddr string,
	protocol transport.ProtocolType,
	timeout time.Duration,
) (*Client, error) {
	return &Client{
		conn:              nil,
		fsm:               fsm.NewFSM(fsm.StateClosed),
		EventChan:         make(chan fsm.Event, eventBufferSize),
		serverAddr:        serverAddr,
		protocol:          protocol,
		connectionTimeout: timeout,
		watchdogTTL:       time.Second * watchdogTTL,
	}, nil
}

func (c *Client) Connect() error {
	conn, err := transport.NewDiameterConnection(c.serverAddr, c.protocol, c.connectionTimeout)
	if err != nil {
		return err
	}
	c.conn = conn
	c.fsm = fsm.NewFSM(fsm.StateClosed)

	// Start event loop in the background
	go c.Run()

	c.EventChan <- EventStart
	return nil
}

// // SendMessage sends a Diameter message to the server.
func (c *Client) SendMessage(msg *message.DiameterMessage) error {
	if c.fsm.GetState() != StateIOpen {
		return fmt.Errorf("client not in open state")
	}
	encodedMsg, err := msg.Encode()
	if err != nil {
		return err
	}
	_, err = c.conn.Write(encodedMsg)
	return err
}

// Disconnect cleanly disconnects from the server.
func (c *Client) Disconnect() error {
	c.EventChan <- EventDisconnect
	return nil
}

// Run listens for and processes client events.
func (c *Client) Run() {
	for event := range c.EventChan {
		c.fsm.Trigger(event)
	}
}
