package connection

import (
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_api/connection/handlers"
	msg "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_api/messaging"
)

// HandleRequest handles requests initiated by the Charge Point
func (c *Connection) HandleRequest(req *msg.OCPP16CallMessage) (res *msg.OCPP16CallResult, err *msg.OCPP16CallError) {
	switch req.Action {
	case "BootNotification":
		res, err = handlers.BootNotification(req)
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
