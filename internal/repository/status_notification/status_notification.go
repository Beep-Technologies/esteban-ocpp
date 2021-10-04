package statusnotification

import (
	"context"

	"gorm.io/gorm"

	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/models"
)

type BaseRepo interface {
	Create(ctx context.Context, cp models.OcppStatusNotification) (models.OcppStatusNotification, error)
	GetLatest(ctx context.Context, cpid int) ([]models.OcppStatusNotification, error)
}

type baseRepo struct {
	db *gorm.DB
}

func NewBaseRepo(db *gorm.DB) BaseRepo {
	return &baseRepo{
		db: db,
	}
}

func (repo baseRepo) Create(ctx context.Context, sn models.OcppStatusNotification) (models.OcppStatusNotification, error) {
	err := repo.db.Table("bb3.ocpp_status_notification").Create(&sn).Error
	if err != nil {
		return models.OcppStatusNotification{}, err
	}

	return sn, nil
}

func (repo baseRepo) GetLatest(ctx context.Context, cpid int) ([]models.OcppStatusNotification, error) {
	res := make([]models.OcppStatusNotification, 0)

	// this entire query might cause issues if two status notifications somehow have
	// the exact same charge_point_id, connector_id and timestamp
	uniqueSubquery := repo.db.Table("bb3.ocpp_status_notification").
		Select("charge_point_id, connector_id, MAX(timestamp) as timestamp").
		Where("charge_point_id = ?", cpid).
		Group("charge_point_id,connector_id").
		Order("connector_id asc")

	err := repo.db.
		Table("(?) as sn_a", uniqueSubquery).
		Joins("natural inner join bb3.ocpp_status_notification as sn_b").
		Select("*").
		Distinct().
		Order("connector_id asc").
		Find(&res).
		Error

	if err != nil {
		return nil, err
	}

	return res, nil
}
