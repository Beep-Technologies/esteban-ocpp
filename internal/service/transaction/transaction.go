package transaction

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/models"
	transactionRepo "github.com/Beep-Technologies/beepbeep3-ocpp/internal/repository/transaction"
)

const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

type Service struct {
	db              *gorm.DB
	transactionRepo transactionRepo.BaseRepo
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db:              db,
		transactionRepo: transactionRepo.NewBaseRepo(db),
	}
}

func (srv Service) CreateTransaction(ctx context.Context, req *rpc.CreateTransactionReq) (*rpc.CreateTransactionResp, error) {
	t, err := srv.transactionRepo.Create(ctx, models.OcppTransaction{
		ChargePointID: req.ChargePointId,
		ConnectorID:   req.ConnectorId,
		State:         "CREATED",
	})

	if err != nil {
		return nil, err
	}

	res := &rpc.CreateTransactionResp{
		Id:              t.ID,
		ChargePointId:   t.ChargePointID,
		ConnectorId:     t.ConnectorID,
		IdTag:           t.IDTag,
		StartTimestamp:  "",
		StopTimestamp:   "",
		StartMeterValue: t.StartMeterValue,
		StopMeterValue:  t.StopMeterValue,
		StopReason:      t.StopReason,
	}

	return res, nil
}

func (srv Service) StartTransaction(ctx context.Context, req *rpc.StartTransactionReq) (*rpc.StartTransactionResp, error) {
	_, err := srv.transactionRepo.Update(
		ctx,
		req.Id,
		[]string{"state", "ongoing", "start_timestamp", "start_meter_value"},
		models.OcppTransaction{
			State:           "STARTED",
			Ongoing:         true,
			StartTimestamp:  time.Now(),
			StartMeterValue: req.StartMeterValue,
		},
	)

	if err != nil {
		return nil, err
	}

	res := &rpc.StartTransactionResp{}

	return res, err
}

func (srv Service) AbortTransaction(ctx context.Context, req *rpc.AbortTransactionReq) (*rpc.AbortTransactionResp, error) {
	_, err := srv.transactionRepo.Update(
		ctx,
		req.Id,
		[]string{"state", "ongoing"},
		models.OcppTransaction{
			State:   "ABORTED",
			Ongoing: false,
		},
	)

	if err != nil {
		return nil, err
	}

	res := &rpc.AbortTransactionResp{}

	return res, err
}

func (srv Service) StopTransaction(ctx context.Context, req *rpc.StopTransactionReq) (*rpc.StopTransactionResp, error) {
	_, err := srv.transactionRepo.Update(
		ctx,
		req.Id,
		[]string{"state", "stopped", "ongoing", "stop_timestamp", "stop_meter_value", "stop_reason"},
		models.OcppTransaction{
			State:          "STOPPED",
			Ongoing:        false,
			StopTimestamp:  time.Now(),
			StopMeterValue: req.StopMeterValue,
			StopReason:     req.StopReason,
		},
	)

	if err != nil {
		return nil, err
	}

	res := &rpc.StopTransactionResp{}

	return res, err
}
