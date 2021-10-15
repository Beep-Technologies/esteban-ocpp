package ocpp16cp

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	msg "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16_messaging"
	ocpp16 "github.com/Beep-Technologies/beepbeep3-ocpp/internal/schemas/ocpp_16"
)

// RemoteStartTransactionOp makes a RemoteStartTransaction call to the charge point
func (c *OCPP16ChargePoint) RemoteStartTransactionOp(cnid int) (*rpc.RemoteStartTransactionResp, error) {
	// check if there is already an ongoing transaction
	tRes, err := c.transactionService.OnGoingTransaction(
		context.Background(),
		&rpc.OngoingTransactionReq{
			ApplicationId:         int32(c.applicationId),
			ChargePointIdentifier: c.chargePointIdentifier,
			ConnectorId:           int32(cnid),
		})

	if err != nil {
		return nil, err
	}

	if tRes.OngoingTransaction {
		return nil, errors.New("there is already an ongoing transaction at the charge point connector")
	}

	// generate a unique id tag for this connection
	// id tags must be 20 characters long
	idTag := uuid.NewString()[0:20]

	// set the current transaction on this connector
	// create transaction row on db
	t, err := c.transactionService.CreateTransaction(
		context.Background(),
		&rpc.CreateTransactionReq{
			ChargePointId:   int32(c.id),
			ConnectorId:     int32(cnid),
			RemoteInitiated: true,
			IdTag:           idTag,
		})

	if err != nil {
		return nil, err
	}

	// make the RemoteStartTransaction call
	m := msg.OCPP16CallMessage{
		MessageTypeID: msg.CALL,
		UniqueID:      uuid.NewString(),
		Action:        "RemoteStartTransaction",
		Payload: &ocpp16.RemoteStartTransactionRequest{
			IdTag:       idTag,
			ConnectorId: cnid,
		},
	}

	err = c.makeCall(m)
	if err != nil {
		return nil, err
	}

	res := &rpc.RemoteStartTransactionResp{
		TransactionId: t.Transaction.Id,
	}

	return res, nil
}
