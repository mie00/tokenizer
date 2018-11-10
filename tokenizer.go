package tokenizer

import (
	"github.com/mie00/tokenizer/tokenizers"
)

var (
	ts = tokenizers.Tokenizers{
		tokenizers.UUID{},
		tokenizers.Email{},
		tokenizers.Boolean{},
		tokenizers.Method{},
		tokenizers.QuotedString{},
		tokenizers.Number{},
		tokenizers.Domain{},
		tokenizers.URL{},
		tokenizers.Null{},
		tokenizers.DateTime{},
		tokenizers.JSONObject{},
		tokenizers.JSONArray{},
		tokenizers.Bytes{},
	}
)

func TokenizeString(in string) string {
	return Tokenize([]byte(in)).(tokenizers.Token).String()
}

func Tokenize(in []byte) interface{} {
	return tokenize(&tokenizers.UnknownToken{in})
}

func tokenize(t tokenizers.Token) tokenizers.Token {
	unknown, ok := t.(*tokenizers.UnknownToken)
	if !ok {
		return t
	}
	in := unknown.Bytes
	for _, t := range ts {
		matchers := t.Matchers()
		for _, m := range matchers {
			if m(in) {
				token := t.Token(in)
				children := token.Children()
				for i := range children {
					children[i] = tokenize(children[i])
				}
				token.SetChildren(children)
				return token
			}
		}
	}
	return &tokenizers.BytesToken{in}
}
