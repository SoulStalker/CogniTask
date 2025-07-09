package fsm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Состояния FSM
const (
	StateIdle = ""

	// состояния при вводе задачи
	StateWaitingTaskText     = "waiting_task_text"
	StateWaitingTaskDate     = "waiting_task_date"
	StateWaitingTaskCategory = "waiting_task_category"

	// состояния при вводе настроек
	StateDeleteAfterDays   = "set_delete_days"
	StateNotificationHours = "set_notify_period"
	StateNotifyFrom        = "set_notify_from"
	StateNotifyTo          = "set_notify_to"
	StateRandom            = "set_random_hour"
)

// FSMData Структура данных FSM
type FSMData struct {
	State           string `json:"state"`
	TaskText        string `json:"task_text,omitempty"`
	TaskDate        string `json:"task_date,omitempty"`
	TaskCategory    string `json:"task_category,omitempty"`
}

// FSMService FSM сервис
type FSMService struct {
	redis   *redis.Client
	timeout time.Duration
}

// NewFSMService Конструктор FSM сервиса
func NewFSMService(redisClient *redis.Client, timeout time.Duration) *FSMService {
	return &FSMService{
		redis:   redisClient,
		timeout: timeout,
	}
}

// GetState Получить состояние пользователя
func (f *FSMService) GetState(ctx context.Context, userID int64) (*FSMData, error) {
	key := f.getUserKey(userID)

	result, err := f.redis.Get(ctx, key).Result()
	if errors.Is(redis.Nil, err) {
		// Пользователь не в состоянии FSM
		return &FSMData{State: StateIdle}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get state: %w", err)
	}

	var data FSMData
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal state: %w", err)
	}

	return &data, nil
}

// SetState Установить состояние пользователя
func (f *FSMService) SetState(ctx context.Context, userID int64, data *FSMData) error {
	key := f.getUserKey(userID)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := f.redis.Set(ctx, key, jsonData, f.timeout).Err(); err != nil {
		return fmt.Errorf("failed to set state: %w", err)
	}

	return nil
}

// ClearState Очистить состояние пользователя
func (f *FSMService) ClearState(ctx context.Context, userID int64) error {
	key := f.getUserKey(userID)

	if err := f.redis.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to clear state: %w", err)
	}

	return nil
}

// IsInState Проверить, находится ли пользователь в состоянии FSM
func (f *FSMService) IsInState(ctx context.Context, userID int64) (bool, error) {
	data, err := f.GetState(ctx, userID)
	if err != nil {
		return false, err
	}

	return data.State != StateIdle, nil
}

// Получить ключ для пользователя в Redis
func (f *FSMService) getUserKey(userID int64) string {
	return fmt.Sprintf("user:%d:fsm", userID)
}
