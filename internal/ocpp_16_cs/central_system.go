package ocpp16cs

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/gorilla/websocket"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	ocpp16cp "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16_cp"
	applicationsrv "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/application"
	chargepointsrv "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/charge_point"
	statusnotificationsrv "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/status_notification"
	transactionsrv "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/transaction"
)

type OCPP16CentralSystem struct {
	logger                    *log.Logger
	chargePoints              map[int]*ocpp16cp.OCPP16ChargePoint
	applicationService        *applicationsrv.Service
	chargePointService        *chargepointsrv.Service
	transactionService        *transactionsrv.Service
	statusNotificationService *statusnotificationsrv.Service
}

func NewOCPP16CentralSystem(
	l *log.Logger,
	aSrv *applicationsrv.Service,
	cpSrv *chargepointsrv.Service,
	trSrv *transactionsrv.Service,
	snSrv *statusnotificationsrv.Service,
) *OCPP16CentralSystem {
	return &OCPP16CentralSystem{
		logger:                    l,
		chargePoints:              make(map[int]*ocpp16cp.OCPP16ChargePoint),
		applicationService:        aSrv,
		chargePointService:        cpSrv,
		transactionService:        trSrv,
		statusNotificationService: snSrv,
	}
}

func (cs *OCPP16CentralSystem) ConnectChargePoint(
	applicationId string,
	chargePointIdentifier string,
	conn *websocket.Conn) error {

	// check that both the application uuid and charge point identifier
	// correspond to an application and a charge point,
	// and that the charge point belongs to the application
	ctx := context.Background()

	// get the application
	ao, aerr := cs.applicationService.GetApplicationByID(ctx, &rpc.GetApplicationByIdReq{
		ApplicationId: applicationId,
	})

	// get the charge point
	cpo, cerr := cs.chargePointService.GetChargePointByIdentifier(ctx, &rpc.GetChargePointByIdentifierReq{
		ApplicationId:         ao.Application.Id,
		ChargePointIdentifier: chargePointIdentifier,
	})

	if aerr != nil || cerr != nil || cpo.ChargePoint.ApplicationId != ao.Application.Id {
		conn.Close()

		if aerr != nil {
			return aerr
		}

		if cerr != nil {
			return cerr
		}

		return errors.New("charge point identifier does not correspond to the application")
	}

	// if there is a valid match, create a new charge point and add it to the map
	id := int(cpo.ChargePoint.Id)
	cs.chargePoints[id] = ocpp16cp.NewOCPP16ChargePoint(
		id,
		chargePointIdentifier,
		ao.Application.Id,
		conn,
		cs.applicationService,
		cs.chargePointService,
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
