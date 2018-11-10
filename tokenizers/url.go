package tokenizers

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
)

type URL struct{}

type URLToken struct {
	Protocol    []byte
	Host        Token
	Port        *int
	Fragment    Token
	QueryString Token
}

func (e *URLToken) Children() []Token {
	return []Token{e.Host, e.Fragment, e.QueryString}
}

func (e *URLToken) SetChildren(ts []Token) error {
	if len(ts) != 3 {
		return InvalidCount{}
	}
	e.Host = ts[0]
	e.Fragment = ts[1]
	e.QueryString = ts[2]
	return nil
}

func (e *URLToken) String() string {
	port := ""
	if e.Port != nil {
		port = fmt.Sprintf(":%d", *e.Port)
	}
	return fmt.Sprintf("%s://%s%s%s%s", e.Protocol, e.Host, port, e.Fragment, e.QueryString)
}

var validURL = regexp.MustCompile(`^(.+)://([^/:]+)((?:[:][0-9]+)?)(/?[^?]*)((?:\?.*)?)$`)

func (URL) Matchers() []Matcher {
	return []Matcher{
		func(in []byte) bool {
			if !validURL.Match(in) {
				return false
			}
			return true
		},
	}
}

func (URL) Token(in []byte) Token {
	var (
		protocol    []byte
		host        []byte
		port        *int
		fragment    []byte
		queryString []byte
	)
	s := bytes.SplitN(in, []byte("://"), 2)
	protocol = s[0]
	s = bytes.SplitN(s[1], []byte("?"), 2)
	if len(s) == 2 {
		queryString = []byte("?" + string(s[1]))
	}
	s = bytes.SplitN(s[0], []byte("/"), 2)
	hostPort := bytes.SplitN(s[0], []byte(":"), 2)
	host = hostPort[0]
	if len(hostPort) == 2 {
		if i, err := strconv.Atoi(string(hostPort[1])); err == nil {
			port = &i
		}
	}
	if len(s) == 2 {
		fragment = []byte("/" + string(s[1]))
	}
	return &URLToken{protocol, &UnknownToken{host}, port, &UnknownToken{fragment}, &UnknownToken{queryString}}
}
