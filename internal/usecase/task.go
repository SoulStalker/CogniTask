package usecase

import (
	"time"

	"github.com/SoulStalker/cognitask/internal/domain"
)

type TaskService struct {
	repo domain.TaskRepository
}

func NewTaskService(repo domain.TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) Add(task domain.Task) (domain.Task, error) {
	return s.repo.Add(task)
}

func (s *TaskService) MarkDone(id uint) (bool, error) {
	return s.repo.MarkDone(id)
}

func (s *TaskService) GetPending() ([]domain.Task, error) {
	return s.repo.GetPending()
}

func (s *TaskService) All() ([]domain.Task, error) {
	return s.repo.All()
}

func (s *TaskService) EditDate(id uint, newDate time.Time) (domain.Task, error) {
	return s.repo.EditDate(id, newDate)
}

func (s *TaskService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *TaskService) GetExpired(deadline time.Time) ([]domain.Task, error) {
	return s.repo.GetExpired(deadline)
}

func (s *TaskService) DeleteOldDone(N_days int) (int64, error) {
	return s.repo.DeleteOldDone(N_days)
}

func (s *TaskService) GetByID(id uint) (domain.Task, error) {
	return s.repo.GetByID(id)
}
