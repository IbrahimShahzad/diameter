package message

import (
	"fmt"
	"net"
	"time"
	"unicode/utf8"
)

// Derived Type
//
//	Address
//	     The Address format is derived from the OctetString AVP Base
//	     Format.  It is a discriminated union, representing, for example a
//	     32-bit (IPv4) [IPV4] or 128-bit (IPv6) [IPV6] address, most
//	     significant octet first.  The first two octets of the Address
//	     AVP represents the AddressType, which contains an Address Family
//	     defined in [IANAADFAM].  The AddressType is used to discriminate
//	     the content and format of the remaining octets.
type Address struct {
	OctetString
	Data   net.IP
	isIPv4 bool
}

func (a *Address) SetData(data interface{}) error {
	switch v := data.(type) {
	case net.IP:
		a.Data = v
		a.isIPv4 = a.Data.To4() != nil
		return nil
	}
	return fmt.Errorf("invalid data type: %T", data)
}

func (a *Address) Length() uint32 {
	if a.isIPv4 {
		return IPAddressTypeLength + IPv4AddressLength
	}
	return IPAddressTypeLength + IPv6AddressLength
}

func (a *Address) Encode() ([]byte, error) {
	length := a.Length()

	buffer := make([]byte, length)
	if a.isIPv4 {
		ip := a.Data.To4()
		if ip == nil {
			return nil, InvalidIPv4AddressError
		}
		buffer[0] = 0
		buffer[1] = AddressFamilyIPv4Byte
		copy(buffer[IPAddressTypeLength:], ip)
	} else {
		ip := a.Data.To16()
		if ip == nil {
			return nil, InvalidIPv6AddressError
		}
		buffer[0] = 0
		buffer[1] = AddressFamilyIPv6Byte
		copy(buffer[IPAddressTypeLength:], ip)
	}
	return append(buffer, make([]byte, getPadding(len(buffer)))...), nil
}

func (a *Address) Decode(data []byte) error {
	// check for ip type in first 2 bytes + length of ip4 address
	if len(data) < IPAddressTypeLength {
		return InvalidAddressLengthError
	}

	addressFamily := (uint16(data[0]) << 8) | uint16(data[1])
	switch addressFamily {
	case uint16(AddressFamilyIPv4Byte):
		if len(data) != IPAddressTypeLength+IPv4AddressLength {
			return InvalidIPv4AddressLengthError
		}
		a.isIPv4 = true
		a.Data = net.IP(data[IPAddressTypeLength : IPAddressTypeLength+IPv4AddressLength])
	case uint16(AddressFamilyIPv6Byte):
		if len(data) != IPAddressTypeLength+IPv6AddressLength {
			return InvalidIPv6AddressLengthError
		}
		a.isIPv4 = false
		a.Data = net.IP(data[IPAddressTypeLength : IPAddressTypeLength+IPv6AddressLength])
	default:
		return UnknownAddressTypeError
	}

	return nil
}

func (a *Address) String() string {
	return a.Data.String()
}

func (a *Address) Type() AVPType {
	return AddressType
}

// UTF8String
//
//	The UTF8String format is derived from the OctetString Basic AVP
//	Format.  This is a human-readable string represented using the
//	ISO/IEC IS 10646-1 character set, encoded as an OctetString using
//	the UTF-8 transformation format [RFC3629].
//
//	Since additional code points are added by amendments to the 10646
//	standard from time to time, implementations MUST be prepared to
//	encounter any code point from 0x00000001 to 0x7fffffff.  Byte
//	sequences that do not correspond to the valid encoding of a code
//	point into UTF-8 charset or are outside this range are prohibited.
//
//	The use of control codes SHOULD be avoided.  When it is necessary
//	to represent a new line, the control code sequence CR LF SHOULD be
//	used.
//
//	The use of leading or trailing white space SHOULD be avoided.
//
//	For code points not directly supported by user interface hardware
//	or software, an alternative means of entry and display, such as
//	hexadecimal, MAY be provided.
//
//	For information encoded in 7-bit US-ASCII, the UTF-8 charset is
//	identical to the US-ASCII charset.
//
//	UTF-8 may require multiple bytes to represent a single character /
//	code point; thus, the length of a UTF8String in octets may be
//	different from the number of characters encoded.
//
//	Note that the AVP Length field of an UTF8String is measured in
//	octets not characters.
type UTF8String struct {
	OctetString
}

func (u *UTF8String) SetData(data interface{}) error {
	if d, ok := data.(string); ok {
		if !isValidUTF8([]byte(d)) {
			return InvalidUTF8StringError
		}
		u.Data = []byte(d)
		return nil
	}
	return fmt.Errorf("invalid data type: %T", data)
}

func (u *UTF8String) Length() uint32 {
	return uint32(len(u.Data))
}

func (u *UTF8String) Encode() ([]byte, error) {
	return []byte(u.Data), nil
}

func (u *UTF8String) Decode(data []byte) error {
	if !isValidUTF8(data) {
		return InvalidUTF8StringError
	}
	return u.OctetString.Decode(data)
}

func (u *UTF8String) String() string {
	return u.OctetString.String()
}

func (u *UTF8String) Type() AVPType {
	return UTF8StringType
}

func isValidUTF8(data []byte) bool {
	return utf8.Valid(data)
}

// Enumerated
//
//	The Enumerated format is derived from the Integer32 Basic AVP
//	Format.  The definition contains a list of valid values and their
//	interpretation and is described in the Diameter application
//	introducing the AVP.
type Enumerated struct {
	Data uint32
}

func (e *Enumerated) SetData(data interface{}) error {
	if d, ok := data.(uint32); ok {
		e.Data = d
		return nil
	}
	return fmt.Errorf("invalid data type: %T", data)
}

func (e *Enumerated) Length() uint32 {
	return int32Length
}

func (e *Enumerated) Encode() ([]byte, error) {
	return encode32(e.Data), nil
}

func (e *Enumerated) Decode(data []byte) error {
	e.Data = decode32(data, uint32(0))
	return nil
}

func (e *Enumerated) String() string {
	return fmt.Sprintf("%d", e.Data)
}

func (e *Enumerated) Type() AVPType {
	return EnumeratedType
}

// Time
//
//	The Time format is derived from the OctetString Basic AVP Format.
//	The string MUST contain four octets, in the same format as the
//	first four bytes are in the NTP timestamp format.  The NTP
//	timestamp format is defined in Section 3 of [RFC5905].
//
//	This represents the number of seconds since 0h on 1 January 1900
//	with respect to the Coordinated Universal Time (UTC).
//
//	On 6h 28m 16s UTC, 7 February 2036, the time value will overflow.
//	Simple Network Time Protocol (SNTP) [RFC5905] describes a
//	procedure to extend the time to 2104.  This procedure MUST be
//	supported by all Diameter nodes.
type Time struct {
	OctetString
}

func (t *Time) SetData(data interface{}) error {
	timeOffset := 2208988800 // Offset for 1 Jan 1900
	switch v := data.(type) {
	case uint32:
		// Directly use NTP seconds
		t.Data = encode32(v)
	case time.Time:
		// Convert from time.Time to NTP seconds
		ntpSeconds := uint32(v.Unix() + int64(timeOffset))
		t.Data = encode32(ntpSeconds)
	case int64, int:
		// Treat as Epoch seconds and convert to NTP
		epochSeconds := int64Value(v)
		ntpSeconds := uint32(epochSeconds + int64(timeOffset))
		t.Data = encode32(ntpSeconds)
	case string:
		// Parse string as a date
		parsedTime, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return fmt.Errorf("invalid date string: %v", err)
		}
		ntpSeconds := uint32(parsedTime.Unix() + int64(timeOffset))
		t.Data = encode32(ntpSeconds)
	default:
		return fmt.Errorf("unsupported data type: %T", v)
	}
	return nil

}

func int64Value(data interface{}) int64 {
	switch v := data.(type) {
	case int:
		return int64(v)
	case int64:
		return v
	default:
		return 0
	}
}

func (t *Time) Length() uint32 {
	return int32Length
}

func (t *Time) Encode() ([]byte, error) {
	if len(t.Data) != int32Length {
		return nil, fmt.Errorf("invalid data length: expected 4, got %d", len(t.Data))
	}
	return t.Data, nil
}

func (t *Time) Decode(data []byte) error {
	if len(data) != int32Length {
		return fmt.Errorf("invalid data length: expected 4, got %d", len(data))
	}
	t.Data = data
	return nil
}

func (t *Time) String() string {
	if len(t.Data) != int32Length {
		return "invalid time format"
	}
	// Convert 4 bytes back to uint32
	seconds := uint32(t.Data[0])<<24 | uint32(t.Data[1])<<16 | uint32(t.Data[2])<<8 | uint32(t.Data[3])
	return fmt.Sprintf("%d", seconds)
}

func (t *Time) Type() AVPType {
	return TimeType
}

// DiameterIdentity
//
//	The DiameterIdentity format is derived from the OctetString Basic
//	AVP Format.
//
//	                  DiameterIdentity  = FQDN/Realm
//
// The DiameterIdentity value is used to uniquely identify either:
//
//   - A Diameter node for purposes of duplicate connection and
//     routing loop detection.
//
//   - A Realm to determine whether messages can be satisfied locally
//     or whether they must be routed or redirected.
//
//     When a DiameterIdentity value is used to identify a Diameter node,
//     the contents of the string MUST be the Fully Qualified Domain Name
//     (FQDN) of the Diameter node.  If multiple Diameter nodes run on
//     the same host, each Diameter node MUST be assigned a unique
//     DiameterIdentity.  If a Diameter node can be identified by several
//     FQDNs, a single FQDN should be picked at startup and used as the
//     only DiameterIdentity for that node, whatever the connection on
//     which it is sent.  In this document, note that DiameterIdentity is
//     in ASCII form in order to be compatible with existing DNS
//     infrastructure.  See Appendix D for interactions between the
//     Diameter protocol and Internationalized Domain Names (IDNs).
type DiameterIdentity struct {
	OctetString
}

//	func (i *DiameterIdentity) SetData(data interface{}) error {
//		if d, ok := data.(string); ok {
//			i.Data = d
//			return nil
//		}
//		return fmt.Errorf("invalid data type: %T", data)
//	}
func SetData(i *DiameterIdentity, data interface{}) error {
	if d, ok := data.(string); ok {
		return i.OctetString.SetData([]byte(d))
	}
	return fmt.Errorf("invalid data type: %T", data)
}

func (d *DiameterIdentity) Length() uint32 {
	return uint32(len(d.Data))
}

func (d *DiameterIdentity) Encode() ([]byte, error) {
	return d.Data, nil
}

func (d *DiameterIdentity) Decode(data []byte) error {
	d.Data = data
	return nil
}

func (d *DiameterIdentity) String() string {
	return string(d.Data)
}

func (d *DiameterIdentity) Type() AVPType {
	return DiameterIdentityType
}

type AppId struct {
	Data uint32
}

func (a *AppId) SetData(data interface{}) error {
	if d, ok := data.(uint32); ok {
		a.Data = d
		return nil
	}
	return fmt.Errorf("invalid data type: %T", data)
}

func (a *AppId) Length() uint32 {
	return int32Length
}

func (a *AppId) Encode() ([]byte, error) {
	return encode32(a.Data), nil
}

func (a *AppId) Decode(data []byte) error {
	a.Data = decode32(data, uint32(0))
	return nil
}

func (a *AppId) String() string {
	return fmt.Sprintf("%d", a.Data)
}

func (a *AppId) Type() AVPType {
	return AppIdType
}

type VendorId struct {
	Data uint32
}

func (v *VendorId) SetData(data interface{}) error {
	if d, ok := data.(uint32); ok {
		v.Data = d
		return nil
	}
	return fmt.Errorf("invalid data type: %T", data)
}

func (v *VendorId) Length() uint32 {
	return int32Length
}

func (v *VendorId) Encode() ([]byte, error) {
	return encode32(v.Data), nil
}

func (v *VendorId) Decode(data []byte) error {
	v.Data = decode32(data, uint32(0))
	return nil
}

func (v *VendorId) String() string {
	return fmt.Sprintf("%d", v.Data)
}

func (v *VendorId) Type() AVPType {
	return VendorIdType
}

// DiameterURI
//
//	The DiameterURI MUST follow the Uniform Resource Identifiers (RFC
//	3986) syntax [RFC3986] rules specified below:
//
//	"aaa://" FQDN [ port ] [ transport ] [ protocol ]
//
//	                ; No transport security
//
//	"aaas://" FQDN [ port ] [ transport ] [ protocol ]
//
//	                ; Transport security used
//
//	FQDN               = < Fully Qualified Domain Name >
//
//	port               = ":" 1*DIGIT
//
//	                ; One of the ports used to listen for
//	                ; incoming connections.
//	                ; If absent, the default Diameter port
//	                ; (3868) is assumed if no transport
//	                ; security is used and port 5658 when
//	                ; transport security (TLS/TCP and DTLS/SCTP)
//	                ; is used.
//
//	transport          = ";transport=" transport-protocol
//
//	                ; One of the transports used to listen
//	                ; for incoming connections.  If absent,
//	                ; the default protocol is assumed to be TCP.
//	                ; UDP MUST NOT be used when the aaa-protocol
//	                ; field is set to diameter.
//
//	transport-protocol = ( "tcp" / "sctp" / "udp" )
//
//	protocol           = ";protocol=" aaa-protocol
//
//	                ; If absent, the default AAA protocol
//	                ; is Diameter.
//
//	aaa-protocol       = ( "diameter" / "radius" / "tacacs+" )
//
//	The following are examples of valid Diameter host identities:
//
//	aaa://host.example.com;transport=tcp
//	aaa://host.example.com:6666;transport=tcp
//	aaa://host.example.com;protocol=diameter
//	aaa://host.example.com:6666;protocol=diameter
//	aaa://host.example.com:6666;transport=tcp;protocol=diameter
//	aaa://host.example.com:1813;transport=udp;protocol=radius
type DiameterURI struct {
	Data string
}

func (i *DiameterURI) SetData(data interface{}) error {
	if d, ok := data.(string); ok {
		i.Data = d
		return nil
	}
	return fmt.Errorf("invalid data type: %T", data)
}

func (d *DiameterURI) Length() uint32 {
	return uint32(len(d.Data))
}

func (d *DiameterURI) Encode() ([]byte, error) {
	return []byte(d.Data), nil
}

func (d *DiameterURI) Decode(data []byte) error {
	d.Data = string(data)
	return nil
}

func (d *DiameterURI) String() string {
	return d.Data
}

func (d *DiameterURI) Type() AVPType {
	return DiameterURIType
}

// IPFilterRule
//
//	The IPFilterRule format is derived from the OctetString Basic AVP
//	Format and uses the ASCII charset.  The rule syntax is a modified
//	subset of ipfw(8) from FreeBSD.  Packets may be filtered based on
//	the following information that is associated with it:
//
//	      Direction                          (in or out)
//	      Source and destination IP address  (possibly masked)
//	      Protocol
//	      Source and destination port        (lists or ranges)
//	      TCP flags
//	      IP fragment flag
//	      IP options
//	      ICMP types
//
// see rfc6733 for more details on the format
type IPFilterRule struct {
	Data uint32
}

func (i *IPFilterRule) SetData(data interface{}) error {
	if d, ok := data.(uint32); ok {
		i.Data = d
		return nil
	}
	return fmt.Errorf("invalid data type: %T", data)
}

func (i *IPFilterRule) Length() uint32 {
	return int32Length
}

func (i *IPFilterRule) Encode() ([]byte, error) {
	return encode32(i.Data), nil
}

func (i *IPFilterRule) Decode(data []byte) error {
	i.Data = decode32(data, uint32(0))
	return nil
}

func (i *IPFilterRule) String() string {
	return fmt.Sprintf("%d", i.Data)
}

func (i *IPFilterRule) Type() AVPType {
	return IPFilterRuleType
}
