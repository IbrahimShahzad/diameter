// Client struct and methods for starting/stopping
package client

import (
	"context"
	"log/slog"
	"time"

	"github.com/IbrahimShahzad/diameter/message"
	fsm "github.com/IbrahimShahzad/diameter/state"
	"github.com/IbrahimShahzad/diameter/transport"
)

const defaultEventBufferSize = 10
const defaultMessageQueueSize = 10
const defaultWatchdogTTL = 10
const defaultConnectionTimeout = 5
const defaultServerAddr = "localhost:3868"

type ClientOptionsFunc func(*ClientOptions)

type ClientOptions struct {
	serverAddr        string
	protocol          transport.ProtocolType
	connectionTimeout time.Duration
	watchdogTTL       time.Duration
	eventBufferSize   int
	messageQueueSize  int
}

func defaultClientOptions() ClientOptions {
	return ClientOptions{
		serverAddr:        defaultServerAddr,
		protocol:          transport.Proto_TCP,
		connectionTimeout: defaultConnectionTimeout,
		watchdogTTL:       defaultWatchdogTTL,
		eventBufferSize:   defaultEventBufferSize,
		messageQueueSize:  defaultMessageQueueSize,
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
		o.connectionTimeout = timeout * time.Second
	}
}

func WithWatchdogTTL(ttl time.Duration) ClientOptionsFunc {
	return func(o *ClientOptions) {
		o.watchdogTTL = ttl
	}
}

func WithEventBufferSize(size int) ClientOptionsFunc {
	return func(o *ClientOptions) {
		o.eventBufferSize = size
	}
}

func WithMessageQueueSize(size int) ClientOptionsFunc {
	return func(o *ClientOptions) {
		o.messageQueueSize = size
	}
}

type Client struct {
	ClientOptions
	ctx          context.Context
	conn         *transport.DiameterConnection
	fsm          *fsm.FSM[message.DiameterMessage]
	EventChan    chan fsm.Event
	messageQueue chan *message.DiameterMessage
}

// NewClient creates a new Client instance with the provided options.
// It initializes the client with default options and then applies any provided ClientOptionsFunc.
// Returns a pointer to the newly created Client and an error if any.
func NewClient(opts ...ClientOptionsFunc) (*Client, error) {
	o := defaultClientOptions()
	for _, optFunc := range opts {
		optFunc(&o)
	}
	return &Client{
		ctx:           context.Background(),
		conn:          nil,
		fsm:           fsm.NewDiameterFSM(),
		EventChan:     make(chan fsm.Event, o.eventBufferSize),
		messageQueue:  make(chan *message.DiameterMessage, o.messageQueueSize),
		ClientOptions: o,
	}, nil
}

func (c *Client) Connect() error {
	slog.Info(
		"Connecting to server.",
		"serverAddr", c.serverAddr,
		"protocol", c.protocol,
		"connectionTimeout", c.connectionTimeout,
	)
	conn, err := transport.NewDiameterConnection(c.serverAddr, c.protocol, c.connectionTimeout)
	if err != nil {
		return err
	}
	c.conn = conn

	// Start event loop in the background
	c.ctx = context.WithValue(c.ctx, "peer", c.conn.RemoteAddr().String())
	c.ctx = context.WithValue(c.ctx, "connection", c.conn)
	c.fsm.Trigger(c.ctx, fsm.ISendConnReq, nil) // this will send CER message

	// now the client is in WAIT_CONN_ACK state
	// wait for response

	readBuf := make([]byte, 1024)

	// wait for response
	readBytes, err := c.conn.Read(readBuf)
	if err != nil || readBytes == 0 {
		// if the err is timeout, then we should trigger another event
		if err.Error() == "i/o timeout" {
			slog.Debug("Timeout while waiting for response")
			c.fsm.Trigger(c.ctx, fsm.Timeout, nil)
		} else {
			slog.Error("Failed to read response", "err", err)
			return err
		}
	}
	slog.Debug("Received response", "buffer", string(readBuf[:readBytes]))

	msg, err := message.DecodeMessage(readBuf[:readBytes])
	if err != nil {
		slog.Error("Failed to parse response", "err", err)
		return err
	}
	c.fsm.Trigger(c.ctx, fsm.RcvCEA, msg) // this will send CEA message

	// clear buffer
	readBuf = make([]byte, 1024)

	return nil
}

// // SendMessage sends a Diameter message to the server.
func (c *Client) SendMessage(msg *message.DiameterMessage) (*message.DiameterMessage, error) {
	c.messageQueue <- msg
	// wait for response
	readBuf := make([]byte, 1024)
	readBytes, err := c.conn.Read(readBuf)
	if err != nil || readBytes == 0 {
		slog.Error("Failed to read response", "err", err)
		return nil, err
	}
	response, err := message.DecodeMessage(readBuf[:readBytes])
	// create diameter message from response
	// response, err := message.ParseDiameterMessage(readBuf[:readBytes])
	if err != nil {
		slog.Error("Failed to read response", "err", err)
		return nil, err
	}
	return response, nil
}

// Disconnect cleanly disconnects from the server.
func (c *Client) Disconnect() error {
	return nil
}
