package map_provider

import (
	"errors"
	"sort"
	"strings"
	"unicode"
)

var (
	ErrMissingArgument = errors.New("missing argument")
	ErrInvalidRawQuery = errors.New("invalid raw query")
)

type PreparedQuery struct {
	raw    string
	params map[int]byte
}

func NewPreparedQuery(raw string) (*PreparedQuery, error) {
	var s strings.Builder
	origin := 0
	params := make(map[int]byte)
	for {
		n := strings.IndexByte(raw[origin:], '%')
		if n == -1 {
			break
		}
		n += origin
		if n+1 == len(raw) {
			return nil, ErrInvalidRawQuery
		}
		s.WriteString(raw[origin:n])
		if raw[n+1] == '%' {
			s.WriteByte('%')
			origin = n + 1
			continue
		}
		params[s.Len()] = toLower(raw[n+1])
		origin = n + 2
	}
	s.WriteString(raw[origin:])
	return &PreparedQuery{
		raw:    s.String(),
		params: params,
	}, nil
}

func (p *PreparedQuery) Format(args map[byte]string) (string, error) {
	var s strings.Builder
	keys := make([]int, 0, len(p.params))
	for k := range p.params {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	origin := 0
	for _, k := range keys {
		r, ok := args[p.params[k]]
		if !ok {
			return "", ErrMissingArgument
		}

		// write up to and including the next parameter
		s.WriteString(p.raw[origin:k])
		s.WriteString(strings.ReplaceAll(r, "'", ""))
		origin = k
	}

	// write the rest of the query
	s.WriteString(p.raw[origin:])
	return s.String(), nil
}

func toLower(a byte) byte {
	return byte(unicode.ToLower(rune(a)))
}
