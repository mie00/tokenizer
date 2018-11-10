package tokenizers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type JSONObject struct{}

type JSONObjectToken struct {
	Values []Token
	Keys   []string
}

func (e *JSONObjectToken) Children() []Token {
	return e.Values
}

func (e *JSONObjectToken) SetChildren(ts []Token) error {
	if len(ts) != len(e.Values) {
		return InvalidCount{}
	}
	e.Values = ts
	return nil
}

func (e *JSONObjectToken) String() string {
	res := []string{}
	for i, k := range e.Keys {
		// TODO: escape
		res = append(res, fmt.Sprintf(`"%s":%s`, k, e.Values[i].String()))
	}
	return fmt.Sprintf(`{%s}`, strings.Join(res, ", "))
}

func (JSONObject) Matchers() []Matcher {
	return []Matcher{
		func(in []byte) bool {
			trimmed := bytes.TrimSpace(in)
			if len(trimmed) < 2 {
				return false
			}
			if trimmed[0] != '{' || trimmed[len(trimmed)-1] != '}' {
				return false
			}
			return true
		},
		func(in []byte) bool {
			var t map[string]json.RawMessage
			err := json.Unmarshal(in, &t)
			return err == nil
		},
	}
}

func (JSONObject) Token(in []byte) Token {
	var t map[string]json.RawMessage
	json.Unmarshal(in, &t)
	ret := JSONObjectToken{
		Values: make([]Token, 0, len(t)),
		Keys:   make([]string, 0, len(t)),
	}
	for k, v := range t {
		ret.Values = append(ret.Values, &UnknownToken{[]byte(v)})
		ret.Keys = append(ret.Keys, k)
	}
	return &ret
}
