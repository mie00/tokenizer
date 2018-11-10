package tokenizers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type JSONArray struct{}

type JSONArrayToken struct {
	Values []Token
}

func (e *JSONArrayToken) Children() []Token {
	return e.Values
}

func (e *JSONArrayToken) SetChildren(ts []Token) error {
	if len(ts) != len(e.Values) {
		return InvalidCount{}
	}
	e.Values = ts
	return nil
}

func (e *JSONArrayToken) String() string {
	res := []string{}
	for _, v := range e.Values {
		res = append(res, v.String())
	}
	return fmt.Sprintf(`[%s]`, strings.Join(res, ", "))
}

func (JSONArray) Matchers() []Matcher {
	return []Matcher{
		func(in []byte) bool {
			trimmed := bytes.TrimSpace(in)
			if len(trimmed) < 2 {
				return false
			}
			if trimmed[0] != '[' || trimmed[len(trimmed)-1] != ']' {
				return false
			}
			return true
		},
		func(in []byte) bool {
			var t []json.RawMessage
			err := json.Unmarshal(in, &t)
			return err == nil
		},
	}
}

func (JSONArray) Token(in []byte) Token {
	var t []json.RawMessage
	json.Unmarshal(in, &t)
	ret := JSONArrayToken{
		Values: make([]Token, 0, len(t)),
	}
	for _, v := range t {
		ret.Values = append(ret.Values, &UnknownToken{[]byte(v)})
	}
	return &ret
}
