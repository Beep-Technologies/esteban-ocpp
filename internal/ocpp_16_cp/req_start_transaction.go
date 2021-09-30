package ocpp16cp

import (
	"context"
	"encoding/json"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	msg "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16_messaging"
	ocpp16 "github.com/Beep-Technologies/beepbeep3-ocpp/internal/schemas/ocpp_16"
)

// startTransaction handles the Start Transaction operation, initiated by the charge point
func (c *OCPP16ChargePoint) startTransaction(req *msg.OCPP16CallMessage) (*msg.OCPP16CallResult, *msg.OCPP16CallError) {
	b := &ocpp16.StartTransactionRequest{}

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

	tid, ok := c.currentConnectorTransactions[b.ConnectorId]

	if !ok {
		return nil, &msg.OCPP16CallError{
			MessageTypeID:    msg.CALLERROR,
			UniqueID:         req.UniqueID,
			ErrorCode:        msg.InternalError,
			ErrorDescription: "there is no transaction on the charge point",
			ErrorDetails:     struct{}{},
		}
	}

	_, err = c.transactionService.StartTransaction(context.Background(), &rpc.StartTransactionReq{
		Id:              int32(tid),
		StartMeterValue: int32(b.MeterStart),
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

	rb := &ocpp16.StartTransactionResponse{
		TransactionId: tid,
		IdTagInfo: &ocpp16.IdTagInfo{
			Status: "Accepted",
		},
	}

	return &msg.OCPP16CallResult{
		MessageTypeID: msg.CALLRESULT,
		UniqueID:      req.UniqueID,
		Payload:       rb,
	}, nil
}
