// AVP struct, encoding/decoding functions
package message

import (
	"errors"
	"fmt"
	"net"

	"github.com/IbrahimShahzad/diameter/utils"
	"golang.org/x/exp/constraints"
)

const (
	AVPHeaderLength      = 8
	AVPHeaderLengthWithV = 12
)

// AVP Flags
const (
	MANDATORY_FLAG = 0x40
	VENDOR_FLAG    = 0x80
	PROTECTED_FLAG = 0x20
)

const (
	AVP_CODE_LENGTH        = 4
	AVP_FLAGS_LENGTH       = 1
	AVP_LENGTH_LENGTH      = 3
	AVP_VENDOR_ID_LENGTH   = 4
	AVP_UNPROTECTED_LENGTH = 8
	AVP_PROTECTED_LENGTH   = 12
	AVP_PROTECTION_LENGTH  = 4
)

// AVP represents a Diameter Attribute-Value Pair.
// Code 4 bytes
// Flags 1 byte:
// - Vendor specific bit - 1 bit
// - Mandatory bit - 1 bit
// - P bit - 1 bit
// - Reserved - 5 bits
// Length 3 bytes
// VendorID 4 bytes
type AVP struct {
	Code      uint32
	Flags     uint8
	AVPlength uint32 // Contains the length of the AVP header and data
	VendorID  uint32 // This is optional
	Data      AVPData
}

type AVPData interface {
	Encode() ([]byte, error)
	Decode(data []byte) error
	Length() uint32
	String() string
}

func (a *AVP) Length() uint32 {
	return a.AVPlength
}

func (a *AVP) String() string {
	return fmt.Sprintf(
		"AVP{Code: %d (%s), Flags: %d, Length: %d, VendorID: %d, Data: %s}",
		a.Code,
		GetAVPNameFromCode(a.Code),
		a.Flags,
		a.AVPlength,
		a.VendorID,
		a.Data.String(),
	)
}

func (a *AVP) getHeaderLength() int {
	if a.isFlagSet(VENDOR_FLAG) {
		return AVPHeaderLengthWithV
	}
	return AVPHeaderLength
}

func (a *AVP) Encode() ([]byte, error) {

	header := make([]byte, a.getHeaderLength())
	byteCount := 0
	copy(header, utils.ToBytes(a.Code, AVP_CODE_LENGTH))
	byteCount += AVP_CODE_LENGTH
	header[byteCount] = a.Flags
	byteCount += AVP_FLAGS_LENGTH
	copy(header[byteCount:], utils.ToBytes(a.AVPlength, AVP_LENGTH_LENGTH))
	byteCount += AVP_LENGTH_LENGTH
	if a.isFlagSet(VENDOR_FLAG) {
		copy(header[byteCount:], utils.ToBytes(a.VendorID, AVP_VENDOR_ID_LENGTH))
	}

	data, err := utils.Encode(a.Data)
	if err != nil {
		return nil, err
	}
	return append(header, data...), nil
}

func (a *AVP) Decode(data []byte) error {
	if len(data) < a.getHeaderLength() {
		return fmt.Errorf("AVP Decode: Insufficient data")
	}

	a.Code = utils.FromBytes(data[0:AVP_CODE_LENGTH])
	byteCount := AVP_CODE_LENGTH

	a.Flags = data[byteCount]
	byteCount += AVP_FLAGS_LENGTH

	a.AVPlength = utils.FromBytes(data[byteCount : byteCount+AVP_LENGTH_LENGTH])
	byteCount += AVP_LENGTH_LENGTH

	if a.isFlagSet(VENDOR_FLAG) {
		a.VendorID = utils.FromBytes(data[byteCount : byteCount+AVP_VENDOR_ID_LENGTH])
	}

	if len(data) < int(a.AVPlength) {
		return fmt.Errorf("AVP Decode: Insufficient data")
	}

	return utils.Decode(a.Data, data[a.getHeaderLength():])
}

func (a *AVP) setFlag(flag uint8) {
	a.Flags |= flag
}

func (a *AVP) isFlagSet(flag uint8) bool {
	return (a.Flags & flag) == flag
}

func (a *AVP) resetFlag(flag uint8) {
	a.Flags &= ^flag
}

// NewAVP creates a new AVP (Attribute-Value Pair) with the given parameters.
// The function supports various types for the value parameter, including net.IP, string, int32, int64, uint32, uint64, *AVP, and *Grouped.
// The function also handles flags for vendor-specific and protected AVPs.
//
// Parameters:
//
//	code: The AVP code.
//	value: The value of the AVP, which can be of various types.
//		-  net.IP for IPAddr AVP.
//		-  string for OctetString AVP.
//		-  int32 for Integer32 AVP.
//		-  int64 for Integer64 AVP.
//		-  uint32 for Unsigned32 AVP.
//		-  uint64 for Unsigned64 AVP.
//		-  *AVP for nested AVPs.
//		-  *Grouped for grouped AVPs.
//	flag: The AVP flags.
//	vendorID: Optional vendor ID(s) for vendor-specific AVPs.
//
// Returns:
//
//	A pointer to the newly created AVP and an error if the creation fails.
func NewAVP[T constraints.Ordered | net.IP](
	code uint32,
	value T,
	flag uint8,
	vendorID ...uint32,
) (*AVP, error) {
	headerLen := AVPHeaderLength
	if flag&VENDOR_FLAG != 0 {
		headerLen = AVPHeaderLengthWithV
	}
	var data AVPData
	switch v := any(value).(type) {
	case net.IP:
		data = &IPAddr{
			isIPv4: v.To4() != nil,
			Data:   v,
		}
	case string:
		var min_length uint32
		if flag&PROTECTED_FLAG != 0 {
			min_length = AVP_PROTECTED_LENGTH
		} else {
			min_length = AVP_UNPROTECTED_LENGTH
		}
		data = &OctetString{Data: []byte(v), min_length: min_length}
	case int32:
		data = &Integer32{Data: v}
	case int64:
		data = &Integer64{Data: v}
	case uint32:
		data = &Unsigned32{Data: v}
	case uint64:
		data = &Unsigned64{Data: v}
	case *AVP:
		data = v
	case *Grouped:
		data = v
	default:
		errors.Join(UnsupportedTypeError, fmt.Errorf("%T", value))
	}

	length := uint32(headerLen) + data.Length()

	// NOTE: Do i need to do this since i am already checking in the switch case
	// for the octet string
	if flag&PROTECTED_FLAG != 0 {
		length += AVP_PROTECTION_LENGTH
	}

	vID := uint32(0)
	if flag&VENDOR_FLAG != 0 {
		if len(vendorID) == 0 {
			return nil, VendorIDRequiredError
		}
		vID = vendorID[0]
		length += AVP_VENDOR_ID_LENGTH
	}

	return &AVP{
		Code:      code,
		Flags:     flag,
		AVPlength: length,
		VendorID:  vID,
		Data:      data,
	}, nil

}

func getPadding(length int) int {
	return (4 - (length % 4)) % 4
}

func extractAVPs(data []byte) ([]*AVP, error) {
	avps := make([]*AVP, 0)
	offset := 0
	for offset < len(data) {
		avp := &AVP{}
		if err := avp.Decode(data[offset:]); err != nil {
			return nil, err
		}
		avps = append(avps, avp)
		// Move to the next AVP
		// Take care of padding
		offset += int(avp.AVPlength) + getPadding(int(avp.AVPlength))

	}
	return avps, nil
}

// get AVP with either name or code based on type of input
// using generic input type to allow for either string or uint32
func (msg *DiameterMessage) GetAVP(input interface{}) *AVP {
	var code uint32
	switch input.(type) {
	case string:
		code = GetAVPCodeFromName(input.(string))
	case uint32:
		code = input.(uint32)
	default:
		return nil
	}

	for _, avp := range msg.AVPs {
		if avp.Code == code {
			return avp
		}
	}
	return nil
}
