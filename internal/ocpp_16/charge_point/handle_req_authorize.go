package chargepoint

import (
	"github.com/mitchellh/mapstructure"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging/schemas"
)

func (cp *OCPP16ChargePoint) handleAuthorize(msg messaging.OCPP16CallMessage) (*messaging.OCPP16CallResult, *messaging.OCPP16CallError) {
	p := &schemas.AuthorizeRequest{}

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

	// get the id tags associated with this charge point
	idTagsRes, err := cp.chargepointService.GetChargePointIdTags(cp.ctx, &rpc.GetChargePointIdTagsReq{
		EntityCode:            cp.entityCode,
		ChargePointIdentifier: cp.chargePointIdentifier,
	})
	if err != nil {
		return nil, &messaging.OCPP16CallError{
			MessageTypeID:    messaging.CALLERROR,
			UniqueID:         msg.UniqueID,
			ErrorCode:        messaging.InternalError,
			ErrorDescription: "",
			ErrorDetails:     struct{}{},
		}
	}

	match := false
	for _, idTag := range idTagsRes.ChargePointIdTags {
		if idTag.IdTag == p.IdTag {
			match = true
			break
		}
	}

	// if requested id tag is not found, return invalid
	// else return accepted
	status := "Invalid"
	if match {
		status = "Accepted"
	}

	return &messaging.OCPP16CallResult{
		MessageTypeID: messaging.CALLRESULT,
		UniqueID:      msg.UniqueID,
		Payload: &schemas.AuthorizeResponse{
			IdTagInfo: &schemas.IdTagInfo{
				Status: status,
			},
		},
	}, nil
}
