package calendar

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elizavetanr/myDays/events"
	"github.com/elizavetanr/myDays/storage"
	"time"
)

var (
	ErrEventNotFound          = errors.New("событие с введенным id не найдено")
	ErrInvalidDuration        = errors.New("некорректный формат интервала")
	ErrEventExpired           = errors.New("событие уже истекло")
	ErrReminderTimeAfterEvent = errors.New("время напоминания позже времени события")
	ErrReminderTimeBeforeNow  = errors.New("время напоминания раньше текущего времени")
	ErrMarshalFailed          = errors.New("сериализация не выполнена")
	ErrUnmarshalFailed        = errors.New("десериализация не выполнена")
	ErrCalendarSaveFailed     = errors.New("сохранение данных в файл не выполнено")
	ErrCalendarLoadFailed     = errors.New("загрузка данных из файла не выполнена")
)

type Calendar struct {
	calendarEvents map[string]*events.Event
	storage        storage.Store
	Notification   chan string
}

func (c *Calendar) Save() error {
	data, err := json.Marshal(c.calendarEvents)
	if err != nil {
		return fmt.Errorf("не удалось сохранить календарь: %w", ErrMarshalFailed)
	}
	err = c.storage.Save(data)
	if err != nil {
		return fmt.Errorf("не удалось сохранить календарь: %w", ErrCalendarSaveFailed)
	}
	return nil
}

func (c *Calendar) Load() error {
	data, err := c.storage.Load()
	if err != nil {
		return fmt.Errorf("не удалось загрузить календарь: %w", ErrCalendarLoadFailed)
	}
	err = json.Unmarshal(data, &c.calendarEvents)
	if err != nil {
		return fmt.Errorf("не удалось загрузить календарь: %w", ErrUnmarshalFailed)
	}
	return nil
}

func NewCalendar(s storage.Store) *Calendar {
	return &Calendar{
		calendarEvents: make(map[string]*events.Event),
		storage:        s,
		Notification:   make(chan string),
	}
}

func (c *Calendar) AddEvent(title string, date string, priority events.Priority) (*events.Event, error) {
	event, err := events.NewEvent(title, date, priority)
	if err != nil {
		return nil, fmt.Errorf("невозможно добавить событие: %w", err)
	}
	c.calendarEvents[event.ID] = event
	return event, nil
}
func (c *Calendar) DeleteEvent(id string) error {
	if !c.idExists(id) {
		return fmt.Errorf("невозможно удалить событие: %w", ErrEventNotFound)
	}

	delete(c.calendarEvents, id)
	return nil
}

func (c *Calendar) EditEvent(id, newTitle, newDate string, priority events.Priority) error {
	if !c.idExists(id) {
		return fmt.Errorf("невозможно отредактировать событие: %w", ErrEventNotFound)
	}
	if err := c.calendarEvents[id].Update(newTitle, newDate, priority); err != nil {
		return fmt.Errorf("невозможно отредактировать событие: %w", err)
	}
	return nil
}

func (c *Calendar) GetEvent() map[string]*events.Event {
	return c.calendarEvents
}

func (c *Calendar) SetEventReminder(id, message, before string) error {
	if !c.idExists(id) {
		return fmt.Errorf("невозможно назначить напоминание событию: %w", ErrEventNotFound)
	}
	reminderAt, err := c.calculateReminderTime(id, before)
	if err != nil {
		return fmt.Errorf("невозможно назначить напоминание событию: %w", err)
	}

	if err := c.calendarEvents[id].AddReminder(message, reminderAt); err != nil {
		return fmt.Errorf("невозможно назначить напоминание событию: %w", err)
	}
	if err := c.calendarEvents[id].StartReminder(c.Notify); err != nil {
		return fmt.Errorf("невозможно запустить добавленное напоминание: %w", err)
	}
	return nil

}

func (c *Calendar) calculateReminderTime(id, before string) (time.Time, error) {
	duration, err := time.ParseDuration(before)
	if err != nil {
		return time.Time{}, ErrInvalidDuration
	}

	e := c.calendarEvents[id]
	eventStartAt := e.StartAt
	reminderAt := eventStartAt.Add(-duration)
	if eventStartAt.Before(time.Now()) {
		return time.Time{}, ErrEventExpired
	}
	if reminderAt.After(eventStartAt) {
		return time.Time{}, ErrReminderTimeAfterEvent
	}
	if reminderAt.Before(time.Now()) {
		return time.Time{}, ErrReminderTimeBeforeNow
	}
	return reminderAt, nil
}

func (c *Calendar) CancelEventReminder(id string) error {
	if !c.idExists(id) {
		return fmt.Errorf("невозможно удалить напоминание у события: %w", ErrEventNotFound)
	}
	e := c.calendarEvents[id]
	err := e.Reminder.Stop()
	e.RemoveReminder()
	return err
}
func (c *Calendar) StartAllReminder() {
	for _, event := range c.calendarEvents {
		if event.Reminder != nil {
			event.StartReminder(c.Notify)
		}
	}
}
func (c *Calendar) Notify(msg string) {
	c.Notification <- msg
}

func (c *Calendar) idExists(id string) bool {
	_, exists := c.calendarEvents[id]
	return exists
}
