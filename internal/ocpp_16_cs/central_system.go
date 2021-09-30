package ocpp16cs

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/gorilla/websocket"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	ocpp16cp "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16_cp"
	chargepointsrv "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/charge_point"
	statusnotificationsrv "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/status_notification"
	transactionsrv "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/transaction"
)

type OCPP16CentralSystem struct {
	logger                    *log.Logger
	chargePoints              map[int]*ocpp16cp.OCPP16ChargePoint
	chargePointService        *chargepointsrv.Service
	transactionService        *transactionsrv.Service
	statusNotificationService *statusnotificationsrv.Service
}

func NewOCPP16CentralSystem(
	l *log.Logger,
	cpSrv *chargepointsrv.Service,
	trSrv *transactionsrv.Service,
	snSrv *statusnotificationsrv.Service,
) *OCPP16CentralSystem {
	return &OCPP16CentralSystem{
		logger:                    l,
		chargePoints:              make(map[int]*ocpp16cp.OCPP16ChargePoint),
		chargePointService:        cpSrv,
		transactionService:        trSrv,
		statusNotificationService: snSrv,
	}
}

func (cs *OCPP16CentralSystem) ConnectChargePoint(chargePointIdentifier string, conn *websocket.Conn) error {
	// get charge point by identifier
	// if there is no matching identifier, immediately disconnect the charge point
	ctx := context.Background()
	req := &rpc.GetChargePointByIdentifierReq{
		ChargePointIdentifier: chargePointIdentifier,
	}

	cpo, err := cs.chargePointService.GetChargePointByIdentifier(ctx, req)
	if err != nil {
		conn.Close()
		return err
	}

	// if there is a matching identifier, create a new charge point and add it to the map
	id := int(cpo.Id)
	cs.chargePoints[id] = ocpp16cp.NewOCPP16ChargePoint(
		id,
		chargePointIdentifier,
		conn,
		cs.transactionService,
		cs.statusNotificationService,
		cs.logger,
	)

	go cs.chargePoints[id].Listen()

	return nil
}

func (cs *OCPP16CentralSystem) GetChargePoint(id int) (*ocpp16cp.OCPP16ChargePoint, error) {
	cp, ok := cs.chargePoints[id]

	if !ok {
		errorMsg := fmt.Sprintf("charge point with id %d not found", id)
		return nil, errors.New(errorMsg)
	}

	return cp, nil
}
