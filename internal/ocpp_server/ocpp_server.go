package ocppserver

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	ocpp16cs "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16_cs"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout:  60 * time.Second,
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: false,
	// Allow connections from all origins
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type OCPPWebSocketServer struct {
	logger              *log.Logger
	ocpp16centralSystem *ocpp16cs.OCPP16CentralSystem
}

func NewOCPPWebSocketServer(l *log.Logger, o16cs *ocpp16cs.OCPP16CentralSystem) (s *OCPPWebSocketServer) {
	return &OCPPWebSocketServer{
		logger:              l,
		ocpp16centralSystem: o16cs,
	}
}

// HttpUpgradeHandler handles the WebSocket opening handshake,
// along with the OCPP-J specific constraints
// after which, it passes the WebSocket connection to the protocol's corresponding handler
func (s *OCPPWebSocketServer) HttpUpgradeHandler(c *gin.Context) {
	w, r := c.Writer, c.Request

	// get the application uuid
	applicationId := c.Param("applicationId")

	// get the entity code
	entityCode := c.Param("entityCode")

	// get the charge point identifier and decode it
	// charge point identifiers are percent-encoded
	chargePointIdentifier := c.Param("chargePointIdentifier")

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

	s.logger.Printf("Client connecting with identifier \"%+v\" and protocol \"%+v\"\n",
		chargePointIdentifier,
		selectedProtocol,
	)

	conn, err := upgrader.Upgrade(w, r, h)
	if err != nil {
		s.logger.Printf("[ERROR] %+v\n", err.Error())
		return
	}

	switch selectedProtocol {
	case "ocpp1.6":
		err := s.ocpp16centralSystem.ConnectChargePoint(
			applicationId,
			entityCode,
			chargePointIdentifier,
			conn)
		if err != nil {
			s.logger.Printf("[ERROR] %+v\n", err.Error())
		}
	case "":
		conn.Close()
	}
}
