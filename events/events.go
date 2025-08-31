package events

import (
	"errors"
	"github.com/elizavetanr/myDays/reminder"
	"github.com/google/uuid"
	"time"
)

var (
	ErrEmptyReminder = errors.New("пустое напоминание")
)

type Event struct {
	ID       string             `json:"id"`
	Title    string             `json:"title"`
	StartAt  time.Time          `json:"date"`
	Priority Priority           `json:"priority"`
	Reminder *reminder.Reminder `json:"reminder"`
}

func NewEvent(title, date string, priority Priority) (*Event, error) {
	startAt, err := ValidateInput(title, date)
	if err != nil {
		return nil, err
	}
	if err := priority.Validate(); err != nil {
		return nil, err
	}
	return &Event{
		ID:       getNextId(),
		Title:    title,
		StartAt:  startAt,
		Priority: priority,
		Reminder: nil}, nil
}

func getNextId() string {
	return uuid.New().String()
}

func (e *Event) Update(newTitle, newDate string, priority Priority) error {
	startAt, err := ValidateInput(newTitle, newDate)
	if err != nil {
		return err
	}

	if err := priority.Validate(); err != nil {
		return err
	}

	e.Title = newTitle
	e.StartAt = startAt
	e.Priority = priority
	return nil
}

func (e *Event) AddReminder(message string, at time.Time) error {
	if len(message) == 0 {
		return ErrEmptyReminder
	}
	e.Reminder = reminder.NewReminder(message, at)
	return nil
}

func (e *Event) StartReminder(Notify func(string)) error {
	return e.Reminder.Start(Notify)
}

func (e *Event) RemoveReminder() {
	e.Reminder = nil
}
