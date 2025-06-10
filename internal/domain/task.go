package domain

import "time"

type Task struct {
	ID          uint      `json:"id"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	// Надо ли добавлять поля с датой создания датой закртия 
	// createdAt ClosedAt например?
}

type TaskRepository interface {
	Add(task Task) (Task, error) // добавляем новую задачу
	Close() (bool, error) // Закрыавет задачу вернет что закрыл или нет и ошибку
	GetAll() []Task // Список всех задач (не закрытых)
	EditDate(task Task, newDate time.Time) (Task, error) 
	Delete(id uint) error // удалить задачу
	// какие еще методы нужны в задаче?
}