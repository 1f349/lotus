package imap

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
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

func (c *Client) Status(folder string) (*imap.MailboxStatus, error) {
	mbox, err := c.ic.Status(folder, imapStatusFlags)
	return mbox, err
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

	outMsg := make([]*imap.Message, 0, limit)
	for msg := range messages {
		outMsg = append(outMsg, msg)
	}
	if err := <-done; err != nil {
		return nil, err
	}
	return outMsg, nil
}
