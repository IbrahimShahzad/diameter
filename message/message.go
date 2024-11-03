// Diameter message struct, encoding/decoding logic
package message

import (
	"fmt"
	"github.com/IbrahimShahzad/diameter/utils"
	"math/rand/v2"
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

func (h *DiameterHeader) Decode(data []byte) error {
	if len(data) < DIAMETER_HEADER_SIZE {
		return InvalidDiameterHeaderLengthError
	}
	// Decode the header fields from the byte slice.
	byteCount := 0
	h.Version = data[byteCount]
	byteCount += DIAMETER_VERSION_SIZE
	if h.Version != DIAMETER_VERSION {
		return InvalidDiameterVersionError
	}

	h.MessageLength = utils.FromBytes(data[byteCount : byteCount+DIAMETER_MESSAGE_SIZE])
	byteCount += DIAMETER_MESSAGE_SIZE

	h.CommandFlags = data[byteCount]
	byteCount += DIAMETER_COMMAND_FLAGS_SIZE

	h.CommandCode = utils.FromBytes(data[byteCount : byteCount+DIAMETER_COMMAND_CODE_SIZE])
	byteCount += DIAMETER_COMMAND_CODE_SIZE

	h.ApplicationID = utils.FromBytes(data[byteCount : byteCount+DIAMETER_APPLICATION_ID_SIZE])
	byteCount += DIAMETER_APPLICATION_ID_SIZE

	h.HopByHopID = utils.FromBytes(data[byteCount : byteCount+DIAMETER_HOP_BY_HOP_ID_SIZE])
	byteCount += DIAMETER_HOP_BY_HOP_ID_SIZE

	h.EndToEndID = utils.FromBytes(data[byteCount : byteCount+DIAMETER_END_TO_END_ID_SIZE])
	byteCount += DIAMETER_END_TO_END_ID_SIZE

	return nil
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
		"DiameterMessage{\nHeader: %s\nAVPs: %s}",
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

func (msg *DiameterMessage) Decode(data []byte) error {
	if len(data) < DIAMETER_HEADER_SIZE {
		return InvalidMessageLengthError
	}

	version := data[0]
	msgStart := DIAMETER_VERSION_SIZE
	msgLenEnd := msgStart + DIAMETER_MESSAGE_SIZE
	msgLen := utils.FromBytes(data[msgStart : msgLenEnd+1])
	if len(data) < int(msgLen) {
		return InvalidMessageLengthError
	}

	cFlags := data[msgLenEnd]

	cmdCodeStart := msgLenEnd + DIAMETER_COMMAND_FLAGS_SIZE
	cmdCodeEnd := cmdCodeStart + DIAMETER_COMMAND_CODE_SIZE
	cmdCode := utils.FromBytes(data[cmdCodeStart : cmdCodeEnd+1])

	appIDEnd := cmdCodeEnd + DIAMETER_APPLICATION_ID_SIZE
	appID := utils.FromBytes(data[cmdCodeEnd:appIDEnd])

	hopByHopIDEnd := appIDEnd + DIAMETER_HOP_BY_HOP_ID_SIZE
	hopByHopID := utils.FromBytes(data[appIDEnd:hopByHopIDEnd])

	endToEndIDEnd := hopByHopIDEnd + DIAMETER_END_TO_END_ID_SIZE
	endToEndID := utils.FromBytes(data[hopByHopIDEnd:endToEndIDEnd])

	if version != DIAMETER_VERSION {
		return InvalidDiameterVersionError
	}

	// Decode the header
	msg.Header = &DiameterHeader{
		Version:       version,
		MessageLength: msgLen,
		CommandFlags:  cFlags,
		CommandCode:   cmdCode,
		ApplicationID: appID,
		HopByHopID:    hopByHopID,
		EndToEndID:    endToEndID,
	}
	// fmt.Println("DiameterMessage Header: ", msg.Header)

	// Decode each AVP
	offset := DIAMETER_HEADER_SIZE
	avps, err := extractAVPs(data[offset:])
	if err != nil {
		return err
	}
	// fmt.Println("DiameterMessage AVPs: ", avps)
	msg.AVPs = avps

	return nil
}

// NewCER generates a Capabilities-Exchange-Request message.
func NewCER(avps ...*AVP) (*DiameterMessage, error) {
	return &DiameterMessage{
		Header: &DiameterHeader{
			Version:       DIAMETER_VERSION,
			CommandFlags:  COMMAND_FLAG_REQUEST, // Set 'R' bit for request
			CommandCode:   COMMAND_CODE_CER,     // CER Command Code
			ApplicationID: 0,                    // Base Protocol Application ID
			HopByHopID:    generateHopByHopID(),
			EndToEndID:    generateEndToEndID(),
		},
		AVPs: avps,
	}, nil
}

func NewDWR(avps ...*AVP) (*DiameterMessage, error) {
	return &DiameterMessage{
		Header: &DiameterHeader{
			Version:       DIAMETER_VERSION,
			CommandFlags:  COMMAND_FLAG_REQUEST, // Set 'R' bit for request
			CommandCode:   COMMAND_CODE_DWR,     // DWR Command Code
			ApplicationID: 0,                    // Base Protocol Application ID
			HopByHopID:    generateHopByHopID(),
			EndToEndID:    generateEndToEndID(),
		},
		AVPs: avps,
	}, nil
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
