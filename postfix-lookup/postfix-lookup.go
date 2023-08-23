package postfix_lookup

import (
	"bufio"
	"bytes"
	_ "embed"
	"errors"
	"os/exec"
	"strings"
)

var ErrInvalidAlias = errors.New("invalid alias")

//go:embed lookup.sh
var lookupScript string

type PostfixLookup struct {
	execCmd func(key string) ([]byte, error)
}

func NewPostfixLookup() *PostfixLookup {
	return &PostfixLookup{
		execCmd: func(key string) ([]byte, error) {
			return exec.Command("bash", "-c", lookupScript, "--", key).Output()
		},
	}
}

func (d *PostfixLookup) Lookup(key string) (string, error) {
	output, err := d.execCmd(key)
	if err != nil {
		return "", err
	}

	s := bufio.NewScanner(bytes.NewReader(output))
	for s.Scan() {
		a := s.Text()
		n := strings.IndexByte(a, '=')
		if n != -1 && a[:n] == "result" {
			return a[n+1:], nil
		}
	}
	if err := s.Err(); err != nil {
		return "", err
	}
	return "", ErrInvalidAlias
}
