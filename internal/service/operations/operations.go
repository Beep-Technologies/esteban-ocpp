package operations

import (
	"context"

	"gorm.io/gorm"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/models"
	ocpp16cs "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16_cs"
	chargepoint "github.com/Beep-Technologies/beepbeep3-ocpp/internal/repository/charge_point"
	statusnotification "github.com/Beep-Technologies/beepbeep3-ocpp/internal/repository/status_notification"
	transaction "github.com/Beep-Technologies/beepbeep3-ocpp/internal/repository/transaction"
)

type Service struct {
	db                  *gorm.DB
	chargePoint         chargepoint.BaseRepo
	statusNotification  statusnotification.BaseRepo
	transaction         transaction.BaseRepo
	ocpp16CentralSystem *ocpp16cs.OCPP16CentralSystem
}

func NewService(db *gorm.DB, o16cs *ocpp16cs.OCPP16CentralSystem) *Service {
	return &Service{
		db:                  db,
		chargePoint:         chargepoint.NewBaseRepo(db),
		statusNotification:  statusnotification.NewBaseRepo(db),
		transaction:         transaction.NewBaseRepo(db),
		ocpp16CentralSystem: o16cs,
	}
}

func (srv Service) RemoteStartTransaction(ctx context.Context, req rpc.RemoteStartTransactionReq) (*rpc.RemoteStartTransactionResp, error) {
	cpObj, err := srv.chargePoint.GetByChargePointIdentifier(ctx, req.CpId)
	if err != nil {
		return nil, err
	}

	cp, err := srv.ocpp16CentralSystem.GetChargePoint(req.CpId)
	if err != nil {
		return nil, err
	}

	err = cp.RemoteStartTransaction(int(req.ConnectorId))
	if err != nil {
		return nil, err
	}

	t, err := srv.transaction.Create(ctx, models.Transaction{
		ChargePointID: cpObj.ID,
		ConnectorID:   int32(req.ConnectorId),
		// TODO: add IDTag
	})
	if err != nil {
		return nil, err
	}

	res := &rpc.RemoteStartTransactionResp{
		TransactionId: t.ID,
	}

	return res, nil
}

func (srv Service) RemoteStopTransaction(ctx context.Context, req rpc.RemoteStopTransactionReq) (*rpc.RemoteStopTransactionResp, error) {
	t, err := srv.transaction.GetByID(ctx, req.TransactionId)
	if err != nil {
		return nil, err
	}

	cpObj, err := srv.chargePoint.GetByID(ctx, t.ChargePointID)
	if err != nil {
		return nil, err
	}

	cp, err := srv.ocpp16CentralSystem.GetChargePoint(cpObj.ChargePointIdentifier)
	if err != nil {
		return nil, err
	}

	err = cp.RemoteStopTransaction(int(t.ID))
	if err != nil {
		return nil, err
	}

	res := &rpc.RemoteStopTransactionResp{}

	return res, nil
}
