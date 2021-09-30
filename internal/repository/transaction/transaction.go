package transaction

import (
	"context"

	"gorm.io/gorm"

	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/models"
)

type BaseRepo interface {
	Create(ctx context.Context, cp models.OcppTransaction) (models.OcppTransaction, error)
	GetByID(ctx context.Context, id int32) (models.OcppTransaction, error)
	GetAllByChargePointID(ctx context.Context, chargePointId int32) ([]models.OcppTransaction, error)
	Update(ctx context.Context, id int32, fields []string, transaction models.OcppTransaction) (models.OcppTransaction, error)
}

type baseRepo struct {
	db *gorm.DB
}

func NewBaseRepo(db *gorm.DB) BaseRepo {
	return &baseRepo{
		db: db,
	}
}

func (repo baseRepo) Create(ctx context.Context, t models.OcppTransaction) (models.OcppTransaction, error) {
	err := repo.db.Table("bb3.ocpp_transaction").Create(&t).Error
	if err != nil {
		return models.OcppTransaction{}, err
	}

	return t, nil
}

func (repo baseRepo) GetByID(ctx context.Context, id int32) (models.OcppTransaction, error) {
	t := models.OcppTransaction{}
	err := repo.db.Table("bb3.ocpp_transaction").Where("id = ?", id).First(&t).Error

	if err != nil {
		return models.OcppTransaction{}, err
	}

	return t, nil
}

func (repo baseRepo) GetAllByChargePointID(ctx context.Context, chargePointId int32) ([]models.OcppTransaction, error) {
	ts := []models.OcppTransaction{}

	err := repo.db.Table("bb3.ocpp_transaction").Where("charge_point_id = ?", chargePointId).Find(ts).Error

	return ts, err
}

func (repo baseRepo) Update(ctx context.Context, id int32, fields []string, t models.OcppTransaction) (models.OcppTransaction, error) {
	err := repo.db.Model(&t).Where("id = ?", id).Select(fields).Updates(t).Error

	if err != nil {
		return models.OcppTransaction{}, err
	}

	to := models.OcppTransaction{}
	err = repo.db.Where("id = ?", id).
		First(&to).Error

	if err != nil {
		return models.OcppTransaction{}, err
	}

	return to, nil
}
