package handlers

import (
	"fmt"

	"github.com/SoulStalker/cognitask/internal/domain"
	"github.com/SoulStalker/cognitask/internal/messages"
)

func formatTaskList(tasks []domain.Task) string {
	var tasksStr = "Задача: дедлайн, создана\n"
	for _, task := range tasks {
		tasksStr += fmt.Sprintf(
			"%s: %s, %s\n",
			task.Description,
			task.Deadline.Format("02-01"),
			task.CreatedAt.Format("02-01"),
		)
	}
	return tasksStr
}

func formatTask(task domain.Task) string {
	return fmt.Sprintf(
		"%s\n\n%s\n%s", 
		messages.BotMessages.TaskAdded, 
		task.Description, task.Deadline.Format("02-01-2006"))
}
