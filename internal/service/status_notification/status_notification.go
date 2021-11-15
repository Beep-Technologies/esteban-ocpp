package statusnotification

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/models"
	chargepoint "github.com/Beep-Technologies/beepbeep3-ocpp/internal/repository/charge_point"
	statusnotification "github.com/Beep-Technologies/beepbeep3-ocpp/internal/repository/status_notification"
)

const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

type Service struct {
	db                 *gorm.DB
	statusnotification statusnotification.BaseRepo
	chargepoint        chargepoint.BaseRepo
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db:                 db,
		statusnotification: statusnotification.NewBaseRepo(db),
		chargepoint:        chargepoint.NewBaseRepo(db),
	}
}

func (srv Service) CreateStatusNotification(ctx context.Context, req *rpc.CreateStatusNotificationReq) (*rpc.CreateStatusNotificationResp, error) {
	// reported timestamp can be empty
	var rts time.Time
	var err error

	if req.Timestamp == "" {
		rts = time.Time{}
	} else {
		rts, err = time.Parse(time.RFC3339Nano, req.Timestamp)
		if err != nil {
			return nil, err
		}
	}

	sn, err := srv.statusnotification.Create(ctx, models.OcppStatusNotification{
		ChargePointID:     req.ChargePointId,
		ConnectorID:       req.ConnectorId,
		ErrorCode:         req.ErrorCode,
		Info:              req.Info,
		Status:            req.Status,
		VendorID:          req.VendorId,
		VendorErrorCode:   req.VendorErrorCode,
		Timestamp:         time.Now(),
		ReportedTimestamp: rts,
	})

	if err != nil {
		return nil, err
	}

	res := &rpc.CreateStatusNotificationResp{
		StatusNotification: &rpc.StatusNotification{
			Id:                sn.ID,
			ChargePointId:     sn.ChargePointID,
			ConnectorId:       sn.ConnectorID,
			ErrorCode:         sn.ErrorCode,
			Info:              sn.Info,
			Status:            sn.Status,
			VendorId:          sn.VendorID,
			VendorErrorCode:   sn.VendorErrorCode,
			Timestamp:         sn.Timestamp.UTC().Format(time.RFC3339Nano),
			ReportedTimestamp: sn.ReportedTimestamp.UTC().Format(time.RFC3339Nano),
		},
	}

	return res, nil
}

func (srv Service) GetLatestStatusNotifications(ctx context.Context, req *rpc.GetLatestStatusNotificationsReq) (*rpc.GetLatestStatusNotificationsResp, error) {
	cpModel, err := srv.chargepoint.GetChargePointByIdentifier(ctx, req.ApplicationId, req.ChargePointIdentifier)

	if err != nil {
		return nil, err
	}

	snModels, err := srv.statusnotification.GetLatestStatusNotifications(ctx, cpModel.ID)
	if err != nil {
		return nil, err
	}

	snRpcs := make([]*rpc.StatusNotification, 0)

	for _, sn := range snModels {
		snRpc := &rpc.StatusNotification{
			Id:                sn.ID,
			ChargePointId:     sn.ChargePointID,
			ConnectorId:       sn.ConnectorID,
			ErrorCode:         sn.ErrorCode,
			Info:              sn.Info,
			Status:            sn.Status,
			VendorId:          sn.VendorID,
			VendorErrorCode:   sn.VendorErrorCode,
			Timestamp:         sn.Timestamp.UTC().Format(RFC3339Milli),
			ReportedTimestamp: sn.ReportedTimestamp.UTC().Format(RFC3339Milli),
		}

		snRpcs = append(snRpcs, snRpc)
	}

	res := &rpc.GetLatestStatusNotificationsResp{
		ConnectorStatus: snRpcs,
	}

	return res, nil
}
