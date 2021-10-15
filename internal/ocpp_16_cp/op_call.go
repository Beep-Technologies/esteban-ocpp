package ocpp16cp

import (
	"errors"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	msg "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16_messaging"
)

// currentCall refers to the current outgoing call sent by the Central System
type currentCall struct {
	m     sync.Mutex
	timer *time.Timer
	call  *msg.OCPP16CallMessage
}

// makes a call to the charge point
// returns an error if there is already a call that has not been responded to / timed out
func (c *OCPP16ChargePoint) makeCall(m msg.OCPP16CallMessage) error {
	c.currentCall.m.Lock()
	defer c.currentCall.m.Unlock()

	// unparse the call message
	p, err := msg.UnparseOCPP16Call(m)
	if err != nil {
		return err
	}

	// return error if there is already a call in progress, to satisfy synchronicity constraints
	if c.currentCall.call != nil {
		return errors.New("call in progress has not been responded to / timed out")
	}

	// send the message
	err = c.conn.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return err
	}

	c.logger.Printf("[CALL: TO %d] %s", c.id, p)

	// set the current call
	c.currentCall.call = &m
	// cancel call after timeout
	c.currentCall.timer = time.AfterFunc(10*time.Second, c.cancelCall) //TODO: shift the timeout config somewhere else

	return nil
}

// cancels a call to the charge point
// this is called after a timeout
func (c *OCPP16ChargePoint) cancelCall() {
	c.currentCall.m.Lock()
	defer c.currentCall.m.Unlock()

	c.currentCall.timer = nil
	c.currentCall.call = nil
}

// handleCall handles calls (requests) initiated by the Charge Point
// there is no need to lock c.currentCall since that refers to the outgoing call
func (c *OCPP16ChargePoint) handleCall(req *msg.OCPP16CallMessage) (res *msg.OCPP16CallResult, err *msg.OCPP16CallError) {
	// go generics couldn't come sooner ....
	switch req.Action {
	case "Authorize":
		res, err = c.authorize(req)
	case "BootNotification":
		res, err = c.bootNotification(req)
	case "StatusNotification":
		res, err = c.statusNotification(req)
	case "Heartbeat":
		res, err = c.heartbeat(req)
	case "MeterValues":
		res, err = c.meterValues(req)
	case "StartTransaction":
		res, err = c.startTransaction(req)
	case "StopTransaction":
		res, err = c.stopTransaction(req)
	default:
		res, err = nil, &msg.OCPP16CallError{
			MessageTypeID:    msg.CALLERROR,
			UniqueID:         req.UniqueID,
			ErrorCode:        msg.NotImplemented,
			ErrorDescription: "",
			ErrorDetails:     struct{}{},
		}
	}

	return res, err
}

// handleResponse handles call results sent by the Charge Point
func (c *OCPP16ChargePoint) handleCallResult(res *msg.OCPP16CallResult) error {
	c.currentCall.m.Lock()
	defer c.currentCall.m.Unlock()

	if c.currentCall.call == nil {
		return errors.New("there is no call in progress. this call result's corresponding call might have timed out")
	}

	// stop the timer which will remove the call
	if c.currentCall.timer != nil {
		c.currentCall.timer.Stop()
	}

	var err error
	switch c.currentCall.call.Action {
	case "RemoteStartTransaction":
		err = c.remoteStartTransaction(res)
	case "RemoteStopTransaction":
		err = c.remoteStopTransaction(res)
	default:
		err = errors.New("this call result's handler has not been implemented")
	}

	// clear the currentCall
	c.currentCall.timer = nil
	c.currentCall.call = nil

	return err
}

// handleError handles call errors sent by the Charge Point
func (c *OCPP16ChargePoint) handleCallError(res *msg.OCPP16CallError) error {
	c.currentCall.m.Lock()
	defer c.currentCall.m.Unlock()

	if c.currentCall.call == nil {
		return errors.New("there is no call in progress. this call error's corresponding call might have timed out")
	}

	// stop the timer which will remove the call
	if c.currentCall.timer != nil {
		c.currentCall.timer.Stop()
	}

	var err error
	switch c.currentCall.call.Action {
	default:
		err = errors.New("this call error's handler has not been implemented")
	}

	// clear the currentCall
	c.currentCall.timer = nil
	c.currentCall.call = nil

	return err
}
