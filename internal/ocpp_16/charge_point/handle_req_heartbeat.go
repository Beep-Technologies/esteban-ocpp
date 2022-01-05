package chargepoint

import (
	"time"

	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging/schemas"
	"github.com/mitchellh/mapstructure"
)

func (cp *OCPP16ChargePoint) handleHeartbeat(msg messaging.OCPP16CallMessage) (*messaging.OCPP16CallResult, *messaging.OCPP16CallError) {
	p := &schemas.HeartbeatRequest{}

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

	now := time.Now().UTC().Format("2006-01-02T15:04:05.000Z07:00")

	return &messaging.OCPP16CallResult{
		MessageTypeID: messaging.CALLRESULT,
		UniqueID:      msg.UniqueID,
		Payload: &schemas.HeartbeatResponse{
			CurrentTime: now,
		},
	}, nil
}
