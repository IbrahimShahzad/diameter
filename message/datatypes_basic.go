package message

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
)

// https://datatracker.ietf.org/doc/html/rfc6733
type AVPType uint32

const (
	OctetStringType AVPType = iota
	Integer32Type
	Integer64Type
	Unsigned32Type
	Unsigned64Type
	Float32Type
	Float64Type
	GroupedType
	AddressType
	UTF8StringType
	EnumeratedType
	TimeType
	DiameterIdentityType
	DiameterURIType
	IPFilterRuleType
	AppIdType
	VendorIdType
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

func encode32[T uint32 | int32](data T) []byte {
	buffer := make([]byte, int32Length)
	for i := 0; i < int32Length; i++ {
		buffer[i] = byte((data >> uint(i*bitsInByte)) & 0xFF)
	}
	return buffer
}

func encode64[T uint64 | int64](data T) []byte {
	buffer := make([]byte, int64Length)
	for i := 0; i < int64Length; i++ {
		buffer[i] = byte((data >> uint(i*bitsInByte)) & 0xFF)
	}
	return buffer
}

func decode32[T uint32 | int32](data []byte, t T) T {
	for i := 0; i < int32Length; i++ {
		t = t | T(data[i])<<uint(bitsInByte*i)
	}
	return t
}

func decode64[T uint64 | int64](data []byte, t T) T {
	for i := 0; i < int64Length; i++ {
		t = t | T(data[i])<<uint(bitsInByte*i)
	}
	return t
}

// Returns whether the data type is derived from the OctetString type
func IsDerivedFromOctetString(data interface{}) bool {
	if data == nil {
		return false
	}

	t := reflect.TypeOf(data)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Type == reflect.TypeOf(OctetString{}) {
			return true
		}
	}

	return false
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

	if d, ok := data.(string); ok {
		o.Data = []byte(d)
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

func (o *OctetString) Type() AVPType {
	return OctetStringType
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
	return encode32(u.Data), nil
}

func (u *Integer32) Decode(data []byte) error {
	u.Data = decode32(data, int32(0))
	return nil
}

func (i *Integer32) String() string {
	return fmt.Sprintf("%d", i.Data)
}

func (i *Integer32) Type() AVPType {
	return Integer32Type
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
	return encode64(u.Data), nil
}

func (u *Integer64) Decode(data []byte) error {
	u.Data = decode64(data, int64(0))
	return nil
}

func (i *Integer64) Length() uint32 {
	return int64Length
}

func (i *Integer64) String() string {
	return fmt.Sprintf("%d", i.Data)
}

func (i *Integer64) Type() AVPType {
	return Integer64Type
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
	return encode32(u.Data), nil
}

func (u *Unsigned32) Length() uint32 {
	return int32Length
}

func (u *Unsigned32) Decode(data []byte) error {
	u.Data = decode32(data, uint32(0))
	return nil
}

func (u *Unsigned32) String() string {
	return fmt.Sprintf("%d", u.Data)
}

func (u *Unsigned32) Type() AVPType {
	return Unsigned32Type
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
	return encode64(u.Data), nil
}

func (u *Unsigned64) Decode(data []byte) error {
	u.Data = decode64(data, uint64(0))
	return nil
}

func (u *Unsigned64) Length() uint32 {
	return int64Length
}

func (u *Unsigned64) String() string {
	return fmt.Sprintf("%d", u.Data)
}

func (u *Unsigned64) Type() AVPType {
	return Unsigned64Type
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

func (f *Float32) Type() AVPType {
	return Float32Type
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

func (f *Float64) Type() AVPType {
	return Float64Type
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
		avp, err := DecodeAVP(data[offset:])
		if err != nil {
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

func (g *Grouped) Type() AVPType {
	return GroupedType
}
