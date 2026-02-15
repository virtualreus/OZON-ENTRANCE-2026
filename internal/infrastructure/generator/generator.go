package generator

import "math/rand"

type ShortGenerator struct{}

func NewShortGenerator() *ShortGenerator {
	return &ShortGenerator{}
}

func (sg *ShortGenerator) GenerateShortLink() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	b := make([]byte, 10)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
