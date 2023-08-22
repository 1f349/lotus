package comma_list_scanner

import (
	"bufio"
	"bytes"
	"fmt"
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
		println("data", fmt.Sprintf("%s", data))
		println("index", bytes.IndexAny(data, " ,"))
		if i := bytes.IndexAny(data, " ,"); i >= 0 {
			// consume all spaces after the comma
			j := i + 1
			for j < len(data) && data[j] == ' ' {
				j++
			}
			return j, bytes.TrimSpace(data[0:i]), nil
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
		return true
	}
	c.err = c.r.Err()
	return false
}

func (c *CommaListScanner) Text() string {
	return c.text
}

func (c *CommaListScanner) Err() error {
	return c.err
}
