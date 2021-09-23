package transaction

import (
	"context"

	"gorm.io/gorm"

	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/models"
)

type BaseRepo interface {
	Create(ctx context.Context, cp models.Transaction) (models.Transaction, error)
	GetByID(ctx context.Context, id int32) (models.Transaction, error)
	Update(ctx context.Context, id int32, fields []string, transaction models.Transaction) (models.Transaction, error)
}

type baseRepo struct {
	db *gorm.DB
}

func NewBaseRepo(db *gorm.DB) BaseRepo {
	return &baseRepo{
		db: db,
	}
}

func (repo baseRepo) Create(ctx context.Context, t models.Transaction) (models.Transaction, error) {
	err := repo.db.Create(&t).Error
	if err != nil {
		return models.Transaction{}, err
	}

	return t, nil
}

func (repo baseRepo) GetByID(ctx context.Context, id int32) (models.Transaction, error) {
	t := models.Transaction{}
	err := repo.db.Where("id = ?", id).First(&t).Error

	if err != nil {
		return models.Transaction{}, err
	}

	return t, nil
}

func (repo baseRepo) Update(ctx context.Context, id int32, fields []string, t models.Transaction) (models.Transaction, error) {
	err := repo.db.Model(&t).Select(fields).Where("id = ?", id).Updates(t).Error

	if err != nil {
		return models.Transaction{}, err
	}

	to := models.Transaction{}
	err = repo.db.Where("id = ?", id).
		First(&to).Error

	if err != nil {
		return models.Transaction{}, err
	}

	return to, nil
}
