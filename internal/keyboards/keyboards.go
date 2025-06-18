package keyboards

import (
	"fmt"
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

// GetTaskControlKeyboard клавиатура для управления задачей
func GetTaskControlKeyboard(taskID uint) *tele.ReplyMarkup {
	kb := &tele.ReplyMarkup{}

	// Inline кнопки с callback data
	btnComplete := kb.Data(BtnComplete, "complete", fmt.Sprintf("%d", taskID))
	btnRandomPic := kb.Data(BtnRandomPic, "random_pic", "")

	// Раскладка
	kb.Inline(
		kb.Row(btnComplete),
		kb.Row(btnRandomPic), // todo или сюда лучше отмену поставить?
	)

	return kb
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
// Парсинг даты из строки (для календаря)
func ParseDate(dateStr string) (time.Time, error) {
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

	return time.Time{}, fmt.Errorf("неподдерживаемый формат даты. Используйте: DD.MM.YYYY или YYYY-MM-DD")
}
