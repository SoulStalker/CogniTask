package keyboards

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/SoulStalker/cognitask/internal/messages"
)

// todo надо подумать может лучше просто принимать цифру текущего месяца
// Парсинг даты из строки
func ParseDate(dateStr string) (time.Time, error) {

	calendarDate, err := mapDateString(dateStr)
	if err != nil {
		log.Println(err)
	} else {
		dateStr = calendarDate
	}

	switch dateStr {
	case BtnToday.Unique:
		dateStr = GetTodayDate()
	case BtnTomorrow.Unique:
		dateStr = GetTomorrowDate()
	case BtnSkipDate.Unique:
		return time.Time{}, nil
	case calendarDate:
		if date, err := time.Parse("02.01.2006", calendarDate); err == nil {
			return date, nil
		}
	}

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

	return time.Time{}, fmt.Errorf(messages.BotMessages.IncompatibleDate)
}

func GetTodayDate() string {
	return time.Now().Format("2006-01-02")
}

func GetTomorrowDate() string {
	return time.Now().AddDate(0, 0, 1).Format("2006-01-02")
}

// Функция преобразует строку формата "day|DAY|YYYY|MM|DD" в "DD.MM.YYYY"
func mapDateString(input string) (string, error) {
	parts := strings.Split(input, "|")
	if len(parts) != 5 {
		return "", fmt.Errorf("неверный формат строки, ожидается day|DAY|YYYY|MM|DD")
	}

	year := parts[2]  // YYYY
	month := parts[3] // MM
	day := parts[4]   // DD

	// Проверяем, что компоненты даты корректны (можно добавить дополнительные проверки)
	if len(year) != 4 || len(month) != 2 || len(day) != 2 {
		return "", fmt.Errorf("некорректный формат даты")
	}

	return fmt.Sprintf("%s.%s.%s", day, month, year), nil
}
