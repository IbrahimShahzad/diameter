// Definitions for common codes (e.g., Result-Code)
package message

import "errors"

// AVP Codes
const (
	AVP_CODE_SESSION_ID          = uint32(263)
	AVP_CODE_ORIGIN_HOST         = uint32(264)
	AVP_CODE_ORIGIN_REALM        = uint32(296)
	AVP_CODE_HOST_IP_ADDRESS     = uint32(257)
	AVP_CODE_VENDOR_ID           = uint32(266)
	AVP_CODE_PRODUCT_NAME        = uint32(269)
	AVP_CODE_ORIGIN_STATE_ID     = uint32(278)
	AVP_CODE_RESULT_CODE         = uint32(268)
	AVP_CODE_ERROR_MESSAGE       = uint32(281)
	AVP_CODE_EXPERIMENTAL_RESULT = uint32(297)
	AVP_CODE_FAILED_AVP          = uint32(279)
)

var AVPCodeToName map[uint32]string = map[uint32]string{
	AVP_CODE_SESSION_ID:          "Session-Id",
	AVP_CODE_ORIGIN_HOST:         "Origin-Host",
	AVP_CODE_ORIGIN_REALM:        "Origin-Realm",
	AVP_CODE_HOST_IP_ADDRESS:     "Host-IP-Address",
	AVP_CODE_VENDOR_ID:           "Vendor-Id",
	AVP_CODE_PRODUCT_NAME:        "Product-Name",
	AVP_CODE_ORIGIN_STATE_ID:     "Origin-State-Id",
	AVP_CODE_RESULT_CODE:         "Result-Code",
	AVP_CODE_ERROR_MESSAGE:       "Error-Message",
	AVP_CODE_EXPERIMENTAL_RESULT: "Experimental-Result",
	AVP_CODE_FAILED_AVP:          "Failed-AVP",
}

func GetAVPCodeFromName(name string) uint32 {
	if code, ok := AVPNameToCode[name]; ok {
		return code
	}
	return 0
}

var AVPNameToCode map[string]uint32 = map[string]uint32{
	"Origin-Host":         AVP_CODE_ORIGIN_HOST,
	"Origin-Realm":        AVP_CODE_ORIGIN_REALM,
	"Host-IP-Address":     AVP_CODE_HOST_IP_ADDRESS,
	"Vendor-Id":           AVP_CODE_VENDOR_ID,
	"Product-Name":        AVP_CODE_PRODUCT_NAME,
	"Origin-State-Id":     AVP_CODE_ORIGIN_STATE_ID,
	"Result-Code":         AVP_CODE_RESULT_CODE,
	"Error-Message":       AVP_CODE_ERROR_MESSAGE,
	"Experimental-Result": AVP_CODE_EXPERIMENTAL_RESULT,
	"Failed-AVP":          AVP_CODE_FAILED_AVP,
}

func GetAVPNameFromCode(code uint32) string {
	return AVPCodeToName[code]
}

type ResultCode uint32

const (
	DIAMETER_MULTI_ROUND_AUTH ResultCode = 1001
	DIAMETER_SUCCESS          ResultCode = 2001
	DIMAETER_LIMITED_SUCCESS  ResultCode = 2002
)
const (
	DIAMETER_COMMAND_UNSUPPORTED ResultCode = 3001 + iota
	DIAMETER_UNABLE_TO_DELIVER
	DIAMETER_REALM_NOT_SERVED
	DIAMETER_TOO_BUSY
	DIAMETER_LOOP_DETECTED
	DIAMETER_REDIRECT_INDICATION
	DIAMETER_APPLICATION_UNSUPPORTED
	DIAMETER_INVALID_HDR_BITS
	DIAMETER_INVALID_AVP_BITS
	DIAMETER_UNKNOWN_PEER
)

const (
	DIAMETER_AUTHENTICATION_REJECTED ResultCode = 4001 + iota
	DIAMETER_OUT_OF_SPACE
	DIAMETER_ELECTION_LOST
)

const (
	DIAMETER_AVP_UNSUPPORTED ResultCode = 5001 + iota
	DIAMETER_UNKNOWN_SESSION_ID
	DIAMETER_AUTHORIZATION_REJECTED
	DIAMETER_INVALID_AVP_VALUE
	DIAMETER_MISSING_AVP
	DIAMETER_RESOURCES_EXCEEDED
	DIAMETER_CONTRADICTING_AVPS
	DIAMETER_AVP_NOT_ALLOWED
	DIAMETER_AVP_OCCURS_TOO_MANY_TIMES
	DIAMETER_NO_COMMON_APPLICATION
	DIAMETER_UNSUPPORTED_VERSION
	DIAMETER_UNABLE_TO_COMPLY
	DIAMETER_INVALID_BIT_IN_HEADER
	DIAMETER_INVALID_AVP_LENGTH
	DIAMETER_INVALID_MESSAGE_LENGTH
	DIAMETER_INVALID_AVP_BIT_COMBO
	DIAMETER_NO_COMMON_SECURITY
)

var ResultCodeToName map[ResultCode]string = map[ResultCode]string{
	DIAMETER_SUCCESS:                   "DIAMETER_SUCCESS",
	DIMAETER_LIMITED_SUCCESS:           "DIMAETER_LIMITED_SUCCESS",
	DIAMETER_MULTI_ROUND_AUTH:          "DIAMETER_MULTI_ROUND_AUTH",
	DIAMETER_COMMAND_UNSUPPORTED:       "DIAMETER_COMMAND_UNSUPPORTED",
	DIAMETER_UNABLE_TO_DELIVER:         "DIAMETER_UNABLE_TO_DELIVER",
	DIAMETER_REALM_NOT_SERVED:          "DIAMETER_REALM_NOT_SERVED",
	DIAMETER_TOO_BUSY:                  "DIAMETER_TOO_BUSY",
	DIAMETER_LOOP_DETECTED:             "DIAMETER_LOOP_DETECTED",
	DIAMETER_REDIRECT_INDICATION:       "DIAMETER_REDIRECT_INDICATION",
	DIAMETER_APPLICATION_UNSUPPORTED:   "DIAMETER_APPLICATION_UNSUPPORTED",
	DIAMETER_INVALID_HDR_BITS:          "DIAMETER_INVALID_HDR_BITS",
	DIAMETER_INVALID_AVP_BITS:          "DIAMETER_INVALID_AVP_BITS",
	DIAMETER_UNKNOWN_PEER:              "DIAMETER_UNKNOWN_PEER",
	DIAMETER_AUTHENTICATION_REJECTED:   "DIAMETER_AUTHENTICATION_REJECTED",
	DIAMETER_OUT_OF_SPACE:              "DIAMETER_OUT_OF_SPACE",
	DIAMETER_ELECTION_LOST:             "DIAMETER_ELECTION_LOST",
	DIAMETER_AVP_UNSUPPORTED:           "DIAMETER_AVP_UNSUPPORTED",
	DIAMETER_UNKNOWN_SESSION_ID:        "DIAMETER_UNKNOWN_SESSION_ID",
	DIAMETER_AUTHORIZATION_REJECTED:    "DIAMETER_AUTHORIZATION_REJECTED",
	DIAMETER_INVALID_AVP_VALUE:         "DIAMETER_INVALID_AVP_VALUE",
	DIAMETER_MISSING_AVP:               "DIAMETER_MISSING_AVP",
	DIAMETER_RESOURCES_EXCEEDED:        "DIAMETER_RESOURCES_EXCEEDED",
	DIAMETER_CONTRADICTING_AVPS:        "DIAMETER_CONTRADICTING_AVPS",
	DIAMETER_AVP_NOT_ALLOWED:           "DIAMETER_AVP_NOT_ALLOWED",
	DIAMETER_AVP_OCCURS_TOO_MANY_TIMES: "DIAMETER_AVP_OCCURS_TOO_MANY_TIMES",
	DIAMETER_NO_COMMON_APPLICATION:     "DIAMETER_NO_COMMON_APPLICATION",
	DIAMETER_UNSUPPORTED_VERSION:       "DIAMETER_UNSUPPORTED_VERSION",
	DIAMETER_UNABLE_TO_COMPLY:          "DIAMETER_UNABLE_TO_COMPLY",
	DIAMETER_INVALID_BIT_IN_HEADER:     "DIAMETER_INVALID_BIT_IN_HEADER",
	DIAMETER_INVALID_AVP_LENGTH:        "DIAMETER_INVALID_AVP_LENGTH",
	DIAMETER_INVALID_MESSAGE_LENGTH:    "DIAMETER_INVALID_MESSAGE_LENGTH",
	DIAMETER_INVALID_AVP_BIT_COMBO:     "DIAMETER_INVALID_AVP_BIT_COMBO",
	DIAMETER_NO_COMMON_SECURITY:        "DIAMETER_NO_COMMON_SECURITY",
}

// 7.1.  Result-Code AVP
//
//	The Result-Code AVP (AVP Code 268) is of type Unsigned32 and
//	indicates whether a particular request was completed successfully or
//	whether an error occurred.  All Diameter answer messages defined in
//	IETF applications MUST include one Result-Code AVP.  A non-successful
//	Result-Code AVP (one containing a non 2xxx value other than
//	DIAMETER_REDIRECT_INDICATION) MUST include the Error-Reporting-Host
//	AVP if the host setting the Result-Code AVP is different from the
//	identity encoded in the Origin-Host AVP.
//
//	The Result-Code data field contains an IANA-managed 32-bit address
//	space representing errors (see Section 11.4).  Diameter provides the
//	following classes of errors, all identified by the thousands digit in
//	the decimal notation:
//
//	   -  1xxx (Informational)
//	   -  2xxx (Success)
//	   -  3xxx (Protocol Errors)
//	   -  4xxx (Transient Failures)
//	   -  5xxx (Permanent Failure)
func GetResultCode(msg *DiameterMessage) (ResultCode, string, error) {
	// TODO: check for error bit
	for _, avp := range msg.AVPs {
		if avp.Code == AVP_CODE_RESULT_CODE {
			if value, ok := avp.Data.(*Unsigned32); ok {
				return ResultCode(value.Data), ResultCodeToName[ResultCode(value.Data)], nil
			}
		}
	}
	// Return a default value or an error if Result-Code AVP is not found
	return ResultCode(0), "", errors.New("Result-Code AVP not found")
}

func ValidateSuccessfulResponse(msg *DiameterMessage) error {
	_, _, err := GetResultCode(msg)
	return err
}