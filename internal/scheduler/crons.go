package scheduler

import (
	"fmt"
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
	notifyStartEntryID cron.EntryID
	notifyStopEntryID  cron.EntryID
	notifyTaskEntryID  cron.EntryID
	mediaEntryID       cron.EntryID
	deleteEntryID      cron.EntryID

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
	s.SendMediaJob()
	s.setupNotificationSchedule()

	s.cr.Start()
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

// setupDeleteSchedule настраивает расписание удаления старых задач
func (s *Scheduler) setupDeleteSchedule() {
	N_days, err := s.settingsUC.GetExpirationDays()
	if err != nil {
		log.Printf("Не смог получить дни из базы: %v", err)
		return
	}

	entryID, err := s.cr.AddFunc("30 4 * * *", func() {
		s.taskUC.RemoveOldTasks(int(N_days))
	})
	if err != nil {
		log.Printf("Не смог обновить базу %v", err)
	}

	s.deleteEntryID = entryID
	log.Printf("Новое расписание установлено. Задачи старше %d дней будут удаляться", N_days)
}

// updateDeleteSchedule обновляет расписание удаления
func (s *Scheduler) updateDeleteSchedule() {
	// сначала надо удалить старое расписание
	if s.deleteEntryID != 0 {
		s.cr.Remove(s.deleteEntryID)
	}

	// теперь можно создать новое
	s.setupDeleteSchedule()
	log.Printf("Расписание по удаление старых данных обновлено")
}

// setupMediaSchedule настраивает расписание отправки медиа
func (s *Scheduler) setupMediaSchedule() {
	hour, err := s.settingsUC.GetRandomHour()
	if err != nil {
		log.Printf("Не удалось получить расписание по отправке медиа из базы: %v", err)
		return
	}
	entryID, err := s.cr.AddFunc(fmt.Sprintf("* %d * * *", hour), s.SendRandomMedia)
	if err != nil {
		log.Printf("Ошибка изменения расписания медиа: %v", err)
		return
	}

	s.mediaEntryID = entryID
	log.Printf("Настроено новое расписание отправки медиа. Теперь отправка в %d часов", hour)
}

// updateMediaSchedule обновляет расписание отправки медиа
func (s *Scheduler) updateMediaSchedule() {
	if s.mediaEntryID != 0 {
		s.cr.Remove(s.mediaEntryID)
	}

	s.setupMediaSchedule()
	log.Println("Расписание по отправке медиа обновлено")
}

// setupNotificationSchedule настраивает расписание уведомлений
func (s *Scheduler) setupNotificationSchedule() {
	interval, startHour, endHour, err := s.settingsUC.GetNotificationData()
	if err != nil {
		log.Printf("Не смог получить данные по оповещениям из базы: %v", err)
		return
	}

	// Удаляем старые задачи уведомлений если есть
	s.removeNotificationTasks()

	var notifyEntryID cron.EntryID

	startJob := RepeatingNotificationJob{
		Interval:  time.Duration(interval) * time.Hour,
		Cron:      s.cr,
		EntryID:   &notifyEntryID,
		Scheduler: *s,
	}

	stopJob := StopNotificationJob{
		Cron:       s.cr,
		EntryIDPtr: &notifyEntryID,
	}

	// Задачи запуска и остановки
	startEntryID, err := s.cr.AddJob(fmt.Sprintf("0 %d * * *", startHour), startJob)
	if err != nil {
		log.Printf("Ошибка добавления задачи запуска уведомлений: %v", err)
		return
	}

	stopEntryID, err := s.cr.AddJob(fmt.Sprintf("0 %d * * *", endHour), stopJob)
	if err != nil {
		log.Printf("Ошибка добавления задачи остановки уведомлений: %v", err)
		return
	}

	s.notifyStartEntryID = startEntryID
	s.notifyStopEntryID = stopEntryID

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), int(startHour), 0, 0, 0, time.Local)
	end := time.Date(now.Year(), now.Month(), now.Day(), int(endHour), 0, 0, 0, time.Local)

	// если бот перезапустился после запуска расписания, запускаем задание
	if now.After(start) && now.Before(end) {
		startJob.Run()
	}

	log.Printf("Обновлено расписание уведомлений: Каждые %d часов с %d:00 до %d:00", interval, startHour, endHour)
}

// updateNotificationSchedule обновляет расписание уведомлений
func (s *Scheduler) updateNotificationSchedule() {
	s.removeNotificationTasks()
	s.setupNotificationSchedule()
	log.Println("Расписание уведомлений обновлено")
}

// removeNotificationTasks удаляет существующие задачи уведомлений
func (s *Scheduler) removeNotificationTasks() {
	if s.notifyStartEntryID != 0 {
		s.cr.Remove(s.notifyStartEntryID)
	}
	if s.notifyStopEntryID != 0 {
		s.cr.Remove(s.notifyStopEntryID)
	}
	if s.notifyTaskEntryID != 0 {
		s.cr.Remove(s.notifyTaskEntryID)
	}
}

// Задача отправляет медиа файлы из базы в заданное время
func (s *Scheduler) SendRandomMedia() {
	recipient := tele.ChatID(s.chatID)

	media, err := s.mediaUC.Random()
	if err != nil {
		log.Printf("Ошибка получения файла: %v", err)
		return
	}

	switch media.Type {
	case "photo":
		photo := &tele.Photo{File: tele.File{FileID: media.Link}}
		_, err = s.bot.Send(recipient, photo)
	case "video":
		video := &tele.Video{File: tele.File{FileID: media.Link}}
		_, err = s.bot.Send(recipient, video)
	}
	if err != nil {
		log.Printf("Error sending media: %v", err)
		s.bot.Send(recipient, "Ошибка при отправке медиа")
	}
}

func (s *Scheduler) SendMediaJob() {
	hour, err := s.settingsUC.GetRandomHour()
	if err != nil {
		log.Println(err)
	}
	s.cr.AddFunc(fmt.Sprintf("10 %d * * *", hour), s.SendRandomMedia)
}

func (s *Scheduler) Notifier() {
	recipient := tele.ChatID(s.chatID)

	currentTasks, err := s.taskUC.GetPending()
	if err != nil {
		log.Printf("Не смог получить задачи: %v", err)
		s.bot.Send(recipient, "Не смог получить задачи: "+err.Error())
		return
	}

	if len(currentTasks) == 0 {
		s.bot.Send(recipient, "Новых задач нет")
		return
	}

	// rows := handlers.FormatTaskList(currentTasks)
	// res, err := s.bot.Send(recipient, messages.BotMessages.YourTasks, &tele.ReplyMarkup{InlineKeyboard: rows})

	// if err != nil {
	// 	log.Printf("Не смог отправить задачи:\n%v, %v", err, res)
	// }
	// todo надо вывести FormatTaskList в другой модуль
}

// Stop останавливает планировщик
func (s *Scheduler) Stop() {
	s.cr.Stop()
	close(s.updateChan)
}
