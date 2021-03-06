package chargepoint

import (
	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging/schemas"
	"github.com/mitchellh/mapstructure"
)

func (cp *OCPP16ChargePoint) handleStopTransaction(msg messaging.OCPP16CallMessage) (*messaging.OCPP16CallResult, *messaging.OCPP16CallError) {
	p := &schemas.StopTransactionRequest{}

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

	// stop the transaction
	t, err := cp.transactionService.StopTransaction(
		cp.ctx,
		&rpc.StopTransactionReq{
			Id:             int32(p.TransactionId),
			StopMeterValue: int32(p.MeterStop),
			StopReason:     p.Reason,
		},
	)

	if err != nil {
		return nil, &messaging.OCPP16CallError{
			MessageTypeID:    messaging.CALLERROR,
			UniqueID:         msg.UniqueID,
			ErrorCode:        messaging.InternalError,
			ErrorDescription: err.Error(),
			ErrorDetails:     struct{}{},
		}
	}

	// make callback
	go cp.makeCallback("StopTransaction", t)

	return &messaging.OCPP16CallResult{
		MessageTypeID: messaging.CALLRESULT,
		UniqueID:      msg.UniqueID,
		Payload:       &schemas.StopTransactionResponse{},
	}, nil
}
