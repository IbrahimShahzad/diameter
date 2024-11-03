// TCP/SCTP connection management for clients/servers
package transport

import (
	"github.com/ishidawataru/sctp"
	"log"
	"net"
	"time"
)

// Protcol specifier
type ProtocolType int

const (
	Proto_TCP ProtocolType = iota
	Proto_SCTP
)

// DiameterConnection manages a network connection (TCP or SCTP) for Diameter
// communication.
type DiameterConnection struct {
	conn         net.Conn
	readTimeout  time.Duration
	writeTimeout time.Duration
	protocol     ProtocolType
}

// NewDiameterConnection establishes a new connection to a server
// (client-side).
func NewDiameterConnection(
	addr string,
	protocol ProtocolType,
	timeout time.Duration,
) (*DiameterConnection, error) {

	var conn net.Conn
	var err error

	switch protocol {
	case Proto_TCP:
		dialer := net.Dialer{Timeout: timeout}
		conn, err = dialer.Dial("tcp", addr)
	case Proto_SCTP:
		conn, err = sctp.DialSCTP("sctp", nil, &sctp.SCTPAddr{IPAddrs: []net.IPAddr{{IP: net.ParseIP(addr)}}})
	}
	if err != nil {
		return nil, err
	}
	return &DiameterConnection{
		conn:     conn,
		protocol: protocol,
	}, nil
}

// Read reads data from the Diameter connection.
func (dc *DiameterConnection) Read(buffer []byte) (int, error) {
	if dc.readTimeout > 0 {
		dc.conn.SetReadDeadline(time.Now().Add(dc.readTimeout))
	}
	n, err := dc.conn.Read(buffer)
	if err != nil {
		return n, err
	}
	return n, nil
}

// Write writes data to the Diameter connection.
func (dc *DiameterConnection) Write(data []byte) (int, error) {
	if dc.writeTimeout > 0 {
		dc.conn.SetWriteDeadline(time.Now().Add(dc.writeTimeout))
	}
	n, err := dc.conn.Write(data)
	if err != nil {
		return n, err
	}
	return n, nil
}

// Close closes the Diameter connection.
func (dc *DiameterConnection) Close() error {
	log.Println("Closing connection.")
	return dc.conn.Close()
}

// SetTimeouts sets read and write timeouts for the connection.
func (dc *DiameterConnection) SetTimeouts(
	readTimeout, writeTimeout time.Duration) {
	dc.readTimeout = readTimeout
	dc.writeTimeout = writeTimeout
}
