package chargepoint

import "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging"

// listenCP is invoked in the OCPP16ChargePoint.Listen() as a goroutine
func (cp *OCPP16ChargePoint) listenCP() {
	for {
		select {
		case <-cp.ctx.Done():
			return
		case msg := <-cp.inCallStream:
			res, err := cp.handleCall(msg)
			if err != nil {
				cp.outCallErrorStream <- *err
				continue
			}
			if res != nil {
				cp.outCallResultStream <- *res
				continue
			}
		}
	}
}

func (cp *OCPP16ChargePoint) handleCall(msg messaging.OCPP16CallMessage) (*messaging.OCPP16CallResult, *messaging.OCPP16CallError) {
	switch msg.Action {
	case "Authorize":
		return cp.handleAuthorize(msg)
	case "BootNotification":
		return cp.handleBootNotification(msg)
	case "Heartbeat":
		return cp.handleHeartbeat(msg)
	case "MeterValues":
		return cp.handleMeterValues(msg)
	case "StartTransaction":
		return cp.handleStartTransaction(msg)
	case "StatusNotification":
		return cp.handleStatusNotification(msg)
	case "StopTransaction":
		return cp.handleStopTransaction(msg)
	}
	return nil, &messaging.OCPP16CallError{
		MessageTypeID:    messaging.CALLERROR,
		UniqueID:         msg.UniqueID,
		ErrorCode:        messaging.InternalError,
		ErrorDescription: "not implemented",
	}
}
