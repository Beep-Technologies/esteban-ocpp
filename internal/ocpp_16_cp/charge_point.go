package ocpp16cp

import (
	"log"

	"github.com/gorilla/websocket"

	msg "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16_messaging"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/application"
	chargepoint "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/charge_point"
	statusnotification "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/status_notification"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/transaction"
)

type OCPP16ChargePoint struct {
	id                        int
	chargePointIdentifier     string
	applicationId             string
	conn                      *websocket.Conn
	status                    msg.OCPP16Status
	currentCall               *currentCall
	applicationService        *application.Service
	chargepointService        *chargepoint.Service
	transactionService        *transaction.Service
	statusNotificationService *statusnotification.Service
	logger                    *log.Logger
}

func NewOCPP16ChargePoint(
	id int,
	chargePointIdentifier string,
	applicationId string,
	conn *websocket.Conn,
	applicationService *application.Service,
	chargepointService *chargepoint.Service,
	transactionService *transaction.Service,
	statusNotificationService *statusnotification.Service,
	logger *log.Logger,
) *OCPP16ChargePoint {
	return &OCPP16ChargePoint{
		id:                        id,
		chargePointIdentifier:     chargePointIdentifier,
		applicationId:             applicationId,
		conn:                      conn,
		status:                    "Available",
		currentCall:               &currentCall{},
		applicationService:        applicationService,
		chargepointService:        chargepointService,
		transactionService:        transactionService,
		statusNotificationService: statusNotificationService,
		logger:                    logger,
	}
}

func (c *OCPP16ChargePoint) GetIdentifier() string {
	return c.chargePointIdentifier
}

// Listen listens for messages from the Charge point, and responds to them
func (c *OCPP16ChargePoint) Listen() {
	// TODO: handle client disconnects. currently a call to go Listen() will never return

	// listen until close
	for {
		err := c.handleMessage()
		if err != nil {
			c.logger.Printf("[ERROR]: %s", err.Error())
			return
		}
	}
}

// handleMessage handles incoming messages from the Charge Point
// and calls the corresponding handler based on type
func (c *OCPP16ChargePoint) handleMessage() error {

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
			c.logger.Printf("[CALL: FROM %d] %s", c.id, p)

			// parse the OCPP16 call message
			req, err := msg.ParseOCPP16Call(p)
			if err != nil {
				return err
			}

			// route to the appropriate message handler, and get the response (and response error)
			res, resErr := c.handleCall(req)

			// if there is a response error, send a OCPP16 call error message and return
			// no error is returned here
			if resErr != nil {
				resErrJson, err := msg.UnparseOCPP16CallError(resErr)
				if err != nil {
					return err
				}

				c.logger.Printf("[CALLERROR: TO %d] %s", c.id, resErrJson)
				c.conn.WriteMessage(websocket.TextMessage, resErrJson)

				return nil
			}

			// else, send the OCPP16 call result and return
			resJson, err := msg.UnparseOCPP16CallResult(res)
			if err != nil {
				return err
			}

			c.logger.Printf("[CALLRESULT: TO %d] %s", c.id, resJson)
			c.conn.WriteMessage(websocket.TextMessage, resJson)
		}
	case msg.CALLRESULT:
		{
			c.logger.Printf("[CALLRESULT: FROM %d] %s", c.id, p)

			// parse the OCPP16 call result
			res, err := msg.ParseOCPP16CallResult(p)
			if err != nil {
				return err
			}

			// route to the appropriate message handler
			err = c.handleCallResult(res)
			if err != nil {
				return err
			}
		}
	case msg.CALLERROR:
		{
			c.logger.Printf("[CALLERROR: FROM %d] %s", c.id, p)

			// parse the OCPP16 call error
			res, err := msg.ParseOCPP16CallError(p)
			if err != nil {
				return err
			}

			// route to the appropriate message handler
			err = c.handleCallError(res)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
