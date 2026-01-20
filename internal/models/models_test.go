package models

import (
	"testing"
)

func TestCommitmentType(t *testing.T) {
	tests := []struct {
		name     string
		input    CommitmentType
		expected string
	}{
		{"Decimal type", CommitmentTypeDecimal, "decimal"},
		{"Percentage type", CommitmentTypePercentage, "percentage"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.input) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.input))
			}
		})
	}
}

func TestUserRole(t *testing.T) {
	tests := []struct {
		name     string
		input    UserRole
		expected string
	}{
		{"User role", RoleUser, "user"},
		{"Admin role", RoleAdmin, "admin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.input) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.input))
			}
		})
	}
}
