package fake

import (
	"fmt"
	"github.com/emersion/go-imap"
	"time"
)

type Mailbox struct {
	Debug    chan []byte
	ImapName string
}

func (m *Mailbox) Name() string {
	return m.ImapName
}

func (m *Mailbox) Info() (*imap.MailboxInfo, error) {
	return &imap.MailboxInfo{
		Attributes: []string{imap.UnmarkedAttr, imap.HasNoChildrenAttr},
		Delimiter:  "/",
		Name:       m.ImapName,
	}, nil
}

func (m *Mailbox) Status(items []imap.StatusItem) (*imap.MailboxStatus, error) {
	return &imap.MailboxStatus{
		Name:     m.ImapName,
		Messages: 1,
	}, nil
}

func (m *Mailbox) SetSubscribed(subscribed bool) error {
	return fmt.Errorf("failed to subscribe")
}

func (m *Mailbox) Check() error {
	return nil
}

func (m *Mailbox) ListMessages(uid bool, seqset *imap.SeqSet, items []imap.FetchItem, ch chan<- *imap.Message) error {
	return fmt.Errorf("failed to list messages")
}

func (m *Mailbox) SearchMessages(uid bool, criteria *imap.SearchCriteria) ([]uint32, error) {
	return nil, fmt.Errorf("failed to search messages")
}

func (m *Mailbox) CreateMessage(flags []string, date time.Time, body imap.Literal) error {
	return fmt.Errorf("failed to create message")
}

func (m *Mailbox) UpdateMessagesFlags(uid bool, seqset *imap.SeqSet, operation imap.FlagsOp, flags []string) error {
	return fmt.Errorf("failed to update message flags")
}

func (m *Mailbox) CopyMessages(uid bool, seqset *imap.SeqSet, dest string) error {
	return fmt.Errorf("failed to copy messages")
}

func (m *Mailbox) Expunge() error {
	return fmt.Errorf("failed to expunge")
}
