package transport

import "errors"

// ErrAcceptTimeout is returned when the Accept timeout is reached for SCTP.
var (
	ErrAcceptTimeout    = errors.New("accept timeout reached")
	UnsupportedProtocol = errors.New("unsupported protocol")
)
