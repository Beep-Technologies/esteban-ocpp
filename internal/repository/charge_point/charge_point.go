package chargepoint

import (
	"context"

	"gorm.io/gorm"

	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/models"
)

type BaseRepo interface {
	Create(ctx context.Context, cp models.OcppChargePoint) (models.OcppChargePoint, error)
	GetChargePoints(ctx context.Context, aid int32) ([]models.OcppChargePoint, error)
	GetChargePointByID(ctx context.Context, cpid int32) (models.OcppChargePoint, error)
	GetChargePointByIdentifier(ctx context.Context, aid int32, cpid string) (models.OcppChargePoint, error)
	Update(ctx context.Context, cpid int32, fields []string, cp models.OcppChargePoint) (models.OcppChargePoint, error)
	CreateIdTag(ctx context.Context, it models.OcppChargePointIDTag) (models.OcppChargePointIDTag, error)
	GetIdTags(ctx context.Context, cpid int32) ([]models.OcppChargePointIDTag, error)
}

type baseRepo struct {
	db *gorm.DB
}

func NewBaseRepo(db *gorm.DB) BaseRepo {
	return &baseRepo{
		db: db,
	}
}

func (repo baseRepo) Create(ctx context.Context, cp models.OcppChargePoint) (models.OcppChargePoint, error) {
	err := repo.db.Table("bb3.ocpp_charge_point").Create(&cp).Error
	if err != nil {
		return models.OcppChargePoint{}, err
	}

	return cp, nil
}

func (repo baseRepo) GetChargePoints(ctx context.Context, aid int32) ([]models.OcppChargePoint, error) {
	cps := make([]models.OcppChargePoint, 0)

	err := repo.db.Table("bb3.ocpp_charge_point").
		Where("application_id = ?", aid).
		Find(&cps).
		Error

	return cps, err
}

func (repo baseRepo) GetChargePointByID(ctx context.Context, cpid int32) (models.OcppChargePoint, error) {
	cp := models.OcppChargePoint{}
	err := repo.db.Table("bb3.ocpp_charge_point").
		Where("id = ?", cpid).
		First(&cp).
		Error

	return cp, err
}

func (repo baseRepo) GetChargePointByIdentifier(ctx context.Context, aid int32, cpid string) (models.OcppChargePoint, error) {
	cp := models.OcppChargePoint{}
	err := repo.db.Table("bb3.ocpp_charge_point").
		Where("application_id = ?", aid).
		Where("charge_point_identifier = ?", cpid).
		First(&cp).Error

	return cp, err
}

func (repo baseRepo) Update(ctx context.Context, cpid int32, fields []string, cp models.OcppChargePoint) (models.OcppChargePoint, error) {
	err := repo.db.Model(&cp).
		Select(fields).Where("id = ?", cpid).
		Updates(cp).
		Error

	if err != nil {
		return models.OcppChargePoint{}, err
	}

	cpo := models.OcppChargePoint{}
	err = repo.db.Table("bb3.ocpp_charge_point").
		Where("id = ?", cpid).
		First(&cpo).Error

	if err != nil {
		return models.OcppChargePoint{}, err
	}

	return cpo, nil
}

func (repo baseRepo) CreateIdTag(ctx context.Context, it models.OcppChargePointIDTag) (models.OcppChargePointIDTag, error) {
	err := repo.db.Table("bb3.ocpp_charge_point_id_tag").Create(&it).Error
	if err != nil {
		return models.OcppChargePointIDTag{}, err
	}

	return it, nil
}

func (repo baseRepo) GetIdTags(ctx context.Context, cpid int32) ([]models.OcppChargePointIDTag, error) {
	ids := make([]models.OcppChargePointIDTag, 0)

	err := repo.db.Table("bb3.ocpp_charge_point_id_tag").
		Where("charge_point_id = ?", cpid).
		Find(&ids).
		Error

	if err != nil {
		return nil, err
	}

	return ids, nil
}
