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
	rts, err := time.Parse(time.RFC3339Nano, req.Timestamp)
	if err != nil {
		return nil, err
	}

	_, err = srv.statusnotification.Create(ctx, models.OcppStatusNotification{
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

	return &rpc.CreateStatusNotificationResp{}, nil
}

func (srv Service) GetLatestStatusNotifications(ctx context.Context, req *rpc.GetLatestStatusNotificationsReq) (*rpc.GetLatestStatusNotificationsResp, error) {
	cpModel, err := srv.chargepoint.GetByID(ctx, req.ChargePointId)
	if err != nil {
		return nil, err
	}

	snModels, err := srv.statusnotification.GetLatest(ctx, int(req.ChargePointId))
	if err != nil {
		return nil, err
	}

	snRpcs := make([]*rpc.StatusNotification, 0)

	for _, sn := range snModels {
		snRpc := &rpc.StatusNotification{
			ConnectorId:     sn.ChargePointID,
			ErrorCode:       sn.ErrorCode,
			Info:            sn.Info,
			Status:          sn.Status,
			Timestamp:       sn.Timestamp.UTC().Format(RFC3339Milli),
			VendorId:        sn.VendorID,
			VendorErrorCode: sn.VendorErrorCode,
		}

		snRpcs = append(snRpcs, snRpc)
	}

	res := &rpc.GetLatestStatusNotificationsResp{
		ChargePointId:           cpModel.ID,
		ChargePointVendor:       cpModel.ChargePointVendor,
		ChargePointModel:        cpModel.ChargePointModel,
		ChargePointSerialNumber: cpModel.ChargePointSerialNumber,
		ChargeBoxSerialNumber:   cpModel.ChargeBoxSerialNumber,
		Iccid:                   cpModel.Iccid,
		Imsi:                    cpModel.Imsi,
		MeterType:               cpModel.MeterType,
		MeterSerialNumber:       cpModel.MeterSerialNumber,
		FirmwareVersion:         cpModel.FirmwareVersion,
		OcppProtocol:            cpModel.OcppProtocol,
		ChargePointIdentifier:   cpModel.ChargePointIdentifier,
		Description:             cpModel.Description,
		LocationLatitude:        cpModel.LocationLatitude,
		LocationLongitude:       cpModel.LocationLongitude,
		ConnectorStatus:         snRpcs,
	}

	return res, nil
}
