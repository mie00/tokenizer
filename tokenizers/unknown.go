package tokenizers

type Unknown struct{}

type UnknownToken struct {
	Bytes []byte
}

func (UnknownToken) Children() []Token {
	return make([]Token, 0)
}

func (UnknownToken) SetChildren([]Token) error {
	return nil
}

func (UnknownToken) String() string {
	return "{UNKNOWN}"
}

func (Unknown) Matchers() []Matcher {
	return []Matcher{
		func(in []byte) bool {
			return false
		},
	}
}

func (Unknown) Token(in []byte) Token {
	return &UnknownToken{in}
}
