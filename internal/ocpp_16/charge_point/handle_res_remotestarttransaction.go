package chargepoint

import (
	"errors"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging/schemas"
	"github.com/mitchellh/mapstructure"
)

func (cp *OCPP16ChargePoint) handleRemoteStartTransaction(call messaging.OCPP16CallMessage, msg messaging.OCPP16CallResult) error {
	p := &schemas.RemoteStartTransactionResponse{}

	// decode the payload into the struct
	err := mapstructure.Decode(msg.Payload, p)
	if err != nil {
		return err
	}

	c := &schemas.RemoteStartTransactionRequest{}
	err = mapstructure.Decode(msg.Payload, c)
	if err != nil {
		return err
	}

	tRes, err := cp.transactionService.GetOngoingTransaction(
		cp.ctx,
		&rpc.GetOngoingTransactionReq{
			EntityCode:            cp.entityCode,
			ChargePointIdentifier: cp.chargePointIdentifier,
			ConnectorId:           int32(c.ConnectorId),
		},
	)
	if err != nil {
		return err
	}

	// if there is no ongoing transaction, throw an error
	if !tRes.OngoingTransaction {
		return errors.New("remote start transaction response received with no ongoing transaction")
	}

	// if rejected, set the transaction status to ABORTED on the database
	if p.Status != "Accepted" {
		// set the transaction status to ABORTED on the database
		_, err = cp.transactionService.AbortTransaction(
			cp.ctx,
			&rpc.AbortTransactionReq{
				Id: int32(tRes.Transaction.Id),
			})
		if err != nil {
			return err
		}

		return errors.New("remote start transaction request was rejected. transaction was aborted")
	}

	return nil
}
