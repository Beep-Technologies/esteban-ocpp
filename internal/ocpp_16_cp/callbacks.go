package ocpp16cp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
)

// callbackEvents is the set of events (which correspond to the receipt of CP-initiated messages)
// which are currently supported
var callbackEvents = map[string]bool{
	"BootNotification":   true,
	"StartTransaction":   true,
	"StopTransaction":    true,
	"StatusNotification": true,
}

func isCallbackEvent(e string) bool {
	_, ok := callbackEvents[e]

	return ok
}

// makeCallback is meant to be called in a separate goroutine
func (c *OCPP16ChargePoint) makeCallback(event string, data interface{}) {
	// ignore if the callback event is not supported
	if !isCallbackEvent(event) {
		c.logger.Printf("[ERROR] Callback event %s is not supported\n", event)
		return
	}

	// get the callback url
	res, err := c.applicationService.GetApplicationCallbacks(
		context.Background(),
		&rpc.GetApplicationCallbacksReq{
			ApplicationId: c.applicationId,
		},
	)

	if err != nil {
		c.logger.Printf("[ERROR] Error while getting callback urls: %s\n", err.Error())
		return
	}

	url := ""

	for _, callback := range res.ApplicationCallbacks {
		if callback.CallbackEvent == event {
			url = callback.CallbackUrl
		}
	}

	if url == "" {
		c.logger.Printf("no callback found for event %s for application with id %s", event, c.applicationId)
		return
	}

	body, err := json.Marshal(map[string]interface{}{
		"charge_point_id":         c.id,
		"charge_point_identifier": c.chargePointIdentifier,
		"callback_event":          event,
		"callback_data":           data,
	})

	if err != nil {
		c.logger.Printf("[ERROR] Error while parsing callback data: %s\n", err.Error())
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		c.logger.Printf("[ERROR] Error while making callback: %s\n", err.Error())
		return
	}

	c.logger.Printf("%s callback to %s, result: %+v", event, url, resp)
}
