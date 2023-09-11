package sendmail

import (
	"bytes"
	"github.com/emersion/go-message"
	"github.com/emersion/go-message/mail"
	"github.com/stretchr/testify/assert"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

var sendTestMessage []byte

func init() {
	var h mail.Header
	h.SetDate(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.Local))
	h.SetSubject("Happy Millennium")
	h.SetAddressList("From", []*mail.Address{{Name: "Test", Address: "test@localhost"}})
	h.SetAddressList("To", []*mail.Address{{Name: "A", Address: "a@localhost"}})
	h.Set("Content-Type", "text/plain; charset=utf-8")
	entity, err := message.New(h.Header, strings.NewReader("Thanks"))
	if err != nil {
		log.Fatal(err)
	}
	out := new(bytes.Buffer)
	if entity.WriteTo(out) != nil {
		log.Fatal(err)
	}
	sendTestMessage = out.Bytes()
}

func TestSmtp_Send(t *testing.T) {
	execCommand = func(name string, arg ...string) *exec.Cmd {
		cs := append([]string{"-test.run=TestHelperProcess", "--", name}, arg...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
		return cmd
	}
	defer func() { execCommand = exec.Command }()

	m := &Mail{From: &mail.Address{Address: "test@localhost"}, Body: sendTestMessage}

	temp, err := os.CreateTemp("", "sendmail")
	if err != nil {
		return
	}

	addr, err := net.ResolveUnixAddr("")
	assert.NoError(t, err)
	listen, err := net.ListenUnix("", addr)
	assert.NoError(t, err)

	s := &SendMail{SendMailCommand: "/tmp/sendmailXXXXX"}
	assert.NoError(t, s.Send(m))
}

func TestSmtpHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

}
