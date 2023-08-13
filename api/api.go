package api

import (
	"encoding/json"
	"github.com/1f349/primrose/imap"
	"github.com/1f349/primrose/smtp"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

type Conf struct {
	Listen string `yaml:"listen"`
}

func SetupApiServer(conf Conf, send *smtp.Smtp, recv *imap.Imap) *http.Server {
	r := httprouter.New()

	// smtp
	r.POST("/message", func(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
		// check body exists
		if req.Body == nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		// parse json body
		var j smtp.Json
		err := json.NewDecoder(req.Body).Decode(&j)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		mail, err := j.PrepareMail()
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if send.Send(mail) != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusAccepted)
	})

	return &http.Server{
		Addr:              conf.Listen,
		Handler:           r,
		ReadTimeout:       time.Minute,
		ReadHeaderTimeout: time.Minute,
		WriteTimeout:      time.Minute,
		IdleTimeout:       time.Minute,
		MaxHeaderBytes:    2500,
	}
}
