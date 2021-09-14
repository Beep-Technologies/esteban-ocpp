package handlers

import (
	"fmt"

	"github.com/gorilla/websocket"

	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/messaging"
)

// handleOCPP16 is called after the WebSocket handshake
// is completed and the agreed-upon protocol is ocpp1.6
func handleOCPP16(conn *websocket.Conn) {
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%s\n", p)
		_, err = messaging.ParseOCPP16Call(p)
		if err != nil {
			fmt.Println(err)
			return
		}

	}
}
