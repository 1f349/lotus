package marshal

import (
	"encoding/json"
	"github.com/emersion/go-imap"
)

var _, _ json.Marshaler = &MessageSliceJson{}, &MessageJson{}

type MessageSliceJson []*imap.Message

func (m MessageSliceJson) MarshalJSON() ([]byte, error) {
	a := make([]MessageJson, len(m))
	for i := range a {
		a[i] = MessageJson(*m[i])
	}
	return json.Marshal(a)
}

type MessageJson imap.Message

func (m MessageJson) MarshalJSON() ([]byte, error) {
	body := make(map[string]imap.Literal, len(m.Body))
	for k, v := range m.Body {
		body[string(k.FetchItem())] = v
	}
	return json.Marshal(map[string]any{
		"SeqNum":        m.SeqNum,
		"Items":         m.Items,
		"Envelope":      m.Envelope,
		"BodyStructure": m.BodyStructure,
		"Flags":         m.Flags,
		"InternalDate":  m.InternalDate,
		"Size":          m.Size,
		"Uid":           m.Uid,
		"$Body":         body,
	})
}
