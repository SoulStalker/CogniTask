package messages

type Messages struct {
    // Commands
    Start string
    Help  string
	Add string
	List string

    
    // Buttons  
    Settings string      // без "Btn"
    CompleteTask string  // "✅ Выполнить"
    RandomPic string     // "🎲 Random Pic"
    
    // Input requests
    InputNewDate string
    InputTaskText string
}

var BotMessages = Messages{
	Start: "Привет! Я твой таск-менеджер\n. Буду сохранять твои задачи и напоминать о них иногда\nБольше информации по команде /help",
	Help: "Доступные команды: /Settings - настройка бота\n...",
	Settings: "⚙️ Настройка",
	Add: "Добавить задачу",
	List: "Список задач",
	CompleteTask: "✅ Выполнить",
	RandomPic: "🎲 Random Pic",
	InputNewDate: "Введи новую дату задачи: ",
    InputTaskText: "Введи описание задачи: ",
}