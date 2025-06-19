package infra

import (
	"github.com/SoulStalker/cognitask/internal/domain"
	"gorm.io/gorm"
)

type GormMediaRepo struct {
	DB *gorm.DB
}

func NewMediaRepo(db *gorm.DB) *GormMediaRepo {
	return &GormMediaRepo{
		DB: db,
	}
}

func (r *GormMediaRepo) Create(media domain.Media) error {
	err := r.DB.Create(&media).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *GormMediaRepo) GetByLink(link string) (domain.Media, error) {
	var media domain.Media
	err := r.DB.Where("link=?", link).First(&media).Error
	if err != nil {
		return domain.Media{}, err
	}
	return media, nil
}

func (r *GormMediaRepo) Delete(media domain.Media) error {
	result := r.DB.Where("link=?", media.Link).Delete(&domain.Media{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormMediaRepo) Random() (domain.Media, error) {
	var media domain.Media
	err := r.DB.Order("random()").First(&media).Error
	if err != nil {
		return domain.Media{}, nil
	}
	return media, nil
}
