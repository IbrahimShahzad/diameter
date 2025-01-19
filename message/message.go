// Diameter message struct, encoding/decoding logic
package message

import (
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/IbrahimShahzad/diameter/utils"
)

const (
	DIAMETER_VERSION_SIZE        = 1
	DIAMETER_MESSAGE_SIZE        = 3
	DIAMETER_COMMAND_FLAGS_SIZE  = 1
	DIAMETER_COMMAND_CODE_SIZE   = 3
	DIAMETER_APPLICATION_ID_SIZE = 4
	DIAMETER_HOP_BY_HOP_ID_SIZE  = 4
	DIAMETER_END_TO_END_ID_SIZE  = 4
	DIAMETER_HEADER_SIZE         = 20
)

// Command Flags
const (
	COMMAND_FLAG_REQUEST       = 0x80
	COMMAND_FLAG_PROXIABLE     = 0x40
	COMMAND_FLAG_ERROR         = 0x20
	COMMAND_FLAG_RETRANSMITTED = 0x10
	COMMAND_FLAG_RESPONSE      = 0x00
)

// Command Codes
const (
	COMMAND_CODE_CER = uint32(257)
	COMMAND_CODE_DWR = uint32(280)
)

func GetCommandNameFromCode(code uint32) string {
	return CommandCodeToName[code]
}

var CommandCodeToName map[uint32]string = map[uint32]string{
	COMMAND_CODE_CER: "Capabilities-Exchange-Request",
	COMMAND_CODE_DWR: "Diameter-Watchdog-Request",
}

const DIAMETER_VERSION = 1

// DiameterHeader represents a basic Diameter message header.
// The header is 20 bytes long and contains the following fields:
//   - Version: 1 byte
//   - MessageLength: 3 bytes
//   - CommandFlags: 1 byte
//   - CommandCode: 3 bytes
//   - ApplicationID: 4 bytes
//   - HopByHopID: 4 bytes
//   - EndToEndID: 4 bytes
type DiameterHeader struct {
	Version       uint8
	MessageLength uint32
	CommandFlags  uint8
	CommandCode   uint32
	ApplicationID uint32
	HopByHopID    uint32
	EndToEndID    uint32
}

func (h *DiameterHeader) String() string {
	return fmt.Sprintf(
		"Version: %d\nMessageLength: %d\nCommandFlags: %d\nCommandCode: %d\nApplicationID: %d\nHopByHopID: %d\nEndToEndID: %d\n",
		h.Version,
		h.MessageLength,
		h.CommandFlags,
		h.CommandCode,
		h.ApplicationID,
		h.HopByHopID,
		h.EndToEndID,
	)
}

func (h *DiameterHeader) Encode() []byte {
	// Allocate a byte slice of 20 bytes to store the header.
	header := make([]byte, DIAMETER_HEADER_SIZE)

	// Encode the header fields into the byte slice.
	byteCount := 0
	header[byteCount] = h.Version
	byteCount += DIAMETER_VERSION_SIZE

	copy(header[byteCount:], utils.ToBytes(h.MessageLength, DIAMETER_MESSAGE_SIZE))
	byteCount += DIAMETER_MESSAGE_SIZE

	header[byteCount] = h.CommandFlags
	byteCount += DIAMETER_COMMAND_FLAGS_SIZE

	copy(header[byteCount:], utils.ToBytes(h.CommandCode, DIAMETER_COMMAND_CODE_SIZE))
	byteCount += DIAMETER_COMMAND_CODE_SIZE

	copy(header[byteCount:], utils.ToBytes(h.ApplicationID, DIAMETER_APPLICATION_ID_SIZE))
	byteCount += DIAMETER_APPLICATION_ID_SIZE

	copy(header[byteCount:], utils.ToBytes(h.HopByHopID, DIAMETER_HOP_BY_HOP_ID_SIZE))
	byteCount += DIAMETER_HOP_BY_HOP_ID_SIZE

	copy(header[byteCount:], utils.ToBytes(h.EndToEndID, DIAMETER_END_TO_END_ID_SIZE))
	byteCount += DIAMETER_END_TO_END_ID_SIZE

	return header
}

func DecodeHeader(data []byte) (*DiameterHeader, error) {

	if len(data) < DIAMETER_HEADER_SIZE {
		log.Printf("DiameterHeader Decode: Insufficient data %d", len(data))
		return nil, InvalidDiameterHeaderLengthError
	}
	// Decode the header fields from the byte slice.
	byteCount := 0
	version := data[byteCount]
	byteCount += DIAMETER_VERSION_SIZE
	if version != DIAMETER_VERSION {
		return nil, InvalidDiameterVersionError
	}

	messageLength := utils.FromBytes(data[byteCount : byteCount+DIAMETER_MESSAGE_SIZE])
	byteCount += DIAMETER_MESSAGE_SIZE

	commandFlags := data[byteCount]
	byteCount += DIAMETER_COMMAND_FLAGS_SIZE

	commandCode := utils.FromBytes(data[byteCount : byteCount+DIAMETER_COMMAND_CODE_SIZE])
	byteCount += DIAMETER_COMMAND_CODE_SIZE

	applicationID := utils.FromBytes(data[byteCount : byteCount+DIAMETER_APPLICATION_ID_SIZE])
	byteCount += DIAMETER_APPLICATION_ID_SIZE

	hopByHopID := utils.FromBytes(data[byteCount : byteCount+DIAMETER_HOP_BY_HOP_ID_SIZE])
	byteCount += DIAMETER_HOP_BY_HOP_ID_SIZE

	endToEndID := utils.FromBytes(data[byteCount : byteCount+DIAMETER_END_TO_END_ID_SIZE])
	byteCount += DIAMETER_END_TO_END_ID_SIZE

	return &DiameterHeader{
		Version:       version,
		MessageLength: messageLength,
		CommandFlags:  commandFlags,
		CommandCode:   commandCode,
		ApplicationID: applicationID,
		HopByHopID:    hopByHopID,
		EndToEndID:    endToEndID,
	}, nil
}

func generateHopByHopID() uint32 {
	return rand.Uint32()
}

func generateEndToEndID() uint32 {
	return rand.Uint32()
}

// DiameterMessage represents a Diameter message with header and AVPs.
type DiameterMessage struct {
	Header *DiameterHeader
	AVPs   []*AVP
}

func (m *DiameterMessage) String() string {
	var avps string
	for _, avp := range m.AVPs {
		avps += avp.String() + "\n"
	}
	return fmt.Sprintf(
		"DiameterMessage{\nHeader: %s\nAVPs: \n%s}",
		m.Header.String(),
		avps,
	)
}

func (msg *DiameterMessage) Encode() ([]byte, error) {
	// Encode the header
	header := msg.Header.Encode()

	// Encode each AVP
	avps := make([]byte, 0)
	for _, avp := range msg.AVPs {
		encoded, err := avp.Encode()
		if err != nil {
			return nil, err
		}
		avps = append(avps, encoded...)
	}

	// Concatenate the header and AVPs
	return append(header, avps...), nil
}

func DecodeMessage(data []byte) (*DiameterMessage, error) {
	if len(data) < DIAMETER_HEADER_SIZE {
		return nil, InvalidMessageLengthError
	}

	// Decode the header
	header, err := DecodeHeader(data)
	if err != nil {
		log.Printf("DiameterMessage Decode: Header Decode error %v", err)
		return nil, err
	}

	// Decode each AVP
	offset := DIAMETER_HEADER_SIZE
	avps, err := extractAVPs(data[offset:])
	if err != nil {
		return nil, err
	}

	return &DiameterMessage{
		Header: header,
		AVPs:   avps,
	}, nil
}

// NewCER generates a Capabilities-Exchange-Request message.
func NewCER(avps ...*AVP) (*DiameterMessage, error) {
	return NewRequest(COMMAND_CODE_CER, avps...)
}

func NewRequest(commandCode uint32, avps ...*AVP) (*DiameterMessage, error) {
	return &DiameterMessage{
		Header: &DiameterHeader{
			Version:       DIAMETER_VERSION,
			CommandFlags:  COMMAND_FLAG_REQUEST, // Set 'R' bit for request
			CommandCode:   commandCode,
			ApplicationID: 0, // Base Protocol Application ID
			HopByHopID:    generateHopByHopID(),
			EndToEndID:    generateEndToEndID(),
			MessageLength: uint32(DIAMETER_HEADER_SIZE + len(avps)),
		},
		AVPs: avps,
	}, nil
}

func NewResponseFromRequest(request *DiameterMessage, avps ...*AVP) (*DiameterMessage, error) {
	msg := &DiameterMessage{
		Header: &DiameterHeader{
			Version:       DIAMETER_VERSION,
			CommandFlags:  COMMAND_FLAG_RESPONSE, // Set 'R' bit for response
			CommandCode:   request.Header.CommandCode,
			ApplicationID: request.Header.ApplicationID,
			HopByHopID:    request.Header.HopByHopID,
			EndToEndID:    request.Header.EndToEndID,
		},
		AVPs: avps,
	}
	msg.Header.MessageLength = uint32(DIAMETER_HEADER_SIZE + len(msg.AVPs))
	return msg, nil
}

func NewDWR(avps ...*AVP) (*DiameterMessage, error) {
	return NewRequest(COMMAND_CODE_DWR, avps...)
}

// read CEA message, check Success or Failure and return AVPs
func ReadCEA(cea DiameterMessage) ([]*AVP, error) {
	if cea.Header.CommandCode != COMMAND_CODE_CER {
		return nil, InvalidCommandCodeError
	}

	// check for mandatory AVPs
	result_code, _, err := GetResultCode(&cea)
	if err != nil {
		return nil, err
	}
	if result_code != DIAMETER_SUCCESS {
		return nil, fmt.Errorf(
			"CEA failed with Result-Code: %d (%s)",
			result_code,
			ResultCodeToName[result_code],
		)
	}
	return cea.AVPs, nil
}
