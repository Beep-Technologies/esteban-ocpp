package ocppserver

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/Beep-Technologies/beepbeep3-iam/pkg/logger"
	ocpp16 "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16"
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
	logger              *zap.Logger
	ocpp16centralSystem ocpp16.CentralSystem
}

func NewOCPPWebSocketServer(l *zap.Logger, o16cs ocpp16.CentralSystem) (s *OCPPWebSocketServer) {
	wsServerLogger := logger.With(
		zap.String("source", "ocpp_ws_server"),
	)

	return &OCPPWebSocketServer{
		logger:              wsServerLogger,
		ocpp16centralSystem: o16cs,
	}
}

// HttpUpgradeHandler handles the WebSocket opening handshake,
// along with the OCPP-J specific constraints
// after which, it passes the WebSocket connection to the protocol's corresponding handler
func (s *OCPPWebSocketServer) HttpUpgradeHandler(c *gin.Context) {
	w, r := c.Writer, c.Request

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

	s.logger.Info(
		"Client connecting",
		zap.String("event", "ws_client_connect"),
		zap.String("charge_point_identifier", chargePointIdentifier),
		zap.String("protocol", selectedProtocol),
	)

	conn, err := upgrader.Upgrade(w, r, h)
	if err != nil {
		s.logger.Error(
			err.Error(),
			zap.String("event", "ws_upgrade_error"),
		)
		return
	}

	switch selectedProtocol {
	case "ocpp1.6":
		err := s.ocpp16centralSystem.ConnectChargePoint(
			entityCode,
			chargePointIdentifier,
			conn,
		)
		if err != nil {
			s.logger.Error(
				err.Error(),
				zap.String("event", "connect_charge_point_error"),
			)
		}
	case "":
		conn.Close()
	}
}
