package tokenizers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
)

type JWT struct{}

type JWTToken struct {
	Bytes []byte
}

func (JWTToken) Children() []Token {
	return make([]Token, 0)
}

func (JWTToken) SetChildren([]Token) error {
	return nil
}

func (JWTToken) String() string {
	return "{JWT}"
}

func (JWT) Matchers() []Matcher {
	return []Matcher{
		func(in []byte) bool {
			return bytes.Count(in, []byte(".")) == 2
		},
		func(in []byte) bool {
			splitted := bytes.SplitN(in, []byte("."), -1)
			first, err := base64.StdEncoding.DecodeString(string(splitted[0]))
			if err != nil {
				return false
			}
			var tmp map[string]interface{}
			err = json.Unmarshal(first, &tmp)
			if err != nil {
				return false
			}
			typIf, ok := tmp["typ"]
			if !ok {
				return false
			}
			if typ, ok := typIf.(string); !ok || typ != "JWT" {
				return false
			}
			second, err := base64.StdEncoding.DecodeString(string(splitted[1]))
			if err != nil {
				return false
			}
			err = json.Unmarshal(second, &tmp)
			if err != nil {
				return false
			}
			return true
		},
	}
}

func (JWT) Token(in []byte) Token {
	return &JWTToken{in}
}
