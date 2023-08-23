package api

import (
	"encoding/json"
	"github.com/1f349/lotus/imap"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

func SetupApiServer(listen string, auth func(callback AuthCallback) httprouter.Handle, send Smtp, recv Imap) *http.Server {
	r := httprouter.New()

	// === ACCOUNT ===
	r.GET("/account", auth(func(rw http.ResponseWriter, req *http.Request, params httprouter.Params, b AuthClaims) {
		// TODO(melon): find users aliases and other account data
	}))

	// === SMTP ===
	r.POST("/message", auth(MessageSender(send)))

	// === IMAP ===
	type statusJson struct {
		Folder string `json:"folder"`
	}
	r.GET("/status", auth(imapClient(recv, func(rw http.ResponseWriter, req *http.Request, params httprouter.Params, cli *imap.Client, t statusJson) {
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
		_ = json.NewEncoder(rw).Encode(messages)
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

type IcCallback[T any] func(rw http.ResponseWriter, req *http.Request, params httprouter.Params, cli *imap.Client, t T)

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
		cb(rw, req, params, cli, t)
	}
}
