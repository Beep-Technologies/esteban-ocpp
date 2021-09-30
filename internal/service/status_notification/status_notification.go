package statusnotification

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/models"
	statusnotification "github.com/Beep-Technologies/beepbeep3-ocpp/internal/repository/status_notification"
)

type Service struct {
	db                 *gorm.DB
	statusnotification statusnotification.BaseRepo
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db:                 db,
		statusnotification: statusnotification.NewBaseRepo(db),
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
