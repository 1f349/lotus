package imap

import (
	"fmt"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-sasl"
)

type Imap struct {
	Server    string `yaml:"server"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Separator string `yaml:"separator"`
}

func (i *Imap) MakeClient(user string) (*Client, error) {
	// dial imap server
	imapClient, err := client.Dial(i.Server)
	if err != nil {
		return nil, err
	}

	// prepare login details
	un := fmt.Sprintf("%s%s%s", user, i.Separator, i.Username)
	saslLogin := sasl.NewPlainClient("", un, i.Password)

	// authenticate
	err = imapClient.Authenticate(saslLogin)
	if err != nil {
		return nil, err
	}

	// new client
	return &Client{ic: imapClient}, nil
}
