package application

import (
	"context"
	"crypto/sha256"
	"encoding/base64"

	"github.com/google/uuid"
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

func (repo BaseRepo) Create(ctx context.Context, a models.OcppApplication) (models.OcppApplication, error) {
	// an application ID should be 8-character alphanumeric string (first section of a uuid)
	a.ID = uuid.NewString()[:8]

	err := repo.db.Table("bb3.ocpp_application").Create(&a).Error
	if err != nil {
		return models.OcppApplication{}, err
	}

	return a, nil
}

func (repo BaseRepo) GetApplicationByID(ctx context.Context, aid string) (models.OcppApplication, error) {
	a := models.OcppApplication{}
	err := repo.db.Table("bb3.ocpp_application").Where("id = ?", aid).First(&a).Error

	return a, err
}

func (repo BaseRepo) GetApiKeyDetails(ctx context.Context, apiKey string) (models.OcppApplicationAPIKey, error) {
	k := models.OcppApplicationAPIKey{}

	h := sha256.Sum256([]byte(apiKey))
	s := base64.URLEncoding.EncodeToString(h[:])

	err := repo.db.Table("bb3.ocpp_application_api_key").
		Where("api_key_hash = ?", s).
		First(&k).
		Error

	return k, err
}

func (repo BaseRepo) CreateCallback(ctx context.Context, a models.OcppApplicationCallback) (models.OcppApplicationCallback, error) {
	err := repo.db.Table("bb3.ocpp_application_callback").Create(&a).Error
	if err != nil {
		return models.OcppApplicationCallback{}, err
	}

	return a, nil
}

func (repo BaseRepo) GetApplicationCallback(ctx context.Context, aid string, callbackEvent string) (models.OcppApplicationCallback, error) {
	a := models.OcppApplicationCallback{}

	err := repo.db.Table("bb3.ocpp_application_callback").
		Where("application_id = ?", aid).
		Where("callback_event = ?", callbackEvent).
		First(&a).
		Error

	return a, err
}

func (repo BaseRepo) GetApplicationCallbacks(ctx context.Context, aid string) ([]models.OcppApplicationCallback, error) {
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
