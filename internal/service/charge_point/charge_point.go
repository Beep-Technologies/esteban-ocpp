package chargepoint

import (
	"context"

	"gorm.io/gorm"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/models"
	chargepoint "github.com/Beep-Technologies/beepbeep3-ocpp/internal/repository/charge_point"
)

type Service struct {
	db          *gorm.DB
	chargePoint chargepoint.BaseRepo
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db:          db,
		chargePoint: chargepoint.NewBaseRepo(db),
	}
}

func (srv Service) CreateChargePoint(ctx context.Context, req *rpc.CreateChargePointReq) (*rpc.CreateChargePointResp, error) {
	cp, err := srv.chargePoint.Create(ctx, models.OcppChargePoint{
		ApplicationID:         req.ApplicationId,
		ChargePointIdentifier: req.ChargePointIdentifier,
		OcppProtocol:          req.OcppProtocol,
	})

	if err != nil {
		return nil, err
	}

	res := &rpc.CreateChargePointResp{
		ChargePoint: &rpc.ChargePoint{
			Id:                      cp.ID,
			ApplicationId:           cp.ApplicationID,
			ChargePointVendor:       cp.ChargePointVendor,
			ChargePointModel:        cp.ChargePointModel,
			ChargePointSerialNumber: cp.ChargePointSerialNumber,
			ChargeBoxSerialNumber:   cp.ChargeBoxSerialNumber,
			Iccid:                   cp.Iccid,
			Imsi:                    cp.Imsi,
			MeterType:               cp.MeterType,
			MeterSerialNumber:       cp.MeterSerialNumber,
			FirmwareVersion:         cp.FirmwareVersion,
			ChargePointIdentifier:   cp.ChargePointIdentifier,
			OcppProtocol:            cp.OcppProtocol,
		},
	}

	return res, nil
}

func (srv Service) GetChargePointByIdentifier(ctx context.Context, req *rpc.GetChargePointByIdentifierReq) (*rpc.GetChargePointByIdentifierResp, error) {
	cp, err := srv.chargePoint.GetChargePointByIdentifier(ctx, req.ApplicationId, req.ChargePointIdentifier)
	if err != nil {
		return nil, err
	}

	res := &rpc.GetChargePointByIdentifierResp{
		ChargePoint: &rpc.ChargePoint{
			Id:                      cp.ID,
			ApplicationId:           cp.ApplicationID,
			ChargePointVendor:       cp.ChargePointVendor,
			ChargePointModel:        cp.ChargePointModel,
			ChargePointSerialNumber: cp.ChargePointSerialNumber,
			ChargeBoxSerialNumber:   cp.ChargeBoxSerialNumber,
			Iccid:                   cp.Iccid,
			Imsi:                    cp.Imsi,
			MeterType:               cp.MeterType,
			MeterSerialNumber:       cp.MeterSerialNumber,
			FirmwareVersion:         cp.FirmwareVersion,
			ChargePointIdentifier:   cp.ChargePointIdentifier,
			OcppProtocol:            cp.OcppProtocol,
		},
	}

	return res, nil
}

func (srv Service) CreateChargePointIdTag(ctx context.Context, req *rpc.CreateChargePointIdTagReq) (*rpc.CreateChargePointIdTagResp, error) {
	cp, err := srv.chargePoint.GetChargePointByIdentifier(ctx, req.ApplicationId, req.ChargePointIdentifier)
	if err != nil {
		return nil, err
	}

	it, err := srv.chargePoint.CreateIdTag(ctx, models.OcppChargePointIDTag{
		ChargePointID: cp.ID,
		IDTag:         req.IdTag,
	})

	if err != nil {
		return nil, err
	}

	res := &rpc.CreateChargePointIdTagResp{
		ChargePointIdTag: &rpc.ChargePointIdTag{
			ChargePointId:         cp.ID,
			ChargePointIdentifier: cp.ChargePointIdentifier,
			IdTag:                 it.IDTag,
		},
	}

	return res, nil
}

func (srv Service) GetChargePointIdTags(ctx context.Context, req *rpc.GetChargePointIdTagsReq) (*rpc.GetChargePointIdTagsResp, error) {
	cp, err := srv.chargePoint.GetChargePointByIdentifier(ctx, req.ApplicationId, req.ChargePointIdentifier)
	if err != nil {
		return nil, err
	}

	its, err := srv.chargePoint.GetIdTags(ctx, cp.ID)
	if err != nil {
		return nil, err
	}

	r := make([]*rpc.ChargePointIdTag, 0)
	for _, it := range its {
		r = append(r, &rpc.ChargePointIdTag{
			ChargePointId:         cp.ID,
			ChargePointIdentifier: cp.ChargePointIdentifier,
			IdTag:                 it.IDTag,
		})
	}

	res := &rpc.GetChargePointIdTagsResp{
		ChargePointIdTags: r,
	}

	return res, nil
}

func (srv Service) UpdateChargePointDetails(ctx context.Context, req *rpc.UpdateChargePointDetailsReq) (*rpc.UpdateChargePointDetailsResp, error) {
	cp, err := srv.chargePoint.Update(
		ctx,
		req.ChargePointId,
		[]string{
			"charge_point_vendor",
			"charge_point_model",
			"charge_point_serial_number",
			"charge_box_serial_number",
			"iccid",
			"imsi",
			"meter_type",
			"meter_serial_number",
			"firmware_version",
		},
		models.OcppChargePoint{
			ID:                      req.ChargePointId,
			ChargePointVendor:       req.ChargePointVendor,
			ChargePointModel:        req.ChargePointModel,
			ChargePointSerialNumber: req.ChargePointSerialNumber,
			ChargeBoxSerialNumber:   req.ChargeBoxSerialNumber,
			Iccid:                   req.Iccid,
			Imsi:                    req.Imsi,
			MeterType:               req.MeterType,
			MeterSerialNumber:       req.MeterSerialNumber,
			FirmwareVersion:         req.FirmwareVersion,
		})

	if err != nil {
		return nil, err
	}

	res := &rpc.UpdateChargePointDetailsResp{
		ChargePoint: &rpc.ChargePoint{
			Id:                      cp.ID,
			ApplicationId:           cp.ApplicationID,
			ChargePointVendor:       cp.ChargePointVendor,
			ChargePointModel:        cp.ChargePointModel,
			ChargePointSerialNumber: cp.ChargePointSerialNumber,
			ChargeBoxSerialNumber:   cp.ChargeBoxSerialNumber,
			Iccid:                   cp.Iccid,
			Imsi:                    cp.Imsi,
			MeterType:               cp.MeterType,
			MeterSerialNumber:       cp.MeterSerialNumber,
			FirmwareVersion:         cp.FirmwareVersion,
			ChargePointIdentifier:   cp.ChargePointIdentifier,
			OcppProtocol:            cp.OcppProtocol,
		},
	}

	return res, nil

}
