package comma_list_scanner

import (
	"bufio"
	"bytes"
	"io"
)

type CommaListScanner struct {
	r    *bufio.Scanner
	text string
	err  error
}

func NewCommaListScanner(r io.Reader) *CommaListScanner {
	s := bufio.NewScanner(r)
	s.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexByte(data, ','); i >= 0 {
			return i + 1, bytes.TrimSpace(data[0:i]), nil
		}
		// If we're at EOF, we have a final non-terminated line. Return it.
		if atEOF {
			return len(data), bytes.TrimSpace(data), nil
		}
		// Request more data.
		return 0, nil, nil
	})
	return &CommaListScanner{r: s}
}

func (c *CommaListScanner) Scan() bool {
	if c.r.Scan() {
		c.text = c.r.Text()
	}
	c.err = c.r.Err()
	return false
}

func (c *CommaListScanner) Text() string {
	return c.text
}
