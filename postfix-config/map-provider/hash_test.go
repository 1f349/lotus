package map_provider

import (
	"bytes"
	_ "embed"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed hash_example.txt
var hashExample []byte

func TestHash_Load(t *testing.T) {
	h := &Hash{r: bytes.NewReader(hashExample), v: make(map[string]string)}
	assert.NoError(t, h.Load())
	assert.Equal(t, map[string]string{
		"root":    "postmaster",
		"this":    "test",
		"is":      "test",
		"an":      "test",
		"example": "test",
	}, h.v)
}
