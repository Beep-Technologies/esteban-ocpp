package ocpp16cs

import (
	"errors"
	"fmt"
	"log"

	"github.com/gorilla/websocket"

	ocpp16cp "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16_cp"
)

type OCPP16CentralSystem struct {
	logger       *log.Logger
	chargePoints map[string]*ocpp16cp.OCPP16ChargePoint
}

func NewOCPP16CentralSystem(l *log.Logger) *OCPP16CentralSystem {
	return &OCPP16CentralSystem{
		logger:       l,
		chargePoints: make(map[string]*ocpp16cp.OCPP16ChargePoint),
	}
}

func (cs *OCPP16CentralSystem) ConnectChargePoint(cpId string, conn *websocket.Conn) error {
	cs.chargePoints[cpId] = ocpp16cp.NewOCPP16ChargePoint(
		cpId,
		conn,
		cs.logger,
	)

	go cs.chargePoints[cpId].Listen()

	return nil
}

func (cs *OCPP16CentralSystem) GetChargePoint(cpId string) (*ocpp16cp.OCPP16ChargePoint, error) {
	cp, ok := cs.chargePoints[cpId]

	if !ok {
		errorMsg := fmt.Sprintf("charge point with id %s not found", cpId)
		return nil, errors.New(errorMsg)
	}

	return cp, nil
}
