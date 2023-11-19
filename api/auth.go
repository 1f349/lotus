package api

import (
	"crypto/subtle"
	"errors"
	"github.com/1f349/mjwt"
	"github.com/1f349/mjwt/auth"
	"github.com/1f349/violet/utils"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var (
	ErrInvalidToken    = errors.New("invalid token")
	ErrInvalidAudClaim = errors.New("invalid audience claim")
)

type AuthClaims mjwt.BaseTypeClaims[auth.AccessTokenClaims]

type AuthCallback func(rw http.ResponseWriter, req *http.Request, params httprouter.Params, b AuthClaims)

// AuthChecker validates the bearer token against a mjwt.Verifier and returns an
// error message or continues to the next handler
type AuthChecker struct {
	Verify mjwt.Verifier
	Aud    string
}

// Middleware is a httprouter.Handle layer to authenticate requests
func (a *AuthChecker) Middleware(cb AuthCallback) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
		// Get bearer token
		bearer := utils.GetBearer(req)
		if bearer == "" {
			apiError(rw, http.StatusForbidden, "Missing bearer token")
			return
		}

		b, err := a.Check(bearer)
		switch {
		case errors.Is(err, ErrInvalidToken):
			apiError(rw, http.StatusForbidden, "Invalid token")
			return
		case errors.Is(err, ErrInvalidAudClaim):
			apiError(rw, http.StatusForbidden, "Invalid audience claim")
			return
		case err != nil:
			apiError(rw, http.StatusForbidden, "Unknown error")
			return
		}

		cb(rw, req, params, b)
	}
}

// Check takes a token and validates whether it is verified and contains the
// correct audience claim
func (a *AuthChecker) Check(token string) (AuthClaims, error) {
	// Read claims from mjwt
	_, b, err := mjwt.ExtractClaims[auth.AccessTokenClaims](a.Verify, token)
	if err != nil {
		return AuthClaims{}, ErrInvalidToken
	}

	// Check aud value
	var validAud bool
	for _, i := range b.Audience {
		if subtle.ConstantTimeCompare([]byte(i), []byte(a.Aud)) == 1 {
			validAud = true
		}
	}
	if !validAud {
		return AuthClaims{}, ErrInvalidAudClaim
	}

	return AuthClaims(b), nil
}
