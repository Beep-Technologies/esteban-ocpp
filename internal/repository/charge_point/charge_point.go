package chargepoint

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

func (repo BaseRepo) Create(ctx context.Context, cp models.OcppChargePoint) (models.OcppChargePoint, error) {
	err := repo.db.Table("bb3.ocpp_charge_point").Create(&cp).Error
	if err != nil {
		return models.OcppChargePoint{}, err
	}

	return cp, nil
}

func (repo BaseRepo) GetChargePoints(ctx context.Context, entity string) ([]models.OcppChargePoint, error) {
	cps := make([]models.OcppChargePoint, 0)

	err := repo.db.Table("bb3.ocpp_charge_point").
		Where("entity_code = ?", entity).
		Find(&cps).
		Error

	return cps, err
}

func (repo BaseRepo) GetChargePoint(ctx context.Context, entity, identifier string) (models.OcppChargePoint, error) {
	cp := models.OcppChargePoint{}

	err := repo.db.Table("bb3.ocpp_charge_point").
		Where("entity_code = ?", entity).
		Where("charge_point_identifier = ?", identifier).
		First(&cp).
		Error

	return cp, err
}

func (repo BaseRepo) Update(ctx context.Context, cpid int32, fields []string, cp models.OcppChargePoint) (models.OcppChargePoint, error) {
	err := repo.db.Table("bb3.ocpp_charge_point").
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

func (repo BaseRepo) CreateIdTag(ctx context.Context, it models.OcppChargePointIDTag) (models.OcppChargePointIDTag, error) {
	err := repo.db.Table("bb3.ocpp_charge_point_id_tag").Create(&it).Error
	if err != nil {
		return models.OcppChargePointIDTag{}, err
	}

	return it, nil
}

func (repo BaseRepo) GetIdTags(ctx context.Context, cpid int32) ([]models.OcppChargePointIDTag, error) {
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

func (repo BaseRepo) GetIdTag(ctx context.Context, cpid int, idTag string) (models.OcppChargePointIDTag, error) {
	idt := models.OcppChargePointIDTag{}

	err := repo.db.Table("bb3.ocpp_charge_point_id_tag").
		Where("charge_point_id = ?", cpid).
		Where("id_tag = ?", idTag).
		First(&idt).
		Error

	if err != nil {
		return models.OcppChargePointIDTag{}, err
	}

	return idt, nil
}
