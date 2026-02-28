package generator_test

import (
	"crypto/rand"
	"errors"
	"shortener/internal/generator"
	"shortener/internal/validator"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratorGenerate(t *testing.T) {
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	size := 10
	gen := generator.NewGenerator(alphabet, size)

	validator, _ := validator.NewValidator(alphabet, size)

	t.Run("classic", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			s, err := gen.Generate()
			assert.NoError(t, err)
			assert.Len(t, s, size)
			assert.True(t, validator.ValidateShortened(s))
		}
	})
}

type failReader struct{}

func (f *failReader) Read(p []byte) (int, error) {
	return 0, errors.New("read error")
}

func TestGeneratorRandError(t *testing.T) {
	gen := generator.NewGenerator("abc", 5)

	origReader := rand.Reader
	defer func() { rand.Reader = origReader }()

	rand.Reader = &failReader{}

	s, err := gen.Generate()
	assert.Error(t, err)
	assert.EqualError(t, err, "read error")
	assert.Empty(t, s)
}
