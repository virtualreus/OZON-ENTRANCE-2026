package generator

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	charset     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	ShortLength = 10
)

var mx = big.NewInt(int64(len(charset)))

type ShortGenerator struct{}

func NewShortGenerator() *ShortGenerator {
	return &ShortGenerator{}
}

func (sg *ShortGenerator) GenerateShortLink() (string, error) {
	b := make([]byte, ShortLength)
	for i := range b {
		n, err := rand.Int(rand.Reader, mx)
		if err != nil {
			return "", fmt.Errorf("failed to generate random char: %w", err)
		}
		b[i] = charset[n.Int64()]
	}
	return string(b), nil
}
