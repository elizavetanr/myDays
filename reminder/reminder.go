package reminder

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrTimeReminderIsUp = errors.New("время напоминания вышло")
	ErrNotExistReminder = errors.New("напоминания не существует")
)

type Reminder struct {
	Message string
	At      time.Time
	Sent    bool
	timer   *time.Timer
}

func NewReminder(message string, at time.Time) *Reminder {
	return &Reminder{
		Message: message,
		At:      at,
		Sent:    false,
	}
}

func (r *Reminder) Send(Notify func(string)) {
	if r.Sent {
		return
	}
	Notify(r.Message)
	r.Sent = true
}

func (r *Reminder) Start(Notify func(string)) error {
	duration := r.At.Sub(time.Now())
	if duration < 0 {
		return fmt.Errorf("невозможно запустить напоминание: %w", ErrTimeReminderIsUp)
	} else {
		r.timer = time.AfterFunc(duration, func() { r.Send(Notify) })
		return nil
	}
}

func (r *Reminder) Stop() error {
	if !r.isExist() {
		return fmt.Errorf("невозможно остановить напоминание: %w", ErrNotExistReminder)
	}
	r.timer.Stop()
	return nil
}

func (r *Reminder) isExist() bool {
	return r != nil
}
