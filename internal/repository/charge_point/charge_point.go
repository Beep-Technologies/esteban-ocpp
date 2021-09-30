package chargepoint

import (
	"context"

	"gorm.io/gorm"

	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/models"
)

type BaseRepo interface {
	Create(ctx context.Context, cp models.OcppChargePoint) (models.OcppChargePoint, error)
	GetByID(ctx context.Context, id int32) (models.OcppChargePoint, error)
	GetByChargePointIdentifier(ctx context.Context, cpId string) (models.OcppChargePoint, error)
	Update(ctx context.Context, id int32, fields []string, chargePoint models.OcppChargePoint) (models.OcppChargePoint, error)
	// TODO: Implement Delete
	// Delete(ctx context.Context, id int32) (err error)
}

type baseRepo struct {
	db *gorm.DB
}

func NewBaseRepo(db *gorm.DB) BaseRepo {
	return &baseRepo{
		db: db,
	}
}

func (repo baseRepo) Create(ctx context.Context, cp models.OcppChargePoint) (models.OcppChargePoint, error) {
	err := repo.db.Table("bb3.ocpp_charge_point").Create(&cp).Error
	if err != nil {
		return models.OcppChargePoint{}, err
	}

	return cp, nil
}

func (repo baseRepo) GetByID(ctx context.Context, id int32) (models.OcppChargePoint, error) {
	cp := models.OcppChargePoint{}
	err := repo.db.Table("bb3.ocpp_charge_point").Where("id = ?", id).First(&cp).Error

	if err != nil {
		return models.OcppChargePoint{}, err
	}

	return cp, nil
}

func (repo baseRepo) GetByChargePointIdentifier(ctx context.Context, cpId string) (models.OcppChargePoint, error) {
	cp := models.OcppChargePoint{}
	err := repo.db.Table("bb3.ocpp_charge_point").Where("charge_point_identifier = ?", cpId).First(&cp).Error

	if err != nil {
		return models.OcppChargePoint{}, err
	}

	return cp, nil
}

func (repo baseRepo) Update(ctx context.Context, id int32, fields []string, cp models.OcppChargePoint) (models.OcppChargePoint, error) {
	err := repo.db.Model(&cp).Select(fields).Where("id = ?", id).Updates(cp).Error

	if err != nil {
		return models.OcppChargePoint{}, err
	}

	cpo := models.OcppChargePoint{}
	err = repo.db.Table("bb3.ocpp_charge_point").
		Where("id = ?", id).
		First(&cpo).Error

	if err != nil {
		return models.OcppChargePoint{}, err
	}

	return cpo, nil
}
