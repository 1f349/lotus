package smtp

import (
	"bytes"
	"github.com/emersion/go-message/mail"
	"github.com/emersion/go-smtp"
)

type Smtp struct {
	Server string `yaml:"server"`
}

type Mail struct {
	From    string
	Deliver []string
	Body    []byte
}

var defaultDialer = smtp.Dial

func (s *Smtp) Send(mail *Mail) error {
	// dial smtp server
	smtpClient, err := defaultDialer(s.Server)
	if err != nil {
		return err
	}

	// use a reader to send bytes
	r := bytes.NewReader(mail.Body)

	// send mail
	return smtpClient.SendMail(mail.From, mail.Deliver, r)
}

func CreateSenderSlice(to, cc, bcc []*mail.Address) []string {
	a := make([]string, 0, len(to)+len(cc)+len(bcc))
	for _, i := range to {
		a = append(a, i.Address)
	}
	for _, i := range cc {
		a = append(a, i.Address)
	}
	for _, i := range bcc {
		a = append(a, i.Address)
	}
	return a
}
