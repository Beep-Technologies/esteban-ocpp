package chargepoint

import (
	"context"

	"gorm.io/gorm"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/repository/address"
	chargepoint "github.com/Beep-Technologies/beepbeep3-ocpp/internal/repository/charge_point"
)

type Service struct {
	db          *gorm.DB
	chargePoint chargepoint.BaseRepo
	address     address.BaseRepo
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db:          db,
		chargePoint: chargepoint.NewBaseRepo(db),
		address:     address.NewBaseRepo(db),
	}
}

func (srv Service) GetChargePointByID(ctx context.Context, req *rpc.GetChargePointByIDReq) (*rpc.GetChargePointByIDResp, error) {
	cp, err := srv.chargePoint.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	res := &rpc.GetChargePointByIDResp{
		Id:                      cp.ID,
		ChargePointVendor:       cp.ChargePointVendor,
		ChargePointModel:        cp.ChargePointModel,
		ChargePointSerialNumber: cp.ChargePointSerialNumber,
		ChargeBoxSerialNumber:   cp.ChargeBoxSerialNumber,
		Iccid:                   cp.Iccid,
		Imsi:                    cp.Imsi,
		MeterType:               cp.MeterType,
		MeterSerialNumber:       cp.MeterSerialNumber,
		FirmwareVersion:         cp.FirmwareVersion,
		OcppProtocol:            cp.OcppProtocol,
		ChargePointIdentifier:   cp.ChargePointIdentifier,
		Description:             cp.Description,
		LocationLatitude:        cp.LocationLatitude,
		LocationLongitude:       cp.LocationLongitude,
		// TODO: add address here
	}

	return res, nil
}

func (srv Service) GetChargePointByIdentifier(ctx context.Context, req *rpc.GetChargePointByIdentifierReq) (*rpc.GetChargePointByIdentifierResp, error) {
	cp, err := srv.chargePoint.GetByChargePointIdentifier(ctx, req.ChargePointIdentifier)
	if err != nil {
		return nil, err
	}

	res := &rpc.GetChargePointByIdentifierResp{
		Id:                      cp.ID,
		ChargePointVendor:       cp.ChargePointVendor,
		ChargePointModel:        cp.ChargePointModel,
		ChargePointSerialNumber: cp.ChargePointSerialNumber,
		ChargeBoxSerialNumber:   cp.ChargeBoxSerialNumber,
		Iccid:                   cp.Iccid,
		Imsi:                    cp.Imsi,
		MeterType:               cp.MeterType,
		MeterSerialNumber:       cp.MeterSerialNumber,
		FirmwareVersion:         cp.FirmwareVersion,
		OcppProtocol:            cp.OcppProtocol,
		ChargePointIdentifier:   cp.ChargePointIdentifier,
		Description:             cp.Description,
		LocationLatitude:        cp.LocationLatitude,
		LocationLongitude:       cp.LocationLongitude,
		// TODO: add address here
	}

	return res, nil
}
