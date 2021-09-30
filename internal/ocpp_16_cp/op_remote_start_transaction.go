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
	_, ok := c.currentConnectorTransactions[cnid]
	if ok {
		return nil, errors.New("there is already an ongoing transaction at this connection")
	}

	// create transaction row on db
	t, err := c.transactionService.CreateTransaction(
		context.Background(),
		&rpc.CreateTransactionReq{
			ChargePointId: int32(c.id),
			ConnectorId:   int32(cnid),
		})

	if err != nil {
		return nil, err
	}

	// set the current transaction on this connector
	c.currentConnectorTransactions[cnid] = int(t.Id)

	c.logger.Printf("inserting transaction, %+v", c.currentConnectorTransactions)

	// make the RemoteStartTransaction call
	m := msg.OCPP16CallMessage{
		MessageTypeID: msg.CALL,
		UniqueID:      uuid.NewString(), // TODO: maybe use an incremementing ID instead of uuid?
		Action:        "RemoteStartTransaction",
		// TODO: figure out how to deal with IdTag
		Payload: &ocpp16.RemoteStartTransactionRequest{
			IdTag:       "TEST",
			ConnectorId: cnid,
		},
	}

	err = c.makeCall(m)
	if err != nil {
		return nil, err
	}

	res := &rpc.RemoteStartTransactionResp{
		TransactionId: t.Id,
	}

	return res, nil
}
