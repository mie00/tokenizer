package tokenizers

import (
	"bytes"
	"regexp"
)

type Email struct{}

type EmailToken struct {
	Username Token
	Domain   Token
}

func (e *EmailToken) Children() []Token {
	return []Token{
		e.Username,
		e.Domain,
	}
}

func (e *EmailToken) SetChildren(ts []Token) error {
	if len(ts) != 2 {
		return InvalidCount{}
	}
	e.Username = ts[0]
	e.Domain = ts[1]
	return nil
}

func (e *EmailToken) String() string {
	return "{EMAIL}"
}

var validEmail = regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`)

func (Email) Matchers() []Matcher {
	return []Matcher{
		func(in []byte) bool {
			if bytes.Count(in, []byte("@")) != 1 {
				return false
			}
			if !validEmail.Match(in) {
				return false
			}
			return true
		},
	}
}

func (Email) Token(in []byte) Token {
	i := bytes.Index(in, []byte("@"))
	return &EmailToken{
		Username: &UnknownToken{in[:i]},
		Domain:   &UnknownToken{in[i+1:]},
	}
}
