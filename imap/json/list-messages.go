package json

import (
	"encoding/json"
	"github.com/emersion/go-imap"
)

type ListMessagesJson []*imap.Message

func (l ListMessagesJson) MarshalJSON() ([]byte, error) {
	a := make([]encodeImapMessage, len(l))
	for i := range a {
		a[i] = encodeImapMessage(*l[i])
	}
	return json.Marshal(a)
}

type encodeImapMessage imap.Message

func (e encodeImapMessage) MarshalJSON() ([]byte, error) {
	body := make(map[string]imap.Literal, len(e.Body))
	for k, v := range e.Body {
		body[string(k.FetchItem())] = v
	}
	return json.Marshal(map[string]any{
		"SeqNum":        e.SeqNum,
		"Items":         e.Items,
		"Envelope":      e.Envelope,
		"BodyStructure": e.BodyStructure,
		"Flags":         e.Flags,
		"InternalDate":  e.InternalDate,
		"Size":          e.Size,
		"Uid":           e.Uid,
		"$Body":         body,
	})
}
