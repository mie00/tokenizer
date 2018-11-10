package tokenizers

import (
	"strconv"
)

type Number struct{}

type NumberToken struct {
	Value float64
}

func (NumberToken) Children() []Token {
	return make([]Token, 0)
}

func (NumberToken) SetChildren([]Token) error {
	return nil
}

func (NumberToken) String() string {
	return "{Number}"
}

func (Number) Matchers() []Matcher {
	return []Matcher{
		func(in []byte) bool {
			_, err := strconv.ParseFloat(string(in), 64)
			return err == nil
		},
	}
}

func (Number) Token(in []byte) Token {
	v, _ := strconv.ParseFloat(string(in), 64)
	return &NumberToken{v}
}
