package diameter

// import (
// 	"diameter/client"
// 	"diameter/message"
// 	"diameter/server"
// )
//
// // Client API
// func NewClient(addr string) *client.Client {
// 	return client.New(addr)
// }
//
// // Server API
// func NewServer(addr string) *server.Server {
// 	return server.New(addr)
// }
//
// // Message/AVP API for constructing Diameter messages
// func NewDiameterMessage(cmdCode uint32, isRequest bool) *message.DiameterMessage {
// 	return message.NewDiameterMessage(cmdCode, isRequest)
// }
//
// func NewAVP(code uint32, value interface{}) *message.AVP {
// 	return message.NewAVP(code, value)
// }
