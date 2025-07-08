package keyboards

import (
	"fmt"
	"time"

	"github.com/SoulStalker/cognitask/internal/messages"
)

// todo надо подумать может лучше просто принимать цифру текущего месяца
// Парсинг даты из строки
func ParseDate(dateStr string) (time.Time, error) {
	
	switch dateStr {
	case BtnToday.Unique:
		dateStr = GetTodayDate()
	case BtnTomorrow.Unique:
		dateStr = GetTomorrowDate()
	case BtnSkipDate.Unique:
		return time.Time{}, nil
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
