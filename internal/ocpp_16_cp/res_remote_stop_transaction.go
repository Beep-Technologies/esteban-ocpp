package ocpp16cp

import (
	"encoding/json"
	"errors"

	msg "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16_messaging"
	ocpp16 "github.com/Beep-Technologies/beepbeep3-ocpp/internal/schemas/ocpp_16"
)

// remoteStopTransaction handles the response from the Charge Point to RemoteStopTransaction
// c.currentCall.call is not nil and locked, so there is no need to acquire the lock or check for nil values
func (c *OCPP16ChargePoint) remoteStopTransaction(res *msg.OCPP16CallResult) error {
	b := &ocpp16.RemoteStartTransactionResponse{}
	p, err := json.Marshal(res.Payload)
	if err != nil {
		return err
	}

	err = b.UnmarshalJSON(p)
	if err != nil {
		return err
	}

	// if rejected, return an error
	if b.Status != "Accepted" {
		return errors.New("remote stop transaction request was rejected")
	}

	return nil
}
