package chargepoint

import (
	"context"

	"gorm.io/gorm"

	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/models"
)

type BaseRepo interface {
	Create(ctx context.Context, cp models.ChargePoint) (models.ChargePoint, error)
	GetByID(ctx context.Context, id int32) (models.ChargePoint, error)
	GetByChargePointIdentifier(ctx context.Context, cpId string) (models.ChargePoint, error)
	// TODO: Implement Update and Delete
	// Update(ctx context.Context, fields []string) (addr models.Address, err error)
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

func (repo baseRepo) Create(ctx context.Context, cp models.ChargePoint) (models.ChargePoint, error) {
	err := repo.db.Create(&cp).Error
	if err != nil {
		return models.ChargePoint{}, err
	}

	return cp, nil
}

func (repo baseRepo) GetByID(ctx context.Context, id int32) (models.ChargePoint, error) {
	cp := models.ChargePoint{}
	err := repo.db.Where("id = ?", id).First(&cp).Error

	if err != nil {
		return models.ChargePoint{}, err
	}

	return cp, nil
}

func (repo baseRepo) GetByChargePointIdentifier(ctx context.Context, cpId string) (models.ChargePoint, error) {
	cp := models.ChargePoint{}
	err := repo.db.Where("charge_point_identifier = ?", cpId).First(&cp).Error

	if err != nil {
		return models.ChargePoint{}, err
	}

	return cp, nil
}
