package api

import (
	"github.com/1f349/lotus/imap"
	"github.com/1f349/lotus/smtp"
)

type Smtp interface {
	Send(mail *smtp.Mail) error
}

type Imap interface {
	MakeClient(user string) (*imap.Client, error)
}
