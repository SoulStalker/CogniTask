package infra

import (
	"time"

	"github.com/SoulStalker/cognitask/internal/domain"
	"gorm.io/gorm"
)

type GormTaskRepo struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *GormTaskRepo {
	return &GormTaskRepo{
		DB: db,
	}
}

func (r *GormTaskRepo) Add(task domain.Task) (domain.Task, error) {
	err := r.DB.Create(&task).Error
	if err != nil {
		return domain.Task{}, err
	}
	return task, nil
}

func (r *GormTaskRepo) MarkDone(id uint) (bool, error) {
	err := r.DB.Model(&domain.Task{}).Where("id=?", id).Updates(domain.Task{
		Closed:   true,
		ClosedAt: time.Now(),
	}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *GormTaskRepo) GetPending() ([]domain.Task, error) {
	var pendingTasks []domain.Task
	err := r.DB.Where("closed=?", false).Find(&pendingTasks).Error
	if err != nil {
		return []domain.Task{}, err
	}
	return pendingTasks, nil
}

func (r *GormTaskRepo) EditDate(id uint, newDate time.Time) (domain.Task, error) {
	err := r.DB.Model(&domain.Task{}).Where("id=?", id).Update("deadline", newDate).Error
	if err != nil {
		return domain.Task{}, err
	}
	var updatedTask domain.Task
	err = r.DB.First(&updatedTask, id).Error
	return updatedTask, err
}

func (r *GormTaskRepo) Delete(id uint) error {
	result := r.DB.Delete(&domain.Task{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormTaskRepo) GetExpired(deadline time.Time) ([]domain.Task, error) {
	var expiredTasks []domain.Task
	err := r.DB.Where("closed = ? AND deadline < ?", false, deadline).Find(&expiredTasks).Error
	if err != nil {
		return []domain.Task{}, err
	}
	return expiredTasks, nil
}

func (r *GormTaskRepo) DeleteOldDone(N_days int) (int64, error) {
	threshold := time.Now().AddDate(0, 0, -N_days)
	result := r.DB.Where("closed=?", true).Where("closed_at < ?", threshold).Delete(&domain.Task{})
	return result.RowsAffected, result.Error
}

func (r *GormTaskRepo) GetByID(id uint) (domain.Task, error) {
	var task domain.Task
	err := r.DB.First(&task, id).Error
	if err != nil {
		return domain.Task{}, err
	}
	return task, nil
}

func (r *GormMediaRepo) All() ([]domain.Task, error) {
	var tasks []domain.Task
	err := r.DB.Find(&tasks).Error
	if err != nil {
		return []domain.Task{}, err
	}
	return tasks, nil
}
