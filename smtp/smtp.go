package smtp

import (
	"os/exec"
)

type Smtp struct {
	Server string `yaml:"server"`
}

type Mail struct {
	From string
	Body []byte
}

var execSendMail = func(from string) *exec.Cmd {
	return exec.Command("/usr/lib/sendmail", "-f", from, "-t")
}

func (s *Smtp) Send(mail *Mail) error {
	// start sendmail caller
	sendMail := execSendMail(mail.From)
	inPipe, err := sendMail.StdinPipe()
	if err != nil {
		return err
	}

	// write message body
	_, err = inPipe.Write(mail.Body)
	if err != nil {
		return err
	}
	return inPipe.Close()
}
