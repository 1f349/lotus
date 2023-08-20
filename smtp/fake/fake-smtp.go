package fake

import (
	"github.com/emersion/go-smtp"
	"io"
	"log"
)

type SmtpBackend struct {
	Debug chan []byte
}

func (f *SmtpBackend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &SmtpSession{f.Debug}, nil
}

type SmtpSession struct {
	Debug chan []byte
}

func (f *SmtpSession) Reset() {}

func (f *SmtpSession) Logout() error { return nil }

func (f *SmtpSession) AuthPlain(username, password string) error { return nil }

func (f *SmtpSession) Mail(from string, opts *smtp.MailOptions) error {
	log.Println("MAIL " + from)
	f.Debug <- []byte("MAIL " + from + "\n")
	return nil
}

func (f *SmtpSession) Rcpt(to string) error {
	f.Debug <- []byte("RCPT " + to + "\n")
	return nil
}

func (f *SmtpSession) Data(r io.Reader) error {
	all, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	f.Debug <- all
	return nil
}
