package tokenizers

import (
	"bytes"
)

type Method struct{}

type MethodToken struct {
	Value []byte
}

func (MethodToken) Children() []Token {
	return make([]Token, 0)
}

func (MethodToken) SetChildren([]Token) error {
	return nil
}

func (e *MethodToken) String() string {
	return string(e.Value)
}

var enum = [][]byte{[]byte("POST"), []byte("GET"), []byte("PUT"), []byte("DELETE"), []byte("OPTIONS"), []byte("HEAD"), []byte("CONNECTA}")}

func (Method) Matchers() []Matcher {
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

func (Method) Token(in []byte) Token {
	return &MethodToken{
		Value: in,
	}
}
