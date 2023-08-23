package smtp

import (
	"bytes"
	"github.com/1f349/lotus/smtp/fake"
	"github.com/emersion/go-message"
	"github.com/emersion/go-message/mail"
	"github.com/emersion/go-smtp"
	"github.com/hydrogen18/memlistener"
	"github.com/stretchr/testify/assert"
	"log"
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
	listener := memlistener.NewMemoryListener()
	serverData := make(chan []byte, 4)
	server := smtp.NewServer(&fake.SmtpBackend{Debug: serverData})
	go func() {
		_ = server.Serve(listener)
	}()

	defaultDialer = func(addr string) (*smtp.Client, error) {
		dial, err := listener.Dial("", "")
		if err != nil {
			return nil, err
		}
		return smtp.NewClient(dial, "localhost")
	}

	s := &Smtp{Server: "localhost:25"}
	err := s.Send(&Mail{From: "test@localhost", Deliver: []string{"a@localhost", "b@localhost"}, Body: sendTestMessage})
	assert.NoError(t, err)
	assert.Equal(t, []byte("MAIL test@localhost\n"), <-serverData)
	assert.Equal(t, []byte("RCPT a@localhost\n"), <-serverData)
	assert.Equal(t, []byte("RCPT b@localhost\n"), <-serverData)
	assert.Equal(t, append(sendTestMessage, '\r', '\n'), <-serverData)
}

func TestCreateSenderSlice(t *testing.T) {
	a := []*mail.Address{{Address: "a@example.com"}, {Address: "b@example.com"}}
	b := []*mail.Address{{Address: "a@example.com"}, {Address: "c@example.com"}}
	c := []*mail.Address{{Address: "a@example.com"}, {Address: "d@example.com"}}
	assert.Equal(t, []string{
		"a@example.com",
		"b@example.com",
		"a@example.com",
		"c@example.com",
		"a@example.com",
		"d@example.com",
	}, CreateSenderSlice(a, b, c))
}
