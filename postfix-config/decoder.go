package postfix_config

import (
	"bufio"
	"fmt"
	configParser "github.com/1f349/lotus/postfix-config/config-parser"
	mapProvider "github.com/1f349/lotus/postfix-config/map-provider"
	"io"
	"strings"
)

type Decoder struct {
	r *configParser.ConfigParser
	v *Config
	t map[string]string
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: configParser.NewConfigParser(r)}
}

func (d *Decoder) Load() error {
	d.v = &Config{}
	for d.r.Scan() {
		k, v := d.r.Pair()
		if d.v.NeedsMapProvider(k) {
			m := mapProvider.SequenceMapProvider{}

			s := bufio.NewScanner(strings.NewReader(v))
			s.Split(bufio.ScanWords)
			for s.Scan() {
				a := s.Text()
				println("a", a)
				if strings.HasPrefix(a, "$") {
					// is variable
				}
				n := strings.IndexByte(a, ':')
				if n == -1 {
					return fmt.Errorf("missing prefix")
				}
			}
			if err := s.Err(); err != nil {
				return err
			}
			d.v.SetKey(k, m)
		}
	}
	return d.r.Err()
}
