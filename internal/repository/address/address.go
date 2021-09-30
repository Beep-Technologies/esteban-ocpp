package address

import (
	"context"

	"gorm.io/gorm"

	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/models"
)

type BaseRepo interface {
	Create(ctx context.Context, addr models.OcppAddress) (models.OcppAddress, error)
	GetByID(ctx context.Context, id int32) (models.OcppAddress, error)
	// TODO: Implement Update and Delete
	// Update(ctx context.Context, fields []string) (addr models.OcppAddress, err error)
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

func (repo baseRepo) Create(ctx context.Context, addr models.OcppAddress) (models.OcppAddress, error) {
	err := repo.db.Table("bb3.ocpp_address").Create(&addr).Error
	if err != nil {
		return models.OcppAddress{}, err
	}

	return addr, nil
}

func (repo baseRepo) GetByID(ctx context.Context, addrID int32) (models.OcppAddress, error) {
	addr := models.OcppAddress{}
	err := repo.db.Table("bb3.ocpp_address").Where("id = ?", addrID).First(&addr).Error

	if err != nil {
		return models.OcppAddress{}, err
	}

	return addr, nil
}
