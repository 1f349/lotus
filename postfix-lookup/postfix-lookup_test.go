package postfix_lookup

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var postfixLookupData = []struct {
	Input  string
	Output string
}{
	{"hi@example.com", "admin@example.com"},
	{"test@example.com", "admin@example.com"},
	{"user@example.org", "admin@example.org"},
	{"user@example.net", ""},
}

func TestDecoder_Load(t *testing.T) {
	p := &PostfixLookup{execCmd: func(key string) ([]byte, error) {
		n := strings.IndexByte(key, '@')
		if n == -1 {
			return []byte{}, nil
		}
		addr := key[n+1:]
		switch addr {
		case "example.com", "example.org":
			return []byte("result=admin@" + addr + "\nadmin@" + addr + "\n"), nil
		}
		return []byte{}, nil
	}}
	for _, i := range postfixLookupData {
		t.Run(i.Input, func(t *testing.T) {
			lookup, err := p.Lookup(i.Input)
			if i.Output == "" && err == nil {
				t.Fatal("expected error for empty output test case")
			}
			if i.Output != "" && err != nil {
				t.Fatal("expected no error for non-empty output test case")
			}
			assert.Equal(t, i.Output, lookup)
		})
	}
}
