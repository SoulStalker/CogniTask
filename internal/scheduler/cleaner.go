package scheduler

import "log"

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
