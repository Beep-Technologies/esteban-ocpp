package ocpp16

import (
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging"
	"github.com/gorilla/websocket"
)

type ChargePoint interface {
	// start listening on the (websocket) connection for messages
	Listen()
	// make a request to remote start transactions
	RemoteStartTransaction(connectorID int) (transactionID int, err error)
	// make a request to remote stop transaction
	RemoteStopTransaction(connectorID int) (err error)
	// make a request to trigger a status notification from the charge point
	TriggerStatusNotification(connectorID int, errorStatusCode messaging.OCPP16ChargePointErrorCode) (err error)
}

type CentralSystem interface {
	// connect a charge point to this central system
	ConnectChargePoint(entityCode, identifier string, conn *websocket.Conn) (err error)
	// get the charge point, if available
	GetChargePoint(entityCode, identifier string) (cp ChargePoint, err error)
}
