package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	postfixLookup "github.com/1f349/lotus/postfix-lookup"
	"github.com/1f349/lotus/sendmail"
	"github.com/MrMelon54/mjwt/auth"
	"github.com/MrMelon54/mjwt/claims"
	"github.com/emersion/go-message/mail"
	"github.com/golang-jwt/jwt/v4"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
	"time"
)

func init() {
	defaultPostfixLookup = func(key string) (string, error) {
		switch key {
		case "noreply@example.com", "admin@example.com":
			return "admin@example.com", nil
		case "user@example.com":
			return "user@example.com", nil
		}
		return "", postfixLookup.ErrInvalidAlias
	}
	timeNow = func() time.Time {
		return time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	}
}

type fakeSmtp struct {
	from *mail.Address
	body []byte
}

func (f *fakeSmtp) Send(mail *sendmail.Mail) error {
	// remove the Bcc header line
	s := bufio.NewScanner(bytes.NewReader(mail.Body))
	b := new(bytes.Buffer)
	b.Grow(len(mail.Body))
	for s.Scan() {
		a := s.Text()
		if strings.HasPrefix(a, "Bcc:") {
			continue
		}
		b.WriteString(a + "\r\n")
	}
	if err := s.Err(); err != nil {
		return err
	}

	// check values are the same
	if mail.From.String() != f.from.String() {
		return fmt.Errorf("test fail: invalid from address")
	}
	if !slices.Equal(b.Bytes(), f.body) {
		return fmt.Errorf("test fail: invalid message body")
	}
	return nil
}

type fakeFailedSmtp struct{}

func (f *fakeFailedSmtp) Send(mail *sendmail.Mail) error {
	return errors.New("sending failed")
}

var messageSenderTestData = []struct {
	req    func() (*http.Request, error)
	smtp   Smtp
	claims AuthClaims
	status int
	output string
}{
	{
		req: func() (*http.Request, error) {
			return http.NewRequest(http.MethodPost, "https://api.example.com/v1/mail/message", nil)
		},
		smtp:   &fakeSmtp{},
		claims: makeFakeAuthClaims("admin@example.com"),
		status: http.StatusBadRequest,
		output: "Missing request body",
	},
	{
		req: func() (*http.Request, error) {
			return http.NewRequest(http.MethodPost, "https://api.example.com/v1/mail/message", strings.NewReader(`{`))
		},
		smtp:   &fakeSmtp{},
		claims: makeFakeAuthClaims("admin@example.com"),
		status: http.StatusBadRequest,
		output: "Invalid JSON body",
	},
	{
		req: func() (*http.Request, error) {
			j, err := json.Marshal(sendmail.Json{
				From:     "noreply2@example.com",
				ReplyTo:  "admin@example.com",
				To:       "user@example.com",
				Subject:  "Test Subject",
				BodyType: "plain",
				Body:     "Plain text",
			})
			if err != nil {
				return nil, err
			}
			return http.NewRequest(http.MethodPost, "https://api.example.com/v1/mail/message", bytes.NewReader(j))
		},
		smtp:   &fakeSmtp{},
		claims: makeFakeAuthClaims("admin@example.com"),
		status: http.StatusBadRequest,
		output: "Invalid sender alias",
	},
	{
		req: func() (*http.Request, error) {
			j, err := json.Marshal(sendmail.Json{
				From:     "user@example.com",
				ReplyTo:  "admin@example.com",
				To:       "user@example.com",
				Subject:  "Test Subject",
				BodyType: "plain",
				Body:     "Plain text",
			})
			if err != nil {
				return nil, err
			}
			return http.NewRequest(http.MethodPost, "https://api.example.com/v1/mail/message", bytes.NewReader(j))
		},
		smtp:   &fakeSmtp{},
		claims: makeFakeAuthClaims("admin@example.com"),
		status: http.StatusBadRequest,
		output: "User does not own sender alias",
	},
	{
		req: func() (*http.Request, error) {
			j, err := json.Marshal(sendmail.Json{
				From:     "noreply@example.com, user2@example.com",
				ReplyTo:  "admin@example.com",
				To:       "user@example.com",
				Subject:  "Test Subject",
				BodyType: "plain",
				Body:     "Plain text",
			})
			if err != nil {
				return nil, err
			}
			return http.NewRequest(http.MethodPost, "https://api.example.com/v1/mail/message", bytes.NewReader(j))
		},
		smtp:   &fakeSmtp{},
		claims: makeFakeAuthClaims("admin@example.com"),
		status: http.StatusBadRequest,
		output: "Invalid mail message: multiple from addresses",
	},
	{
		req: func() (*http.Request, error) {
			j, err := json.Marshal(sendmail.Json{
				From:     "noreply@example.com",
				ReplyTo:  "admin@example.com",
				To:       "user@example.com",
				Subject:  "Test Subject",
				BodyType: "no",
				Body:     "Plain text",
			})
			if err != nil {
				return nil, err
			}
			return http.NewRequest(http.MethodPost, "https://api.example.com/v1/mail/message", bytes.NewReader(j))
		},
		smtp:   &fakeSmtp{},
		claims: makeFakeAuthClaims("admin@example.com"),
		status: http.StatusBadRequest,
		output: "Invalid mail message: invalid body type",
	},
	{
		req: func() (*http.Request, error) {
			j, err := json.Marshal(sendmail.Json{
				From:     "noreply@example.com",
				ReplyTo:  "admin@example.com",
				To:       "a <user@example.com",
				Subject:  "Test Subject",
				BodyType: "no",
				Body:     "Plain text",
			})
			if err != nil {
				return nil, err
			}
			return http.NewRequest(http.MethodPost, "https://api.example.com/v1/mail/message", bytes.NewReader(j))
		},
		smtp:   &fakeSmtp{},
		claims: makeFakeAuthClaims("admin@example.com"),
		status: http.StatusBadRequest,
		output: "Invalid mail message: mail: unclosed angle-addr",
	},
	{
		req: func() (*http.Request, error) {
			j, err := json.Marshal(sendmail.Json{
				From:     "noreply@example.com",
				ReplyTo:  "admin@example.com",
				To:       "a <user>",
				Subject:  "Test Subject",
				BodyType: "no",
				Body:     "Plain text",
			})
			if err != nil {
				return nil, err
			}
			return http.NewRequest(http.MethodPost, "https://api.example.com/v1/mail/message", bytes.NewReader(j))
		},
		smtp:   &fakeSmtp{},
		claims: makeFakeAuthClaims("admin@example.com"),
		status: http.StatusBadRequest,
		output: "Invalid mail message: mail: missing @ in addr-spec",
	},
	{
		req: func() (*http.Request, error) {
			j, err := json.Marshal(sendmail.Json{
				From:     "noreply@example.com",
				ReplyTo:  "admin@example.com",
				To:       "user@example.com",
				Subject:  "Test Subject",
				BodyType: "plain",
				Body:     "Plain text",
			})
			if err != nil {
				return nil, err
			}
			return http.NewRequest(http.MethodPost, "https://api.example.com/v1/mail/message", bytes.NewReader(j))
		},
		smtp:   &fakeFailedSmtp{},
		claims: makeFakeAuthClaims("admin@example.com"),
		status: http.StatusInternalServerError,
		output: "Failed to send mail",
	},
	{
		req: func() (*http.Request, error) {
			j, err := json.Marshal(sendmail.Json{
				From:     "noreply@example.com",
				ReplyTo:  "admin@example.com",
				To:       "user@example.com",
				Cc:       "user2@example.com",
				Bcc:      "user3@example.com",
				Subject:  "Test Subject",
				BodyType: "plain",
				Body:     "Some plain text",
			})
			if err != nil {
				return nil, err
			}
			return http.NewRequest(http.MethodPost, "https://api.example.com/v1/mail/message", bytes.NewReader(j))
		},
		smtp: &fakeSmtp{
			from: &mail.Address{Address: "noreply@example.com"},
			body: []byte("Mime-Version: 1.0\r\n" +
				"Content-Type: text/plain; charset=utf-8\r\n" +
				"Cc: <user2@example.com>\r\n" +
				"To: <user@example.com>\r\n" +
				"Reply-To: <admin@example.com>\r\n" +
				"From: <noreply@example.com>\r\n" +
				"Subject: Test Subject\r\n" +
				"Date: Sat, 01 Jan 2000 00:00:00 +0000\r\n" +
				"\r\n" +
				"Some plain text\r\n"),
		},
		claims: makeFakeAuthClaims("admin@example.com"),
		status: http.StatusAccepted,
		output: "",
	},
}

func makeFakeAuthClaims(subject string) AuthClaims {
	return struct {
		jwt.RegisteredClaims
		ClaimType string
		Claims    auth.AccessTokenClaims
	}{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "Test",
			Subject:  subject,
			Audience: jwt.ClaimStrings{"mail.example.com"},
		},
		ClaimType: "access-token",
		Claims:    auth.AccessTokenClaims{Perms: claims.NewPermStorage()},
	}
}

func TestMessageSender(t *testing.T) {
	for _, i := range messageSenderTestData {
		rec := httptest.NewRecorder()
		req, err := i.req()
		assert.NoError(t, err)
		MessageSender(i.smtp)(rec, req, httprouter.Params{}, i.claims)

		res := rec.Result()
		assert.Equal(t, i.status, res.StatusCode)
		assert.NotNil(t, res.Body)
		all, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		if i.output == "" {
			assert.Equal(t, "", string(all))
		} else {
			assert.Equal(t, "{\"error\":\""+i.output+"\"}\n", string(all))
		}
	}
}
