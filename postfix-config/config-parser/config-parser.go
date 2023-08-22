package config_parser

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

var ErrInvalidConfigLine = errors.New("invalid config line")

type ConfigParser struct {
	s    *bufio.Scanner
	pair [2]string
	err  error
}

func NewConfigParser(r io.Reader) *ConfigParser {
	return &ConfigParser{s: bufio.NewScanner(r)}
}

func (c *ConfigParser) Scan() bool {
scanAgain:
	if !c.s.Scan() {
		return false
	}
	text := strings.TrimSpace(c.s.Text())
	if text == "" || strings.HasPrefix(text, "#") {
		goto scanAgain
	}
	n := strings.IndexByte(text, '=')
	if n < 2 || n+2 >= len(text) || text[n-1] != ' ' || text[n+1] != ' ' {
		c.err = ErrInvalidConfigLine
		return false
	}
	c.pair = [2]string{text[:n-1], text[n+2:]}
	return true
}

func (c *ConfigParser) Pair() (string, string) {
	return strings.TrimSpace(c.pair[0]), strings.TrimSpace(c.pair[1])
}

func (c *ConfigParser) Err() error {
	if c.err != nil {
		return c.err
	}
	return c.s.Err()
}
