package chargepoint

import (
	"context"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging/schemas"
	"github.com/mitchellh/mapstructure"
)

func (cp *OCPP16ChargePoint) handleStatusNotification(msg messaging.OCPP16CallMessage) (*messaging.OCPP16CallResult, *messaging.OCPP16CallError) {
	p := &schemas.StatusNotificationRequest{}

	// decode the payload into the struct
	err := mapstructure.Decode(msg.Payload, p)
	if err != nil {
		return nil, &messaging.OCPP16CallError{
			MessageTypeID:    messaging.CALLERROR,
			UniqueID:         msg.UniqueID,
			ErrorCode:        messaging.FormationViolation,
			ErrorDescription: "",
			ErrorDetails:     struct{}{},
		}
	}

	// make callback
	go cp.makeCallback("StatusNotification", p)

	// if the status notification shows that the connector is available,
	// if there is ongoing transaction on the database, set it as abnormally stopped
	if p.Status == "Available" {
		tRes, err := cp.transactionService.GetOngoingTransaction(
			cp.ctx,
			&rpc.GetOngoingTransactionReq{
				EntityCode:            cp.entityCode,
				ChargePointIdentifier: cp.chargePointIdentifier,
				ConnectorId:           int32(p.ConnectorId),
			})

		if err != nil {
			return nil, &messaging.OCPP16CallError{
				MessageTypeID:    messaging.CALLERROR,
				UniqueID:         msg.UniqueID,
				ErrorCode:        messaging.InternalError,
				ErrorDescription: err.Error(),
				ErrorDetails:     struct{}{},
			}
		}

		if tRes.OngoingTransaction {
			cp.transactionService.AbnormalStopTransaction(
				context.Background(),
				&rpc.AbnormalStopTransactionReq{
					Id: tRes.Transaction.Id,
				},
			)
		}
	}

	_, err = cp.statusNotificationService.CreateStatusNotification(
		cp.ctx,
		&rpc.CreateStatusNotificationReq{
			EntityCode:            cp.entityCode,
			ChargePointIdentifier: cp.chargePointIdentifier,
			ConnectorId:           int32(p.ConnectorId),
			ErrorCode:             p.ErrorCode,
			Info:                  p.Info,
			Status:                p.Status,
			Timestamp:             p.Timestamp,
			VendorId:              p.VendorId,
			VendorErrorCode:       p.VendorErrorCode,
		},
	)

	if err != nil {
		return nil, &messaging.OCPP16CallError{
			MessageTypeID:    messaging.CALLERROR,
			UniqueID:         msg.UniqueID,
			ErrorCode:        messaging.FormationViolation,
			ErrorDescription: err.Error(),
			ErrorDetails:     struct{}{},
		}
	}

	return &messaging.OCPP16CallResult{
		MessageTypeID: messaging.CALLRESULT,
		UniqueID:      msg.UniqueID,
		Payload:       &schemas.StatusNotificationResponse{},
	}, nil
}
