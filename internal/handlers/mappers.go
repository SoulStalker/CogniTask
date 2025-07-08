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
			Data: fmt.Sprint(task.ID),
		}
		if !task.Closed {
			btn.Text = fmt.Sprintf("✅ %s : %s", task.Description, task.Deadline.Format("02-01"))
		} else {
			btn.Unique = keyboards.BtnDelete.Unique
			btn.Text = fmt.Sprintf("❌ %s : %s", task.Description, task.ClosedAt.Format("02-01"))
		}
		rows = append(rows, []tele.InlineButton{btn})
	}
	cancelBtn := tele.InlineButton{Unique: "cancel", Text: keyboards.BtnCancel.Unique}
	rows = append(rows, []tele.InlineButton{cancelBtn})

	return rows
}

func formatTask(task domain.Task) string {
	return fmt.Sprintf(
		"%s\n\n%s\n%s",
		messages.BotMessages.TaskAdded,
		task.Description, task.Deadline.Format("02-01-2006"))
}
