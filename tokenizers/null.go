package tokenizers

import (
	"bytes"
)

type Null struct{}

type NullToken struct{}

func (NullToken) Children() []Token {
	return make([]Token, 0)
}

func (NullToken) SetChildren([]Token) error {
	return nil
}

func (e *NullToken) String() string {
	return "{Null}"
}

func (Null) Matchers() []Matcher {
	return []Matcher{
		func(in []byte) bool {
			return bytes.Compare(bytes.ToLower(in), []byte("nil")) == 0 || bytes.Compare(bytes.ToLower(in), []byte("null")) == 0
		},
	}
}

func (Null) Token(in []byte) Token {
	return &NullToken{}
}
