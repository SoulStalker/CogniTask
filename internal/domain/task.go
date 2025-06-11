package domain

import "time"

type Task struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Description string    `gorm:"size:200; not null" json:"description"`
	Deadline    time.Time `json:"deadline"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	Closed      bool      `json:"closed"`
	ClosedAt    time.Time `json:"closed_at"`
}

type TaskRepository interface {
	Add(task Task) (Task, error)    // добавляем новую задачу
	MarkDone(id uint) (bool, error) // Закрыавет задачу вернет что закрыл или нет и ошибку
	GetPending() ([]Task, error)    // Список всех задач (не закрытых)
	EditDate(task Task, newDate time.Time) (Task, error)
	Delete(id uint) error // удалить задачу
	GetExpired(deadline time.Time) ([]Task, error)
	DeleteOldDone(N_days int) error
	GetByID(id uint) (Task, error)
}
