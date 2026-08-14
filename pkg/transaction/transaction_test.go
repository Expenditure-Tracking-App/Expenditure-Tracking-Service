package transaction

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateAmount(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        float32
		wantErr     bool
		errContains string
	}{
		{"given integer input when parsed then returns float", "100", 100.0, false, ""},
		{"given decimal input when parsed then returns float", "99.99", 99.99, false, ""},
		{"given negative input when parsed then returns negative float", "-50", -50.0, false, ""},
		{"given string input when parsed then returns error", "abc", 0, true, "invalid amount"},
		{"given empty input when parsed then returns error", "", 0, true, "invalid amount"},
		{"given zero input when parsed then returns zero", "0", 0.0, false, ""},
		{"given large number when parsed then returns float", "999999.99", 999999.99, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateAmount(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestValidateBool(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        bool
		wantErr     bool
		errContains string
	}{
		{"given true string when parsed then returns true", "true", true, false, ""},
		{"given false string when parsed then returns false", "false", false, false, ""},
		{"given yes string when parsed then returns error", "yes", false, true, "invalid boolean value"},
		{"given no string when parsed then returns error", "no", false, true, "invalid boolean value"},
		{"given invalid string when parsed then returns error", "maybe", false, true, "invalid boolean value"},
		{"given empty string when parsed then returns error", "", false, true, "invalid boolean value"},
		{"given numeric one when parsed then returns true", "1", true, false, ""},
		{"given numeric zero when parsed then returns false", "0", false, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateBool(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestProcessDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		validate func(t *testing.T, result string)
	}{
		{
			name:  "given t shorthand when processed then returns today",
			input: "t",
			validate: func(t *testing.T, result string) {
				today := time.Now().Format("2006-01-02")
				assert.Equal(t, today, result)
			},
		},
		{
			name:     "given valid DDMMYY when processed then returns ISO format",
			input:    "15.04.25",
			wantErr:  false,
			validate: func(t *testing.T, result string) { assert.Equal(t, "2025-04-15", result) },
		},
		{
			name:     "given another valid date when processed then returns ISO format",
			input:    "01.01.24",
			wantErr:  false,
			validate: func(t *testing.T, result string) { assert.Equal(t, "2024-01-01", result) },
		},
		{
			name:  "given invalid format when processed then returns today",
			input: "not-a-date",
			validate: func(t *testing.T, result string) {
				today := time.Now().Format("2006-01-02")
				assert.Equal(t, today, result)
			},
		},
		{
			name:  "given empty string when processed then returns today",
			input: "",
			validate: func(t *testing.T, result string) {
				today := time.Now().Format("2006-01-02")
				assert.Equal(t, today, result)
			},
		},
		{
			name:     "given leap year date when processed then returns ISO format",
			input:    "29.02.24",
			wantErr:  false,
			validate: func(t *testing.T, result string) { assert.Equal(t, "2024-02-29", result) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessDate(tt.input)
			tt.validate(t, result)
		})
	}
}
