package usecase

import "github.com/SoulStalker/cognitask/internal/domain"

type MediaService struct {
	repo domain.MediaRepository
}

func NewMediaService(repo domain.MediaRepository) *MediaService {
	return &MediaService{
		repo: repo,
	}
}

func (s *MediaService) Create(media domain.Media) error {
	return s.repo.Create(media)
}

func (s *MediaService) Delete(media domain.Media) error {
	return s.repo.Delete(media)
}

func (s *MediaService) GetByLink(link string) (domain.Media, error) {
	return s.repo.GetByLink(link)
}

func (s *MediaService) Random() (domain.Media, error) {
	return s.repo.Random()
}
