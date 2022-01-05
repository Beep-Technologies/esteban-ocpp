package chargepoint

import (
	"fmt"

	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging"
)

func (cp *OCPP16ChargePoint) handleBootNotification(msg messaging.OCPP16CallMessage) (*messaging.OCPP16CallResult, *messaging.OCPP16CallError) {
	// STUB
	fmt.Println("bootnotification")
	return nil, nil
}
