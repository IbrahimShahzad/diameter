package message

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
)

// TODO: Fix the encoding decoding functions for the derived types

// https://datatracker.ietf.org/doc/html/rfc6733

const (
	int32Length = 4
	int64Length = 8
	bitsInByte  = 8
)

const (
	IPAddressTypeLength = 2
	IPv4AddressLength   = 4
	IPv6AddressLength   = 16
)

const (
	AddressFamilyIPv4Byte = byte(0x01) // 0x01 for IPv4
	AddressFamilyIPv6Byte = byte(0x02) // 0x02 for IPv6
)

func encode32[T uint32 | int32](data T) ([]byte, error) {
	buffer := make([]byte, int32Length)
	for i := 0; i < int32Length; i++ {
		buffer[int32Length-1-i] = byte((data >> uint(i*bitsInByte)) & 0xFF)
	}
	return buffer, nil
}

func encode64[T uint64 | int64](data T) ([]byte, error) {
	buffer := make([]byte, int64Length)
	for i := 0; i < int64Length; i++ {
		buffer[int64Length-1-i] = byte((data >> uint(i*bitsInByte)) & 0xFF)
	}
	return buffer, nil
}

func decode32[T uint32 | int32](data []byte, t T) (T, error) {
	for i := 0; i < int32Length; i++ {
		t = t | T(data[i])<<uint(bitsInByte*i)
	}
	return t, nil
}

func decode64[T uint64 | int64](data []byte, t T) (T, error) {
	for i := 0; i < int64Length; i++ {
		t = t | T(data[i])<<uint(bitsInByte*i)
	}
	return t, nil
}

// Basic AVP Data Types

// OctetString
//
//	The data contains arbitrary data of variable length.  Unless
//	otherwise noted, the AVP Length field MUST be set to at least 8
//	(12 if the 'V' bit is enabled).  AVP values of this type that are
//	not a multiple of 4 octets in length are followed by the necessary
//	padding so that the next AVP (if any) will start on a 32-bit
//	boundary.
type OctetString struct {
	Data       []byte
	min_length uint32
}

func (o *OctetString) SetData(data interface{}) error {
	if d, ok := data.([]byte); ok {
		o.Data = d
		return nil
	}
	return fmt.Errorf("invalid data type: %T", data)
}

func (o *OctetString) Length() uint32 {
	return uint32(len(o.Data))
}

func (o *OctetString) Encode() ([]byte, error) {
	length := len(o.Data)
	if length < int(o.min_length) {
		length = int(o.min_length)
	}
	buffer := make([]byte, length+getPadding(length))
	copy(buffer, o.Data)
	return buffer, nil
}

func (o *OctetString) Decode(data []byte) error {
	o.Data = data
	return nil
}

func (o *OctetString) String() string {
	return string(o.Data)
}

// Integer32
//
//	32-bit signed value, in network byte order.  The AVP Length field
//	MUST be set to 12 (16 if the 'V' bit is enabled).
type Integer32 struct {
	Data int32
}

func (i *Integer32) SetData(data interface{}) error {
	if d, ok := data.(int32); ok {
		i.Data = d
		return nil
	}
	return fmt.Errorf("invalid data type: %T", data)
}

func (i *Integer32) Length() uint32 {
	return int32Length
}

func (u *Integer32) Encode() ([]byte, error) {
	return encode32(u.Data)
}

func (u *Integer32) Decode(data []byte) error {
	var err error
	u.Data, err = decode32(data, int32(0))
	return err
}

func (i *Integer32) String() string {
	return fmt.Sprintf("%d", i.Data)
}

// Integer64
//
//	64-bit signed value, in network byte order.  The AVP Length field
//	MUST be set to 16 (20 if the 'V' bit is enabled).
type Integer64 struct {
	Data int64
}

func (i *Integer64) SetData(data interface{}) error {
	if d, ok := data.(int64); ok {
		i.Data = d
		return nil
	}
	return fmt.Errorf("invalid data type: %T", data)
}

func (u *Integer64) Encode() ([]byte, error) {
	return encode64(u.Data)
}

func (u *Integer64) Decode(data []byte) error {
	var err error
	u.Data, err = decode64(data, int64(0))
	return err
}

func (i *Integer64) Length() uint32 {
	return int64Length
}

func (i *Integer64) String() string {
	return fmt.Sprintf("%d", i.Data)
}

// Unsigned32
//
//	32-bit unsigned value, in network byte order.  The AVP Length
//	field MUST be set to 12 (16 if the 'V' bit is enabled).
type Unsigned32 struct {
	Data uint32
}

func (u *Unsigned32) SetData(data interface{}) error {
	if d, ok := data.(uint32); ok {
		u.Data = d
		return nil
	}
	return fmt.Errorf("invalid data type: %T", data)
}

func (u *Unsigned32) Encode() ([]byte, error) {
	return encode32(u.Data)
}

func (u *Unsigned32) Length() uint32 {
	return int32Length
}

func (u *Unsigned32) Decode(data []byte) error {
	var err error
	u.Data, err = decode32(data, uint32(0))
	return err
}

func (u *Unsigned32) String() string {
	return fmt.Sprintf("%d", u.Data)
}

// Unsigned64
//
//	64-bit unsigned value, in network byte order.  The AVP Length
//	field MUST be set to 16 (20 if the 'V' bit is enabled).
type Unsigned64 struct {
	Data uint64
}

func (u *Unsigned64) SetData(data interface{}) error {
	if d, ok := data.(uint64); ok {
		u.Data = d
		return nil
	}
	return fmt.Errorf("invalid data type: %T", data)
}

func (u *Unsigned64) Encode() ([]byte, error) {
	return encode64(u.Data)
}

func (u *Unsigned64) Decode(data []byte) error {
	var err error
	u.Data, err = decode64(data, uint64(0))
	return err
}

func (u *Unsigned64) Length() uint32 {
	return int64Length
}

func (u *Unsigned64) String() string {
	return fmt.Sprintf("%d", u.Data)
}

// Float32
//
//	This represents floating point values of single precision as
//	described by [FLOATPOINT].  The 32-bit value is transmitted in
//	network byte order.  The AVP Length field MUST be set to 12 (16 if
//	the 'V' bit is enabled).
type Float32 struct {
	Data float32
}

func (f *Float32) SetData(data interface{}) error {
	if d, ok := data.(float32); ok {
		f.Data = d
		return nil
	}
	return fmt.Errorf("invalid data type: %T", data)
}

func (f *Float32) Length() uint32 {
	return int32Length
}

func (f *Float32) Encode() ([]byte, error) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, math.Float32bits(f.Data))
	return buf, nil
}

func (f *Float32) Decode(data []byte) error {
	if len(data) != 4 {
		return fmt.Errorf("invalid data length: %d", len(data))
	}
	f.Data = math.Float32frombits(binary.BigEndian.Uint32(data))
	return nil
}

func (f *Float32) String() string {
	return fmt.Sprintf("%f", f.Data)
}

// Float64
//
//	This represents floating point values of double precision as
//	described by [FLOATPOINT].  The 64-bit value is transmitted in
//	network byte order.  The AVP Length field MUST be set to 16 (20 if
//	the 'V' bit is enabled).
type Float64 struct {
	Data float64
}

func (f *Float64) SetData(data interface{}) error {
	if d, ok := data.(float64); ok {
		f.Data = d
		return nil
	}
	return fmt.Errorf("invalid data type: %T", data)
}

func (f *Float64) Length() uint32 {
	return int64Length
}

func (f *Float64) Encode() ([]byte, error) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, math.Float64bits(f.Data))
	return buf, nil
}

func (f *Float64) Decode(data []byte) error {
	if len(data) != 8 {
		return fmt.Errorf("invalid data length: %d", len(data))
	}
	f.Data = math.Float64frombits(binary.BigEndian.Uint64(data))
	return nil
}

func (f *Float64) String() string {
	return fmt.Sprintf("%f", f.Data)
}

// Grouped
//
//	The Data field is specified as a sequence of AVPs.  These AVPs are
//	concatenated -- including their headers and padding -- in the
//	order in which they are specified and the result encapsulated in
//	the Data field.  The AVP Length field is set to 8 (12 if the 'V'
//	bit is enabled) plus the total length of all included AVPs,
//	including their headers and padding.  Thus, the AVP Length field
//	of an AVP of type Grouped is always a multiple of 4.
type Grouped struct {
	AVPs []*AVP
}

func (g *Grouped) SetData(data interface{}) error {
	if d, ok := data.([]*AVP); ok {
		g.AVPs = d
		return nil
	}
	return fmt.Errorf("invalid data type: %T", data)
}

func (g *Grouped) Length() uint32 {
	length := uint32(0)
	for _, avp := range g.AVPs {
		length += avp.Length()
	}
	return length
}

func (g *Grouped) Encode() ([]byte, error) {
	buffer := make([]byte, 0)
	for _, avp := range g.AVPs {
		encoded, err := avp.Encode()
		if err != nil {
			return nil, err
		}
		buffer = append(buffer, encoded...)
	}
	return buffer, nil
}

func (g *Grouped) Decode(data []byte) error {
	offset := 0
	for offset < len(data) {
		avp := &AVP{}
		if err := avp.Decode(data[offset:]); err != nil {
			return err
		}
		g.AVPs = append(g.AVPs, avp)
		offset += int(avp.AVPlength)
	}
	return nil
}

func (g *Grouped) String() string {
	var str string
	for _, avp := range g.AVPs {
		str += avp.String() + "\n"
	}
	return str
}

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
	Data   net.IP
	isIPv4 bool
}

func (i *Address) SetData(data interface{}) error {
	if d, ok := data.(net.IP); ok {
		i.Data = d
		i.isIPv4 = i.Data.To4() != nil
		return nil
	}
	return fmt.Errorf("invalid data type: %T", data)
}

func (i *Address) Length() uint32 {
	if i.isIPv4 {
		return IPAddressTypeLength + IPv4AddressLength
	}
	return IPAddressTypeLength + IPv6AddressLength
}

func (i *Address) Encode() ([]byte, error) {

	buffer := make([]byte, i.Length())
	if i.isIPv4 {
		ip := i.Data.To4()
		if ip == nil {
			return nil, InvalidIPv4AddressError
		}
		copy(buffer, []byte{AddressFamilyIPv4Byte})
		copy(buffer[AddressFamilyIPv4Byte:], ip)
	} else {
		ip := i.Data.To16()
		if ip == nil {
			return nil, InvalidIPv6AddressError
		}
		copy(buffer, []byte{AddressFamilyIPv6Byte})
		copy(buffer[AddressFamilyIPv6Byte:], ip)
	}
	return append(buffer, make([]byte, getPadding(len(buffer)))...), nil
}

func (i *Address) Decode(data []byte) error {
	// check for ip type in first 2 bytes + length of ip4 address
	if len(data) < IPAddressTypeLength {
		return InvalidAddressLengthError
	}

	if data[0] == 0 && data[1] == 1 {
		if len(data) != IPAddressTypeLength+IPv4AddressLength {
			return InvalidIPv4AddressLengthError
		}
		i.isIPv4 = true
		i.Data = net.IP(data[IPAddressTypeLength : IPAddressTypeLength+IPv4AddressLength])
	} else if data[0] == 0 && data[1] == 2 {
		if len(data) != IPAddressTypeLength+IPv6AddressLength {
			return InvalidIPv6AddressLengthError
		}
		i.isIPv4 = false
		i.Data = net.IP(data[IPAddressTypeLength : IPAddressTypeLength+IPv6AddressLength])
	} else {
		return UnknownAddressTypeError
	}
	return nil
}

func (i *Address) String() string {
	return i.Data.String()
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
	Data string
}

func (u *UTF8String) SetData(data interface{}) error {
	if d, ok := data.(string); ok {
		u.Data = d
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
	u.Data = string(data)
	return nil
}

func (u *UTF8String) String() string {
	return u.Data
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
	return encode32(e.Data)
}

func (e *Enumerated) Decode(data []byte) error {
	var err error
	e.Data, err = decode32(data, uint32(0))
	return err
}

func (e *Enumerated) String() string {
	return fmt.Sprintf("%d", e.Data)
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
	Data uint32
}

func (t *Time) SetData(data interface{}) error {
	if d, ok := data.(uint32); ok {
		t.Data = d
		return nil
	}
	return fmt.Errorf("invalid data type: %T", data)
}

func (t *Time) Length() uint32 {
	return int32Length
}

func (t *Time) Encode() ([]byte, error) {
	return encode32(t.Data)
}

func (t *Time) Decode(data []byte) error {
	var err error
	t.Data, err = decode32(data, uint32(0))
	return err
}

func (t *Time) String() string {
	return fmt.Sprintf("%d", t.Data)
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
	Data string
}

func (i *DiameterIdentity) SetData(data interface{}) error {
	if d, ok := data.(string); ok {
		i.Data = d
		return nil
	}
	return fmt.Errorf("invalid data type: %T", data)
}

func (d *DiameterIdentity) Length() uint32 {
	return uint32(len(d.Data))
}

func (d *DiameterIdentity) Encode() ([]byte, error) {
	return []byte(d.Data), nil
}

func (d *DiameterIdentity) Decode(data []byte) error {
	d.Data = string(data)
	return nil
}

func (d *DiameterIdentity) String() string {
	return d.Data
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
	return encode32(a.Data)
}

func (a *AppId) Decode(data []byte) error {
	var err error
	a.Data, err = decode32(data, uint32(0))
	return err
}

func (a *AppId) String() string {
	return fmt.Sprintf("%d", a.Data)
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
	return encode32(v.Data)
}

func (v *VendorId) Decode(data []byte) error {
	var err error
	v.Data, err = decode32(data, uint32(0))
	return err
}

func (v *VendorId) String() string {
	return fmt.Sprintf("%d", v.Data)
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
	return encode32(i.Data)
}

func (i *IPFilterRule) Decode(data []byte) error {
	var err error
	i.Data, err = decode32(data, uint32(0))
	return err
}

func (i *IPFilterRule) String() string {
	return fmt.Sprintf("%d", i.Data)
}
