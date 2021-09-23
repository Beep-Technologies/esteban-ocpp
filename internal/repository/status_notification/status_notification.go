package statusnotification

import (
	"context"

	"gorm.io/gorm"

	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/models"
)

type BaseRepo interface {
	Create(ctx context.Context, cp models.StatusNotification) (models.StatusNotification, error)
}

type baseRepo struct {
	db *gorm.DB
}

func NewBaseRepo(db *gorm.DB) BaseRepo {
	return &baseRepo{
		db: db,
	}
}

func (repo baseRepo) Create(ctx context.Context, sn models.StatusNotification) (models.StatusNotification, error) {
	err := repo.db.Create(&sn).Error
	if err != nil {
		return models.StatusNotification{}, err
	}

	return sn, nil
}
