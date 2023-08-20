package fake

import (
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
)

type Backend struct {
	Debug    chan []byte
	Username string
	Password string
}

func (i *Backend) Login(connInfo *imap.ConnInfo, username, password string) (backend.User, error) {
	if username != i.Username || password != i.Password {
		return nil, fmt.Errorf("invalid user")
	}
	return &User{i.Debug, username}, nil
}
