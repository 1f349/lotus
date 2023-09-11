package api

import (
	"encoding/json"
	"errors"
	"fmt"
	postfixLookup "github.com/1f349/lotus/postfix-lookup"
	"github.com/1f349/lotus/sendmail"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

var defaultPostfixLookup = postfixLookup.NewPostfixLookup().Lookup
var timeNow = time.Now

// MessageSender is the internal handler for `POST /message` requests
// the access token is already validated at this point
func MessageSender(send Smtp) func(rw http.ResponseWriter, req *http.Request, params httprouter.Params, b AuthClaims) {
	return func(rw http.ResponseWriter, req *http.Request, params httprouter.Params, b AuthClaims) {
		// check body exists
		if req.Body == nil {
			apiError(rw, http.StatusBadRequest, "Missing request body")
			return
		}

		// parse json body
		var j sendmail.Json
		err := json.NewDecoder(req.Body).Decode(&j)
		if err != nil {
			apiError(rw, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		// prepare the mail for sending
		mail, err := j.PrepareMail(timeNow())
		if err != nil {
			apiError(rw, http.StatusBadRequest, fmt.Sprintf("Invalid mail message: %s", err))
			return
		}

		// this looks up the underlying account for the sender alias
		lookup, err := defaultPostfixLookup(mail.From.Address)

		// the alias does not exist
		if errors.Is(err, postfixLookup.ErrInvalidAlias) {
			apiError(rw, http.StatusBadRequest, "Invalid sender alias")
			return
		}

		// the alias lookup failed to run
		if err != nil {
			apiError(rw, http.StatusInternalServerError, "Sender alias lookup failed")
			return
		}

		// the alias does not match the logged-in user
		if lookup != b.Subject {
			apiError(rw, http.StatusBadRequest, "User does not own sender alias")
			return
		}

		// try sending the mail
		if err := send.Send(mail); err != nil {
			apiError(rw, http.StatusInternalServerError, "Failed to send mail")
			log.Printf("Failed to send mail: %#v: %s\n", mail, err)
			return
		}

		rw.WriteHeader(http.StatusAccepted)
	}
}
