package message

import (
	"fmt"
	"net"
)

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

type OctetString struct {
	Data       []byte
	min_length uint32
}

func (o *OctetString) Length() uint32 {
	return uint32(len(o.Data))
}

// OctetString
//
//	The data contains arbitrary data of variable length.  Unless
//	otherwise noted, the AVP Length field MUST be set to at least 8
//	(12 if the 'V' bit is enabled).  AVP Values of this type that are
//	not a multiple of four-octets in length is followed by the
//	necessary padding so that the next AVP (if any) will start on a
//	32-bit boundary.
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

type Integer32 struct {
	Data int32
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

type Integer64 struct {
	Data int64
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

type Unsigned32 struct {
	Data uint32
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

type Unsigned64 struct {
	Data uint64
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

type Grouped struct {
	AVPs []*AVP
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
type IPAddr struct {
	Data   net.IP
	isIPv4 bool
}

func (i *IPAddr) Length() uint32 {
	if i.isIPv4 {
		return IPAddressTypeLength + IPv4AddressLength
	}
	return IPAddressTypeLength + IPv6AddressLength
}

func (i *IPAddr) Encode() ([]byte, error) {

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

func (i *IPAddr) Decode(data []byte) error {
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

func (i *IPAddr) String() string {
	return i.Data.String()
}
