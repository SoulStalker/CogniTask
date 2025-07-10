package usecase

import (
	"errors"

	"github.com/SoulStalker/cognitask/internal/domain"
	"github.com/SoulStalker/cognitask/internal/messages"
	"gorm.io/gorm"
)

type MediaService struct {
	repo domain.MediaRepository
}

func NewMediaService(repo domain.MediaRepository) *MediaService {
	return &MediaService{
		repo: repo,
	}
}

func (s *MediaService) Create(media domain.Media) error {
	existed, err := s.GetByLink(media.Link)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.repo.Create(media)
		}
	}
	if existed.Link != "" {
		return errors.New(messages.BotMessages.FileExisted)
	}
	return err
}

func (s *MediaService) Delete(media domain.Media) error {
	return s.repo.Delete(media)
}

func (s *MediaService) GetByLink(link string) (domain.Media, error) {
	return s.repo.GetByLink(link)
}

func (s *MediaService) Random() (domain.Media, error) {
	media, err := s.repo.Random()

	// если все отправлено, сбрасываем статусы и начинаем заново
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = s.repo.ClearStatus()
		if err != nil {
			return domain.Media{}, err
		}
		return s.repo.Random()
	}
	return media, err
}
