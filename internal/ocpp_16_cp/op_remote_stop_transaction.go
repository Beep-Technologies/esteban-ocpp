package ocpp16cp

import (
	"errors"

	"github.com/google/uuid"

	msg "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16_messaging"
	ocpp16 "github.com/Beep-Technologies/beepbeep3-ocpp/internal/schemas/ocpp_16"
)

// RemoteStopTransaction makes a RemoteStopTransaction call to the charge point
func (c *OCPP16ChargePoint) RemoteStopTransaction(transactionId int) error {
	if c.currentTransaction == nil || c.currentTransaction.transactionId != transactionId {
		return errors.New("no ongoing transaction with the id")
	}

	m := msg.OCPP16CallMessage{
		MessageTypeID: msg.CALL,
		UniqueID:      uuid.NewString(),
		Action:        "RemoteStopTransaction",
		Payload: &ocpp16.RemoteStopTransactionRequest{
			TransactionId: transactionId,
		},
	}

	err := c.makeCall(m)
	if err != nil {
		return err
	}

	return nil
}
