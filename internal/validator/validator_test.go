package validator_test

import (
	"testing"

	"shortener/internal/validator"

	"github.com/stretchr/testify/assert"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantURL   string
		wantValid bool
	}{
		{
			name:      "ok",
			url:       "http://example.com",
			wantURL:   "http://example.com",
			wantValid: true,
		},
		{
			name:      "slash",
			url:       "https://example.com/",
			wantURL:   "https://example.com",
			wantValid: true,
		},
		{
			name:      "spaces",
			url:       "  https://example.com  ",
			wantURL:   "https://example.com",
			wantValid: true,
		},
		{
			name:      "invalid url",
			url:       "invalid",
			wantURL:   "",
			wantValid: false,
		},
		{
			name:      "empty",
			url:       "   ",
			wantURL:   "",
			wantValid: false,
		},
		{
			name:      "invalid scheme",
			url:       "ftp://example.com",
			wantURL:   "",
			wantValid: false,
		},
	}

	validator, _ := validator.NewValidator("alphabet", 2)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURL, gotValid := validator.ValidateURL(tt.url)

			assert.Equal(t, tt.wantValid, gotValid)
			assert.Equal(t, tt.wantURL, gotURL)
		})
	}
}

func TestValidateShortened(t *testing.T) {
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	v, _ := validator.NewValidator(alphabet, 10)

	tests := []struct {
		name      string
		shortened string
		wantValid bool
	}{
		{
			name:      "ok",
			shortened: "abcABC123_",
			wantValid: true,
		},
		{
			name:      "short",
			shortened: "abc123",
			wantValid: false,
		},
		{
			name:      "long",
			shortened: "abcABC123__",
			wantValid: false,
		},
		{
			name:      "invalid character",
			shortened: "abcABC!23_",
			wantValid: false,
		},
		{
			name:      "empty",
			shortened: "",
			wantValid: false,
		},
		{
			name:      "space",
			shortened: "abc ABC123",
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValid := v.ValidateShortened(tt.shortened)
			assert.Equal(t, tt.wantValid, gotValid)
		})
	}
}
