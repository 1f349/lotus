package comma_list_scanner

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var testCommaList = []string{
	"hello, wow-this-is-cool, amazing",
	"hello, wow-this-is-cool",
	"hello, wow-this-is-cool, ",
	"hello, wow-this-is-cool,",
	",hello, wow-this-is-cool",
	",hello, wow-this-is-cool,",
	"hello, wow-this-is-cool,,,",
}

func TestNewCommaListScanner(t *testing.T) {
	for _, i := range testCommaList {
		t.Run(i, func(t *testing.T) {
			// use comma list scanner
			s := NewCommaListScanner(strings.NewReader(i))
			n := strings.Count(i, ",")
			a := make([]string, 0, n+1)
			for s.Scan() {
				a = append(a, s.Text())
			}
			assert.NoError(t, s.Err())

			// test against splitting and trimming strings
			b := strings.Split(i, ",")
			for i := 0; i < len(b); i++ {
				c := strings.TrimSpace(b[i])
				if c == "" {
					b = append(b[0:i], b[i+1:]...)
					i--
				} else {
					b[i] = c
				}
			}
			assert.Equal(t, b, a)
		})
	}
}
