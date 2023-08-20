package imap

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"time"
)

var imapStatusFlags = []imap.StatusItem{
	imap.StatusMessages,
	imap.StatusRecent,
	imap.StatusUidNext,
	imap.StatusUidValidity,
	imap.StatusUnseen,
}

type Client struct {
	ic *client.Client
}

func (c *Client) Append(name string, flags []string, date time.Time, msg imap.Literal) error {
	return c.ic.Append(name, flags, date, msg)
}

func (c *Client) Copy(seqset *imap.SeqSet, dest string) error {
	return c.ic.Copy(seqset, dest)
}

func (c *Client) Create(name string) error {
	return c.ic.Create(name)
}

func (c *Client) Delete(name string) error {
	return c.ic.Delete(name)
}

func (c *Client) Fetch(folder string, start, end, limit uint32) ([]*imap.Message, error) {
	// select the mailbox
	mbox, err := c.ic.Select(folder, false)
	if err != nil {
		return nil, err
	}

	// setup fetch range
	if end > mbox.Messages {
		end = mbox.Messages
	}
	if end-start > limit {
		start = end - (limit - 1)
	}
	seqSet := new(imap.SeqSet)
	seqSet.AddRange(start, end)

	messages := make(chan *imap.Message, limit)
	done := make(chan error, 1)
	go func() {
		done <- c.ic.Fetch(seqSet, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	out := make([]*imap.Message, 0, limit)
	for msg := range messages {
		out = append(out, msg)
	}
	if err := <-done; err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) List(ref, name string) ([]*imap.MailboxInfo, error) {
	infos := make(chan *imap.MailboxInfo, 1)
	done := make(chan error, 1)
	go func() {
		done <- c.ic.List(ref, name, infos)
	}()

	out := make([]*imap.MailboxInfo, 0)
	for info := range infos {
		out = append(out, info)
	}
	if err := <-done; err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) Move(seqset *imap.SeqSet, dest string) error {
	return c.ic.Move(seqset, dest)
}

func (c *Client) Noop() error {
	return c.ic.Noop()
}

func (c *Client) Rename(existingName, newName string) error {
	return c.ic.Rename(existingName, newName)
}

func (c *Client) Search(criteria *imap.SearchCriteria) ([]uint32, error) {
	return c.ic.Search(criteria)
}

func (c *Client) Status(name string) (*imap.MailboxStatus, error) {
	mbox, err := c.ic.Status(name, imapStatusFlags)
	return mbox, err
}
