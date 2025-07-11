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

	s.cr.Start()
}

func (s *Scheduler) deleteOldTasks() {
	s.cr.AddFunc("0 5 * * *", func() { log.Println("В пять утра удалил старые задачи", time.Now()) })
}

func (s *Scheduler) SendMedia(hour uint) {
	s.cr.AddFunc(fmt.Sprintf("0 %d * * *", hour),
		func() {
			log.Printf("В %d часов отправил файл", hour)
			log.Println(time.Now())
		})
}

func (s *Scheduler) Notifier(hour uint) cron.EntryID {
	entryID, err := s.cr.AddFunc(fmt.Sprintf("%d * * * *", hour), func() { // потом на часы переделать
		log.Printf("Уведомления каждые %d часа", hour)
	})
	if err != nil {
		return 0
	}
	return entryID
}

func (s *Scheduler) StartNotifications(hour uint) cron.EntryID {
	s.cr.AddFunc(
		fmt.Sprintf("* %d * * *", hour), func() { log.Println("Начинаю слать уведомления", time.Now()) },
	)
	return s.Notifier(3)
}

func (s *Scheduler) StopNotifications(hour uint, entry cron.EntryID) {
	s.cr.AddFunc(
		fmt.Sprintf("* %d * * *", hour), func() { log.Println("В пять утра удалил старые задачи", time.Now()) },
	)
	s.cr.Remove(entry)
}
