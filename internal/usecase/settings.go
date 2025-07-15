package usecase

import "github.com/SoulStalker/cognitask/internal/domain"

type SettingsService struct {
	repo domain.SettingsRepository
}

func NewSettingsService(repo domain.SettingsRepository) *SettingsService {
	return &SettingsService{
		repo: repo,
	}
}

func (s *SettingsService) All() (domain.Settings, error) {
	return s.repo.All()
}

func (s *SettingsService) SetDeleteDays(days uint) error {
	return s.repo.SetDeleteDays(days)
}

func (s *SettingsService) SetNotificationHours(hours uint) error {
	return s.repo.SetNotificationHours(hours)
}

func (s *SettingsService) SetNotifyFrom(hours uint) error {
	return s.repo.SetNotifyFrom(hours)
}

func (s *SettingsService) SetNotifyTo(hours uint) error {
	return s.repo.SetNotifyTo(hours)
}

func (s *SettingsService) SetRandomHour(hour uint) error {
	return s.repo.SetRandomHour(hour)
}

func (s *SettingsService) GetNotificationInterval() (uint, error) {
	return s.repo.GetNotificationInterval()
}

func (s *SettingsService) GetExpirationDays() (uint, error) {
	return s.repo.GetExpirationDays()
}
