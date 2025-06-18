package keyboards

import (
	"fmt"
	"log"
	"time"

	tele "gopkg.in/telebot.v3"
)

const (
	// Кнопки выбора даты
	BtnToday    = "📅 Сегодня"
	BtnTomorrow = "🌅 Завтра"
	BtnCalendar = "🗓️ Выбрать"
	BtnSkipDate = "⏭️ Пропустить"
	BtnCancel   = "🚫 Отмена"

	// Кнопки управления задачами
	BtnComplete  = "✅ Выполнить"
	BtnDelete    = "🗑 Удалить"
	BtnEditDate  = "📅 Изменить дату"
	BtnRandomPic = "🎲 Random Pic"
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
	BtnCancel := kb.Text(BtnCancel)

	// Раскладка кнопок три ряда: Сегодня | Завтра, Выбрать, Пропустить | Отмена
	kb.Reply(
		kb.Row(btnToday, btnTomorrow),
		kb.Row(btnCalendar),
		kb.Row(btnSkip, BtnCancel),
	)

	return kb
}

// CreateTaskKeyboard клавиатура для управления задачей
func CreateTaskKeyboard(taskID uint) *tele.ReplyMarkup {
	kb := &tele.ReplyMarkup{}

	// Inline кнопки с callback data
	btnComplete := kb.Data(BtnComplete, "complete", fmt.Sprintf("%d", taskID))
	btnRandomPic := kb.Data(BtnRandomPic, "random_pic", "")
	btnDelete := kb.Data(BtnDelete, "delete_task", fmt.Sprintf("%d", taskID))
	btnEditDate := kb.Data(BtnEditDate, "edit_date", fmt.Sprintf("%d", taskID))

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

// GetRemoveKeyboard удаляет reply клаву
func GetRemoveKeyboard() *tele.ReplyMarkup {
	return &tele.ReplyMarkup{
		RemoveKeyboard: true,
	}
}

func GetTodayDate() string {
	return time.Now().Format("2006-01-02")
}

func GetTomorrowDate() string {
	return time.Now().AddDate(0, 0, 1).Format("2006-01-02")
}

// todo надо подумать может лучше просто принимать цифру текущего месяца
// Парсинг даты из строки
func ParseDate(dateStr string) (time.Time, error) {
	switch dateStr {
	case BtnToday:
		dateStr = GetTodayDate()
	case BtnTomorrow:
		dateStr = GetTomorrowDate()
	}

	log.Printf("Выбрана дата для парсинга: %s", dateStr)
	// Поддерживаемые форматы
	formats := []string{
		"2006-01-02", // 2024-12-25
		"02.01.2006", // 25.12.2024
		"02/01/2006", // 25/12/2024
		"02-01-2006", // 25-12-2024
	}

	for _, format := range formats {
		if date, err := time.Parse(format, dateStr); err == nil {
			return date, nil
		}
	}

	return time.Time{}, fmt.Errorf("неподдерживаемый формат даты.\nИспользуйте: DD.MM.YYYY или YYYY-MM-DD")
}
