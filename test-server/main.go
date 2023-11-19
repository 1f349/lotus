package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/1f349/mjwt"
	"github.com/1f349/mjwt/auth"
	"github.com/1f349/mjwt/claims"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/rs/cors"
)

func main() {
	log.Println("Starting test server")
	signer, err := mjwt.NewMJwtSignerFromFileOrCreate("Test SSO Service", "private.key.local", rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	go clientServer()
	go ssoServer(signer)
	go apiServer(signer)
	done := make(chan struct{})
	<-done
}

func clientServer() {
	r := http.NewServeMux()
	r.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		http.ServeFile(rw, req, "index.html")
	})
	log.Println("[Client Server]", http.ListenAndServe(":2020", r))
}

func ssoServer(signer mjwt.Signer) {
	r := http.NewServeMux()
	r.HandleFunc("/popup", func(w http.ResponseWriter, r *http.Request) {
		ps := claims.NewPermStorage()
		ps.Set("mail:inbox=admin@localhost")
		accessToken, err := signer.GenerateJwt("81b99bd7-bf74-4cc2-9133-80ed2393dfe6", uuid.NewString(), jwt.ClaimStrings{"d0555671-df9d-42d0-a4d6-94b694251f0b"}, 15*time.Minute, auth.AccessTokenClaims{
			Perms: ps,
		})
		if err != nil {
			http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="en">
<head>
  <title>Test SSO Service</title>
  <script>
    let loginData = {
      target: "http://localhost:2020",
      userinfo: {
        "aud": "d0555671-df9d-42d0-a4d6-94b694251f0b",
        "email": "admin@localhost",
        "email_verified": true,
        "name": "Admin",
        "preferred_username": "admin",
        "sub": "81b99bd7-bf74-4cc2-9133-80ed2393dfe6",
        "picture": "http://localhost:5173/1f349.svg",
        "updated_at": 0
      },
      tokens: {
        access: "%s",
        refresh: "%s",
      },
    };
    window.addEventListener("load", function () {
      setTimeout(function() {
        window.opener.postMessage(loginData, loginData.target);
      },2000);
    });
  </script>
</head>
<body>
<header>
  <h1>Test SSO Service</h1>
</header>
<main id="mainBody">Loading...</main>
</body>
</html>
`, accessToken, "")
	})
	log.Println("[SSO Server]", http.ListenAndServe(":9090", r))
}

var serveApiCors = cors.New(cors.Options{
	AllowedOrigins: []string{"*"}, // allow all origins for api requests
	AllowedHeaders: []string{"Content-Type", "Authorization"},
	AllowedMethods: []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
	},
	AllowCredentials: true,
})

func apiServer(verify mjwt.Verifier) {
	r := http.NewServeMux()
	r.Handle("/v1/violet/route", hasPerm(verify, "violet:route", func(rw http.ResponseWriter, req *http.Request) {
		m := make([]map[string]any, 0, 40)
		for i := 0; i < 20; i++ {
			m = append(m, map[string]any{
				"src":    uuid.NewString() + ".example.com",
				"dst":    "127.0.0.1:8080",
				"desc":   "This is a test description",
				"flags":  181,
				"active": true,
			})
		}
		for i := 0; i < 20; i++ {
			m = append(m, map[string]any{
				"src":    uuid.NewString() + ".example.org",
				"dst":    "127.0.0.1:8085",
				"desc":   "This is a test description",
				"flags":  17,
				"active": true,
			})
		}
		json.NewEncoder(rw).Encode(m)
	}))
	r.Handle("/v1/violet/redirect", hasPerm(verify, "violet:redirect", func(rw http.ResponseWriter, req *http.Request) {
		m := make([]map[string]any, 0, 40)
		for i := 0; i < 20; i++ {
			m = append(m, map[string]any{
				"src":    uuid.NewString() + ".example.com",
				"dst":    "test1.example.com",
				"desc":   "This is a test description",
				"flags":  1,
				"code":   308,
				"active": true,
			})
		}
		for i := 0; i < 20; i++ {
			m = append(m, map[string]any{
				"src":    uuid.NewString() + ".example.org",
				"dst":    "test2.example.org",
				"desc":   "This is a test description",
				"flags":  3,
				"code":   307,
				"active": true,
			})
		}
		json.NewEncoder(rw).Encode(m)
	}))
	r.Handle("/v1/orchid/owned", hasPerm(verify, "orchid:cert", func(rw http.ResponseWriter, req *http.Request) {
		m := make(map[int]any, 41)
		for i := 0; i < 20; i++ {
			u := uuid.NewString()
			m[i] = map[string]any{
				"id":           i + 1,
				"auto_renew":   true,
				"active":       true,
				"renewing":     false,
				"renew_failed": false,
				"not_after":    "2024-02-06T11:52:05Z",
				"updated_at":   "2023-11-08T07:32:08Z",
				"domains": []string{
					u + ".example.com",
					"*." + u + ".example.com",
				},
			}
		}
		for i := 0; i < 20; i++ {
			u := uuid.NewString()
			m[i+20] = map[string]any{
				"id":           i + 21,
				"auto_renew":   false,
				"active":       false,
				"renewing":     false,
				"renew_failed": false,
				"not_after":    "2024-02-06T11:52:05Z",
				"updated_at":   "2023-11-08T07:32:08Z",
				"domains": []string{
					u + ".example.org",
					"*." + u + ".example.org",
				},
			}
		}
		u := uuid.NewString()
		m[40] = map[string]any{
			"id":           41,
			"auto_renew":   false,
			"active":       false,
			"renewing":     false,
			"renew_failed": true,
			"not_after":    "2024-02-06T11:52:05Z",
			"updated_at":   "2023-11-08T07:32:08Z",
			"domains": []string{
				u + ".example.org",
				"*." + u + ".example.org",
			},
		}
		json.NewEncoder(rw).Encode(m)
	}))
	r.Handle("/v1/sites", hasPerm(verify, "sites:manage", func(rw http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			defer req.Body.Close()
			dec := json.NewDecoder(req.Body)
			var m map[string]string
			if err := dec.Decode(&m); err != nil {
				http.Error(rw, err.Error(), http.StatusBadRequest)
				return
			}
			switch m["submit"] {
			case "secret":
				rw.WriteHeader(http.StatusOK)
				fmt.Fprintf(rw, "{\"secret\":\"%s\"}\n", uuid.NewString())
				return
			case "delete-branch":
				rw.WriteHeader(http.StatusOK)
			}
			return
		}
		m := make([]any, 0, 40)
		for i := 0; i < 20; i++ {
			m = append(m, map[string]any{
				"domain":   uuid.NewString() + ".example.com",
				"branches": []string{"", "beta"},
			})
		}
		for i := 0; i < 20; i++ {
			m = append(m, map[string]any{
				"domain":   uuid.NewString() + ".example.org",
				"branches": []string{"", "alpha"},
			})
		}
		json.NewEncoder(rw).Encode(m)
	}))

	logger := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		log.Println("[API Server]", req.URL.String())
		r.ServeHTTP(rw, req)
	})
	log.Println("[API Server]", http.ListenAndServe(":9095", serveApiCors.Handler(logger)))
}

func hasPerm(verify mjwt.Verifier, perm string, next func(rw http.ResponseWriter, req *http.Request)) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		a := req.Header.Get("Authorization")
		if !strings.HasPrefix(a, "Bearer ") {
			http.Error(rw, "Missing bearer authorization", http.StatusForbidden)
			return
		}
		_, b, err := mjwt.ExtractClaims[auth.AccessTokenClaims](verify, a[len("Bearer "):])
		if err != nil {
			http.Error(rw, "Invalid token", http.StatusForbidden)
			log.Println("Invalid token:", err)
			return
		}
		if !b.Claims.Perms.Has("violet:route") {
			http.Error(rw, "Missing permission", http.StatusForbidden)
			return
		}
		next(rw, req)
	})
}
