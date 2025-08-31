package events

import (
	"errors"
	"github.com/araddon/dateparse"
	"regexp"
	"time"
)

var (
	ErrInvalidTitle = errors.New("некорректное имя события")
	ErrInvalidDate  = errors.New("некорректный формат даты")
)

func IsValidTitle(title string) bool {
	pattern := "^[а-яА-Я0-9a-zA-Z //./+]{3,50}$"
	matched, err := regexp.MatchString(pattern, title)
	if err != nil {
		return false
	}
	return matched
}

func ValidateInput(title, date string) (time.Time, error) {
	if !IsValidTitle(title) {
		return time.Time{}, ErrInvalidTitle
	}
	at, err := dateparse.ParseLocal(date)
	if err != nil {
		return time.Time{}, ErrInvalidDate
	}
	return at, err
}
