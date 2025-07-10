package messages

import tele "gopkg.in/telebot.v3"

var Commands = []tele.Command{
	{Text: "start", Description: "Запустить бота"},
	{Text: "help", Description: "Помощь"},
	{Text: "random", Description: "Случайная мотивация"},
}

type Messages struct {
	// Commands
	Start string
	Help  string
	Add   string
	List  string

	// Input requests
	InputNewDate   string
	InputTaskText  string
	ErrorSomeError string
	ErrorTryAgain  string
	ChooseAction   string

	// Replies
	TaskAdded      string
	FileSaved      string
	FileDeleted    string
	NotImplemented string
	UnknownType    string
	NoOpenTasks    string
	NoTasks        string
	YourTasks      string
	TaskCompleted  string
	TaskDeleted    string
	UnknownText    string

	IncompatibleDate string
	FileExisted      string
}

var BotMessages = Messages{
	Start: "👋 Привет! Я твой таск‑менеджер\nБуду сохранять твои задачи и напоминать о них иногда\nБольше информации по команде /help",
	Help:  "ℹ️ Доступные команды:\n/start — запуск бота\n/settings — настройка бота",

	Add:  "➕ Добавить задачу",
	List: "📋 Список задач",

	InputNewDate:   "📅 Введи новую дату задачи:",
	InputTaskText:  "✏️ Введи описание задачи:",
	ErrorSomeError: "⚠️ Произошла ошибка. Попробуйте позже.",
	ErrorTryAgain:  "⚠️ Произошла ошибка. Попробуйте заново.",
	ChooseAction:   "🔽 Выбери действие:\n",

	TaskAdded:        "☑️ Задача успешно добавлена",
	FileSaved:        "📁 Файл успешно сохранен",
	FileDeleted:      "🗑️ Файл удалён",
	UnknownType:      "❓ Неизвестный тип файла",
	NoOpenTasks:      "✔️ Нет открытых задач",
	NoTasks:          "📝 Нет задач",
	YourTasks:        "📌 Задачи",
	TaskCompleted:    "✅ Задача завершена",
	TaskDeleted:      "❌ Задача удалена",
	UnknownText:      "❓ Unknown text: ",
	IncompatibleDate: "📅 Неподдерживаемый формат даты.\nИспользуйте DD.MM.YYYY или YYYY-MM-DD",

	FileExisted: "‼️ Такой файл уже существует",
}
