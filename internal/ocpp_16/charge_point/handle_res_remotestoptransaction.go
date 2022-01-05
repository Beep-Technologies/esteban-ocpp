package chargepoint

import (
	"errors"

	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging/schemas"
	"github.com/mitchellh/mapstructure"
)

func (cp *OCPP16ChargePoint) handleRemoteStopTransaction(call messaging.OCPP16CallMessage, msg messaging.OCPP16CallResult) error {
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

	// if rejected, return an error
	if p.Status != "Accepted" {
		return errors.New("remote stop transaction request was rejected")
	}

	return nil
}
