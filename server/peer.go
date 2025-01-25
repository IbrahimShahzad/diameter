package server

import (
	"context"
	"log"

	"github.com/IbrahimShahzad/diameter/message"
	fsm "github.com/IbrahimShahzad/diameter/state"
	"github.com/IbrahimShahzad/diameter/transport"
)

type Peer struct {
	conn         *transport.DiameterConnection
	fsm          *fsm.FSM[message.DiameterMessage]
	EventChan    chan fsm.Event
	messageQueue chan *message.DiameterMessage
}

func (p *Peer) handleMessage(msg *message.DiameterMessage) []byte {
	// Process the message
	switch msg.Header.CommandCode {
	case message.COMMAND_CODE_CER:
		// Handle CER message
		ctx := context.Background()
		ctx = context.WithValue(ctx, "peer", p.conn.RemoteAddr().String())
		ctx = context.WithValue(ctx, "connection", p.conn)
		p.fsm.Trigger(ctx, fsm.RConnCER, msg)

	case message.COMMAND_CODE_DWR:
		// Handle DWR message
		log.Println("Received DWR message")
	default:
		// Handle unknown message
	}
	return nil
}
