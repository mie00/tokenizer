package tokenizers

import (
	"bytes"
)

type QuotedString struct{}

type QuotedStringToken struct {
	Value Token
	Quote byte
}

func (e *QuotedStringToken) Children() []Token {
	return []Token{e.Value}
}
func (e *QuotedStringToken) SetChildren(ts []Token) error {
	if len(ts) != 1 {
		return InvalidCount{}
	}
	e.Value = ts[0]
	return nil
}
func (e *QuotedStringToken) String() string {
	inner := e.Value.String()
	ret := make([]byte, 0, len(inner)+2)
	ret = append(ret, e.Quote)
	ret = append(ret, inner...)
	ret = append(ret, e.Quote)
	return string(ret)
}

func (QuotedString) Matchers() []Matcher {
	return []Matcher{
		func(in []byte) bool {
			if len(in) < 2 {
				return false
			}
			if in[0] != in[len(in)-1] || (in[0] != '"' && in[0] != '\'') {
				return false
			}
			// TODO: check inner escapes
			return true
		},
	}
}

func (QuotedString) Token(in []byte) Token {
	return &QuotedStringToken{
		Value: &UnknownToken{bytes.Replace(in[1:len(in)-1], []byte(`\\`), []byte(`\`), -1)},
		Quote: in[0],
	}
}
