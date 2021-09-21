package ocpp16cp

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	msg "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16_messaging"
)

type currentCall struct {
	m        sync.Mutex
	timer    *time.Timer
	uniqueID string
}

type currentTransaction struct {
	transactionId int
	// taskId        string // TODO: find a better name for this
}

type OCPP16ChargePoint struct {
	cpId               string
	conn               *websocket.Conn
	status             msg.OCPP16Status
	currentCall        *currentCall
	currentTransaction *currentTransaction
	logger             *log.Logger
}

func NewOCPP16ChargePoint(cpId string, conn *websocket.Conn, logger *log.Logger) *OCPP16ChargePoint {
	return &OCPP16ChargePoint{
		cpId:               cpId,
		conn:               conn,
		status:             "Available",
		currentCall:        &currentCall{},
		currentTransaction: nil,
		logger:             logger,
	}
}

// Listen listens for messages from the Charge point, and responds to them
func (c *OCPP16ChargePoint) Listen() error {
	// TODO: handle client disconnects. currently a call to go Listen() will never return

	// listen until close
	for {
		err := c.handleMessage()
		if err != nil {
			return err
		}

	}
}

// handleMessage handles incoming messages from the Charge Point
// and calls the corresponding handler based on type
func (c *OCPP16ChargePoint) handleMessage() error {
	{
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			return err
		}

		msgType, err := msg.GetOCPP16MessageType(p)
		if err != nil {
			return err
		}

		switch msgType {
		case msg.CALL:
			{
				c.logger.Printf("[CALL: FROM %s] %s", c.cpId, p)
				req, err := msg.ParseOCPP16Call(p)
				if err != nil {
					return err
				}

				res, resErr := c.handleRequest(req)

				if resErr != nil {
					resErrJson, err := msg.UnparseOCPP16CallError(resErr)
					if err != nil {
						return err
					}

					c.logger.Printf("[CALLERROR: TO %s] %s", c.cpId, resErrJson)
					c.conn.WriteMessage(websocket.TextMessage, resErrJson)

					errorMsg := fmt.Sprintf("%+v", resErr)
					return errors.New(errorMsg)
				}

				resJson, err := msg.UnparseOCPP16CallResult(res)
				if err != nil {
					return err
				}

				c.logger.Printf("[CALLRESULT: TO %s] %s", c.cpId, resJson)
				c.conn.WriteMessage(websocket.TextMessage, resJson)
			}
		case msg.CALLRESULT:
			{
				_, err := msg.ParseOCPP16CallResult(p)
				if err != nil {
					return err
				}

				c.logger.Printf("[CALLRESULT: FROM %s] %s", c.cpId, p)
			}
		case msg.CALLERROR:
			{
				_, err := msg.ParseOCPP16CallError(p)
				if err != nil {
					return err
				}

				c.logger.Printf("[CALLERROR: FROM %s] %s", c.cpId, p)
			}
		}
	}

	return nil
}

// handleRequest handles requests initiated by the Charge Point
func (c *OCPP16ChargePoint) handleRequest(req *msg.OCPP16CallMessage) (res *msg.OCPP16CallResult, err *msg.OCPP16CallError) {
	// go generics couldn't come sooner ....
	switch req.Action {
	case "BootNotification":
		res, err = c.bootNotification(req)
	case "StatusNotification":
		res, err = c.statusNotification(req)
	case "Heartbeat":
		res, err = c.heartbeat(req)
	case "MeterValues":
		res, err = c.meterValues(req)
	case "StartTransaction":
		res, err = c.startTransaction(req)
	case "StopTransaction":
		res, err = c.stopTransaction(req)
	default:
		res, err = nil, &msg.OCPP16CallError{
			MessageTypeID:    msg.CALLERROR,
			UniqueID:         req.UniqueID,
			ErrorCode:        msg.NotImplemented,
			ErrorDescription: "",
			ErrorDetails:     struct{}{},
		}
	}

	return res, err
}
