package map_provider

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type Hash struct {
	r io.Reader
	v map[string]string
}

var _ MapProvider = &Hash{}

func NewHashMapProvider(filename string) (*Hash, error) {
	open, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return &Hash{open, make(map[string]string)}, nil
}

func (h *Hash) Load() error {
	scanner := bufio.NewScanner(h.r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(text, "#") {
			continue
		}

		n := strings.IndexByte(text, ':')
		key := strings.TrimSpace(text[:n])
		values := strings.Split(text[n+1:], ",")
		for _, i := range values {
			k := strings.TrimSpace(i)
			h.v[k] = key
		}
	}
	return scanner.Err()
}

func (h *Hash) Find(name string) (string, bool) {
	v, ok := h.v[name]
	return v, ok
}
