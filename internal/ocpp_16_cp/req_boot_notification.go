package ocpp16cp

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	msg "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16_messaging"
	ocpp16 "github.com/Beep-Technologies/beepbeep3-ocpp/internal/schemas/ocpp_16"
)

// convenience function to get current system time
func getCurrentTime() string {
	RFC3339Milli := "2006-01-02T15:04:05.000Z07:00"
	return time.Now().UTC().Format(RFC3339Milli)
}

// bootNotification handles the Boot Notification operation, initiated by the charge point
func (c *OCPP16ChargePoint) bootNotification(req *msg.OCPP16CallMessage) (*msg.OCPP16CallResult, *msg.OCPP16CallError) {
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

	// make callback
	go c.makeCallback("BootNotification", b)

	// TODO: figure out if it is necessary to handle errors
	c.chargepointService.UpdateChargePointDetails(
		context.Background(),
		&rpc.UpdateChargePointDetailsReq{
			ChargePointId:           int32(c.id),
			ChargePointVendor:       b.ChargePointVendor,
			ChargePointModel:        b.ChargePointModel,
			ChargePointSerialNumber: b.ChargeBoxSerialNumber,
			ChargeBoxSerialNumber:   b.ChargeBoxSerialNumber,
			Iccid:                   b.Iccid,
			Imsi:                    b.Imsi,
			MeterType:               b.MeterType,
			MeterSerialNumber:       b.MeterSerialNumber,
			FirmwareVersion:         b.FirmwareVersion,
		})

	// TODO: move heartbeat interval config somewhere else
	rb := &ocpp16.BootNotificationResponse{
		Status:      "Accepted",
		CurrentTime: getCurrentTime(),
		Interval:    20,
	}

	return &msg.OCPP16CallResult{
		MessageTypeID: msg.CALLRESULT,
		UniqueID:      req.UniqueID,
		Payload:       rb,
	}, nil
}
