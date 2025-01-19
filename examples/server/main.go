package main

import (
	ds "github.com/IbrahimShahzad/diameter/server"
	"log"
)

func main() {
	server := ds.NewServer(
		ds.WithServerAddr("localhost:3868"),
		ds.WithTCP(),
		ds.WithConnectionTimeout(0),
	)
	log.Printf("Listening on %s", server.Addr())
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
