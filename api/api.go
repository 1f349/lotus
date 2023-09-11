package api

import (
	"bytes"
	"encoding/json"
	"github.com/1f349/lotus/imap"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{}

func SetupApiServer(listen string, auth *AuthChecker, send Smtp, recv Imap) *http.Server {
	r := httprouter.New()

	// === ACCOUNT ===
	r.GET("/identities", auth.Middleware(func(rw http.ResponseWriter, req *http.Request, params httprouter.Params, b AuthClaims) {
		// TODO(melon): find users aliases and other account data
	}))

	// === SMTP ===
	r.POST("/smtp", auth.Middleware(MessageSender(send)))

	r.GET("/imap", func(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
		// upgrade to websocket conn and defer close
		c, err := upgrader.Upgrade(rw, req, nil)
		if err != nil {
			log.Println("[Imap] Failed to upgrade to websocket:", err)
			return
		}
		defer c.Close()

		// set a really short deadline to refuse unauthenticated clients
		deadline := time.Now().Add(5 * time.Second)
		_ = c.SetReadDeadline(deadline)
		_ = c.SetWriteDeadline(deadline)

		// close on all possible errors, assume we are being attacked
		mt, msg, err := c.ReadMessage()
		if err != nil {
			return
		}
		if mt != websocket.TextMessage {
			return
		}
		if len(msg) >= 2000 {
			return
		}

		// parse token from message
		var tokenMsg struct {
			Token string `json:"token"`
		}
		dec := json.NewDecoder(bytes.NewReader(msg))
		dec.DisallowUnknownFields()
		err = dec.Decode(&tokenMsg)
		if err != nil {
			_ = c.WriteJSON(map[string]string{"error": "Authentication missing"})
			return
		}

		// get a "possible" auth token value
		// exit on empty token value
		if tokenMsg.Token == "" {
			_ = c.WriteJSON(map[string]string{"error": "Authentication missing"})
			return
		}

		// check the token
		authUser, err := auth.Check(tokenMsg.Token)
		if err != nil {
			_ = c.WriteJSON(map[string]string{"error": "Authentication invalid"})
			return
		}

		// open imap client
		client, err := recv.MakeClient(authUser.Subject)
		if err != nil {
			_ = c.WriteJSON(map[string]string{"error": "Making client failed"})
			return
		}

		// auth was ok
		err = c.WriteJSON(map[string]string{"auth": "ok"})
		if err != nil {
			return
		}

		for {
			// authenticated users get longer to reply
			// a simple ping/pong setup bypasses this
			d := time.Now().Add(5 * time.Minute)
			_ = c.SetReadDeadline(d)
			_ = c.SetWriteDeadline(d)

			// read incoming message
			var m struct {
				Action string   `json:"action"`
				Args   []string `json:"args"`
			}
			err := c.ReadJSON(&m)
			if err != nil {
				_ = c.WriteJSON(map[string]string{"error": "Invalid input"})
				return
			}

			// handle action
			j, err := client.HandleWS(m.Action, m.Args)
			if err != nil {
				_ = c.WriteJSON(map[string]string{"error": "Action failed"})
				return
			}

			// write outgoing message
			err = c.WriteJSON(j)
			if err != nil {
				_ = c.WriteJSON(map[string]string{"error": "Invalid output"})
				return
			}
		}
	})

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
