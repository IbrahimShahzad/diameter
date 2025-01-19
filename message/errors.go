package message

import "errors"

// Diameter errors
var (
	InvalidDiameterVersionError      = errors.New("invalid version")
	InvalidDiameterHeaderLengthError = errors.New("invalid header length")
)

// datatype errors
var (
	UnsupportedTypeError = errors.New("unsupported type")
)

// AVP errors
var (
	VendorIDRequiredError = errors.New("VendorID is required for vendor specific AVP")
)

// IPAddr errors
var (
	InvalidIPv4AddressError       = errors.New("invalid IPv4 address")
	InvalidIPv6AddressError       = errors.New("invalid IPv6 address")
	InvalidIPv4AddressLengthError = errors.New("invalid IPv4 address length")
	InvalidIPv6AddressLengthError = errors.New("invalid IPv6 address length")
	UnknownAddressTypeError       = errors.New("unknown address type")
	InvalidAddressLengthError     = errors.New("invalid address length")
)

// Decoding errors
var (
	InvalidMessageLengthError = errors.New("invalid message length for decoding")
)

var (
	InvalidCommandCodeError = errors.New("invalid command code")
)

var (
	InvalidUTF8StringError = errors.New("invalid UTF-8 string")
)
