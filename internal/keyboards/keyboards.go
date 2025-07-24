package keyboards

import (
	"fmt"
	"strconv"
	"time"

	tele "gopkg.in/telebot.v3"
)

// Глобальные кнопки для регистрации хендлеров
var (
	BtnCancel = &tele.InlineButton{Unique: "cancel", Text: "🚫 Отмена"}
	// кнопки основной клавы
	BtnSettings  = &tele.InlineButton{Unique: "settings", Text: "Настройки"}
	BtnPending   = &tele.InlineButton{Unique: "pending", Text: "Текущие задачи"}
	BtnAll       = &tele.InlineButton{Unique: "all_tasks", Text: "Все задачи"}
	BtnRandomPic = &tele.InlineButton{Unique: "random_pic", Text: "🎲 Random Pic"}

	// кнопки задач
	BtnComplete = &tele.InlineButton{Unique: "complete_task", Text: "✅ Выполнить"}
	BtnDelete   = &tele.InlineButton{Unique: "delete_task", Text: "🗑 Удалить"}
	BtnEditDate = &tele.InlineButton{Unique: "edit_date", Text: "📅 Изменить дату"}
	BtnAdd      = &tele.InlineButton{Unique: "add", Text: "Новая задача"}

	// кнопки добавления задачи
	BtnToday    = &tele.InlineButton{Unique: "today", Text: "📅 Сегодня"}
	BtnTomorrow = &tele.InlineButton{Unique: "tomorrow", Text: "🌅 Завтра"}
	BtnCalendar = &tele.InlineButton{Unique: "choose", Text: "🗓️ Выбрать"}
	BtnSkipDate = &tele.InlineButton{Unique: "skip", Text: "⏭️ Пропустить"}

	// кнопки настроек
	BtnAutoDelete    = &tele.InlineButton{Unique: "setDeleteDays", Text: "🗑️ Авто‑удаление"}
	BtnNotifications = &tele.InlineButton{Unique: "setNotifications", Text: "⏰ Уведомления"}
	BtnNotifyFrom    = &tele.InlineButton{Unique: "setNotifyFrom", Text: "📅 Начало"}
	BtnNotifyTo      = &tele.InlineButton{Unique: "setNotifyTo", Text: "📅 Конец"}
	BtnRandomHour    = &tele.InlineButton{Unique: "setRandomHour", Text: "💡 Мотивация"}
)

// GetDateSelectionKeyboard клавиатура для выбора даты задачи
func GetDateSelectionKeyboard() *tele.ReplyMarkup {
	kb := &tele.ReplyMarkup{}

	// Кнопки
	btnToday := kb.Data(BtnToday.Text, BtnToday.Unique)
	btnTomorrow := kb.Data(BtnTomorrow.Text, BtnTomorrow.Unique)
	btnCalendar := kb.Data(BtnCalendar.Text, BtnCalendar.Unique)
	btnSkip := kb.Data(BtnSkipDate.Text, BtnSkipDate.Unique)
	btnCancel := kb.Data(BtnCancel.Text, BtnCancel.Unique)

	// Раскладка кнопок три ряда: Сегодня | Завтра, Выбрать, Пропустить | Отмена
	kb.Inline(
		kb.Row(btnToday, btnTomorrow),
		kb.Row(btnCalendar, btnSkip),
		kb.Row(btnCancel),
	)

	return kb
}

// CreateMainKeyboard основная клавиатура
func CreateMainKeyboard() *tele.ReplyMarkup {
	kb := &tele.ReplyMarkup{
		// ResizeKeyboard:  true,
		// OneTimeKeyboard: true,
	}

	btnAdd := kb.Data(BtnAdd.Text, BtnAdd.Unique)
	btnSettings := kb.Data(BtnSettings.Text, BtnSettings.Unique)
	btnPending := kb.Data(BtnPending.Text, BtnPending.Unique)
	btnAll := kb.Data(BtnAll.Text, BtnAll.Unique)
	btnCancel := kb.Data(BtnCancel.Text, BtnCancel.Unique)

	kb.Inline(
		kb.Row(btnAdd, btnSettings),
		kb.Row(btnPending, btnAll),
		kb.Row(btnCancel),
	)

	return kb
}

// CreateTaskKeyboard клавиатура для управления задачей
func CreateTaskKeyboard(taskID uint) *tele.ReplyMarkup {
	kb := &tele.ReplyMarkup{}

	// Inline кнопки с callback data
	btnComplete := kb.Data(BtnComplete.Text, BtnComplete.Unique, fmt.Sprintf("%d", taskID))
	btnRandomPic := kb.Data(BtnRandomPic.Text, BtnRandomPic.Unique, "")
	btnDelete := kb.Data(BtnDelete.Text, BtnDelete.Unique, fmt.Sprintf("%d", taskID))
	btnEditDate := kb.Data(BtnEditDate.Text, BtnEditDate.Unique, fmt.Sprintf("%d", taskID))

	// Раскладка
	kb.Inline(
		kb.Row(btnComplete),
		kb.Row(btnEditDate, btnDelete),
		kb.Row(btnRandomPic),
	)

	return kb
}

// CreateSettingsKeyboard создает клавиатуру для настроек
func CreateSettingsKeyboard() *tele.ReplyMarkup {
	kb := &tele.ReplyMarkup{}

	btnAutoDelete := kb.Data(BtnAutoDelete.Text, BtnAutoDelete.Unique)
	btnNotifications := kb.Data(BtnNotifications.Text, BtnNotifications.Unique)
	btnNotifyFrom := kb.Data(BtnNotifyFrom.Text, BtnNotifyFrom.Unique)
	btnNotifyTo := kb.Data(BtnNotifyTo.Text, BtnNotifyTo.Unique)
	btnRandomHour := kb.Data(BtnRandomHour.Text, BtnRandomHour.Unique)

	cancel := kb.Data(BtnCancel.Text, BtnCancel.Unique)

	kb.Inline(
		kb.Row(btnAutoDelete, btnNotifications),
		kb.Row(btnNotifyFrom, btnNotifyTo),
		kb.Row(btnRandomHour),
		kb.Row(cancel),
	)

	return kb
}

// CreateMainKeyboard основная клавиатура
func CreateCancelKeyboard() *tele.ReplyMarkup {
	kb := &tele.ReplyMarkup{}
	btnCancel := kb.Data(BtnCancel.Text, BtnCancel.Unique)
	kb.Inline(kb.Row(btnCancel))

	return kb
}

func CreateHoursKeyboard(rowsCount int) *tele.ReplyMarkup {
	kb := &tele.ReplyMarkup{}
	if rowsCount < 1 {
		rowsCount = 1
	}

	var btns []tele.Btn
	for i := 1; i <= 24; i++ {
		btns = append(btns, kb.Data(fmt.Sprint(i), fmt.Sprint(i)))
	}

	cancel := kb.Data(BtnCancel.Text, BtnCancel.Unique)

	perRow := (len(btns) + rowsCount - 1) / rowsCount
	rows := kb.Split(perRow, btns)

	rows = append(rows, kb.Row(cancel))

	kb.Inline(rows...)
	return kb
}

/* ---------- построение календаря ---------- */

func BuildKeyboard(year int, month time.Month) *tele.ReplyMarkup {
	const (
		btnWidth = 7 // 7 дней в строке
	)

	loc := time.Local
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, loc)
	daysInMonth := firstOfMonth.AddDate(0, 1, -1).Day()
	weekdayOffset := int(firstOfMonth.Weekday())
	if weekdayOffset == 0 { // в Go Sunday==0, хотим, чтобы Monday==0
		weekdayOffset = 6
	} else {
		weekdayOffset--
	}

	m := &tele.ReplyMarkup{}

	// навигация
	prev := firstOfMonth.AddDate(0, -1, 0)
	next := firstOfMonth.AddDate(0, 1, 0)
	btnPrev := m.Data("◀︎", "nav_prev",
		fmt.Sprintf("NAV|%d|%d", prev.Year(), int(prev.Month())))
	btnNext := m.Data("▶︎", "nav_next",
		fmt.Sprintf("NAV|%d|%d", next.Year(), int(next.Month())))
	m.Inline(
		m.Row(btnPrev, btnNext),
	)

	// пустые ячейки до 1-го числа
	cells := make([]tele.Btn, weekdayOffset)
	for i := range cells {
		cells[i] = m.Data(" ", "empty", "x") // «пустая» кнопка
	}

	// сами дни
	for d := 1; d <= daysInMonth; d++ {
		date := fmt.Sprintf("DAY|%d|%02d|%02d", year, month, d)
		cells = append(cells, m.Data(strconv.Itoa(d), "day", date))
	}

	// разбиваем в строки по 7 кнопок
	for len(cells)%btnWidth != 0 {
		cells = append(cells, m.Data(" ", "empty", "x"))
	}
	for i := 0; i < len(cells); i += btnWidth {
		m.InlineKeyboard = append(m.InlineKeyboard, cells[i:i+btnWidth])
	}
	return m
}
