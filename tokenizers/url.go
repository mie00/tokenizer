package tokenizers

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type URL struct{}

type URLToken struct {
	Protocol          []byte
	Host              Token
	Port              *int
	Fragment          []Token
	QueryStringKeys   [][]byte
	QueryStringValues []Token
}

func (e *URLToken) Children() []Token {
	ret := []Token{e.Host}
	ret = append(ret, e.Fragment...)
	ret = append(ret, e.QueryStringValues...)
	return ret
}

func (e *URLToken) SetChildren(ts []Token) error {
	if len(ts) != 1+len(e.Fragment)+len(e.QueryStringValues) {
		return InvalidCount{}
	}
	e.Host = ts[0]
	e.Fragment = ts[1 : len(e.Fragment)+1]
	e.QueryStringValues = ts[len(e.Fragment)+1 : len(ts)]
	return nil
}

func (e *URLToken) String() string {
	port := ""
	if e.Port != nil {
		port = fmt.Sprintf(":%d", *e.Port)
	}
	qs := make([]string, len(e.QueryStringKeys))
	for i, k := range e.QueryStringKeys {
		qs[i] = fmt.Sprintf("%s=%s", k, e.QueryStringValues[i].String())
	}
	qss := ""
	if len(qs) > 0 {
		qss = fmt.Sprintf("?%s", strings.Join(qs, "&"))
	}
	fragment := make([]string, len(e.Fragment))
	for i, k := range e.Fragment {
		fragment[i] = fmt.Sprintf("%s", k)
	}
	fragments := ""
	if len(e.Fragment) > 0 {
		fragments = fmt.Sprintf("/%s", strings.Join(fragment, "/"))
	}
	return fmt.Sprintf("%s://%s%s%s%s", e.Protocol, e.Host, port, fragments, qss)
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
		queryString = s[1]
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
	fragments := []Token{}
	for _, f := range bytes.Split(fragment, []byte("/"))[1:] {
		fragments = append(fragments, &UnknownToken{f})
	}
	queryStringKeys := [][]byte{}
	queryStringValues := []Token{}
	for _, f := range bytes.Split(queryString, []byte("&")) {
		if kv := bytes.SplitN(f, []byte("="), 2); len(kv) == 2 {
			queryStringKeys = append(queryStringKeys, kv[0])
			queryStringValues = append(queryStringValues, &UnknownToken{kv[1]})
		}
	}
	return &URLToken{protocol, &UnknownToken{host}, port, fragments, queryStringKeys, queryStringValues}
}
