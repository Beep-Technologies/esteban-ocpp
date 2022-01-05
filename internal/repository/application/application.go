package application

import (
	"context"

	"gorm.io/gorm"

	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/models"
)

type BaseRepo struct {
	db *gorm.DB
}

func NewBaseRepo(db *gorm.DB) BaseRepo {
	return BaseRepo{
		db: db,
	}
}

func (repo BaseRepo) GetApplicationByID(ctx context.Context, aid string) (models.OcppApplication, error) {
	a := models.OcppApplication{}
	err := repo.db.Table("bb3.ocpp_application").Where("id = ?", aid).First(&a).Error

	return a, err
}

func (repo BaseRepo) CreateCallback(ctx context.Context, a models.OcppApplicationCallback) (models.OcppApplicationCallback, error) {
	err := repo.db.Table("bb3.ocpp_application_callback").Create(&a).Error
	if err != nil {
		return models.OcppApplicationCallback{}, err
	}

	return a, nil
}

func (repo BaseRepo) GetApplicationCallbacks(ctx context.Context, entity, callbackEvent string) ([]models.OcppApplicationCallback, error) {
	callbacks := make([]models.OcppApplicationCallback, 0)

	err := repo.db.Table("bb3.ocpp_application_callback").
		Joins("join bb3.ocpp_application_entity on bb3.ocpp_application_callback.application_id = bb3.ocpp_application_entity.application_id").
		Where("entity_code = ?", entity).
		Where("callback_event = ?", callbackEvent).
		Find(&callbacks).
		Error

	if err != nil {
		return nil, err
	}

	return callbacks, nil
}

func (repo BaseRepo) UpdateCallback(ctx context.Context, aid string, fields []string, a models.OcppApplicationCallback) (models.OcppApplicationCallback, error) {
	err := repo.db.Table("bb3.ocpp_application_callback").
		Select(fields).Where("id = ?", aid).
		Updates(a).
		Error

	if err != nil {
		return models.OcppApplicationCallback{}, err
	}

	ao := models.OcppApplicationCallback{}
	err = repo.db.Model(&ao).Where("id = ?", aid).First(&ao).Error

	if err != nil {
		return models.OcppApplicationCallback{}, err
	}

	return ao, nil
}

func (repo BaseRepo) DeleteCallback(ctx context.Context, acid int32) error {
	ac := models.OcppApplicationCallback{
		ID: acid,
	}

	err := repo.db.Table("bb3.ocpp_application_callback").Delete(&ac).Error

	return err
}
