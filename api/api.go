package api

import (
	"encoding/json"
	"github.com/1f349/lotus/imap"
	"github.com/1f349/lotus/imap/marshal"
	imap2 "github.com/emersion/go-imap"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

func SetupApiServer(listen string, auth func(callback AuthCallback) httprouter.Handle, send Smtp, recv Imap) *http.Server {
	r := httprouter.New()

	// === ACCOUNT ===
	r.GET("/identities", auth(func(rw http.ResponseWriter, req *http.Request, params httprouter.Params, b AuthClaims) {
		// TODO(melon): find users aliases and other account data
	}))

	// === SMTP ===
	r.POST("/message", auth(MessageSender(send)))

	r.Handle(http.MethodConnect, "/", func(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
		
	})

	// === IMAP ===
	type mailboxStatusJson struct {
		Folder string `json:"folder"`
	}
	r.GET("/mailbox/status", auth(imapClient(recv, func(rw http.ResponseWriter, req *http.Request, params httprouter.Params, cli *imap.Client, t mailboxStatusJson) error {
		status, err := cli.Status(t.Folder)
		if err != nil {
			return err
		}
		return json.NewEncoder(rw).Encode(status)
	})))

	type mailboxListJson struct {
		Folder  string `json:"folder"`
		Pattern string `json:"pattern"`
	}
	r.GET("/mailbox/list", auth(imapClient(recv, func(rw http.ResponseWriter, req *http.Request, params httprouter.Params, cli *imap.Client, t mailboxListJson) error {
		list, err := cli.List(t.Folder, t.Pattern)
		if err != nil {
			return err
		}
		return json.NewEncoder(rw).Encode(list)
	})))

	type mailboxCreateJson struct {
		Name string `json:"name"`
	}
	r.POST("/mailbox/create", auth(imapClient(recv, func(rw http.ResponseWriter, req *http.Request, params httprouter.Params, cli *imap.Client, t mailboxCreateJson) error {
		err := cli.Create(t.Name)
		if err != nil {
			return err
		}
		return json.NewEncoder(rw).Encode(map[string]string{"Status": "OK"})
	})))

	type messagesListJson struct {
		Folder string `json:"folder"`
		Start  uint32 `json:"start"`
		End    uint32 `json:"end"`
		Limit  uint32 `json:"limit"`
	}
	r.GET("/list-messages", auth(imapClient(recv, func(rw http.ResponseWriter, req *http.Request, params httprouter.Params, cli *imap.Client, t messagesListJson) error {
		messages, err := cli.Fetch(t.Folder, t.Start, t.End, t.Limit)
		if err != nil {
			return err
		}
		return json.NewEncoder(rw).Encode(marshal.MessageSliceJson(messages))
	})))

	type messagesSearchJson struct {
		Folder       string `json:"folder"`
		SeqNum       imap2.SeqSet
		Uid          imap2.SeqSet
		Since        time.Time
		Before       time.Time
		SentSince    time.Time
		SentBefore   time.Time
		Body         []string
		Text         []string
		WithFlags    []string
		WithoutFlags []string
		Larger       uint32
		Smaller      uint32
	}
	r.GET("/search-messages", auth(imapClient(recv, func(rw http.ResponseWriter, req *http.Request, params httprouter.Params, cli *imap.Client, t messagesSearchJson) error {
		status, err := cli.Search(&imap2.SearchCriteria{
			SeqNum:       t.SeqNum,
			Uid:          t.Uid,
			Since:        time.Time{},
			Before:       time.Time{},
			SentSince:    time.Time{},
			SentBefore:   time.Time{},
			Header:       nil,
			Body:         nil,
			Text:         nil,
			WithFlags:    nil,
			WithoutFlags: nil,
			Larger:       0,
			Smaller:      0,
			Not:          nil,
			Or:           nil,
		})
		if err != nil {
			return err
		}
		return json.NewEncoder(rw).Encode(status)
	})))
	r.POST("/update-messages-flags", auth(imapClient(recv, func(rw http.ResponseWriter, req *http.Request, params httprouter.Params, cli *imap.Client, t mailboxStatusJson) {
		status, err := cli.Status(t.Folder)
		if err != nil {
			rw.WriteHeader(http.StatusForbidden)
			return
		}
		_ = json.NewEncoder(rw).Encode(status)
	})))
	r.GET("/list-messages", auth(imapClient(recv, func(rw http.ResponseWriter, req *http.Request, params httprouter.Params, cli *imap.Client, t statusJson) {
		messages, err := cli.Fetch(t.Folder, 1, 100, 100)
		if err != nil {
			rw.WriteHeader(http.StatusForbidden)
			return
		}
		err = json.NewEncoder(rw).Encode(marshal.ListMessagesJson(messages))
		if err != nil {
			log.Println("list-messages json encode error:", err)
		}
	})))
	r.GET("/search-messages", auth(imapClient(recv, func(rw http.ResponseWriter, req *http.Request, params httprouter.Params, cli *imap.Client, t statusJson) {
		status, err := cli.Status(t.Folder)
		if err != nil {
			rw.WriteHeader(http.StatusForbidden)
			return
		}
		_ = json.NewEncoder(rw).Encode(status)
	})))
	r.POST("/create-message", auth(imapClient(recv, func(rw http.ResponseWriter, req *http.Request, params httprouter.Params, cli *imap.Client, t statusJson) {
		status, err := cli.Status(t.Folder)
		if err != nil {
			rw.WriteHeader(http.StatusForbidden)
			return
		}
		_ = json.NewEncoder(rw).Encode(status)
	})))
	r.POST("/update-messages-flags", auth(imapClient(recv, func(rw http.ResponseWriter, req *http.Request, params httprouter.Params, cli *imap.Client, t statusJson) {
		status, err := cli.Status(t.Folder)
		if err != nil {
			rw.WriteHeader(http.StatusForbidden)
			return
		}
		_ = json.NewEncoder(rw).Encode(status)
	})))
	r.POST("/copy-messages", auth(imapClient(recv, func(rw http.ResponseWriter, req *http.Request, params httprouter.Params, cli *imap.Client, t statusJson) {
		status, err := cli.Status(t.Folder)
		if err != nil {
			rw.WriteHeader(http.StatusForbidden)
			return
		}
		_ = json.NewEncoder(rw).Encode(status)
	})))

	return &http.Server{
		Addr:              listen,
		Handler:           r,
		ReadTimeout:       time.Minute,
		ReadHeaderTimeout: time.Minute,
		WriteTimeout:      time.Minute,
		IdleTimeout:       time.Minute,
		MaxHeaderBytes:    2500,
	}
}

// apiError outputs a generic JSON error message
func apiError(rw http.ResponseWriter, code int, m string) {
	rw.WriteHeader(code)
	_ = json.NewEncoder(rw).Encode(map[string]string{
		"error": m,
	})
}

type IcCallback[T any] func(rw http.ResponseWriter, req *http.Request, params httprouter.Params, cli *imap.Client, t T) error

func imapClient[T any](recv Imap, cb IcCallback[T]) AuthCallback {
	return func(rw http.ResponseWriter, req *http.Request, params httprouter.Params, b AuthClaims) {
		if req.Body == nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		var t T
		if json.NewDecoder(req.Body).Decode(&t) != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		cli, err := recv.MakeClient(b.Subject)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = cb(rw, req, params, cli, t)
		if err != nil {
			log.Println("[ImapClient] Error:", err)
		}
	}
}
