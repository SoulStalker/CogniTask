package usecase

import "github.com/SoulStalker/cognitask/internal/domain"

type SettingsService struct {
	repo domain.TaskRepository
}

func NewSettingsService(repo domain.TaskRepository) *SettingsService {
	return &SettingsService{
		repo: repo,
	}
}