package api

import (
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

	r.Handle(http.MethodConnect, "/imap", func(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
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

		// get a "possible" auth token value
		authToken := string(msg)

		// wait for authToken or error
		// exit on empty reply
		if authToken == "" {
			return
		}

		// check the token
		authUser, err := auth.Check(authToken)
		if err != nil {
			// exit on error
			return
		}

		_ = authUser

		client, err := recv.MakeClient(authUser.Subject)
		if err != nil {
			_ = c.WriteJSON(map[string]string{"Error": "Making client failed"})
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
				// errors should close the connection
				return
			}

			// handle action
			j, err := client.HandleWS(m.Action, m.Args)
			if err != nil {
				// errors should close the connection
				return
			}

			// write outgoing message
			err = c.WriteJSON(j)
			if err != nil {
				// errors should close the connection
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
