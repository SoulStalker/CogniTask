package messages

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
	TaskAdded string
}

var BotMessages = Messages{
	Start: "Привет! Я твой таск-менеджер\nБуду сохранять твои задачи и напоминать о них иногда\nБольше информации по команде /help",
	Help:  "Доступные команды:\n/start - запуск бота\n/settings - настройка бота\n/add - добавить задачу\n/pending - список не выполненных задач",

	Add:  "Добавить задачу",
	List: "Список задач",

	InputNewDate:   "Введи новую дату задачи: ",
	InputTaskText:  "Введи описание задачи: ",
	ErrorSomeError: "Произошла ошибка. Попробуйте позже.",
	ErrorTryAgain:  "Произошла ошибка. Попробуйте заново.",
	TaskAdded:      "☑️ Задача успешно добавлена",
	ChooseAction:   "Выбери действие:\n",
}
