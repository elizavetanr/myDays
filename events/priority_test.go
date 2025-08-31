package events

import (
	"testing"
)

func TestValidatePriority(t *testing.T) {
	var p Priority = "very low"
	err := p.Validate()
	if err == nil {
		t.Error("Expected an error for priority, got none")
	}
	p = "low"
	err = p.Validate()
	if err != nil {
		t.Errorf("Expected no error for correct priority, got %v\"", err)
	}
}
