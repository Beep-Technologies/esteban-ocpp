package operation

import (
	"context"

	"gorm.io/gorm"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	ocpp16cs "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16_cs"
	chargepoint "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/charge_point"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/transaction"
)

type Service struct {
	db                 *gorm.DB
	chargePointService *chargepoint.Service
	transactionService *transaction.Service
	// statusNotificationService *statusnotification.Service
	ocpp16CentralSystem *ocpp16cs.OCPP16CentralSystem
}

func NewService(
	db *gorm.DB,
	cpSrv *chargepoint.Service,
	trSrv *transaction.Service,
	o16cs *ocpp16cs.OCPP16CentralSystem) *Service {
	return &Service{
		db:                  db,
		chargePointService:  cpSrv,
		transactionService:  trSrv,
		ocpp16CentralSystem: o16cs,
	}
}

func (srv Service) RemoteStartTransaction(
	ctx context.Context,
	req *rpc.RemoteStartTransactionReq,
) (*rpc.RemoteStartTransactionResp, error) {
	// TODO: figure out how this interfaces with the database
	// for now we handle this locally

	// get the charge point
	cp, err := srv.ocpp16CentralSystem.GetChargePoint(int(req.ChargePointId))
	if err != nil {
		return nil, err
	}

	// make start transaction call (return error if there is a transaction ongoing)
	res, err := cp.RemoteStartTransactionOp(int(req.ConnectorId))

	return res, err
}

func (srv Service) RemoteStopTransaction(
	ctx context.Context,
	req *rpc.RemoteStopTransactionReq,
) (*rpc.RemoteStopTransactionResp, error) {
	// TODO: figure out how this interfaces with the database

	// get the charge point
	cp, err := srv.ocpp16CentralSystem.GetChargePoint(int(req.ChargePointId))
	if err != nil {
		return nil, err
	}

	res, err := cp.RemoteStopTransactionOp(int(req.TransactionId))

	return res, err
}
