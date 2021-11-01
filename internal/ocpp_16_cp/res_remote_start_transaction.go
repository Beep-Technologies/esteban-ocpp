package ocpp16cp

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	msg "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16_messaging"
	ocpp16 "github.com/Beep-Technologies/beepbeep3-ocpp/internal/schemas/ocpp_16"
)

// remoteStartTransaction handles the response from the Charge Point to RemoteStartTransaction
func (c *OCPP16ChargePoint) remoteStartTransaction(res *msg.OCPP16CallResult) error {
	b := &ocpp16.RemoteStartTransactionResponse{}
	p, err := json.Marshal(res.Payload)
	if err != nil {
		return err
	}

	err = b.UnmarshalJSON(p)
	if err != nil {
		return err
	}

	// get the connector id
	br := &ocpp16.RemoteStartTransactionRequest{}
	pr, err := json.Marshal(c.currentCall.call.Payload)
	if err != nil {
		return err
	}

	err = br.UnmarshalJSON(pr)
	if err != nil {
		return err
	}

	tRes, err := c.transactionService.GetOngoingTransaction(context.Background(), &rpc.GetOngoingTransactionReq{
		ApplicationId:         c.applicationId,
		ChargePointIdentifier: c.chargePointIdentifier,
		ConnectorId:           int32(br.ConnectorId),
	})
	if err != nil {
		return err
	}

	// if rejected, set the transaction status to ABORTED on the database
	if b.Status != "Accepted" {
		// set the transaction status to ABORTED on the database
		_, err = c.transactionService.AbortTransaction(context.Background(), &rpc.AbortTransactionReq{
			Id: int32(tRes.Transaction.Id),
		})
		if err != nil {
			return err
		}

		return errors.New("remote start transaction request was rejected. transaction was aborted")
	}

	return nil
}
