package ocpp16cp

import (
	"context"
	"encoding/json"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	msg "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16_messaging"
	ocpp16 "github.com/Beep-Technologies/beepbeep3-ocpp/internal/schemas/ocpp_16"
)

// heartbeat handles the Heartbeat operation, initiated by the charge point
func (c *OCPP16ChargePoint) authorize(req *msg.OCPP16CallMessage) (*msg.OCPP16CallResult, *msg.OCPP16CallError) {
	b := &ocpp16.AuthorizeRequest{}

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

	res, err := c.chargepointService.GetChargePointIdTag(context.Background(), &rpc.GetChargePointIdTagReq{
		ChargePointIdentifier: c.chargePointIdentifier,
		ApplicationId:         c.applicationId,
		IdTag:                 b.IdTag,
	})

	if err != nil {
		return nil, &msg.OCPP16CallError{
			MessageTypeID:    msg.CALLERROR,
			UniqueID:         req.UniqueID,
			ErrorCode:        msg.FormationViolation,
			ErrorDescription: "",
			ErrorDetails:     struct{}{},
		}
	}

	rb := &ocpp16.AuthorizeResponse{
		IdTagInfo: &ocpp16.IdTagInfo{},
	}

	if res.ChargePointIdTag != nil {
		rb.IdTagInfo.Status = "Accepted"
	} else {
		rb.IdTagInfo.Status = "Invalid"
	}

	return &msg.OCPP16CallResult{
		MessageTypeID: msg.CALLRESULT,
		UniqueID:      req.UniqueID,
		Payload:       rb,
	}, nil
}
