package events

import (
	"testing"
)

func TestValidation(t *testing.T) {
	_, err := ValidateInput("Поступление заказа в пункт ozon кошачий корм 2шт", "2025-08-15")
	if err != nil {
		t.Errorf("Expected no error for correct title and date, got %v\"", err)
	}
	_, err = ValidateInput("& Поступление заказа в пункт ozon кошачий корм 2шт", "2025-08-15")
	if err == nil {
		t.Error("Expected an error for char '&' in title, got none")
	}
	_, err = ValidateInput("Поступление заказа в пункт ozon кошачий корм 2шт", "2025-23-01")
	if err == nil {
		t.Error("Expected an error for date, got none")
	}
}
