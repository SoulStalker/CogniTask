package domain

type Settings struct {
	DeleteAfterDays   uint // Удалять закрытые задачи после N дней
	NotificationHours uint // Интервал уведомлений в часах
	NotifyFrom        uint // Уведомления от
	NotifyTo          uint // Уведомления до
	RandomHour        uint // во сколько присылать мотивацию
}

type SettingsRepository interface {
	All() (Settings, error)
	SetDeleteDays(days uint) (Settings, error)
	SetNotificationHours(hours uint) (Settings, error)
	SetNotifyFrom(hours uint) (Settings, error)
	SetNotifyTo(hours uint) (Settings, error)
	SetRandomHour(hour uint) (Settings, error)
}
