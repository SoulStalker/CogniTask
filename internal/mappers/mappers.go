package mappers

import (
	"fmt"
	"time"

	"github.com/SoulStalker/cognitask/internal/domain"
	"github.com/SoulStalker/cognitask/internal/keyboards"
	"github.com/SoulStalker/cognitask/internal/messages"
	tele "gopkg.in/telebot.v3"
)

var russianShortMonths = map[time.Month]string{
	time.January:   "янв.",
	time.February:  "фев.",
	time.March:     "мар.", // Обратите внимание, для короткой формы иногда используются другие склонения или формы
	time.April:     "апр.",
	time.May:       "мая",  // Короткая форма для мая часто совпадает с полной
	time.June:      "июня", // Аналогично
	time.July:      "июля", // Аналогично
	time.August:    "авг.",
	time.September: "сен.",
	time.October:   "окт.",
	time.November:  "ноя.",
	time.December:  "дек.",
}

func FormatTaskList(tasks []domain.Task) [][]tele.InlineButton {
	var rows [][]tele.InlineButton

	for _, task := range tasks {
		btn := tele.InlineButton{
			Unique: keyboards.BtnComplete.Unique,
			Data:   fmt.Sprint(task.ID),
		}
		var taskDate string
		if !task.Deadline.IsZero() {
			t := task.Deadline
			day := t.Day()
			monthName := russianShortMonths[t.Month()]

			// taskDate = task.Deadline.Format("02-01")
			taskDate = fmt.Sprintf("%d %s", day, monthName)
		}
		if !task.Closed {
			btn.Text = fmt.Sprintf("✅ %s : %s", task.Description, taskDate)
		} else {
			btn.Unique = keyboards.BtnDelete.Unique
			btn.Text = fmt.Sprintf("❌ %s : %s", task.Description, taskDate)
		}
		rows = append(rows, []tele.InlineButton{btn})
	}
	cancelBtn := tele.InlineButton{Unique: keyboards.BtnCancel.Unique, Text: keyboards.BtnCancel.Text}
	rows = append(rows, []tele.InlineButton{cancelBtn})

	return rows
}

func FormatTask(task domain.Task) string {
	return fmt.Sprintf(
		"%s\n\n%s\n%s",
		messages.BotMessages.TaskAdded,
		task.Description, task.Deadline.Format("02-01-2006"))
}
