package tokenizers

import (
	"bytes"
	"time"

	"github.com/araddon/dateparse"
)

type DateTime struct{}

type DateTimeToken struct {
	Value  time.Time
	Format []byte
}

func (DateTimeToken) Children() []Token {
	return make([]Token, 0)
}

func (DateTimeToken) SetChildren([]Token) error {
	return nil
}

func (DateTimeToken) String() string {
	return "{DateTime}"
}

func (DateTime) Matchers() []Matcher {
	return []Matcher{
		func(in []byte) bool {
			if len(bytes.TrimSpace(in)) == 0 {
				return false
			}
			_, err := dateparse.ParseAny(string(in))
			return err == nil
		},
	}
}

func (DateTime) Token(in []byte) Token {
	value, _ := dateparse.ParseAny(string(in))
	// TODO: format
	return &DateTimeToken{value, []byte("")}
}
