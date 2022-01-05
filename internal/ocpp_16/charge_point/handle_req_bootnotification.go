package chargepoint

import (
	"time"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging/schemas"
	"github.com/mitchellh/mapstructure"
)

func (cp *OCPP16ChargePoint) handleBootNotification(msg messaging.OCPP16CallMessage) (*messaging.OCPP16CallResult, *messaging.OCPP16CallError) {
	p := &schemas.BootNotificationRequest{}

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

	// make callback
	go cp.makeCallback("BootNotification", p)

	cp.chargepointService.UpdateChargePointDetails(
		cp.ctx,
		&rpc.UpdateChargePointDetailsReq{
			EntityCode:              cp.entityCode,
			ChargePointIdentifier:   cp.chargePointIdentifier,
			ChargePointVendor:       p.ChargePointVendor,
			ChargePointModel:        p.ChargePointModel,
			ChargePointSerialNumber: p.ChargeBoxSerialNumber,
			ChargeBoxSerialNumber:   p.ChargeBoxSerialNumber,
			Iccid:                   p.Iccid,
			Imsi:                    p.Imsi,
			MeterType:               p.MeterType,
			MeterSerialNumber:       p.MeterSerialNumber,
			FirmwareVersion:         p.FirmwareVersion,
		},
	)

	now := time.Now().UTC().Format("2006-01-02T15:04:05.000Z07:00")

	return &messaging.OCPP16CallResult{
		MessageTypeID: messaging.CALLRESULT,
		UniqueID:      msg.UniqueID,
		Payload: &schemas.BootNotificationResponse{
			Status:      "Accepted",
			CurrentTime: now,
			Interval:    20,
		},
	}, nil
}
