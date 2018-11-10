package tokenizers

import (
	"regexp"
)

type Domain struct{}

type DomainToken struct {
	Domain []byte
}

func (e *DomainToken) Children() []Token {
	return []Token{}
}

func (e *DomainToken) SetChildren(ts []Token) error {
	return nil
}

func (e *DomainToken) String() string {
	return "{DOMAIN}"
}

var validDomain = regexp.MustCompile(`^[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`)

func (Domain) Matchers() []Matcher {
	return []Matcher{
		func(in []byte) bool {
			if !validDomain.Match(in) {
				return false
			}
			return true
		},
	}
}

func (Domain) Token(in []byte) Token {
	return &DomainToken{in}
}
