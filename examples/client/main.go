package main

import (
	dc "github.com/IbrahimShahzad/diameter/client"
	"log"
)

func main() {
	client, err := dc.NewClient(
		dc.WithTCP(),
		dc.WithConnectionTimeout(5),
		dc.WithServerAddr("localhost:3868"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	err = client.Connect()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	client.Disconnect()
}

// transport example
// package transport
//
// import (
// 	"log"
// 	"net"
// 	"time"
// )
//
// // DiameterListener manages incoming Diameter connections on the server side.
// type DiameterListener struct {
// 	listener   net.Listener
// 	addr       string
// 	acceptTimeout time.Duration
// }
//
// // NewDiameterListener creates a new listener on the specified address.
// func NewDiameterListener(addr string, acceptTimeout time.Duration) (*DiameterListener, error) {
// 	ln, err := net.Listen("tcp", addr) // TCP is used; SCTP can be integrated similarly.
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &DiameterListener{listener: ln, addr: addr, acceptTimeout: acceptTimeout}, nil
// }
//
// // Accept waits for and returns the next incoming connection.
// func (dl *DiameterListener) Accept() (*DiameterConnection, error) {
// 	if dl.acceptTimeout > 0 {
// 		dl.listener.(*net.TCPListener).SetDeadline(time.Now().Add(dl.acceptTimeout))
// 	}
// 	conn, err := dl.listener.Accept()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &DiameterConnection{conn: conn}, nil
// }
//
// // Close closes the listener, stopping it from accepting any more connections.
// func (dl *DiameterListener) Close() error {
// 	log.Printf("Shutting down listener on %s\n", dl.addr)
// 	return dl.listener.Close()
// }
//
// // Addr returns the address the listener is listening on.
// func (dl *DiameterListener) Addr() net.Addr {
// 	return dl.listener.Addr()
// }
