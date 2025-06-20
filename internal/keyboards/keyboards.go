package keyboards

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

const (
	// Кнопки выбора даты
	BtnToday    = "📅 Сегодня"
	BtnTomorrow = "🌅 Завтра"
	BtnCalendar = "🗓️ Выбрать"
	BtnSkipDate = "⏭️ Пропустить"
	// BtnCancel   = "🚫 Отмена"

)

// Глобальные кнопки для регистрации хендлеров
var (
	BtnComplete  = &tele.InlineButton{Unique: "complete_task", Text: "✅ Выполнить"}
	BtnDelete    = &tele.InlineButton{Unique: "delete_task", Text: "🗑 Удалить"}
	BtnEditDate  = &tele.InlineButton{Unique: "edit_date", Text: "📅 Изменить дату"}
	BtnCancel    = &tele.InlineButton{Unique: "cancel", Text: "🚫 Отмена"}
	BtnAdd       = &tele.InlineButton{Unique: "add", Text: "Новая задача"}
	BtnSettings  = &tele.InlineButton{Unique: "settings", Text: "Настройки"}
	BtnPending   = &tele.InlineButton{Unique: "pending", Text: "Текущие задачи"}
	BtnAll       = &tele.InlineButton{Unique: "all_tasks", Text: "Все задачи"}
	BtnRandomPic = &tele.InlineButton{Unique: "random_pic", Text: "🎲 Random Pic"}
)

// GetDateSelectionKeyboard клавиатура для выбора даты задачи
func GetDateSelectionKeyboard() *tele.ReplyMarkup {
	kb := &tele.ReplyMarkup{
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}

	// Кнопки
	btnToday := kb.Text(BtnToday)
	btnTomorrow := kb.Text(BtnTomorrow)
	btnCalendar := kb.Text(BtnCalendar)
	btnSkip := kb.Text(BtnSkipDate)
	// BtnCancel := kb.Text(BtnCancel)

	// Раскладка кнопок три ряда: Сегодня | Завтра, Выбрать, Пропустить | Отмена
	kb.Reply(
		kb.Row(btnToday, btnTomorrow),
		kb.Row(btnCalendar),
		kb.Row(btnSkip),
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
	markup := &tele.ReplyMarkup{}

	// Кнопки настроек (заготовка для будущего)
	btnNotifications := markup.Data("🔔 Уведомления", "settings_notifications")
	btnAutoDelete := markup.Data("🗑 Авто-удаление", "settings_auto_delete")

	markup.Inline(
		markup.Row(btnNotifications),
		markup.Row(btnAutoDelete),
	)

	return markup
}
