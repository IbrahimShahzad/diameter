// Client struct and methods for starting/stopping
// FIXME: Fix client errors
package client

import (
	"log"
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
	for _, optFunc := range opts {
		optFunc(&o)
	}
	return &Client{
		conn:          nil,
		fsm:           fsm.NewFSM(StateClosed),
		EventChan:     make(chan fsm.Event, o.eventBufferSize),
		messageQueue:  make(chan *message.DiameterMessage, o.messageQueueSize),
		ClientOptions: o,
	}, nil
}

func (c *Client) Connect() error {
	log.Println("Connecting to server.")
	conn, err := transport.NewDiameterConnection(c.serverAddr, c.protocol, c.connectionTimeout)
	if err != nil {
		return err
	}
	c.fsm = fsm.NewFSM(StateClosed)
	c.conn = conn

	// Start event loop in the background
	go c.Run()

	c.EventChan <- EventStart

	// diameter.CreateAVP(diameter.GetAVPCodeFromName("Result-Code"), uint32(2001), diameter.MANDATORY_FLAG), // Success
	// the client sends a CER message to the server
	originHost, err := message.NewAVP(message.AVP_ORIGIN_HOST, "client.example.com", message.MANDATORY_FLAG)
	if err != nil {
		log.Fatalf("Failed to create AVP: %v", err)
	}
	originRealm, err := message.NewAVP(message.AVP_ORIGIN_REALM, "example.com", message.MANDATORY_FLAG)
	if err != nil {
		log.Fatalf("Failed to create AVP: %v", err)
	}

	msgCER, err := message.NewCER(
		originHost,
		originRealm,
	)
	if err != nil {
		log.Fatalf("Failed to create CER message: %v", err)
	}
	log.Printf("Sending CER message: %v\n", msgCER)
	resp, err := c.SendMessage(msgCER)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
	log.Printf("Received response: %v", resp)

	return nil
}

// // SendMessage sends a Diameter message to the server.
func (c *Client) SendMessage(msg *message.DiameterMessage) (*message.DiameterMessage, error) {
	c.messageQueue <- msg
	// wait for response
	readBuf := make([]byte, 1024)
	readBytes, err := c.conn.Read(readBuf)
	if err != nil || readBytes == 0 {
		log.Printf("Failed to read response: %v", err)
		return nil, err
	}
	response, err := message.DecodeMessage(readBuf[:readBytes])
	// create diameter message from response
	// response, err := message.ParseDiameterMessage(readBuf[:readBytes])
	if err != nil {
		log.Printf("Failed to decode response: %v", err)
		return nil, err
	}
	return response, nil
}

// Disconnect cleanly disconnects from the server.
func (c *Client) Disconnect() error {
	c.EventChan <- EventDisconnect
	return nil
}

// Run listens for and processes client events.
func (c *Client) Run() {

	readBuf := make([]byte, 1024)
	for {
		select {
		case msg := <-c.messageQueue:
			log.Printf("received message on queue: %v", msg)
			if msg == nil {
				return
			}
			bytes, err := msg.Encode()
			if err != nil {
				log.Printf("Failed to encode message: %v", err)
				continue
			}
			c.conn.Write(bytes)
			// wait for response
			readBytes, err := c.conn.Read(readBuf)
			if err != nil || readBytes == 0 {
				log.Printf("Failed to read response: %v", err)
				continue
			}
			log.Printf("Received response: %v", string(readBuf[:readBytes]))
			// clear buffer
			readBuf = make([]byte, 1024)
		case event := <-c.EventChan:
			log.Printf("received event: %v", event)
			c.fsm.Trigger(event, nil)
		}
	}
}
