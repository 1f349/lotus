package sendmail

import (
	"github.com/emersion/go-message/mail"
	"os/exec"
)

type Smtp struct {
	SendMailCommand string `json:"send_mail_command"`
}

type Mail struct {
	From *mail.Address
	Body []byte
}

var execCommand = exec.Command

func (s *Smtp) Send(mail *Mail) error {
	// start sendmail caller
	if s.SendMailCommand == "" {
		s.SendMailCommand = "/usr/sbin/sendmail"
	}
	sendMail := execCommand(s.SendMailCommand, "-f", mail.From.Address, "-t")
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

	// run command
	return sendMail.Run()
}
