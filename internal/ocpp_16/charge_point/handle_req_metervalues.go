package chargepoint

import (
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging/schemas"
	"github.com/mitchellh/mapstructure"
)

func (cp *OCPP16ChargePoint) handleMeterValues(msg messaging.OCPP16CallMessage) (*messaging.OCPP16CallResult, *messaging.OCPP16CallError) {
	p := &schemas.MeterValuesRequest{}

	// decode the payload into the struct
	err := mapstructure.Decode(msg.Payload, p)
	if err != nil {
		return nil, &messaging.OCPP16CallError{
			MessageTypeID:    messaging.CALLERROR,
			UniqueID:         msg.UniqueID,
			ErrorCode:        messaging.FormationViolation,
			ErrorDescription: "",
			ErrorDetails:     struct{}{},
		}
	}

	return &messaging.OCPP16CallResult{
		MessageTypeID: messaging.CALLRESULT,
		UniqueID:      msg.UniqueID,
		Payload:       &schemas.MeterValuesResponse{},
	}, nil
}
