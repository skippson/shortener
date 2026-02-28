package validator

import (
	"errors"
	netURL "net/url"
	"strings"
)

type Validator struct {
	letters map[rune]struct{}
	len     int
}

func NewValidator(alphabet string, size int) (*Validator, error) {
	if alphabet == "" || size <= 0 {
		return nil, errors.New("invalid options")
	}

	m := make(map[rune]struct{}, len(alphabet))

	for _, r := range alphabet {
		m[r] = struct{}{}
	}

	return &Validator{
		letters: m,
		len:     size,
	}, nil
}

func (v *Validator) ValidateURL(url string) (string, bool) {
	url = strings.TrimSpace(url)
	if url == "" {
		return "", false
	}

	parsed, err := netURL.ParseRequestURI(url)
	if err != nil {
		return "", false
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", false
	}

	normalized := strings.TrimRight(parsed.String(), "/")
	return normalized, true
}

func (v *Validator) ValidateShortened(shortened string) bool {
	if len(shortened) != v.len {
		return false
	}

	for _, r := range shortened {
		if _, ok := v.letters[r]; !ok {
			return false
		}
	}

	return true
}
