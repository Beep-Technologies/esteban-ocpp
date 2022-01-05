package operation

import (
	"context"

	"gorm.io/gorm"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	ocpp16 "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16"
	chargepoint "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/charge_point"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/transaction"
)

type Service struct {
	db                 *gorm.DB
	chargePointService *chargepoint.Service
	transactionService *transaction.Service
	// statusNotificationService *statusnotification.Service
	ocpp16CentralSystem ocpp16.CentralSystem
}

func NewService(
	db *gorm.DB,
	cpSrv *chargepoint.Service,
	trSrv *transaction.Service,
	o16cs ocpp16.CentralSystem,
) *Service {
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
	// get the charge point
	cp, err := srv.ocpp16CentralSystem.GetChargePoint(req.EntityCode, req.ChargePointIdentifier)
	if err != nil {
		return nil, err
	}

	// call remote start transaction
	tid, err := cp.RemoteStartTransaction(int(req.ConnectorId))
	if err != nil {
		return nil, err
	}

	res := &rpc.RemoteStartTransactionResp{
		TransactionId: int32(tid),
	}

	return res, err
}

func (srv Service) RemoteStopTransaction(
	ctx context.Context,
	req *rpc.RemoteStopTransactionReq,
) (*rpc.RemoteStopTransactionResp, error) {
	// get the charge point
	cp, err := srv.ocpp16CentralSystem.GetChargePoint(req.EntityCode, req.ChargePointIdentifier)
	if err != nil {
		return nil, err
	}

	// call remote start transaction
	err = cp.RemoteStopTransaction(int(req.ConnectorId))
	if err != nil {
		return nil, err
	}

	res := &rpc.RemoteStopTransactionResp{}

	return res, err
}
