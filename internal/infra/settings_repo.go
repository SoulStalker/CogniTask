package infra

import (
	"github.com/SoulStalker/cognitask/internal/domain"
	"gorm.io/gorm"
)

type GormSettingsRepo struct {
	DB *gorm.DB
}

func NewSettingsRepository(db *gorm.DB) *GormSettingsRepo {
	return &GormSettingsRepo{
		DB: db,
	}
}

func (r *GormSettingsRepo) All() (domain.Settings, error) {
	var settings domain.Settings
	err := r.DB.Find(&settings).Error
	if err != nil {
		return domain.Settings{}, err
	}
	return settings, nil
}

func (r *GormSettingsRepo) SetDeleteDays(days uint) (domain.Settings, error) {

}

func (r *GormSettingsRepo) SetNotificationHours(hours uint) (domain.Settings, error) {

}

func (r *GormSettingsRepo) SetNotifyFrom(hours uint) (domain.Settings, error) {

}

func (r *GormSettingsRepo) SetNotifyTo(hours uint) (domain.Settings, error) {

}

func (r *GormSettingsRepo) SetRandomHour(hour uint) (domain.Settings, error) {

}
