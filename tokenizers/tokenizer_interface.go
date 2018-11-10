package tokenizers

type Token interface {
	Children() []Token
	SetChildren([]Token) error
	String() string
}

type Matcher func([]byte) bool

type Tokenizer interface {
	Matchers() []Matcher
	Token([]byte) Token
}

type Tokenizers []Tokenizer
