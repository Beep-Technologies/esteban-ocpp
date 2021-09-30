package statusnotification

import (
	"context"

	"gorm.io/gorm"

	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/models"
)

type BaseRepo interface {
	Create(ctx context.Context, cp models.OcppStatusNotification) (models.OcppStatusNotification, error)
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
