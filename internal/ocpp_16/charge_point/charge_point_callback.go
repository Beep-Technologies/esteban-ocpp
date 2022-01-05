package chargepoint

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	"go.uber.org/zap"
)

func (cp *OCPP16ChargePoint) makeCallback(event string, data interface{}) {
	// get callback urls
	callbacks, err := cp.applicationService.GetApplicationCallbacks(
		cp.ctx,
		&rpc.GetApplicationCallbacksReq{
			EntityCode:    cp.entityCode,
			CallbackEvent: event,
		},
	)
	if err != nil {
		return
	}

	body, _ := json.Marshal(map[string]interface{}{
		"charge_point_id":         cp.id,
		"entity_code":             cp.entityCode,
		"charge_point_identifier": cp.chargePointIdentifier,
		"callback_event":          event,
		"callback_data":           data,
	})

	for _, callback := range callbacks.ApplicationCallbacks {
		url := callback.CallbackUrl
		res, err := http.Post(url, "application/json", bytes.NewBuffer(body))

		resString := ""
		errString := ""

		if err != nil {
			errString = err.Error()
		} else {
			resBody, _ := ioutil.ReadAll(res.Body)
			defer res.Body.Close()
			resString = string(resBody)
		}

		cp.logger.Info(
			"POST "+url,
			zap.String("event", "application_callback"),
			zap.String("callback_event", event),
			zap.String("callback_url", url),
			zap.String("callback_res", resString),
			zap.String("callback_err", errString),
		)
	}
}
