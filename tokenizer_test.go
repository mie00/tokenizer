package tokenizer

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/d4l3k/messagediff"
	"github.com/davecgh/go-spew/spew"
	"github.com/mie00/tokenizer/tokenizers"
)

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
				Protocol:    []byte("https"),
				Host:        &tokenizers.DomainToken{[]byte("github.com")},
				Port:        intPointer(443),
				Fragment:    &tokenizers.BytesToken{[]byte("/mie00")},
				QueryString: &tokenizers.BytesToken{[]byte("?something=true")},
			},
			"https://{DOMAIN}:443/mie00?something=true",
		},
		&InOut{
			"https://github.com/mie00?something=true",
			&tokenizers.URLToken{
				Protocol:    []byte("https"),
				Host:        &tokenizers.DomainToken{[]byte("github.com")},
				Port:        nil,
				Fragment:    &tokenizers.BytesToken{[]byte("/mie00")},
				QueryString: &tokenizers.BytesToken{[]byte("?something=true")},
			},
			"https://{DOMAIN}/mie00?something=true",
		},
		&InOut{
			"https://github.com/mie00/pulls?something=true",
			&tokenizers.URLToken{
				Protocol:    []byte("https"),
				Host:        &tokenizers.DomainToken{[]byte("github.com")},
				Port:        nil,
				Fragment:    &tokenizers.BytesToken{[]byte("/mie00/pulls")},
				QueryString: &tokenizers.BytesToken{[]byte("?something=true")},
			},
			"https://{DOMAIN}/mie00/pulls?something=true",
		},
		&InOut{
			"https://github.com/mie00/pulls",
			&tokenizers.URLToken{
				Protocol:    []byte("https"),
				Host:        &tokenizers.DomainToken{[]byte("github.com")},
				Port:        nil,
				Fragment:    &tokenizers.BytesToken{[]byte("/mie00/pulls")},
				QueryString: &tokenizers.BytesToken{nil},
			},
			"https://{DOMAIN}/mie00/pulls",
		},
		&InOut{
			"https://github.com/",
			&tokenizers.URLToken{
				Protocol:    []byte("https"),
				Host:        &tokenizers.DomainToken{[]byte("github.com")},
				Port:        nil,
				Fragment:    &tokenizers.BytesToken{[]byte("/")},
				QueryString: &tokenizers.BytesToken{nil},
			},
			"https://{DOMAIN}/",
		},
		&InOut{
			"https://github.com/?something=true",
			&tokenizers.URLToken{
				Protocol:    []byte("https"),
				Host:        &tokenizers.DomainToken{[]byte("github.com")},
				Port:        nil,
				Fragment:    &tokenizers.BytesToken{[]byte("/")},
				QueryString: &tokenizers.BytesToken{[]byte("?something=true")},
			},
			"https://{DOMAIN}/?something=true",
		},
		&InOut{
			"https://github.com?something=true",
			&tokenizers.URLToken{
				Protocol:    []byte("https"),
				Host:        &tokenizers.DomainToken{[]byte("github.com")},
				Port:        nil,
				Fragment:    &tokenizers.BytesToken{nil},
				QueryString: &tokenizers.BytesToken{[]byte("?something=true")},
			},
			"https://{DOMAIN}?something=true",
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
			&tokenizers.MethodToken{[]byte("POST")},
			"POST",
		},
		&InOut{
			"GET",
			&tokenizers.MethodToken{[]byte("GET")},
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
	}
	for _, inout := range inouts {
		out := Tokenize([]byte(inout.in))
		if !reflect.DeepEqual(inout.out, out) {
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
