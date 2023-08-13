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
	from    string
	deliver []string
	body    []byte
}

func (s *Smtp) Send(mail *Mail) error {
	// dial smtp server
	smtpClient, err := smtp.Dial(s.Server)
	if err != nil {
		return err
	}

	// use a reader to send bytes
	r := bytes.NewReader(mail.body)

	// send mail
	return smtpClient.SendMail(mail.from, mail.deliver, r)
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
