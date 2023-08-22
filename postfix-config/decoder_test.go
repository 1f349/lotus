package postfix_config

import (
	"bytes"
	_ "embed"
	configParser "github.com/1f349/lotus/postfix-config/config-parser"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

//go:embed example.cf
var exampleConfig []byte

func TestDecoder_Load(t *testing.T) {
	// get working directory
	wd, err := os.Getwd()
	assert.NoError(t, err)

	// read example config
	b := bytes.NewReader(exampleConfig)
	d := &Decoder{
		r:        configParser.NewConfigParser(b),
		basePath: filepath.Join(wd, "test-data"),
	}
	assert.NoError(t, d.Load())
}
