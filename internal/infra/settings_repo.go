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

func (r *GormSettingsRepo) SetDeleteDays(days uint) error {
	err := r.DB.Model(&domain.Settings{}).Where("id=?", 1).Updates(domain.Settings{DeleteAfterDays: days}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *GormSettingsRepo) SetNotificationHours(hours uint) error {
	err := r.DB.Model(&domain.Settings{}).Where("id=?", 1).Updates(domain.Settings{NotificationHours: hours}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *GormSettingsRepo) SetNotifyFrom(hours uint) error {
	err := r.DB.Model(&domain.Settings{}).Where("id=?", 1).Updates(domain.Settings{NotifyFrom: hours}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *GormSettingsRepo) SetNotifyTo(hours uint) error {
	err := r.DB.Model(&domain.Settings{}).Where("id=?", 1).Updates(domain.Settings{NotifyTo: hours}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *GormSettingsRepo) SetRandomHour(hour uint) error {
	err := r.DB.Model(&domain.Settings{}).Where("id=?", 1).Updates(domain.Settings{RandomHour: hour}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *GormSettingsRepo) GetNotificationInterval() (uint, error) {
	var setting domain.Settings
	err := r.DB.First(&setting).Error
	if err != nil {
		return 0, err
	}
	return setting.NotificationHours, nil
}

func (r *GormSettingsRepo) GetExpirationDays() (uint, error) {
	var setting domain.Settings
	err := r.DB.First(&setting).Error
	if err != nil {
		return 0, err
	}
	return setting.DeleteAfterDays, nil
}

func (r *GormSettingsRepo) GetRandomHour() (uint, error) {
	var setting domain.Settings
	err := r.DB.First(&setting).Error
	if err != nil {
		return 0, err
	}
	return setting.RandomHour, nil
}

func (r *GormSettingsRepo) GetNotificationData() (uint, uint, uint, error) {
	var setting domain.Settings
	err := r.DB.First(&setting).Error
	if err != nil {
		return 0, 0, 0, err
	}
	return setting.NotificationHours, setting.NotifyFrom, setting.NotifyTo, nil
}
