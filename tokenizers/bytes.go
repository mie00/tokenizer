package tokenizers

type Bytes struct{}

type BytesToken struct {
	Bytes []byte
}

func (BytesToken) Children() []Token {
	return make([]Token, 0)
}
func (BytesToken) SetChildren([]Token) error {
	return nil
}

func (e *BytesToken) String() string {
	return string(e.Bytes)
}

func (Bytes) Matchers() []Matcher {
	return []Matcher{
		func(in []byte) bool {
			return true
		},
	}
}

func (Bytes) Token(in []byte) Token {
	return &BytesToken{in}
}
