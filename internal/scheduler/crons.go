package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

// надо добавить расписания для автоудаления сообщений (проверять раз в день и удалять все что старше заданных дней)
// расписание для моитиваций
// расписание для запуска уведомлений
// расписание для остановки уведомлений
// либо перенсти тикер в крон

type Scheduler struct {
	cr *cron.Cron
}

func NewScheduler(cr *cron.Cron) *Scheduler {
	return &Scheduler{
		cr: cr,
	}
}

// InitDefaultSchedule запускает расписания по настройками из базы
func (s *Scheduler) InitDefaultSchedule() {
	// авто-удаление старых задач
	s.deleteOldTasks()
	// расписаний мотиваций
	s.SendMedia(17)
	// можно начинать слать задачи с интервалом
	// s.cr.AddFunc("26 16 * * *", func() { log.Println("Начинаю слать уведомления", time.Now()) })
	// Останавливаем уведомления
	s.cr.AddFunc("27 16 * * *", func() { log.Println("Перестаю слать уведомления", time.Now()) })
	// сами уведомления каждые 3 часа
	// s.cr.AddFunc("@every 3m", func() { log.Println("Уведомления каждые три часа", time.Now()) })
	// notifierID := s.Notifier(3)
	var notifyEntryID cron.EntryID

	startJob := RepeatingNotificationJob{
		Interval: 3 * time.Second,
		Cron:     s.cr,
		EntryID:  &notifyEntryID,
	}

	stopJob := StopNotificationJob{
		Cron:       s.cr,
		EntryIDPtr: &notifyEntryID,
	}

	s.cr.AddJob("00 9 * * *", startJob)
	s.cr.AddJob("18 16 * * *", stopJob)

	now := time.Now()

	start := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, time.Local)
	end := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, time.Local)

	// если бот перезапустился после запуска расписания, запускаем задание
	if now.After(start) && now.Before(end) {
		startJob.Run()
	}

	s.cr.Start()
}

// Задача по удаление старых записей из базы
func (s *Scheduler) deleteOldTasks() {
	s.cr.AddFunc("0 5 * * *", func() { log.Println("В пять утра удалил старые задачи", time.Now()) })
}

// Задача отправляет медиа файлы из базы в заданное время

func (s *Scheduler) SendMedia(hour uint) {
	s.cr.AddFunc(fmt.Sprintf("0 %d * * *", hour),
		func() {
			log.Printf("В %d часов отправил файл", hour)
			log.Println(time.Now())
		})
}
