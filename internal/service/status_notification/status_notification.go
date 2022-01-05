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
	statusNotification statusnotification.BaseRepo
	chargePoint        chargepoint.BaseRepo
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db:                 db,
		statusNotification: statusnotification.NewBaseRepo(db),
		chargePoint:        chargepoint.NewBaseRepo(db),
	}
}

func (srv Service) CreateStatusNotification(ctx context.Context, req *rpc.CreateStatusNotificationReq) (*rpc.CreateStatusNotificationResp, error) {
	// get charge point
	cp, err := srv.chargePoint.GetChargePoint(ctx, req.EntityCode, req.ChargePointIdentifier)
	if err != nil {
		return nil, err
	}

	// reported timestamp can be empty
	var rts time.Time
	if req.Timestamp == "" {
		rts = time.Time{}
	} else {
		rts, err = time.Parse(time.RFC3339Nano, req.Timestamp)
		if err != nil {
			return nil, err
		}
	}

	sn, err := srv.statusNotification.Create(ctx, models.OcppStatusNotification{
		ChargePointID:     cp.ID,
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
	// get charge point
	cp, err := srv.chargePoint.GetChargePoint(ctx, req.EntityCode, req.ChargePointIdentifier)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	sns, err := srv.statusNotification.GetLatestStatusNotifications(ctx, cp.ID)
	if err != nil {
		return nil, err
	}

	snsRes := make([]*rpc.StatusNotification, 0)

	for _, sn := range sns {
		resStatusNotification := &rpc.StatusNotification{
			Id:                sn.ID,
			ConnectorId:       sn.ConnectorID,
			ErrorCode:         sn.ErrorCode,
			Info:              sn.Info,
			Status:            sn.Status,
			VendorId:          sn.VendorID,
			VendorErrorCode:   sn.VendorErrorCode,
			Timestamp:         sn.Timestamp.UTC().Format(RFC3339Milli),
			ReportedTimestamp: sn.ReportedTimestamp.UTC().Format(RFC3339Milli),
		}

		snsRes = append(snsRes, resStatusNotification)
	}

	res := &rpc.GetLatestStatusNotificationsResp{
		ConnectorStatus: snsRes,
	}

	return res, nil
}
