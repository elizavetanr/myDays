package events

import "errors"

type Priority string

var (
	ErrInvalidPriority = errors.New("некорректный приоритет")
)

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

func (p Priority) Validate() error {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return nil
	default:
		return ErrInvalidPriority
	}
}
