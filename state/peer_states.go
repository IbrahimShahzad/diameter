package state

import (
	"context"
	"log/slog"

	"github.com/IbrahimShahzad/diameter/message"
	"github.com/IbrahimShahzad/diameter/transport"
)

const (
	Closed            State = "Closed"
	WaitConnectionAck       = "Wait-Conn-Ack"
	WaitICEA                = "Wait-I-CEA"
	Elect                   = "Elect"
	WaitReturns             = "Wait-Returns"
	ROpen                   = "R-Open"
	IOpen                   = "I-Open"
	Closing                 = "Closing"
)

const (
	Start       Event = "Start"         // The Diameter application has signaled that a connection should be initiated with the peer.
	RConnCER          = "R-Conn-CER"    // An acknowledgement is received stating that the transport connection has been established, and the associated CER has arrived.
	RcvConnAck        = "Rcv-Conn-Ack"  // A positive acknowledgement is received confirming that the transport connection is established.
	RcvConnNack       = "Rcv-Conn-Nack" // A negative acknowledgement was received stating that the transport connection was not established.
	Timeout           = "Timeout"       // An application-defined timer has expired while waiting for some event.
	RcvCER            = "Rcv-CER"       // A CER message from the peer was received.
	RcvCEA            = "Rcv-CEA"       // A CEA message from the peer was received.
	RcvNonCEA         = "Rcv-Non-CEA"   // A message, other than a CEA, from the peer was received.
	PeerDisc          = "Peer-Disc"     // A disconnection indication from the peer was received.
	RcvDPR            = "Rcv-DPR"       // A DPR message from the peer was received.
	RcvDPA            = "Rcv-DPA"       // A DPA message from the peer was received.
	WinElection       = "Win-Election"  // An election was held, and the local node was the winner.
	SendMessage       = "Send-Message"  // A message is to be sent.
	RcvMessage        = "Rcv-Message"   // A message other than CER, CEA, DPR, DPA, DWR, or DWA was received.
	Stop              = "Stop"          // The Diameter application has signaled that a connection should be terminated (e.g., on system shutdown).
)

const (
	ISendConnReq Event = "I-Send-Conn-Req" // A transport connection is initiated with the peer.
	DError             = "Diameter-Error"  // An error has occurred in the Diameter protocol.
)

type Action[T any] struct {
	Name string
	Fn   ActionFunc[T]
}

var SendConnReq = Action[message.DiameterMessage]{
	Name: "SendConnReq",
	Fn: func(ctx context.Context, args *message.DiameterMessage) (*message.DiameterMessage, error) {
		// Code to send a CER message
		sessionIDString := "1234567890"
		originHostString := "client.example.com"
		originRealmString := "example.com"

		ctx = context.WithValue(ctx, message.AVP_SESSION_ID, sessionIDString)
		sessionID, err := message.NewAVP(message.AVP_SESSION_ID, sessionIDString, message.MANDATORY_FLAG)
		if err != nil {
			return args, err
		}

		ctx = context.WithValue(ctx, message.AVP_ORIGIN_HOST, originHostString)
		originHost, err := message.NewAVP(message.AVP_ORIGIN_HOST, originHostString, message.MANDATORY_FLAG)
		if err != nil {
			return args, err
		}

		ctx = context.WithValue(ctx, message.AVP_ORIGIN_REALM, originRealmString)
		originRealm, err := message.NewAVP(message.AVP_ORIGIN_REALM, originRealmString, message.MANDATORY_FLAG)
		if err != nil {
			return args, err
		}

		return message.NewCER(
			sessionID,
			originHost,
			originRealm,
		)

	},
}

// The incoming connection associated with the R-Conn-CER is accepted as the responder connection.
var AcceptConn = Action[message.DiameterMessage]{
	Name: "AcceptConn",
	Fn: func(ctx context.Context, args *message.DiameterMessage) (*message.DiameterMessage, error) {
		slog.Info("Accepted Connection")
		return args, nil
	},
}

// The incoming connection associated with the R-Conn-CER is disconnected.
var RejectConn = Action[message.DiameterMessage]{
	Name: "RejectConn",
	Fn: func(ctx context.Context, args *message.DiameterMessage) (*message.DiameterMessage, error) {
		// Code to reject the connection
		return args, nil
	},
}

// The CER associated with the R-Conn-CER is processed.
var ProcessCER = Action[message.DiameterMessage]{
	Name: "ProcessCER",
	Fn: func(ctx context.Context, args *message.DiameterMessage) (*message.DiameterMessage, error) {
		// Code to process the CER
		peerAddr := ctx.Value("peer")
		slog.Info("Processing CER", "peer", peerAddr)
		return args, nil
	},
}

// A CER message is sent to the peer.
var SendConnAck = Action[message.DiameterMessage]{
	Name: "SendConnAck",
	Fn: func(ctx context.Context, args *message.DiameterMessage) (*message.DiameterMessage, error) {
		// Code to send a CER message
		return args, nil
	},
}

// A CEA message is sent to the peer.
var SendCEA = Action[message.DiameterMessage]{
	Name: "SendCEA",
	Fn: func(ctx context.Context, args *message.DiameterMessage) (*message.DiameterMessage, error) {
		// Code to send a CEA message
		resultAVP, err := message.NewAVP(message.AVP_RESULT_CODE, uint32(2001), message.MANDATORY_FLAG)
		if err != nil {
			slog.Debug("Error creating AVP", "name", "result", "error", err)
			return args, err
		}

		orignHostAVP, err := message.NewAVP(message.AVP_ORIGIN_HOST, "localhost", message.MANDATORY_FLAG)
		if err != nil {
			slog.Debug("Error creating AVP", "name", "originHost", "error", err)
			return args, err
		}

		orignRealmAVP, err := message.NewAVP(message.AVP_ORIGIN_REALM, "example.ims.com", message.MANDATORY_FLAG)
		if err != nil {
			slog.Debug("Error creating AVP", "name", "originRealm", "error", err)
			return args, err
		}

		vendorIDAVP, err := message.NewAVP(message.AVP_VENDOR_ID, uint32(10415), message.MANDATORY_FLAG)
		if err != nil {
			slog.Debug("Error creating AVP", "name", "vendorID", "error", err)
			return args, err
		}

		productNameAVP, err := message.NewAVP(message.AVP_PRODUCT_NAME, "Diameter Server", message.MANDATORY_FLAG)
		if err != nil {
			slog.Debug("Error creating AVP", "name", "productName", "error", err)
			return args, err
		}

		slog.Debug("Sending Capabilities-Exchange-Answer (CEA) in response to CER.")
		return message.NewResponseFromRequest(args,
			resultAVP,
			orignHostAVP,
			orignRealmAVP,
			vendorIDAVP,
			productNameAVP)
	},
}

// If necessary, the connection is shut down, and any local resources are freed.
var Cleanup = Action[message.DiameterMessage]{
	Name: "Cleanup",
	Fn: func(ctx context.Context, args *message.DiameterMessage) (*message.DiameterMessage, error) {
		// Code to clean up resources
		return args, nil
	},
}

// The transport layer connection is disconnected
var DiameterError = Action[message.DiameterMessage]{
	Name: "DiameterError",
	Fn: func(ctx context.Context, args *message.DiameterMessage) (*message.DiameterMessage, error) {
		// Code to handle an error
		return args, nil
	},
}

// A received CEA is processed.
var ProcessCEA = Action[message.DiameterMessage]{
	Name: "ProcessCEA",
	Fn: func(ctx context.Context, args *message.DiameterMessage) (*message.DiameterMessage, error) {
		// Code to process a CEA message
		slog.Info("Processing CEA message.")
		resultCode, _, err := message.GetResultCode(args)
		if err != nil {
			slog.Error("Error getting Result-Code", "error", err)
			return args, err
		}
		slog.Info("Result-Code", "code", resultCode)

		// TODO:
		// sessionIDorig := ctx.Value(message.AVP_SESSION_ID)
		// match session ID
		// sessionID, _, err := message.GetSessionID(args)
		// if err != nil {
		// 	slog.Error("Error getting Session-ID", "error", err)
		// 	return args, err
		// }
		// if sessionID != sessionIDorig {
		//	slog.Error("Session-ID mismatch", "sessionID", sessionID, "expected", sessionIDorig)
		//	return args, fmt.Errorf("Session-ID mismatch")
		// }
		return args, nil
	},
}

// A DPR message is sent to the peer.
var SendDPR = Action[message.DiameterMessage]{
	Name: "SendDPR",
	Fn: func(ctx context.Context, args *message.DiameterMessage) (*message.DiameterMessage, error) {
		// Code to send a DPR message
		return args, nil
	},
}

// A DPA message is sent to the peer.
var SendDPA = Action[message.DiameterMessage]{
	Name: "SendDPA",
	Fn: func(ctx context.Context, args *message.DiameterMessage) (*message.DiameterMessage, error) {
		// Code to send a DPA message
		return args, nil
	},
}

// The transport layer connection is disconnected, and local resources are freed.
var Disconnect = Action[message.DiameterMessage]{
	Name: "Disconnect",
	Fn: func(ctx context.Context, args *message.DiameterMessage) (*message.DiameterMessage, error) {
		// Code to disconnect the connection
		return args, nil
	},
}

// An election occurs
var Election = Action[message.DiameterMessage]{
	Name: "Election",
	Fn: func(ctx context.Context, args *message.DiameterMessage) (*message.DiameterMessage, error) {
		// Code to handle an election
		return args, nil
	},
}

// A message is sent.
var SendDiameterMessage = Action[message.DiameterMessage]{
	Name: "SendMessage",
	Fn: func(ctx context.Context, args *message.DiameterMessage) (*message.DiameterMessage, error) {
		// Code to send a message
		conn := ctx.Value("connection").(*transport.DiameterConnection)
		slog.Info("Sending Diameter message.")
		encodedMsg, err := args.Encode()
		if err != nil {
			slog.Error("Error encoding message", "error", err)
			return args, err
		}
		_, err = conn.Write(encodedMsg)
		if err != nil {
			slog.Error("Error sending message", "error", err)
			return args, err
		}
		return args, nil
	},
}

// A DWR message is sent.
var SendDWR = Action[message.DiameterMessage]{
	Name: "SendDWR",
	Fn: func(ctx context.Context, args *message.DiameterMessage) (*message.DiameterMessage, error) {
		// Code to send a DWR message
		return args, nil
	},
}

// A DWA message is sent.
var SendDWA = Action[message.DiameterMessage]{
	Name: "SendDWA",
	Fn: func(ctx context.Context, args *message.DiameterMessage) (*message.DiameterMessage, error) {
		// Code to send a DWA message
		return args, nil
	},
}

// The DWR message is serviced.
var ProcessDWR = Action[message.DiameterMessage]{
	Name: "ProcessDWR",
	Fn: func(ctx context.Context, args *message.DiameterMessage) (*message.DiameterMessage, error) {
		// Code to process a DWR message
		return args, nil
	},
}

// The DWA message is serviced.
var ProcessDWA = Action[message.DiameterMessage]{
	Name: "ProcessDWA",
	Fn: func(ctx context.Context, args *message.DiameterMessage) (*message.DiameterMessage, error) {
		// Code to process a DWA message
		return args, nil
	},
}

// A message is serviced.
var ProcessMessage = Action[message.DiameterMessage]{
	Name: "ProcessMessage",
	Fn: func(ctx context.Context, args *message.DiameterMessage) (*message.DiameterMessage, error) {
		// Code to process a DWR message
		return args, nil
	},
}

// TODO: implement state transitions
// Closed --> WaitConnAck: Start / I-Snd-Conn-Req
// Closed --> ROpen: R-Conn-CER / R-Accept, Process-CER, R-Snd-CEA
// WaitConnAck --> WaitICEA: I-Rcv-Conn-Ack / I-Snd-CER
// WaitConnAck --> Closed: I-Rcv-Conn-Nack / Cleanup
// WaitConnAck --> WaitConnAckElect: R-Conn-CER / R-Accept, Process-CER
// WaitConnAck --> Closed: Timeout / Error
// WaitICEA --> IOpen: I-Rcv-CEA / Process-CEA
// WaitICEA --> WaitReturns: R-Conn-CER / R-Accept, Process-CER, Elect
// WaitICEA --> Closed: I-Peer-Disc / I-Disc
// WaitICEA --> Closed: I-Rcv-Non-CEA / Error
// WaitICEA --> Closed: Timeout / Error
// WaitConnAckElect --> WaitReturns: I-Rcv-Conn-Ack / I-Snd-CER, Elect
// WaitConnAckElect --> ROpen: I-Rcv-Conn-Nack / R-Snd-CEA
// WaitConnAckElect --> WaitConnAck: R-Peer-Disc / R-Disc
// WaitConnAckElect --> WaitConnAckElect: R-Conn-CER / R-Reject
// WaitConnAckElect --> Closed: Timeout / Error
// WaitReturns --> ROpen: Win-Election / I-Disc, R-Snd-CEA
// WaitReturns --> ROpen: I-Peer-Disc / I-Disc, R-Snd-CEA
// WaitReturns --> IOpen: I-Rcv-CEA / R-Disc
// WaitReturns --> WaitICEA: R-Peer-Disc / R-Disc
// WaitReturns --> WaitReturns: R-Conn-CER / R-Reject
// WaitReturns --> Closed: Timeout / Error
// ROpen --> ROpen: Send-Message / R-Snd-Message
// ROpen --> ROpen: R-Rcv-Message / Process
// ROpen --> ROpen: R-Rcv-DWR / Process-DWR, R-Snd-DWA
// ROpen --> ROpen: R-Rcv-DWA / Process-DWA
// ROpen --> ROpen: R-Conn-CER / R-Reject
// ROpen --> Closing: Stop / R-Snd-DPR
// ROpen --> Closing: R-Rcv-DPR / R-Snd-DPA
// ROpen --> Closed: R-Peer-Disc / R-Disc
// IOpen --> IOpen: Send-Message / I-Snd-Message
// IOpen --> IOpen: I-Rcv-Message / Process
// IOpen --> IOpen: I-Rcv-DWR / Process-DWR, I-Snd-DWA
// IOpen --> IOpen: I-Rcv-DWA / Process-DWA
// IOpen --> IOpen: R-Conn-CER / R-Reject
// IOpen --> Closing: Stop / I-Snd-DPR
// IOpen --> Closing: I-Rcv-DPR / I-Snd-DPA
// IOpen --> Closed: I-Peer-Disc / I-Disc
// Closing --> Closed: I-Rcv-DPA / I-Disc
// Closing --> Closed: R-Rcv-DPA / R-Disc
// Closing --> Closed: Timeout / Error
// Closing --> Closed: I-Peer-Disc / I-Disc
// Closing --> Closed: R-Peer-Disc / R-Disc

func NewDiameterFSM() *FSM[message.DiameterMessage] {
	// Initial State (Closed)
	fsm := NewFSM[message.DiameterMessage](Closed)
	fsm.RegisterState(WaitConnectionAck)
	fsm.RegisterState(WaitICEA)
	fsm.RegisterState(Elect)
	fsm.RegisterState(WaitReturns)
	fsm.RegisterState(ROpen)
	fsm.RegisterState(IOpen)
	fsm.RegisterState(Closing)

	// Starts from Closed state.
	// If an I-Snd-Conn-Req event occurs, it moves to Wait-Conn-Ack.
	// If a R-Conn-CER event occurs (incoming connection), it transitions to R-Open.
	fsm.AddTransition(Closed, WaitICEA, ISendConnReq, []Action[message.DiameterMessage]{
		SendConnReq,
		SendDiameterMessage,
	})
	fsm.AddTransition(Closed, ROpen, RConnCER, []Action[message.DiameterMessage]{
		AcceptConn,
		ProcessCER,
		SendCEA,
		SendDiameterMessage,
	})

	// Wait-Conn-Ack State:
	// Awaits acknowledgment after initiating a connection.
	// On receiving I-Rcv-Conn-Ack, it transitions to Wait-I-CEA.
	// If a timeout occurs, it returns to Closed.
	fsm.AddTransition(WaitConnectionAck, WaitICEA, RcvConnAck, []Action[message.DiameterMessage]{
		SendCEA,
		SendDiameterMessage,
	})
	fsm.AddTransition(WaitConnectionAck, Closed, Timeout, []Action[message.DiameterMessage]{DiameterError})

	// Wait-I-CEA State:
	// Awaits peer connection response.
	// Upon I-Rcv-CEA, it transitions to I-Open.
	// Errors or disconnections transition back to Closed.
	fsm.AddTransition(WaitICEA, IOpen, RcvCEA, []Action[message.DiameterMessage]{ProcessCEA})
	fsm.AddTransition(WaitICEA, Closed, PeerDisc, []Action[message.DiameterMessage]{Disconnect})
	fsm.AddTransition(WaitICEA, Closed, DError, []Action[message.DiameterMessage]{DiameterError})

	// I-Open / R-Open:
	// In I-Open, it can send or receive messages, and handle disconnections or peer requests.
	// In R-Open, similar operations occur for responder scenarios.
	// Both move to Closing when a stop event happens.
	fsm.AddTransition(IOpen, Closing, Stop, []Action[message.DiameterMessage]{Cleanup})
	fsm.AddTransition(ROpen, Closing, Stop, []Action[message.DiameterMessage]{Cleanup})

	// Closing State:
	// Awaits disconnection confirmation (DPA).
	// Transitions to Closed on acknowledgment or timeout.
	fsm.AddTransition(Closing, Closed, RcvDPA, []Action[message.DiameterMessage]{Cleanup})
	fsm.AddTransition(Closing, Closed, Timeout, []Action[message.DiameterMessage]{Cleanup})

	// Election Phase (Wait-Conn-Ack/Elect & Wait-Returns):
	// If multiple connection attempts occur, elections decide the controlling node.
	// The winner transitions to R-Open, and the losing node goes back to waiting or closes.
	// TODO: Implement election transitions.
	return fsm
}
