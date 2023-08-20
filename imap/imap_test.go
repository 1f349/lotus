package imap

import (
	"github.com/1f349/lotus/imap/fake"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-imap/server"
	"github.com/hydrogen18/memlistener"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestImap_MakeClient(t *testing.T) {
	listener := memlistener.NewMemoryListener()
	serverData := make(chan []byte, 4)
	srv := server.New(&fake.Backend{Debug: serverData, Username: "a@localhost*master@localhost", Password: "1234"})
	srv.AllowInsecureAuth = true
	go func() {
		_ = srv.Serve(listener)
	}()

	defaultDialer = func(addr string) (*client.Client, error) {
		dial, err := listener.Dial("", "")
		if err != nil {
			return nil, err
		}
		return client.New(dial)
	}

	i := &Imap{Server: "localhost", Username: "master@localhost", Password: "1234", Separator: "*"}
	cli, err := i.MakeClient("a@localhost")
	assert.NoError(t, err)
	status, err := cli.Status("INBOX")
	assert.NoError(t, err)
	assert.Equal(t, "INBOX", status.Name)
	assert.Equal(t, uint32(1), status.Messages)
}
