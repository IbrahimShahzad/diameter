package main

import (
	dc "github.com/IbrahimShahzad/diameter/client"
	"log"
)

func main() {
	client, err := dc.NewClient(
		dc.WithTCP(),
		dc.WithConnectionTimeout(5),
		dc.WithServerAddr("localhost:3868"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	err = client.Connect()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	client.Disconnect()
}
