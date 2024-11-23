// Client struct and methods for starting/stopping
package client

import (
	"time"

	"github.com/IbrahimShahzad/diameter/message"
	fsm "github.com/IbrahimShahzad/diameter/state"
	"github.com/IbrahimShahzad/diameter/transport"
)

const eventBufferSize = 10
const messageQueueSize = 10
const watchdogTTL = 10

type ClientOptionsFunc func(*ClientOptions)

type ClientOptions struct {
	serverAddr        string
	protocol          transport.ProtocolType
	connectionTimeout time.Duration
	watchdogTTL       time.Duration
}

func defaultClientOptions() ClientOptions {
	return ClientOptions{
		serverAddr:        "localhost:3868",
		protocol:          transport.Proto_TCP,
		connectionTimeout: 5 * time.Second,
		watchdogTTL:       watchdogTTL,
	}
}

func WithServerAddr(serverAddr string) ClientOptionsFunc {
	return func(o *ClientOptions) {
		o.serverAddr = serverAddr
	}
}

func WithSCTP() ClientOptionsFunc {
	return func(o *ClientOptions) {
		o.protocol = transport.Proto_SCTP
	}
}

func WithTCP() ClientOptionsFunc {
	return func(o *ClientOptions) {
		o.protocol = transport.Proto_TCP
	}
}

func WithConnectionTimeout(timeout time.Duration) ClientOptionsFunc {
	return func(o *ClientOptions) {
		o.connectionTimeout = timeout
	}
}

func WithWatchdogTTL(ttl time.Duration) ClientOptionsFunc {
	return func(o *ClientOptions) {
		o.watchdogTTL = ttl
	}
}

type Client struct {
	ClientOptions
	conn         *transport.DiameterConnection
	fsm          *fsm.FSM
	EventChan    chan fsm.Event
	messageQueue chan *message.DiameterMessage
}

// NewClient creates a new Client instance with the provided options.
// It initializes the client with default options and then applies any provided ClientOptionsFunc.
// Returns a pointer to the newly created Client and an error if any.
func NewClient(opts ...ClientOptionsFunc) (*Client, error) {
	o := defaultClientOptions()
	for _, opt := range opts {
		opt(&o)
	}
	return &Client{
		conn:          nil,
		fsm:           fsm.NewFSM(fsm.StateClosed),
		EventChan:     make(chan fsm.Event, eventBufferSize),
		messageQueue:  make(chan *message.DiameterMessage, messageQueueSize),
		ClientOptions: o,
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
	c.messageQueue <- msg
	return nil
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
