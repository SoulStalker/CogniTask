package scheduler

import (
	"log"
	"time"

	"github.com/SoulStalker/cognitask/internal/usecase"
	"github.com/robfig/cron/v3"
	tele "gopkg.in/telebot.v3"
)

type ScheduleUpdateType int

const (
	UpdateNotifications ScheduleUpdateType = iota
	UpdateMediaSchedule
	UpdateDeleteSchedule
)

type ScheduleUpdate struct {
	Type ScheduleUpdateType
	Data interface{} // можем передать дополнительные данные если нужно
}

type Scheduler struct {
	cr         *cron.Cron
	settingsUC *usecase.SettingsService
	taskUC     *usecase.TaskService
	mediaUC    *usecase.MediaService
	bot        *tele.Bot
	chatID     int64

	// Для отслеживания активных задач
	notifyEntryID cron.EntryID
	mediaEntryID  cron.EntryID
	deleteEntryID cron.EntryID

	// Канал для обновления расписания
	updateChan chan ScheduleUpdate
}

func NewScheduler(
	cr *cron.Cron,
	settingsUC *usecase.SettingsService,
	taskUC *usecase.TaskService,
	mediaUC *usecase.MediaService,
	bot *tele.Bot,
	chatID int64,
) *Scheduler {
	return &Scheduler{
		cr:         cr,
		settingsUC: settingsUC,
		taskUC:     taskUC,
		mediaUC:    mediaUC,
		bot:        bot,
		chatID:     chatID,
		updateChan: make(chan ScheduleUpdate, 10),
	}
}

// GetUpdateChannel возвращает канал для отправки обновлений
func (s *Scheduler) GetUpdateChannel() chan<- ScheduleUpdate {
	return s.updateChan
}

// InitDefaultSchedule запускает расписания по настройками из базы
func (s *Scheduler) InitDefaultSchedule() {
	// горутина будет обрабатывать обновления
	go s.handleUpdates()

	// инициализация всех расписаний
	s.setupDeleteSchedule()
	s.setupMediaSchedule()
	s.setupNotificationSchedule()

	s.cr.Start()

	s.PrintAllTasks()
}

// handleUpdate обрабатывает обновления расписаний
func (s *Scheduler) handleUpdates() {
	for update := range s.updateChan {
		switch update.Type {
		case UpdateNotifications:
			s.updateNotificationSchedule()
		case UpdateMediaSchedule:
			s.updateMediaSchedule()
		case UpdateDeleteSchedule:
			s.updateDeleteSchedule()
		}
	}
}

// SetBot устанавливает бота в планировщик
func (s *Scheduler) SetBot(bot *tele.Bot) {
	s.bot = bot
}

// setupNotificationSchedule настраивает расписание уведомлений
func (s *Scheduler) setupNotificationSchedule() {
	interval, startHour, endHour, err := s.settingsUC.GetNotificationData()
	if err != nil {
		log.Printf("Не смог получить данные по оповещениям из базы: %v\n", err)
		return
	}

	// Удаляем старые задачи уведомлений если есть
	s.cr.Remove(s.notifyEntryID)

	notifyJob := RepeatingNotificationJob{
		Interval:  time.Duration(interval) * time.Hour, 
		Cron:      s.cr,
		EntryID:   &s.notifyEntryID,
		Scheduler: *s,
	}

	if s.notificationTime() {
		notifyJob.Run()
	}

	s.notifyEntryID = *notifyJob.EntryID

	log.Printf("Обновлено расписание уведомлений: Каждые %d часов с %d:00 до %d:00\n", interval, startHour, endHour)
	s.PrintAllTasks()
}

// updateNotificationSchedule обновляет расписание уведомлений
func (s *Scheduler) updateNotificationSchedule() {
	log.Printf("Задача %d удалена", s.notifyEntryID)
	s.cr.Remove(s.notifyEntryID)

	s.setupNotificationSchedule()
	log.Println("Расписание уведомлений обновлено")
}

// notificationTime возвращает можно ли в этом промежутке времени слать уведомления
func (s *Scheduler) notificationTime() bool {
	_, startHour, endHour, err := s.settingsUC.GetNotificationData()

	if err != nil {
		return false
	}

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), int(startHour), 0, 0, 0, time.Local)
	end := time.Date(now.Year(), now.Month(), now.Day(), int(endHour), 0, 0, 0, time.Local)

	if now.After(start) && now.Before(end) {
		return true
	}
	return false
}

// Stop останавливает планировщик
func (s *Scheduler) Stop() {
	s.cr.Stop()
	close(s.updateChan)
}


// PrintAllTasks cлужебный метод для вывода задач в консоль
func (s *Scheduler) PrintAllTasks() {
	for _, entry := range s.cr.Entries() {
		log.Printf("Задача %d: следующее выполнение %v", entry.ID, entry.Next)
	}
}
