package handlers

import (
	"encoding/json"
	"time"

	msg "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_api/messaging"
	ocpp16 "github.com/Beep-Technologies/beepbeep3-ocpp/internal/schemas/ocpp_16"
)

const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

// BootNotification handles the Boot Notification operation, initiated by the charge point
func BootNotification(req *msg.OCPP16CallMessage) (*msg.OCPP16CallResult, *msg.OCPP16CallError) {
	b := &ocpp16.BootNotificationRequest{}

	p, err := json.Marshal(req.Payload)
	if err != nil {
		return nil, &msg.OCPP16CallError{
			MessageTypeID:    msg.CALLERROR,
			UniqueID:         req.UniqueID,
			ErrorCode:        msg.FormationViolation,
			ErrorDescription: "",
			ErrorDetails:     struct{}{},
		}
	}

	err = b.UnmarshalJSON(p)
	if err != nil {
		return nil, &msg.OCPP16CallError{
			MessageTypeID:    msg.CALLERROR,
			UniqueID:         req.UniqueID,
			ErrorCode:        msg.FormationViolation,
			ErrorDescription: "",
			ErrorDetails:     struct{}{},
		}
	}

	currentTime := time.Now().UTC().Format(RFC3339Milli)

	rb := &ocpp16.BootNotificationResponse{
		Status:      "Accepted",
		CurrentTime: currentTime,
		Interval:    300,
	}

	return &msg.OCPP16CallResult{
		MessageTypeID: msg.CALLRESULT,
		UniqueID:      req.UniqueID,
		Payload:       rb,
	}, nil
}
