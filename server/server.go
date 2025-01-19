// Server struct and methods for starting/stopping
package server

import (
	"log"
	"time"

	"github.com/IbrahimShahzad/diameter/message"
	fsm "github.com/IbrahimShahzad/diameter/state"
	"github.com/IbrahimShahzad/diameter/transport"
)

const defaultServerAddr = "localhost:3868"
const defaultWatchdogTTL = 10
const defaultConnectionTimeout = 5 * time.Second
const defaultEventBufferSize = 10
const defaultMessageQueueSize = 10

// multiple connections with different peers
// each with its own FSM

type Server struct {
	ServerOptions
	peers map[string]*Peer
}

type ServerOptionsFunc func(*ServerOptions)

type ServerOptions struct {
	serverAddr            string
	protocol              transport.ProtocolType
	connectionTimeout     time.Duration
	watchdogTTL           time.Duration
	eventBufferSize       int
	messageQueueSize      int
	supportedApplications []uint32
}

func defaultServerOptions() ServerOptions {
	return ServerOptions{
		serverAddr:            defaultServerAddr,
		protocol:              transport.Proto_TCP,
		connectionTimeout:     defaultConnectionTimeout,
		watchdogTTL:           defaultWatchdogTTL,
		eventBufferSize:       defaultEventBufferSize,
		messageQueueSize:      defaultMessageQueueSize,
		supportedApplications: []uint32{},
	}
}

func WithServerAddr(serverAddr string) ServerOptionsFunc {
	return func(o *ServerOptions) {
		o.serverAddr = serverAddr
	}
}

func WithSCTP() ServerOptionsFunc {
	return func(o *ServerOptions) {
		o.protocol = transport.Proto_SCTP
	}
}

func WithTCP() ServerOptionsFunc {
	return func(o *ServerOptions) {
		o.protocol = transport.Proto_TCP
	}
}

func WithConnectionTimeout(timeout time.Duration) ServerOptionsFunc {
	return func(o *ServerOptions) {
		o.connectionTimeout = timeout
	}
}

func WithWatchdogTTL(ttl time.Duration) ServerOptionsFunc {
	return func(o *ServerOptions) {
		o.watchdogTTL = ttl
	}
}

func WithEventBufferSize(size int) ServerOptionsFunc {
	return func(o *ServerOptions) {
		o.eventBufferSize = size
	}
}

func WithMessageQueueSize(size int) ServerOptionsFunc {
	return func(o *ServerOptions) {
		o.messageQueueSize = size
	}
}

func WithSupportedApplications(apps ...uint32) ServerOptionsFunc {
	return func(o *ServerOptions) {
		o.supportedApplications = apps
	}
}

func NewServer(opts ...ServerOptionsFunc) *Server {
	options := defaultServerOptions()
	for _, optFunc := range opts {
		optFunc(&options)
	}

	return &Server{
		peers:         make(map[string]*Peer),
		ServerOptions: options,
	}
}

func (s *Server) AddNewPeer(conn *transport.DiameterConnection) {
	s.peers[conn.RemoteAddr().String()] = &Peer{
		conn:         conn,
		fsm:          fsm.NewFSM(StateClosed),
		EventChan:    make(chan fsm.Event, s.eventBufferSize),
		messageQueue: make(chan *message.DiameterMessage, s.messageQueueSize),
	}
	// initialize the FSM
	s.peers[conn.RemoteAddr().String()].InitializeFSM()
}

func (s *Server) ListenAndServe() error {
	listener, err := transport.NewDiameterListener(s.serverAddr, s.protocol, s.connectionTimeout)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		s.AddNewPeer(conn)
		go s.handlePeer(s.peers[conn.RemoteAddr().String()])
	}
}

func (s *Server) handlePeer(p *Peer) {
	defer p.conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := p.conn.Read(buffer)
		if err != nil {
			return
		}

		log.Printf("Received message: %s", string(buffer[:n]))

		// Parse the message
		msg, err := message.DecodeMessage(buffer[:n])
		if err != nil {
			log.Printf("Error parsing message: %v", err)
			return
		}
		// Handle the message
		p.handleMessage(msg)
	}
}

func (s *Server) Addr() string {
	return s.serverAddr
}
