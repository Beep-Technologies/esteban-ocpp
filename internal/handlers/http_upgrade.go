package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout:  60 * time.Second,
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: false,
}

// HttpUpgradeHandler handles the WebSocket opening handshake,
// along with the OCPP-J specific constraints
// after which, it passes the WebSocket connection to the protocol's corresponding handler
func HttpUpgradeHandler(w http.ResponseWriter, r *http.Request) {
	// get the charge point identifier and decode it
	// charge point identifiers are percent-encoded, mux decodes it by default
	vars := mux.Vars(r)
	cpId := vars["chargePointIdentifier"]

	if cpId == "" {
		http.Error(w, "Invalid Charge Point Identifier", http.StatusNotFound)
		return
	}

	// server-offered subprotocols
	serverProtocols := []string{"ocpp1.6"}

	// Check client-offered subprotocols
	ps := r.Header.Get("Sec-WebSocket-Protocol")
	ps = strings.ToLower(ps)
	clientProtocols := strings.Split(ps, ",")

	// agreed-upon subprotocol
	selectedProtocol := ""

	// check for a matching subprotocol
	for _, clientProtocol := range clientProtocols {
		clientProtocol = strings.TrimSpace(clientProtocol)

		for _, serverProtocol := range serverProtocols {
			if clientProtocol == serverProtocol {
				selectedProtocol = serverProtocol
				break
			}
		}

		if selectedProtocol != "" {
			break
		}
	}

	// if no matching protocols, complete the WebSocket handshake
	// without a Sec-WebSocket-Protocol header and then immediately close the WebSocket connection
	// if there is a matching protocol, complete the WebSocket handshake
	// and begin handling OCPP messages
	h := http.Header{}
	h.Set("Sec-WebSocket-Protocol", selectedProtocol)

	if selectedProtocol == "" {
		h = nil
	}

	fmt.Printf("Client connecting with identifier \"%+v\" and protocol \"%+v\"\n",
		cpId,
		selectedProtocol,
	)

	conn, err := upgrader.Upgrade(w, r, h)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		return
	}

	switch selectedProtocol {
	case "ocpp1.6":
		handleOCPP16(conn)
	case "":
		conn.Close()
	}
}
