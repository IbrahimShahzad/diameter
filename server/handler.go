// Request/response handling and connection managemen
package server

import (
	"log"

	"github.com/IbrahimShahzad/diameter/message"
)

func (p *Peer) generateCEA(msg *message.DiameterMessage) *message.DiameterMessage {
	// This function is called when a CER message is received
	log.Println("Received CER message")
	// Send CEA message
	resultAVP, err := message.NewAVP(message.AVP_RESULT_CODE, uint32(2001), message.MANDATORY_FLAG)
	if err != nil {
		log.Printf("Error creating AVP: %v", err)
		return nil
	}

	orignHostAVP, err := message.NewAVP(message.AVP_ORIGIN_HOST, "localhost", message.MANDATORY_FLAG)
	if err != nil {
		log.Printf("Error creating AVP: %v", err)
		return nil
	}

	orignRealmAVP, err := message.NewAVP(message.AVP_ORIGIN_REALM, "example.ims.com", message.MANDATORY_FLAG)
	if err != nil {
		log.Printf("Error creating AVP: %v", err)
		return nil
	}

	vendorIDAVP, err := message.NewAVP(message.AVP_VENDOR_ID, uint32(10415), message.MANDATORY_FLAG)
	if err != nil {
		log.Printf("Error creating AVP: %v", err)
		return nil
	}

	productNameAVP, err := message.NewAVP(message.AVP_PRODUCT_NAME, "Diameter Server", message.MANDATORY_FLAG)
	if err != nil {
		log.Printf("Error creating AVP: %v", err)
		return nil
	}

	cea, err := message.NewResponseFromRequest(msg,
		resultAVP,
		orignHostAVP,
		orignRealmAVP,
		vendorIDAVP,
		productNameAVP)
	if err != nil {
		log.Printf("Error creating CEA message: %v", err)
		return nil
	}

	log.Println("Sending Capabilities-Exchange-Answer (CEA) in response to CER.")
	// send CEA message over channel
	return cea
}
