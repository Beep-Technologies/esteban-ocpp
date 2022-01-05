package chargepoint

import (
	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging/schemas"
	"github.com/mitchellh/mapstructure"
)

func (cp *OCPP16ChargePoint) handleStartTransaction(msg messaging.OCPP16CallMessage) (*messaging.OCPP16CallResult, *messaging.OCPP16CallError) {
	p := &schemas.StartTransactionRequest{}

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

	// check if there is an ongoing transaction
	otRes, err := cp.transactionService.GetOngoingTransaction(cp.ctx, &rpc.GetOngoingTransactionReq{
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

	// if there is no ongoing transaction, this should not be a remote-initiated transaction
	if !otRes.OngoingTransaction {
		// create transaction
		t, err := cp.transactionService.CreateTransaction(cp.ctx, &rpc.CreateTransactionReq{
			EntityCode:            cp.entityCode,
			ChargePointIdentifier: cp.chargePointIdentifier,
			ConnectorId:           int32(p.ConnectorId),
			RemoteInitiated:       false,
			IdTag:                 p.IdTag,
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

		// start transaction
		_, err = cp.transactionService.StartTransaction(cp.ctx, &rpc.StartTransactionReq{
			Id:              t.Transaction.Id,
			StartMeterValue: int32(p.MeterStart),
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

		return &messaging.OCPP16CallResult{
			MessageTypeID: messaging.CALLRESULT,
			UniqueID:      msg.UniqueID,
			Payload: &schemas.StartTransactionResponse{
				TransactionId: int(t.Transaction.Id),
				IdTagInfo: &schemas.IdTagInfo{
					Status: "Accepted",
				},
			},
		}, nil
	}

	// else it should be a remote-initiated transaction
	// check if the idTags match, reject if they dont
	if p.IdTag != otRes.Transaction.IdTag {
		return &messaging.OCPP16CallResult{
			MessageTypeID: messaging.CALLRESULT,
			UniqueID:      msg.UniqueID,
			Payload: &schemas.StartTransactionResponse{
				IdTagInfo: &schemas.IdTagInfo{
					Status: "Rejected",
				},
			},
		}, nil
	}

	t, err := cp.transactionService.StartTransaction(
		cp.ctx,
		&rpc.StartTransactionReq{
			Id:              int32(otRes.Transaction.Id),
			StartMeterValue: int32(p.MeterStart),
		})

	// make callback
	go cp.makeCallback("StartTransaction", t)

	if err != nil {
		return nil, &messaging.OCPP16CallError{
			MessageTypeID:    messaging.CALLERROR,
			UniqueID:         msg.UniqueID,
			ErrorCode:        messaging.InternalError,
			ErrorDescription: err.Error(),
			ErrorDetails:     struct{}{},
		}
	}

	rb := &schemas.StartTransactionResponse{
		TransactionId: int(otRes.Transaction.Id),
		IdTagInfo: &schemas.IdTagInfo{
			Status: "Accepted",
		},
	}

	return &messaging.OCPP16CallResult{
		MessageTypeID: messaging.CALLRESULT,
		UniqueID:      msg.UniqueID,
		Payload:       rb,
	}, nil
}
