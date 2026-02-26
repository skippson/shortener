package generator

import (
	"crypto/rand"
	"math/big"
)

type Generator struct {
	alphabet string
	len      int
}

func NewGenerator(alphabet string, len int) *Generator {
	return &Generator{
		alphabet: alphabet,
		len:      len,
	}
}

func (g *Generator) Generate() (string, error) {
	b := make([]byte, g.len)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(g.alphabet))))
		if err != nil {
			return "", err
		}

		b[i] = g.alphabet[n.Int64()]
	}

	return string(b), nil
}
