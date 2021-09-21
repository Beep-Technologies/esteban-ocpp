package ocpp16cp

import (
	"github.com/google/uuid"

	msg "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16_messaging"
	ocpp16 "github.com/Beep-Technologies/beepbeep3-ocpp/internal/schemas/ocpp_16"
)

// RemoteStartTransaction makes a RemoteStartTransaction call to the charge point
func (c *OCPP16ChargePoint) RemoteStartTransaction(connectorId int) error {
	m := msg.OCPP16CallMessage{
		MessageTypeID: msg.CALL,
		UniqueID:      uuid.NewString(),
		Action:        "RemoteStartTransaction",
		// TODO: figure out how to deal with IdTag
		Payload: &ocpp16.RemoteStartTransactionRequest{
			IdTag:       "TEST",
			ConnectorId: connectorId,
		},
	}

	err := c.makeCall(m)
	if err != nil {
		return err
	}

	return nil
}
