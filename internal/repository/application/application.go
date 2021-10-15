package application

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/models"
)

type BaseRepo interface {
	Create(ctx context.Context, a models.OcppApplication) (models.OcppApplication, error)
	GetApplicationByID(ctx context.Context, aid int32) (models.OcppApplication, error)
	GetApplicationByUUID(ctx context.Context, auuid string) (models.OcppApplication, error)
	CreateCallback(ctx context.Context, a models.OcppApplicationCallback) (models.OcppApplicationCallback, error)
	GetApplicationCallback(ctx context.Context, aid int32, callbackEvent string) (models.OcppApplicationCallback, error)
	GetApplicationCallbacks(ctx context.Context, aid int32) ([]models.OcppApplicationCallback, error)
	DeleteCallback(ctx context.Context, aid int32) error
}

type baseRepo struct {
	db *gorm.DB
}

func NewBaseRepo(db *gorm.DB) BaseRepo {
	return &baseRepo{
		db: db,
	}
}

func (repo baseRepo) Create(ctx context.Context, a models.OcppApplication) (models.OcppApplication, error) {
	err := repo.db.Table("bb3.ocpp_application").Create(models.OcppApplication{
		UUID: uuid.NewString(),
		Name: a.Name,
	}).Error
	if err != nil {
		return models.OcppApplication{}, err
	}

	return a, nil
}

func (repo baseRepo) GetApplicationByID(ctx context.Context, aid int32) (models.OcppApplication, error) {
	a := models.OcppApplication{}
	err := repo.db.Table("bb3.ocpp_application").Where("id = ?", aid).First(&a).Error

	if err != nil {
		return models.OcppApplication{}, err
	}

	return a, nil
}

func (repo baseRepo) GetApplicationByUUID(ctx context.Context, auuid string) (models.OcppApplication, error) {
	a := models.OcppApplication{}
	err := repo.db.Table("bb3.ocpp_application").Where("uuid = ?", auuid).First(&a).Error

	if err != nil {
		return models.OcppApplication{}, err
	}

	return a, nil
}

func (repo baseRepo) CreateCallback(ctx context.Context, a models.OcppApplicationCallback) (models.OcppApplicationCallback, error) {
	err := repo.db.Table("bb3.ocpp_application_callback").Create(&a).Error
	if err != nil {
		return models.OcppApplicationCallback{}, err
	}

	return a, nil
}

func (repo baseRepo) GetApplicationCallback(ctx context.Context, aid int32, callbackEvent string) (models.OcppApplicationCallback, error) {
	a := models.OcppApplicationCallback{}

	err := repo.db.Table("bb3.ocpp_application_callback").
		Where("application_id = ?", aid).
		Where("callback_event = ?", callbackEvent).
		First(&a).
		Error

	if err != nil {
		return a, err
	}

	return a, nil
}

func (repo baseRepo) GetApplicationCallbacks(ctx context.Context, aid int32) ([]models.OcppApplicationCallback, error) {
	as := make([]models.OcppApplicationCallback, 0)

	err := repo.db.Table("bb3.ocpp_application_callback").
		Where("application_id = ?", aid).
		Find(&as).
		Error

	if err != nil {
		return nil, err
	}

	return as, nil
}

func (repo baseRepo) DeleteCallback(ctx context.Context, aid int32) error {
	a := models.OcppApplicationCallback{
		ID: aid,
	}

	err := repo.db.Table("bb3.ocpp_application_callback").Delete(&a).Error

	return err
}
