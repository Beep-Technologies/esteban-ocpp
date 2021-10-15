package transaction

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/models"
)

type BaseRepo interface {
	Create(ctx context.Context, cp models.OcppTransaction) (models.OcppTransaction, error)
	GetByID(ctx context.Context, tid int32) (models.OcppTransaction, error)
	GetAllByChargePointID(ctx context.Context, cpid int32) ([]models.OcppTransaction, error)
	GetAllByChargePointIDStates(ctx context.Context, cpid int32, states []string) ([]models.OcppTransaction, error)
	GetByChargePointIDConnectorStates(ctx context.Context, cpid int32, cnid int32, states []string) (models.OcppTransaction, error)
	Update(ctx context.Context, tid int32, fields []string, transaction models.OcppTransaction) (models.OcppTransaction, error)
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

func (repo baseRepo) GetByID(ctx context.Context, tid int32) (models.OcppTransaction, error) {
	t := models.OcppTransaction{}
	err := repo.db.Table("bb3.ocpp_transaction").
		Where("id = ?", tid).
		First(&t).
		Error

	if err != nil {
		return models.OcppTransaction{}, err
	}

	return t, nil
}

func (repo baseRepo) GetAllByChargePointID(ctx context.Context, cpid int32) ([]models.OcppTransaction, error) {
	ts := []models.OcppTransaction{}

	err := repo.db.Table("bb3.ocpp_transaction").
		Where("charge_point_id = ?", cpid).
		Find(ts).
		Error

	return ts, err
}

func (repo baseRepo) GetAllByChargePointIDStates(ctx context.Context, cpid int32, states []string) ([]models.OcppTransaction, error) {
	ts := []models.OcppTransaction{}

	// TODO: order by
	err := repo.db.Table("bb3.ocpp_transaction").
		Where("charge_point_id = ?", cpid).
		Where("state IN ?", states).
		Find(&ts).
		Error

	return ts, err
}

func (repo baseRepo) GetByChargePointIDConnectorStates(ctx context.Context, cpid int32, cnid int32, states []string) (models.OcppTransaction, error) {
	t := models.OcppTransaction{}

	fmt.Printf("%d %d %+v", cpid, cnid, states)

	// TODO: order by
	err := repo.db.Table("bb3.ocpp_transaction").
		Where("charge_point_id = ?", cpid).
		Where("connector_id = ?", cnid).
		Where("state IN ?", states).
		First(&t).
		Error

	return t, err
}

func (repo baseRepo) Update(ctx context.Context, tid int32, fields []string, t models.OcppTransaction) (models.OcppTransaction, error) {
	err := repo.db.Model(&t).
		Where("id = ?", tid).
		Select(fields).
		Updates(&t).
		Error

	if err != nil {
		return models.OcppTransaction{}, err
	}

	to := models.OcppTransaction{}
	err = repo.db.Where("id = ?", tid).
		First(&to).Error

	if err != nil {
		return models.OcppTransaction{}, err
	}

	return to, nil
}
