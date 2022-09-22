package transaction

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/models"
	chargepoint "github.com/Beep-Technologies/beepbeep3-ocpp/internal/repository/charge_point"
	transaction "github.com/Beep-Technologies/beepbeep3-ocpp/internal/repository/transaction"
	"github.com/Beep-Technologies/beepbeep3-ocpp/pkg/util"
)

const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

type Service struct {
	db          *gorm.DB
	transaction transaction.BaseRepo
	chargePoint chargepoint.BaseRepo
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db:          db,
		chargePoint: chargepoint.NewBaseRepo(db),
		transaction: transaction.NewBaseRepo(db),
	}
}

func (srv Service) GetOngoingTransaction(ctx context.Context, req *rpc.GetOngoingTransactionReq) (*rpc.GetOngoingTransactionResp, error) {
	// get charge point
	cp, err := srv.chargePoint.GetChargePoint(ctx, req.EntityCode, req.ChargePointIdentifier)
	if err != nil {
		return nil, err
	}

	t, err := srv.transaction.GetByChargePointIDConnectorStates(ctx, cp.ID, req.ConnectorId, []string{"CREATED", "STARTED"})
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		res := &rpc.GetOngoingTransactionResp{
			OngoingTransaction: false,
		}
		return res, nil
	}

	if err != nil {
		return nil, err
	}

	to := rpc.Transaction{}
	err = util.ConvertCopyStruct(&to, &t, map[string]util.ConverterFunc{})
	if err != nil {
		return nil, err
	}

	res := &rpc.GetOngoingTransactionResp{
		OngoingTransaction: true,
		Transaction:        &to,
	}

	return res, nil
}

func (srv Service) CreateTransaction(ctx context.Context, req *rpc.CreateTransactionReq) (*rpc.CreateTransactionResp, error) {
	// get charge point
	cp, err := srv.chargePoint.GetChargePoint(ctx, req.EntityCode, req.ChargePointIdentifier)
	if err != nil {
		return nil, err
	}

	t, err := srv.transaction.Create(ctx, models.OcppTransaction{
		ChargePointID:   cp.ID,
		ConnectorID:     req.ConnectorId,
		State:           "CREATED",
		RemoteInitiated: req.RemoteInitiated,
		IDTag:           req.IdTag,
	})

	if err != nil {
		return nil, err
	}

	_, err = srv.chargePoint.CreateIdTag(ctx, models.OcppChargePointIDTag{
		ChargePointID: cp.ID,
		IDTag:         req.IdTag,
	})

	if err != nil {
		return nil, err
	}

	res := &rpc.CreateTransactionResp{
		Transaction: &rpc.Transaction{
			Id:              t.ID,
			ChargePointId:   t.ChargePointID,
			ConnectorId:     t.ConnectorID,
			IdTag:           t.IDTag,
			State:           t.State,
			RemoteInitiated: t.RemoteInitiated,
			StartTimestamp:  t.StartTimestamp.UTC().Format(RFC3339Milli),
			StopTimestamp:   t.StopTimestamp.UTC().Format(RFC3339Milli),
			StartMeterValue: t.StartMeterValue,
			StopMeterValue:  t.StopMeterValue,
			StopReason:      t.StopReason,
		},
	}

	return res, nil
}

func (srv Service) StartTransaction(ctx context.Context, req *rpc.StartTransactionReq) (*rpc.StartTransactionResp, error) {
	t, err := srv.transaction.Update(
		ctx,
		req.Id,
		[]string{"state", "ongoing", "start_timestamp", "start_meter_value"},
		models.OcppTransaction{
			State:           "STARTED",
			StartTimestamp:  time.Now(),
			StartMeterValue: req.StartMeterValue,
		},
	)

	if err != nil {
		return nil, err
	}

	res := &rpc.StartTransactionResp{
		Transaction: &rpc.Transaction{
			Id:              t.ID,
			ChargePointId:   t.ChargePointID,
			ConnectorId:     t.ConnectorID,
			IdTag:           t.IDTag,
			State:           t.State,
			RemoteInitiated: t.RemoteInitiated,
			StartTimestamp:  t.StartTimestamp.UTC().Format(RFC3339Milli),
			StopTimestamp:   t.StopTimestamp.UTC().Format(RFC3339Milli),
			StartMeterValue: t.StartMeterValue,
			StopMeterValue:  t.StopMeterValue,
			StopReason:      t.StopReason,
		},
	}

	return res, err
}

func (srv Service) AbortTransaction(ctx context.Context, req *rpc.AbortTransactionReq) (*rpc.AbortTransactionResp, error) {
	t, err := srv.transaction.Update(
		ctx,
		req.Id,
		[]string{"state", "ongoing"},
		models.OcppTransaction{
			State: "ABORTED",
		},
	)

	if err != nil {
		return nil, err
	}

	res := &rpc.AbortTransactionResp{
		Transaction: &rpc.Transaction{
			Id:              t.ID,
			ChargePointId:   t.ChargePointID,
			ConnectorId:     t.ConnectorID,
			IdTag:           t.IDTag,
			State:           t.State,
			RemoteInitiated: t.RemoteInitiated,
			StartTimestamp:  t.StartTimestamp.UTC().Format(RFC3339Milli),
			StopTimestamp:   t.StopTimestamp.UTC().Format(RFC3339Milli),
			StartMeterValue: t.StartMeterValue,
			StopMeterValue:  t.StopMeterValue,
			StopReason:      t.StopReason,
		},
	}

	return res, err
}

func (srv Service) AbnormalStopTransaction(ctx context.Context, req *rpc.AbnormalStopTransactionReq) (*rpc.AbnormalStopTransactionResp, error) {
	t, err := srv.transaction.Update(
		ctx,
		req.Id,
		[]string{"state", "ongoing"},
		models.OcppTransaction{
			State: "ABNORMAL_STOPPED",
		},
	)

	if err != nil {
		return nil, err
	}

	res := &rpc.AbnormalStopTransactionResp{
		Transaction: &rpc.Transaction{
			Id:              t.ID,
			ChargePointId:   t.ChargePointID,
			ConnectorId:     t.ConnectorID,
			IdTag:           t.IDTag,
			State:           t.State,
			RemoteInitiated: t.RemoteInitiated,
			StartTimestamp:  t.StartTimestamp.UTC().Format(RFC3339Milli),
			StopTimestamp:   t.StopTimestamp.UTC().Format(RFC3339Milli),
			StartMeterValue: t.StartMeterValue,
			StopMeterValue:  t.StopMeterValue,
			StopReason:      t.StopReason,
		},
	}

	return res, err
}

func (srv Service) StopTransaction(ctx context.Context, req *rpc.StopTransactionReq) (*rpc.StopTransactionResp, error) {
	t, err := srv.transaction.Update(
		ctx,
		req.Id,
		[]string{"state", "stopped", "ongoing", "stop_timestamp", "stop_meter_value", "stop_reason"},
		models.OcppTransaction{
			State:          "STOPPED",
			StopTimestamp:  time.Now(),
			StopMeterValue: req.StopMeterValue,
			StopReason:     req.StopReason,
		},
	)

	if err != nil {
		return nil, err
	}

	res := &rpc.StopTransactionResp{
		Transaction: &rpc.Transaction{
			Id:              t.ID,
			ChargePointId:   t.ChargePointID,
			ConnectorId:     t.ConnectorID,
			IdTag:           t.IDTag,
			State:           t.State,
			StartTimestamp:  t.StartTimestamp.UTC().Format(RFC3339Milli),
			StopTimestamp:   t.StopTimestamp.UTC().Format(RFC3339Milli),
			StartMeterValue: t.StartMeterValue,
			StopMeterValue:  t.StopMeterValue,
			StopReason:      t.StopReason,
		},
	}

	return res, err
}
