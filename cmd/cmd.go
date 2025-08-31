package cmd

import (
	"errors"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/elizavetanr/myDays/calendar"
	"github.com/elizavetanr/myDays/events"
	"github.com/elizavetanr/myDays/logger"
	"github.com/elizavetanr/myDays/reminder"
	"github.com/google/shlex"
	"os"
	"strings"
	"sync"
)

var (
	ErrLoggerFailed = errors.New("ошибка создания лога")
)

type Cmd struct {
	calendar *calendar.Calendar
	log      []string
	mu       sync.Mutex
}

func NewCmd(c *calendar.Calendar) *Cmd {
	return &Cmd{
		calendar: c,
		log:      []string{},
	}
}
func (c *Cmd) executor(input string) {
	input = strings.TrimSpace(input)
	var output string
	if input == "" {
		output = "Вы ввели пустую строку. 'help' для списка команд"
		c.logIOHistory(output)
		return
	}
	c.logIOHistory(input)

	parts, err := shlex.Split(input)
	if err != nil {
		output = "Некорректный ввод команды"
		c.logIOHistory(output)
		c.logError("Ошибка парсинга ввода: " + err.Error())
	}

	cmd := strings.ToLower(parts[0])
	c.calendar.StartAllReminder()
	switch cmd {
	case "add":
		if len(parts) < 4 {
			output = "Формат: add \"название события\" \"дата и время\" \"приоритет\""
			c.logIOHistory(output)
			return
		}

		title := parts[1]
		date := parts[2]
		priority := events.Priority(parts[3])

		event, err := c.calendar.AddEvent(title, date, priority)

		if err != nil {
			switch {
			case errors.Is(err, events.ErrInvalidPriority):
				output = fmt.Sprintf("Некорректный приоритет. Возможные приоритеты: \"%s\", \"%s\", \"%s\"",
					events.PriorityLow, events.PriorityMedium, events.PriorityHigh)
			case errors.Is(err, events.ErrInvalidTitle):
				output = "Некорректное название события. Длина названия от 3 до 50 символов." +
					"\nНазвание может состоять из букв русского и английского алфавита, цифр, пробелов и точек."
			case errors.Is(err, events.ErrInvalidDate):
				output = "Некорректный формат даты. Пример правильного формата: \"2025-10-11 15:00\""
			}
			c.logError(err.Error())
		} else {
			output = "Событие добавлено"
			c.logInfo(fmt.Sprintf("Добавлено событие: ID - %s Title - %s Date - %s Priority - %s ",
				event.ID, event.Title, event.StartAt.Format("02.01.2006  15:04:05"), string(event.Priority)))
		}

	case "update":
		if len(parts) < 5 {
			output = "Формат: update \"ID события\" \"название события\" \"дата и время\" \"приоритет\""
			c.logIOHistory(output)
			return
		}
		ID := parts[1]
		title := parts[2]
		date := parts[3]
		priority := events.Priority(parts[4])
		err = c.calendar.EditEvent(ID, title, date, priority)
		if err != nil {
			switch {
			case errors.Is(err, calendar.ErrEventNotFound):
				output = "Событие с введенным id не найдено"
			case errors.Is(err, events.ErrInvalidPriority):
				output = fmt.Sprintf("Некорректный приоритет. Возможные приоритеты: \"%s\", \"%s\", \"%s\"",
					events.PriorityLow, events.PriorityMedium, events.PriorityHigh)
			case errors.Is(err, events.ErrInvalidTitle):
				output = "Некорректное название события. Длина названия от 3 до 50 символов." +
					"\nНазвание может состоять из букв русского и английского алфавита, цифр, пробелов и точек."
			case errors.Is(err, events.ErrInvalidDate):
				output = "Некорректный формат даты. Пример правильного формата: \"2025-10-11 15:00\""
			}
			c.logError(err.Error())
		} else {
			output = "Событие изменено"
			c.logInfo(fmt.Sprintf("Изменено событие с ID - %s: Title - %s Date - %s Priority - %s ",
				ID, title, date, priority))
		}

	case "remove":
		if len(parts) < 2 {
			output = "Формат: remove \"ID события\""
			c.logIOHistory(output)
			return
		}
		ID := parts[1]
		err = c.calendar.DeleteEvent(ID)
		if errors.Is(err, calendar.ErrEventNotFound) {
			output = "Событие с введенным id не найдено"
			c.logError(err.Error())
		} else {
			output = "Событие удалено"
			c.logInfo(fmt.Sprintf("Удалено событие с ID - %s", ID))
		}
	case "list":
		calendarEvents := c.calendar.GetEvent()
		if len(calendarEvents) == 0 {
			output = "Список событий пуст"
			c.logIOHistory(output)
			return
		}
		output = ""
		for _, event := range calendarEvents {
			dateStr := event.StartAt.Format("2006-01-02 15:04")
			outputEvent := event.Title + " - " + dateStr + " - ID: " + event.ID
			output += outputEvent + "\n"
		}
	case "add_reminder":
		if len(parts) < 4 {
			output = "Формат: add_reminder \"ID события\" \"текст напоминания\" \"интервал до события\""
			c.logIOHistory(output)
			return
		}
		ID := parts[1]
		message := parts[2]
		before := parts[3]
		err = c.calendar.SetEventReminder(ID, message, before)
		if err != nil {
			switch {
			case errors.Is(err, reminder.ErrTimeReminderIsUp):
				output = "Нельзя запустить напоминание с истекшим временем"
			case errors.Is(err, calendar.ErrEventNotFound):
				output = "Событие с введенным id не найдено"
			case errors.Is(err, calendar.ErrInvalidDuration):
				output = "Некорректный ввод интервала. Примеры правильного ввода: \"2h45m\", \"1.5h\", \"120m\""
			case errors.Is(err, calendar.ErrEventExpired):
				output = "Нельзя добавить напоминание прошедшему событию"
			case errors.Is(err, calendar.ErrReminderTimeAfterEvent):
				output = "Нельзя добавить напоминание после начала события"
			case errors.Is(err, calendar.ErrReminderTimeBeforeNow):
				output = "Нельзя добавить напоминание раньше текущего времени"
			}
			c.logError(err.Error())
		} else {
			output = "Напоминание добавлено и запущено"
			c.logInfo(fmt.Sprintf("Добавлено напоминание к событию с ID - %s: Message - %s Before - %s",
				ID, message, before))
		}

	case "remove_reminder":
		if len(parts) < 2 {
			output = "Формат: remove_reminder \"ID события\""
			c.logIOHistory(output)
			return
		}
		ID := parts[1]
		err = c.calendar.CancelEventReminder(ID)
		if err != nil {
			switch {
			case errors.Is(err, calendar.ErrEventNotFound):
				output = "Событие с введенным id не найдено"
			case errors.Is(err, reminder.ErrNotExistReminder):
				output = "У этого события не существует напоминания"
			}
			c.logError(err.Error())
		} else {
			output = "Напоминание удалено"
			c.logInfo(fmt.Sprintf("Удалено напоминание у события с ID - %s", ID))
		}
	case "help":
		output = "Доступные команды:" +
			"\nДобавление события: add \"название события\" \"дата и время\" \"приоритет\"" +
			"\nРедактирование события: update \"ID события\" \"название события\" \"дата и время\" \"приоритет\"" +
			"\nУдаление события: remove \"ID события\"" +
			"\nДобавление напоминания: add_reminder \"ID события\" \"текст напоминания\" \"интервал до события\"" +
			"\nУдаление напоминания: remove_reminder \"ID события\"" +
			"\nВывести список всех событий: list" +
			"\nВывести список всех команд: help" +
			"\nВывести логи: log" +
			"\nВыход из приложения: exit"
	case "log":
		c.showLogIOHistory()
	case "exit":
		err = c.calendar.Save()
		if err != nil {
			output = "Сохранение не выполнено"
			c.logError(err.Error())
		} else {
			output = "Сохранено"
			c.logInfo("Выполнено сохранение календаря")
		}
		c.logInfo("Приложение закрыто")
		os.Exit(0)
	default:
		output = "Неизвестная команда. Введите 'help' для списка команд"
	}
	c.logIOHistory(output)
}

func (c *Cmd) completer(d prompt.Document) []prompt.Suggest {
	suggestions := []prompt.Suggest{
		{Text: "add", Description: "Добавить событие"},
		{Text: "list", Description: "Показать все события"},
		{Text: "remove", Description: "Удалить событие"},
		{Text: "add_reminder", Description: "Добавить напоминание"},
		{Text: "remove_reminder", Description: "Удалить напоминание"},
		{Text: "help", Description: "Показать справку"},
		{Text: "log", Description: "Показать логи"},
		{Text: "exit", Description: "Выйти из программы"},
	}
	return prompt.FilterHasPrefix(suggestions, d.GetWordAfterCursor(), true)
}

func (c *Cmd) Run() {
	p := prompt.New(
		c.executor,
		c.completer,
		prompt.OptionPrefix("> "),
	)
	go func() {
		for msg := range c.calendar.Notification {
			c.logIOHistory(msg)
			c.logInfo(fmt.Sprintf("Пользователю выведено напоминание: %s", msg))
		}
	}()
	c.logInfo("Приложение запущено")
	p.Run()
	close(c.calendar.Notification)
}

func (c *Cmd) logIOHistory(log string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.log = append(c.log, log)
	fmt.Println(log)
}

func (c *Cmd) showLogIOHistory() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, log := range c.log {
		fmt.Println(log)
	}
}

func (c *Cmd) logError(log string) {
	if err := logger.Error("Ошибка: " + log); err != nil {
		fmt.Println(ErrLoggerFailed)
		c.logIOHistory(ErrLoggerFailed.Error())
	}
}

func (c *Cmd) logInfo(log string) {
	if err := logger.Info(log); err != nil {
		fmt.Println(ErrLoggerFailed)
		c.logIOHistory(ErrLoggerFailed.Error())
	}
}
