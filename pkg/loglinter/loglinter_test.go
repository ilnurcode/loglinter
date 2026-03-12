package loglinter

import (
	"testing"
)

func TestAnalyzerRules(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		wantErr  bool
		errMatch string
	}{
		// Правило 1: начинается со строчной буквы
		{
			name:     "lowercase_start",
			message:  "starting server",
			wantErr:  false,
			errMatch: "",
		},
		{
			name:     "uppercase_start",
			message:  "Starting server",
			wantErr:  true,
			errMatch: "lowercase",
		},
		// Правило 2: только английский
		{
			name:     "english_only",
			message:  "server started",
			wantErr:  false,
			errMatch: "",
		},
		{
			name:     "russian_not_allowed",
			message:  "запуск сервера",
			wantErr:  true,
			errMatch: "English",
		},
		// Правило 3: спецсимволы и эмодзи
		{
			name:     "no_special_chars",
			message:  "server started",
			wantErr:  false,
			errMatch: "",
		},
		{
			name:     "emoji_not_allowed",
			message:  "server started 🚀",
			wantErr:  true,
			errMatch: "special characters",
		},
		{
			name:     "exclamation_marks_not_allowed",
			message:  "connection failed!!!",
			wantErr:  true,
			errMatch: "special characters",
		},
		{
			name:     "dots_not_allowed",
			message:  "something went wrong...",
			wantErr:  true,
			errMatch: "special characters",
		},
		// Правило 4: чувствительные данные
		{
			name:     "no_sensitive_data",
			message:  "user authenticated",
			wantErr:  false,
			errMatch: "",
		},
		{
			name:     "password_not_allowed",
			message:  "user password: secret",
			wantErr:  true,
			errMatch: "sensitive",
		},
		{
			name:     "token_not_allowed",
			message:  "token: abc123",
			wantErr:  true,
			errMatch: "sensitive",
		},
		{
			name:     "api_key_not_allowed",
			message:  "api_key=secret",
			wantErr:  true,
			errMatch: "sensitive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Тестовая проверка логики функций
			var hasError bool

			if !startsWithLowercase(tt.message) && tt.errMatch == "lowercase" {
				hasError = true
			}
			if !isEnglishOnly(tt.message) && tt.errMatch == "English" {
				hasError = true
			}
			if hasSpecialCharsOrEmoji(tt.message) && tt.errMatch == "special characters" {
				hasError = true
			}
			if containsSensitiveData(tt.message) && tt.errMatch == "sensitive" {
				hasError = true
			}

			if tt.wantErr && !hasError {
				t.Errorf("expected error for %q, but got none", tt.message)
			}
		})
	}
}
