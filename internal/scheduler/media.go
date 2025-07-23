package scheduler

import (
	"fmt"
	"log"

	tele "gopkg.in/telebot.v3"
)

// setupMediaSchedule настраивает расписание отправки медиа
func (s *Scheduler) setupMediaSchedule() {
	hour, err := s.settingsUC.GetRandomHour()
	if err != nil {
		log.Printf("Не удалось получить расписание по отправке медиа из базы: %v", err)
		return
	}
	entryID, err := s.cr.AddFunc(fmt.Sprintf("10 %d * * *", hour), s.SendRandomMedia)
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
	log.Printf("Расписание по отправке медиа обновлено\n")
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
