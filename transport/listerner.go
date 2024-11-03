// Listener for incoming connections
package transport

import (
	"log"
	"net"
	"time"

	"github.com/ishidawataru/sctp"
)

// DiameterListener manages incoming Diameter connections on the server side.
type DiameterListener struct {
	listener      net.Listener
	addr          string
	acceptTimeout time.Duration
	protocol      ProtocolType
}

// NewDiameterListener creates a new listener on the specified address.
func NewDiameterListener(addr string, protocol ProtocolType, acceptTimeout time.Duration) (*DiameterListener, error) {
	var listener net.Listener
	var err error

	switch protocol {
	case Proto_TCP:
		listener, err = net.Listen("tcp", addr)
	case Proto_SCTP:
		listener, err = sctp.ListenSCTP("sctp", &sctp.SCTPAddr{IPAddrs: []net.IPAddr{{IP: net.ParseIP(addr)}}})
	}

	if err != nil {
		return nil, err
	}
	return &DiameterListener{
		listener:      listener,
		addr:          addr,
		acceptTimeout: acceptTimeout,
		protocol:      protocol,
	}, nil
}

// Accept waits for and returns the next incoming connection, applying a timeout if specified.
func (dl *DiameterListener) Accept() (*DiameterConnection, error) {
	// If TCP, apply the standard SetDeadline for accept timeout.
	if dl.protocol == Proto_TCP {
		if dl.acceptTimeout > 0 {
			dl.listener.(*net.TCPListener).SetDeadline(time.Now().Add(dl.acceptTimeout))
		}
		conn, err := dl.listener.Accept()
		if err != nil {
			return nil, err
		}
		return &DiameterConnection{conn: conn, protocol: dl.protocol}, nil
	}

	// For SCTP, implement a custom timeout mechanism.
	if dl.protocol == Proto_SCTP {
		connChan := make(chan net.Conn)
		errChan := make(chan error)

		// Start a goroutine to accept the connection.
		go func() {
			conn, err := dl.listener.Accept()
			if err != nil {
				errChan <- err
				return
			}
			connChan <- conn
		}()

		// Wait for either a connection or a timeout.
		select {
		case conn := <-connChan:
			return &DiameterConnection{conn: conn, protocol: dl.protocol}, nil
		case err := <-errChan:
			return nil, err
		case <-time.After(dl.acceptTimeout):
			return nil, ErrAcceptTimeout
		}
	}

	return nil, UnsupportedProtocol
}

// Close closes the listener, stopping it from accepting any more connections.
func (dl *DiameterListener) Close() error {
	log.Printf("Shutting down listener on %s\n", dl.addr)
	return dl.listener.Close()
}

// Addr returns the address the listener is listening on.
func (dl *DiameterListener) Addr() net.Addr {
	return dl.listener.Addr()
}
