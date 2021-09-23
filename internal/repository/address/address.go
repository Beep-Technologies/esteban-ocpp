package address

import (
	"context"

	"gorm.io/gorm"

	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/models"
)

type BaseRepo interface {
	Create(ctx context.Context, addr models.Address) (models.Address, error)
	GetByID(ctx context.Context, id int32) (models.Address, error)
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

func (repo baseRepo) Create(ctx context.Context, addr models.Address) (models.Address, error) {
	err := repo.db.Create(&addr).Error
	if err != nil {
		return models.Address{}, err
	}

	return addr, nil
}

func (repo baseRepo) GetByID(ctx context.Context, addrID int32) (models.Address, error) {
	addr := models.Address{}
	err := repo.db.Where("id = ?", addrID).First(&addr).Error

	if err != nil {
		return models.Address{}, err
	}

	return addr, nil
}
