package handlers

import (
	"fmt"

	"github.com/SoulStalker/cognitask/internal/domain"
	"github.com/SoulStalker/cognitask/internal/keyboards"
	"github.com/SoulStalker/cognitask/internal/messages"
	tele "gopkg.in/telebot.v3"
)

func formatTaskList(tasks []domain.Task) [][]tele.InlineButton {
 	var rows [][]tele.InlineButton

	for _, task := range tasks {
		btn := tele.InlineButton{
			Unique: keyboards.BtnComplete.Unique,
			Text:   fmt.Sprintf("✅ %s срок: %s", task.Description, task.Deadline.Format("02-01")),
			Data:   fmt.Sprint(task.ID),
		}
		rows = append(rows, []tele.InlineButton{btn})
	}
	cancelBtn := tele.InlineButton{Unique: "cancel", Text: keyboards.BtnCancel}
	rows = append(rows, []tele.InlineButton{cancelBtn})

	return rows
}

func formatTask(task domain.Task) string {
	return fmt.Sprintf(
		"%s\n\n%s\n%s",
		messages.BotMessages.TaskAdded,
		task.Description, task.Deadline.Format("02-01-2006"))
}
