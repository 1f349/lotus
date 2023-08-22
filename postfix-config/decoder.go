package postfix_config

import (
	"errors"
	"fmt"
	commaListScanner "github.com/1f349/lotus/postfix-config/comma-list-scanner"
	configParser "github.com/1f349/lotus/postfix-config/config-parser"
	mapProvider "github.com/1f349/lotus/postfix-config/map-provider"
	"io"
	"path/filepath"
	"strings"
)

type Decoder struct {
	r        *configParser.ConfigParser
	temp     map[string]string
	basePath string
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: configParser.NewConfigParser(r)}
}

func (d *Decoder) Load() error {
	for d.r.Scan() {
		k, v := d.r.Pair()
		d.temp[k] = v
	}
	if err := d.r.Err(); err != nil {
		return err
	}

		switch d.value.ParseProvider(k) {
		case "comma":
			m := mapProvider.SequenceMapProvider{}

			s := commaListScanner.NewCommaListScanner(strings.NewReader(v))
			for s.Scan() {
				a := s.Text()
				println("a", a)
				if strings.HasPrefix(a, "$") {
					m = append(m, &mapProvider.Variable{Name: a[1:]})
					continue
				}

				v2, err := d.createValue(a)
				if err != nil {
					return err
				}
				m = append(m, v2)
			}
			if err := s.Err(); err != nil {
				return err
			}
			d.value.SetKey(k, m)
		case "union":
			if !strings.HasPrefix(v, "unionmap:{") || !strings.HasSuffix(v, "}") {
				return errors.New("key requires a union map")
			}
			v = v[len("unionmap:{") : len(v)-1]

			m := mapProvider.SequenceMapProvider{}
			s := commaListScanner.NewCommaListScanner(strings.NewReader(v))
			for s.Scan() {
				a := s.Text()
				v2, err := d.createValue(a)
				if err != nil {
					return err
				}
				m = append(m, v2)
			}
		default:
			return fmt.Errorf("key '%s' has no defined parse provider", k)
		}
	}
	return d.r.Err()
}

func (d *Decoder) createValue(a string) (mapProvider.MapProvider, error) {
	n := strings.IndexByte(a, ':')
	if n == -1 {
		return nil, fmt.Errorf("missing prefix")
	}

	namespace := a[:n]
	value := a[n+1:]
	switch namespace {
	case "mysql":
		if !filepath.IsAbs(value) {
			value = filepath.Join(d.basePath, value)
		}
		provider, err := mapProvider.NewMySqlMapProvider(value)
		if err != nil {
			return nil, err
		}
		return provider, nil
	case "hash":
		if !filepath.IsAbs(value) {
			value = filepath.Join(d.basePath, value)
		}
		provider, err := mapProvider.NewHashMapProvider(value)
		if err != nil {
			return nil, err
		}
		return provider, nil
	}
	return nil, errors.New("invalid provider namespace")
}
