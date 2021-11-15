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
			ErrorDescription: err.Error(),
			ErrorDetails:     struct{}{},
		}
	}

	err = b.UnmarshalJSON(p)
	if err != nil {
		return nil, &msg.OCPP16CallError{
			MessageTypeID:    msg.CALLERROR,
			UniqueID:         req.UniqueID,
			ErrorCode:        msg.FormationViolation,
			ErrorDescription: err.Error(),
			ErrorDetails:     struct{}{},
		}
	}

	// make callback
	go c.makeCallback("StatusNotification", b)

	// if the status notification shows that the connector is available,
	// if there is ongoing transaction on the database, set it as abnormally stopped
	if b.Status == "Available" {
		tRes, err := c.transactionService.GetOngoingTransaction(context.Background(), &rpc.GetOngoingTransactionReq{
			ApplicationId:         c.applicationId,
			ChargePointIdentifier: c.chargePointIdentifier,
			ConnectorId:           int32(b.ConnectorId),
		})

		if err != nil {
			return nil, &msg.OCPP16CallError{
				MessageTypeID:    msg.CALLERROR,
				UniqueID:         req.UniqueID,
				ErrorCode:        msg.InternalError,
				ErrorDescription: err.Error(),
				ErrorDetails:     struct{}{},
			}
		}

		if tRes.OngoingTransaction {
			c.transactionService.AbnormalStopTransaction(
				context.Background(),
				&rpc.AbnormalStopTransactionReq{
					Id: tRes.Transaction.Id,
				},
			)
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
			ErrorDescription: err.Error(),
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
