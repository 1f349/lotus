package config_parser

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var configParserData = []struct {
	Input  string
	Values [][2]string
}{
	{
		"a = a",
		[][2]string{{"a", "a"}},
	},
	{
		"     a = a    ",
		[][2]string{{"a", "a"}},
	},
	{
		"   # this is a comment\n  a = a, b\nb = c, d",
		[][2]string{{"a", "a, b"}, {"b", "c, d"}},
	},
}

func TestConfigParser(t *testing.T) {
	for _, i := range configParserData {
		t.Run(i.Input, func(t *testing.T) {
			a := NewConfigParser(strings.NewReader(i.Input))
			n := 0
			for a.Scan() {
				assert.False(t, n >= len(i.Values))
				assert.Equal(t, i.Values[n], a.pair)
				n++
			}
		})
	}
}
