package ocpp16cp

import (
	"context"
	"encoding/json"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	msg "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16_messaging"
	ocpp16 "github.com/Beep-Technologies/beepbeep3-ocpp/internal/schemas/ocpp_16"
)

// stopTransaction handles the Stop Transaction operation, initiated by the charge point
func (c *OCPP16ChargePoint) stopTransaction(req *msg.OCPP16CallMessage) (*msg.OCPP16CallResult, *msg.OCPP16CallError) {
	b := &ocpp16.StopTransactionRequest{}

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

	connectorId := c.GetTransactionConnectorID(b.TransactionId)
	transactionExists := connectorId != 0

	if !transactionExists {
		return nil, &msg.OCPP16CallError{
			MessageTypeID:    msg.CALLERROR,
			UniqueID:         req.UniqueID,
			ErrorCode:        msg.InternalError,
			ErrorDescription: "there is no transaction on the charge point",
			ErrorDetails:     struct{}{},
		}
	}

	_, err = c.transactionService.StopTransaction(context.Background(), &rpc.StopTransactionReq{
		Id:             int32(b.TransactionId),
		StopMeterValue: int32(b.MeterStop),
		StopReason:     b.Reason,
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

	delete(c.currentConnectorTransactions, connectorId)

	rb := &ocpp16.StopTransactionResponse{}

	return &msg.OCPP16CallResult{
		MessageTypeID: msg.CALLRESULT,
		UniqueID:      req.UniqueID,
		Payload:       rb,
	}, nil
}
