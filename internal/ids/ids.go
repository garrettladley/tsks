package ids

import (
	"crypto/rand"
	"strings"
)

const (
	defaultAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	defaultSize     = 21
)

type Option func(*Options)

func WithAlphabet(alphabet string) Option {
	return func(o *Options) {
		o.Alphabet = alphabet
	}
}

func WithSize(size int) Option {
	return func(o *Options) {
		o.Size = size
	}
}

type Options struct {
	Alphabet string
	Size     int
}

var defaultOptions = Options{
	Alphabet: defaultAlphabet,
	Size:     defaultSize,
}

func New(prefix string, opts ...Option) string {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	return new(prefix, options)
}

func new(prefix string, options Options) string {
	if prefix != "" {
		return prefix + "_" + generateRandomString(options.Alphabet, options.Size)
	}
	return generateRandomString(options.Alphabet, options.Size)
}

func generateRandomString(alphabet string, size int) string {
	var sb strings.Builder
	sb.Grow(size)

	var (
		alphabetLen = len(alphabet)
		bytes       = make([]byte, size)
	)

	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}

	for i := range size {
		sb.WriteByte(alphabet[int(bytes[i])%alphabetLen])
	}

	return sb.String()
}
