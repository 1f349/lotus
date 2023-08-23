package smtp

import (
	"github.com/emersion/go-message/mail"
	"os/exec"
)

type Smtp struct {
	Server string `yaml:"server"`
}

type Mail struct {
	From *mail.Address
	Body []byte
}

var execSendMail = func(from string) *exec.Cmd {
	return exec.Command("/usr/lib/sendmail", "-f", from, "-t")
}

func (s *Smtp) Send(mail *Mail) error {
	// start sendmail caller
	sendMail := execSendMail(mail.From.String())
	inPipe, err := sendMail.StdinPipe()
	if err != nil {
		return err
	}

	// write message body
	_, err = inPipe.Write(mail.Body)
	if err != nil {
		return err
	}
	err = inPipe.Close()
	if err != nil {
		return err
	}

	return sendMail.Run()
}
