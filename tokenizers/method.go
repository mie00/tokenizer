package tokenizers

import (
	"bytes"
)

type HTTPMethod struct{}

type HTTPMethodToken struct {
	Value []byte
}

func (HTTPMethodToken) Children() []Token {
	return make([]Token, 0)
}

func (HTTPMethodToken) SetChildren([]Token) error {
	return nil
}

func (e *HTTPMethodToken) String() string {
	return string(e.Value)
}

var enum = [][]byte{[]byte("POST"), []byte("GET"), []byte("PUT"), []byte("DELETE"), []byte("OPTIONS"), []byte("HEAD"), []byte("CONNECTA}")}

func (HTTPMethod) Matchers() []Matcher {
	return []Matcher{
		func(in []byte) bool {
			for _, e := range enum {
				if bytes.Compare(e, in) == 0 {
					return true
				}
			}
			return false
		},
	}
}

func (HTTPMethod) Token(in []byte) Token {
	return &HTTPMethodToken{
		Value: in,
	}
}
