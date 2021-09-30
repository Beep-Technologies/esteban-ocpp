package ocpp16cp

import (
	"context"
	"encoding/json"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	msg "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16_messaging"
	ocpp16 "github.com/Beep-Technologies/beepbeep3-ocpp/internal/schemas/ocpp_16"
)

// statusNotification handles the Status Notification operation, initiated by the charge point
func (c *OCPP16ChargePoint) statusNotification(req *msg.OCPP16CallMessage) (*msg.OCPP16CallResult, *msg.OCPP16CallError) {
	b := &ocpp16.StatusNotificationRequest{}

	p, err := json.Marshal(req.Payload)
	if err != nil {
		return nil, &msg.OCPP16CallError{
			MessageTypeID:    msg.CALLERROR,
			UniqueID:         req.UniqueID,
			ErrorCode:        msg.FormationViolation,
			ErrorDescription: "",
			ErrorDetails:     struct{}{},
		}
	}

	err = b.UnmarshalJSON(p)
	if err != nil {
		return nil, &msg.OCPP16CallError{
			MessageTypeID:    msg.CALLERROR,
			UniqueID:         req.UniqueID,
			ErrorCode:        msg.FormationViolation,
			ErrorDescription: "",
			ErrorDetails:     struct{}{},
		}
	}

	c.status = msg.OCPP16Status(b.Status)

	_, err = c.statusNotificationService.CreateStatusNotification(context.Background(), &rpc.CreateStatusNotificationReq{
		ChargePointId:   int32(c.id),
		ConnectorId:     int32(b.ConnectorId),
		ErrorCode:       b.ErrorCode,
		Info:            b.Info,
		Status:          b.Status,
		Timestamp:       b.Timestamp,
		VendorId:        b.VendorId,
		VendorErrorCode: b.VendorErrorCode,
	})

	if err != nil {
		return nil, &msg.OCPP16CallError{
			MessageTypeID:    msg.CALLERROR,
			UniqueID:         req.UniqueID,
			ErrorCode:        msg.FormationViolation,
			ErrorDescription: "",
			ErrorDetails:     struct{}{},
		}
	}

	rb := &ocpp16.StatusNotificationResponse{}

	return &msg.OCPP16CallResult{
		MessageTypeID: msg.CALLRESULT,
		UniqueID:      req.UniqueID,
		Payload:       rb,
	}, nil
}
