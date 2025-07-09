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