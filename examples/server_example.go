package main

// import (
// 	"diameter"
// 	"log"
// )
//
// func main() {
// 	server := diameter.NewServer(":3868")
// 	err := server.ListenAndServe()
// 	if err != nil {
// 		log.Fatalf("Failed to start server: %v", err)
// 	}
// }

// transport example

// //package main
//
// import (
// 	"diameter/transport"
// 	"log"
// 	"time"
// )
//
// func main() {
// 	listener, err := transport.NewDiameterListener(":3868", 5*time.Second)
// 	if err != nil {
// 		log.Fatalf("Failed to start listener: %v", err)
// 	}
// 	defer listener.Close()
// 	log.Printf("Listening on %s", listener.Addr())
//
// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			log.Printf("Error accepting connection: %v", err)
// 			continue
// 		}
//
// 		go handleConnection(conn)
// 	}
// }
//
// func handleConnection(conn *transport.DiameterConnection) {
// 	defer conn.Close()
// 	log.Println("Handling new connection.")
//
// 	buffer := make([]byte, 1024)
// 	for {
// 		n, err := conn.Read(buffer)
// 		if err != nil {
// 			log.Printf("Error reading from connection: %v", err)
// 			return
// 		}
// 		log.Printf("Received message: %s", string(buffer[:n]))
//
// 		// Example of sending a response
// 		response := []byte("Diameter Response Message")
// 		_, err = conn.Write(response)
// 		if err != nil {
// 			log.Printf("Error writing to connection: %v", err)
// 			return
// 		}
// 	}
// }
