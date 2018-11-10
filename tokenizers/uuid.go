package tokenizers

import (
	"bytes"
	"regexp"
)

type UUID struct{}

type UUIDToken struct {
	UUID []byte
	Orig []byte
}

func (UUIDToken) Children() []Token {
	return make([]Token, 0)
}

func (UUIDToken) SetChildren([]Token) error {
	return nil
}

func (e *UUIDToken) String() string {
	return "{UUID}"
}

var validUUID = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)

func (UUID) Matchers() []Matcher {
	return []Matcher{
		func(in []byte) bool {
			if len(in) != 36 {
				return false
			}
			if !validUUID.Match(in) {
				return false
			}
			return true
		},
	}
}

func (UUID) Token(in []byte) Token {
	return &UUIDToken{
		UUID: bytes.ToLower(in),
		Orig: in,
	}
}
