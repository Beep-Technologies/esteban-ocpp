package ocpp16cp

import (
	"errors"
	"time"

	"github.com/gorilla/websocket"

	msg "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16_messaging"
)

// makes a call to the charge point
// returns an error if there is already a call that has not been responded to / timed out
func (c *OCPP16ChargePoint) makeCall(m msg.OCPP16CallMessage) error {
	c.currentCall.m.Lock()
	defer c.currentCall.m.Unlock()

	p, err := msg.UnparseOCPP16Call(m)
	if err != nil {
		return err
	}

	if c.currentCall.uniqueID != "" {
		return errors.New("call in progress has not been responded to / timed out")
	}

	err = c.conn.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return err
	}

	c.logger.Printf("[CALL: TO %s] %s", c.cpId, p)

	// TODO: Move timer configs to a dedicated config source
	// cancel call after timeout
	time.AfterFunc(10*time.Second, c.cancelCall)

	return nil
}

// cancels a call to the charge point
func (c *OCPP16ChargePoint) cancelCall() {
	c.currentCall.m.Lock()
	defer c.currentCall.m.Unlock()

	c.currentCall.timer = nil
	c.currentCall.uniqueID = ""
}
