package domain

type Settings struct {
	ID                uint `gorm:"primaryKey"`
	DeleteAfterDays   uint // Удалять закрытые задачи после N дней
	NotificationHours uint // Интервал уведомлений в часах
	NotifyFrom        uint // Уведомления от
	NotifyTo          uint // Уведомления до
	RandomHour        uint // во сколько присылать мотивацию
}

type SettingsRepository interface {
	All() (Settings, error)
	SetDeleteDays(days uint) error
	SetNotificationHours(hours uint) error
	SetNotifyFrom(hours uint) error
	SetNotifyTo(hours uint) error
	SetRandomHour(hour uint) error
	Interval() (uint, error)
	DeleteOldDataDays(N_days uint) (uint, error)
}
