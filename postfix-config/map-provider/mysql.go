package map_provider

import (
	"database/sql"
	configParser "github.com/1f349/lotus/postfix-config/config-parser"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"os"
	"regexp"
	"strings"
)

var checkUatD = regexp.MustCompile("^[^@]+@[^@]+$")

type MySql struct {
	r     io.Reader
	db    *sql.DB
	query *PreparedQuery
}

var _ MapProvider = &MySql{}

func NewMySqlMapProvider(filename string) (*MySql, error) {
	open, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return &MySql{r: open}, nil
}

func (m *MySql) Load() error {
	p := configParser.NewConfigParser(m.r)
	c := mysql.NewConfig()
	var q string
	for p.Scan() {
		k, v := p.Pair()
		switch k {
		case "user":
			c.User = v
		case "password":
			c.Passwd = v
		case "hosts":
			c.Net = "tcp"
			c.Addr = v
		case "dbname":
			c.DBName = v
		case "query":
			q = v
		}
	}
	if err := p.Err(); err != nil {
		return err
	}

	q2, err := NewPreparedQuery(q)
	if err != nil {
		return err
	}
	m.query = q2

	// try opening connection
	db, err := sql.Open("mysql", c.FormatDSN())
	if err != nil {
		return err
	}
	m.db = db

	return db.Ping()
}

func (m *MySql) Find(name string) (string, bool) {
	format, err := m.query.Format(genQueryArgs(name))
	return format, err == nil
}

// genQueryArgs converts an input key into the % encoded parameters
//
// %s - full input key
// %u - user part of user@domain or full input key
// %d - domain part of user@domain or missing parameter
// %[1-9] - replaced with the most significant component of the input key's domain
//          for `user@mail.example.com` %1 = com, %2 = example, %3 = mail
//          otherwise they are missing parameters
func genQueryArgs(name string) map[byte]string {
	args := make(map[byte]string)
	args['s'] = name
	args['u'] = name
	if checkUatD.MatchString(name) {
		n := strings.IndexByte(name, '@')
		args['u'] = name[:n]
		args['d'] = name[n+1:]

		genDomainArgs(args, name[n+1:])
	}
	return args
}

// genDomainArgs replaces with the most significant component of the input key's
// domain for `user@mail.example.com` %1 = com, %2 = example, %3 = mail,
// otherwise they are missing parameters
func genDomainArgs(args map[byte]string, s string) {
	i, l := byte(1), len(s)
	for {
		n := strings.LastIndexByte(s, '.')
		if n == -1 {
			break
		}
		args[(i + '0')] = s[n+1 : l]
		l = n
	}
}
