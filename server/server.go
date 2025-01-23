// Server struct and methods for starting/stopping
package server

import (
	"log/slog"
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
const defaultMessageReadSize = 1024

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
		fsm:          fsm.NewDiameterFSM(),
		EventChan:    make(chan fsm.Event, s.eventBufferSize),
		messageQueue: make(chan *message.DiameterMessage, s.messageQueueSize),
	}
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

	buffer := make([]byte, defaultMessageReadSize)
	for {
		n, err := p.conn.Read(buffer)
		if err != nil {
			slog.Error("Error reading from connection", "err", err)
			return
		}

		slog.Debug("Received", "message", string(buffer[:n]))

		// Parse the message
		msg, err := message.DecodeMessage(buffer[:n])
		if err != nil {
			slog.Error("Error parsing", "message", err)
			return
		}
		// Handle the message
		p.handleMessage(msg)
	}
}

func (s *Server) Addr() string {
	return s.serverAddr
}
