package ocpp16cp

import (
	"context"

	"github.com/google/uuid"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"

	msg "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16_messaging"
	ocpp16 "github.com/Beep-Technologies/beepbeep3-ocpp/internal/schemas/ocpp_16"
)

// RemoteStopTransaction makes a RemoteStopTransaction call to the charge point
func (c *OCPP16ChargePoint) RemoteStopTransactionOp(tid int) (*rpc.RemoteStopTransactionResp, error) {
	_, err := c.transactionService.GetTransactionById(context.Background(), &rpc.GetTransactionByIdReq{Id: int32(tid)})
	if err != nil {
		return nil, err
	}

	m := msg.OCPP16CallMessage{
		MessageTypeID: msg.CALL,
		UniqueID:      uuid.NewString(),
		Action:        "RemoteStopTransaction",
		Payload: &ocpp16.RemoteStopTransactionRequest{
			TransactionId: tid,
		},
	}

	err = c.makeCall(m)
	if err != nil {
		return nil, err
	}

	return &rpc.RemoteStopTransactionResp{}, err
}
