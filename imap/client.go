package imap

import (
	"errors"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"strconv"
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

var ErrInvalidArguments = errors.New("invalid arguments")

func (c *Client) HandleWS(action string, args []string) (map[string]any, error) {
	switch action {
	case "copy":
		// TODO: implementation
	case "create":
		// TODO: implementation
	case "delete":
		// TODO: implementation
	case "list":
		if len(args) != 2 {
			return nil, ErrInvalidArguments
		}

		// do list
		list, err := c.list(args[0], args[1])
		if err != nil {
			return nil, err
		}
		return map[string]any{"type": "list", "value": list}, nil
	case "fetch":
		if len(args) != 4 {
			return nil, ErrInvalidArguments
		}

		// parse numeric parameters
		arg1i, err := strconv.Atoi(args[1])
		if err != nil {
			return nil, err
		}
		arg2i, err := strconv.Atoi(args[2])
		if err != nil {
			return nil, err
		}
		arg3i, err := strconv.Atoi(args[3])
		if err != nil {
			return nil, err
		}

		// do fetch
		fetch, err := c.fetch(args[0], uint32(arg1i), uint32(arg2i), uint32(arg3i))
		if err != nil {
			return nil, err
		}
		return map[string]any{"type": "fetch", "value": fetch}, nil
	case "move":
		// TODO: implementation
	case "rename":
		// TODO: implementation
	case "search":
		// TODO: implementation
	case "status":
		// TODO: implementation
	}
	return map[string]any{"error": "Not implemented"}, nil
}

func (c *Client) fetch(folder string, start, end, limit uint32) ([]*imap.Message, error) {
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
		done <- c.ic.Fetch(seqSet, []imap.FetchItem{imap.FetchEnvelope, imap.FetchUid, imap.FetchFlags}, messages)
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

func (c *Client) list(ref, name string) ([]*imap.MailboxInfo, error) {
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
