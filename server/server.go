// Server struct and methods for starting/stopping
package server

import (
	"time"

	fsm "github.com/IbrahimShahzad/diameter/state"
	"github.com/IbrahimShahzad/diameter/transport"
)

type Server struct {
	conn              *transport.DiameterConnection
	fsm               *fsm.FSM
	EventChan         chan fsm.Event
	serverAddr        string
	protocol          transport.ProtocolType
	connectionTimeout time.Duration
	watchdogTTL       time.Duration
}
