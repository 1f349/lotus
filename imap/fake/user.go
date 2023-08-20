package fake

import (
	"fmt"
	"github.com/emersion/go-imap/backend"
)

type User struct {
	Debug    chan []byte
	ImapUser string
}

func (i *User) Username() string {
	return i.ImapUser
}

func (i *User) ListMailboxes(subscribed bool) ([]backend.Mailbox, error) {
	return []backend.Mailbox{}, nil
}

func (i *User) GetMailbox(name string) (backend.Mailbox, error) {
	return &Mailbox{i.Debug, name}, nil
}

func (i *User) CreateMailbox(name string) error {
	return fmt.Errorf("failed to create mailbox")
}

func (i *User) DeleteMailbox(name string) error {
	return fmt.Errorf("failed to delete mailbox")
}

func (i *User) RenameMailbox(existingName, newName string) error {
	return fmt.Errorf("failed to rename mailbox")
}

func (i *User) Logout() error {
	return nil
}
