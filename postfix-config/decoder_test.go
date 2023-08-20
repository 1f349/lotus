package postfix_config

import (
	"bytes"
	_ "embed"
	configParser "github.com/1f349/lotus/postfix-config/config-parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed example.cf
var exampleConfig []byte

func TestDecoder_Load(t *testing.T) {
	b := bytes.NewReader(exampleConfig)
	d := &Decoder{r: configParser.NewConfigParser(b)}
	assert.NoError(t, d.Load())
}
