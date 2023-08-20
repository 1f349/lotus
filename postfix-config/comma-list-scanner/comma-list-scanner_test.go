package comma_list_scanner

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var testCommaList = []struct {
	text string
	out  []string
}{
	{"hello, wow this is cool, amazing", []string{"hello", "wow this is cool", "amazing"}},
	{"hello, wow this is cool, amazing", []string{"hello", "wow this is cool", "amazing"}},
}

func TestNewCommaListScanner(t *testing.T) {
	for _, i := range testCommaList {
		t.Run(i.text, func(t *testing.T) {
			s := NewCommaListScanner(strings.NewReader(i.text))
			n := 0
			for s.Scan() {
				assert.Equal(t, i.out[n], s.Text())
			}
		})
	}
}
