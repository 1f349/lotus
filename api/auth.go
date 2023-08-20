package api

import (
	"crypto/subtle"
	"github.com/1f349/violet/utils"
	"github.com/MrMelon54/mjwt"
	"github.com/MrMelon54/mjwt/auth"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type AuthClaims mjwt.BaseTypeClaims[auth.AccessTokenClaims]

type AuthCallback func(rw http.ResponseWriter, req *http.Request, params httprouter.Params, b AuthClaims)

type authChecker struct {
	verify mjwt.Verifier
	aud    string
	cb     AuthCallback
}

func (a *authChecker) Handle(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Get bearer token
	bearer := utils.GetBearer(req)
	if bearer == "" {
		apiError(rw, http.StatusForbidden, "Missing bearer token")
		return
	}

	// Read claims from mjwt
	_, b, err := mjwt.ExtractClaims[auth.AccessTokenClaims](a.verify, bearer)
	if err != nil {
		apiError(rw, http.StatusForbidden, "Invalid token")
		return
	}

	var validAud bool
	for _, i := range b.Audience {
		if subtle.ConstantTimeCompare([]byte(i), []byte(a.aud)) == 1 {
			validAud = true
		}
	}
	if !validAud {
		apiError(rw, http.StatusForbidden, "Invalid audience claim")
		return
	}

	a.cb(rw, req, params, AuthClaims(b))
}

// CheckAuth validates the bearer token against a mjwt.Verifier and returns an
// error message or continues to the next handler
func CheckAuth(verify mjwt.Verifier, aud string) func(cb AuthCallback) httprouter.Handle {
	return func(cb AuthCallback) httprouter.Handle {
		return (&authChecker{
			verify: verify,
			aud:    aud,
			cb:     cb,
		}).Handle
	}
}
