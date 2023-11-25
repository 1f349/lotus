package imap

import (
	"encoding/json"
	"errors"
	"github.com/1f349/lotus/imap/marshal"
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

var ErrInvalidArguments = errors.New("invalid arguments")

func (c *Client) HandleWS(action string, args json.RawMessage) (map[string]any, error) {
	switch action {
	case "copy":
		// TODO: implementation
	case "create":
		// TODO: implementation
	case "delete":
		// TODO: implementation
	case "list":
		var listArgs []string
		err := json.Unmarshal(args, &listArgs)
		if err != nil {
			return nil, err
		}

		if len(listArgs) != 2 {
			return nil, ErrInvalidArguments
		}

		// do list
		list, err := c.list(listArgs[0], listArgs[1])
		if err != nil {
			return nil, err
		}
		return map[string]any{"type": "list", "value": list}, nil
	case "fetch":
		var fetchArgs struct {
			Sync  uint64 `json:"sync"`
			Path  string `json:"path"`
			Start uint32 `json:"start"`
			End   uint32 `json:"end"`
			Limit uint32 `json:"limit"`
		}
		err := json.Unmarshal(args, &fetchArgs)
		if err != nil {
			return nil, err
		}

		if fetchArgs.Sync == 0 || len(fetchArgs.Path) == 0 || fetchArgs.Start == 0 || fetchArgs.End == 0 || fetchArgs.Limit == 0 {
			return nil, ErrInvalidArguments
		}

		// do fetch
		fetch, err := c.fetch(fetchArgs.Path, fetchArgs.Start, fetchArgs.End, fetchArgs.Limit)
		if err != nil {
			return nil, err
		}
		return map[string]any{"type": "fetch", "sync": fetchArgs.Sync, "value": marshal.MessageSliceJson(fetch)}, nil
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
		done <- c.ic.Fetch(seqSet, []imap.FetchItem{imap.FetchEnvelope, imap.FetchUid, imap.FetchFlags, imap.FetchInternalDate}, messages)
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
