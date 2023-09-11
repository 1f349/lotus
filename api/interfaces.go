package api

import (
	"github.com/1f349/lotus/imap"
	"github.com/1f349/lotus/sendmail"
)

type Smtp interface {
	Send(mail *sendmail.Mail) error
}

type Imap interface {
	MakeClient(user string) (*imap.Client, error)
}
