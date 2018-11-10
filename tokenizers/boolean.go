package tokenizers

import (
	"bytes"
)

type Boolean struct{}

type BooleanToken struct {
	Value bool
}

func (BooleanToken) Children() []Token {
	return make([]Token, 0)
}

func (BooleanToken) SetChildren([]Token) error {
	return nil
}

func (e *BooleanToken) String() string {
	return "{Boolean}"
}

func (Boolean) Matchers() []Matcher {
	return []Matcher{
		func(in []byte) bool {
			return bytes.Compare(bytes.ToLower(in), []byte("true")) == 0 || bytes.Compare(bytes.ToLower(in), []byte("false")) == 0
		},
	}
}

func (Boolean) Token(in []byte) Token {
	return &BooleanToken{
		Value: bytes.Compare(bytes.ToLower(in), []byte("true")) == 0,
	}
}
