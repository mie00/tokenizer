package tokenizer

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/d4l3k/messagediff"
	"github.com/davecgh/go-spew/spew"
	"github.com/mie00/tokenizer/tokenizers"
)

type InOutMessage struct {
	InOut   *InOut
	Message string
}

type InOut struct {
	in  string
	out tokenizers.Token
	str string
}

func intPointer(in int) *int {
	return &in
}

// func timeMust(t time.Time, err error) time.Time {
// 	if err != nil {
// 		panic(err)
// 	}
// 	return t
// }
func TestMatcher(t *testing.T) {
	ioms := []*InOutMessage{
		&InOutMessage{
			InOut: &InOut{
				"asdasdsad",
				nil,
				"asdasdsad",
			},
			Message: "failed on both",
		},
		&InOutMessage{
			InOut: &InOut{
				"asd.asd.sad",
				nil,
				"asd.asd.sad",
			},
			Message: "failed second only",
		},
		&InOutMessage{
			InOut: &InOut{
				"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJrZXkiOiJ2YWwiLCJpYXQiOjE0MjI2MDU0NDV9.eUiabuiKv-8PYk2AkGY4Fb5KMZeorYBLw261JPQD5lM.asd",
				nil,
				"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJrZXkiOiJ2YWwiLCJpYXQiOjE0MjI2MDU0NDV9.eUiabuiKv-8PYk2AkGY4Fb5KMZeorYBLw261JPQD5lM.asd",
			},
			Message: "failed on first only",
		},
		&InOutMessage{
			InOut: &InOut{
				"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJrZXkiOiJ2YWwiLCJpYXQiOjE0MjI2MDU0NDV9.eUiabuiKv-8PYk2AkGY4Fb5KMZeorYBLw261JPQD5lM",
				nil,
				"{JWT}",
			},
			Message: "succeed",
		},
	}
	for _, iom := range ioms {
		inout := iom.InOut
		out := Tokenize([]byte(inout.in))
		if out.(tokenizers.Token).String() != inout.str {
			t.Errorf("case: %s should have %s, out: %s, was expecting: %s", inout.in, iom.Message, out.(tokenizers.Token).String(), inout.str)
		}
	}
}
func TestTokenize(t *testing.T) {
	inouts := []*InOut{
		&InOut{
			"CA761232-ED42-11CE-BACD-00AA0057B223",
			&tokenizers.UUIDToken{
				UUID: []byte("ca761232-ed42-11ce-bacd-00aa0057b223"),
				Orig: []byte("CA761232-ED42-11CE-BACD-00AA0057B223"),
			},
			"{UUID}",
		},
		&InOut{
			"mie",
			&tokenizers.BytesToken{
				Bytes: []byte("mie"),
			},
			"mie",
		},
		&InOut{
			"mohamed@elawadi.net",
			&tokenizers.EmailToken{
				Username: &tokenizers.BytesToken{[]byte("mohamed")},
				Domain:   &tokenizers.DomainToken{[]byte("elawadi.net")},
			},
			"{EMAIL}",
		},
		&InOut{
			"https://github.com:443/mie00?something=true",
			&tokenizers.URLToken{
				Protocol:          []byte("https"),
				Host:              &tokenizers.DomainToken{[]byte("github.com")},
				Port:              intPointer(443),
				Fragment:          []tokenizers.Token{&tokenizers.BytesToken{[]byte("mie00")}},
				QueryStringKeys:   [][]byte{[]byte("something")},
				QueryStringValues: []tokenizers.Token{&tokenizers.BooleanToken{true}},
			},
			"https://github.com:443/mie00?something={Boolean}",
		},
		&InOut{
			"https://github.com/mie00?something=true",
			&tokenizers.URLToken{
				Protocol:          []byte("https"),
				Host:              &tokenizers.DomainToken{[]byte("github.com")},
				Port:              nil,
				Fragment:          []tokenizers.Token{&tokenizers.BytesToken{[]byte("mie00")}},
				QueryStringKeys:   [][]byte{[]byte("something")},
				QueryStringValues: []tokenizers.Token{&tokenizers.BooleanToken{true}},
			},
			"https://github.com/mie00?something={Boolean}",
		},
		&InOut{
			"https://github.com/mie00/pulls?something=true",
			&tokenizers.URLToken{
				Protocol:          []byte("https"),
				Host:              &tokenizers.DomainToken{[]byte("github.com")},
				Port:              nil,
				Fragment:          []tokenizers.Token{&tokenizers.BytesToken{[]byte("mie00")}, &tokenizers.BytesToken{[]byte("pulls")}},
				QueryStringKeys:   [][]byte{[]byte("something")},
				QueryStringValues: []tokenizers.Token{&tokenizers.BooleanToken{true}},
			},
			"https://github.com/mie00/pulls?something={Boolean}",
		},
		&InOut{
			"https://github.com/mie00/pulls",
			&tokenizers.URLToken{
				Protocol:          []byte("https"),
				Host:              &tokenizers.DomainToken{[]byte("github.com")},
				Port:              nil,
				Fragment:          []tokenizers.Token{&tokenizers.BytesToken{[]byte("mie00")}, &tokenizers.BytesToken{[]byte("pulls")}},
				QueryStringKeys:   [][]byte{},
				QueryStringValues: []tokenizers.Token{},
			},
			"https://github.com/mie00/pulls",
		},
		&InOut{
			"https://github.com/",
			&tokenizers.URLToken{
				Protocol:          []byte("https"),
				Host:              &tokenizers.DomainToken{[]byte("github.com")},
				Port:              nil,
				Fragment:          []tokenizers.Token{&tokenizers.BytesToken{[]byte("")}},
				QueryStringKeys:   [][]byte{},
				QueryStringValues: []tokenizers.Token{},
			},
			"https://github.com/",
		},
		&InOut{
			"https://github.com/?something=true",
			&tokenizers.URLToken{
				Protocol:          []byte("https"),
				Host:              &tokenizers.DomainToken{[]byte("github.com")},
				Port:              nil,
				Fragment:          []tokenizers.Token{&tokenizers.BytesToken{[]byte("")}},
				QueryStringKeys:   [][]byte{[]byte("something")},
				QueryStringValues: []tokenizers.Token{&tokenizers.BooleanToken{true}},
			},
			"https://github.com/?something={Boolean}",
		},
		&InOut{
			"https://github.com?something=true",
			&tokenizers.URLToken{
				Protocol:          []byte("https"),
				Host:              &tokenizers.DomainToken{[]byte("github.com")},
				Port:              nil,
				Fragment:          []tokenizers.Token{},
				QueryStringKeys:   [][]byte{[]byte("something")},
				QueryStringValues: []tokenizers.Token{&tokenizers.BooleanToken{true}},
			},
			"https://github.com?something={Boolean}",
		},
		&InOut{
			"http://192.168.1.11:8887/2.0/discovery/delivery/search?orderingEnabled=0&compact=0&area=37493b8a-2224-11e8-924e-0242ac110011&page=1&pageSize=20",
			nil,
			"http://192.168.1.11:8887/{Number}/discovery/delivery/search?orderingEnabled={Number}&compact={Number}&area={UUID}&page={Number}&pageSize={Number}",
		},
		&InOut{
			"true",
			&tokenizers.BooleanToken{true},
			"{Boolean}",
		},
		&InOut{
			"True",
			&tokenizers.BooleanToken{true},
			"{Boolean}",
		},
		&InOut{
			"false",
			&tokenizers.BooleanToken{false},
			"{Boolean}",
		},
		&InOut{
			"POST",
			&tokenizers.HTTPMethodToken{[]byte("POST")},
			"POST",
		},
		&InOut{
			"GET",
			&tokenizers.HTTPMethodToken{[]byte("GET")},
			"GET",
		},
		&InOut{
			"\"true\"",
			&tokenizers.QuotedStringToken{&tokenizers.BooleanToken{true}, '"'},
			"\"{Boolean}\"",
		},
		&InOut{
			`{"mie": "true"}`,
			&tokenizers.JSONObjectToken{[]tokenizers.Token{&tokenizers.QuotedStringToken{&tokenizers.BooleanToken{true}, '"'}}, []string{"mie"}},
			`{"mie":"{Boolean}"}`,
		},
		&InOut{
			`{"mie": {"asd": "true"}, "b": 1}`,
			&tokenizers.JSONObjectToken{
				[]tokenizers.Token{
					&tokenizers.JSONObjectToken{
						[]tokenizers.Token{
							&tokenizers.QuotedStringToken{&tokenizers.BooleanToken{true}, '"'},
						},
						[]string{"asd"},
					},
					&tokenizers.NumberToken{1},
				},
				[]string{"mie", "b"},
			},
			`{"mie":{"asd":"{Boolean}"}, "b":{Number}}`,
		},
		&InOut{
			`[{"mie": {"asd": "true"}, "b": 1}]`,
			&tokenizers.JSONArrayToken{
				[]tokenizers.Token{
					&tokenizers.JSONObjectToken{
						[]tokenizers.Token{
							&tokenizers.JSONObjectToken{
								[]tokenizers.Token{
									&tokenizers.QuotedStringToken{&tokenizers.BooleanToken{true}, '"'},
								},
								[]string{"asd"},
							},
							&tokenizers.NumberToken{1},
						},
						[]string{"mie", "b"},
					},
				},
			},
			`[{"mie":{"asd":"{Boolean}"}, "b":{Number}}]`,
		},
		&InOut{
			"null",
			&tokenizers.NullToken{},
			"{Null}",
		},
		// &InOut{
		// 	"1-1-2005",
		// 	&tokenizers.DateTimeToken{Value: timeMust(time.Parse("02-01-2006", "01-01-2005")), Format: []byte("")},
		// 	"{DateTime}",
		// },
		&InOut{
			"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJrZXkiOiJ2YWwiLCJpYXQiOjE0MjI2MDU0NDV9.eUiabuiKv-8PYk2AkGY4Fb5KMZeorYBLw261JPQD5lM",
			nil,
			"{JWT}",
		},
	}
	for _, inout := range inouts {
		out := Tokenize([]byte(inout.in))
		if inout.out != nil && !reflect.DeepEqual(inout.out, out) {
			t.Errorf("error in case: in: %s, out: %#+v, was expecting: %#+v", inout.in, out, inout.out)
			spew.Dump("got", out, "wants", inout.out)
			diff, _ := messagediff.PrettyDiff(inout.out, out)
			fmt.Println(diff)
		} else {
			if out.(tokenizers.Token).String() != inout.str {
				t.Errorf("error in case: in: %s, out: %s, was expecting: %s", inout.in, out.(tokenizers.Token).String(), inout.str)
			}
		}
	}
}
