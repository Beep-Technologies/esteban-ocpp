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

	otRes, err := c.transactionService.OnGoingTransaction(context.Background(), &rpc.OngoingTransactionReq{
		ApplicationId:         int32(c.applicationId),
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

	// if there is no ongoing transaction, this should not be a remote-initiated transaction
	if !otRes.OngoingTransaction {
		// TODO: allow non remote-initiated transactions
		return &msg.OCPP16CallResult{
			MessageTypeID: msg.CALLRESULT,
			UniqueID:      req.UniqueID,
			Payload: &ocpp16.StartTransactionResponse{
				IdTagInfo: &ocpp16.IdTagInfo{
					Status: "Rejected",
				},
			},
		}, nil
	}

	// else it should be a remote-initiated transaction
	tRes, err := c.transactionService.GetTransactionById(context.Background(),
		&rpc.GetTransactionByIdReq{Id: otRes.TransactionId},
	)
	if err != nil {
		return nil, &msg.OCPP16CallError{
			MessageTypeID:    msg.CALLERROR,
			UniqueID:         req.UniqueID,
			ErrorCode:        msg.InternalError,
			ErrorDescription: err.Error(),
			ErrorDetails:     struct{}{},
		}
	}

	// check if the idTags match, reject if they dont
	if b.IdTag != tRes.Transaction.IdTag {
		return &msg.OCPP16CallResult{
			MessageTypeID: msg.CALLRESULT,
			UniqueID:      req.UniqueID,
			Payload: &ocpp16.StartTransactionResponse{
				IdTagInfo: &ocpp16.IdTagInfo{
					Status: "Rejected",
				},
			},
		}, nil
	}

	_, err = c.transactionService.StartTransaction(context.Background(), &rpc.StartTransactionReq{
		Id:              int32(tRes.Transaction.Id),
		StartMeterValue: int32(b.MeterStart),
	})

	// make callback
	c.makeCallback("StartTransaction", b)

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
		TransactionId: int(tRes.Transaction.Id),
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
