// Package hooks wraps the different Postmark webhook routes and their logic
package hooks

import (
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/agladfield/postcart/pkg/jdb"
	"github.com/agladfield/postcart/pkg/shared/env"
)

var (
	postmarkServerUser     string
	postmarkServerPassword string
)

var (
	errBadPostmarkServerUser     = errors.New("empty postmark server username provided")
	errBadPostmarkServerPassword = errors.New("empty postmark server password provided")
)

const (
	hooksErrFmrStr       = "webhooks err: %w"
	hooksConfigErrFmtStr = "webhooks config err: %w"
)

func Configure() error {
	// double check hook basic auth credentials
	postmarkServerUser = env.PostmarkServerUsername()
	if postmarkServerUser == "" {
		return fmt.Errorf(hooksConfigErrFmtStr, errBadPostmarkServerUser)
	}

	postmarkServerPassword = env.PostmarkServerPassword()
	if postmarkServerPassword == "" {
		return fmt.Errorf(hooksConfigErrFmtStr, errBadPostmarkServerPassword)
	}

	return nil
}

func Routes() {
	http.HandleFunc("/hooks/bounce", postmarkOnlyRoute(bounceHandler))
	http.HandleFunc("/hooks/inbound", postmarkOnlyRoute(inboundHandler))
	http.HandleFunc("/hooks/delivered", postmarkOnlyRoute(deliveredHandler))
	http.HandleFunc("/hooks/spam", postmarkOnlyRoute(spamComplaintHandler))
}

// postmarkOnlyRoute wraps the webhook http handlers to only allow postmark to make the requests.
// I did not in my example here but you could also whitelist the postmark webhook ips found here:
// https://postmarkapp.com/support/article/800-ips-for-firewalls#webhooks
func postmarkOnlyRoute(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(postmarkServerUser))
			expectedPasswordHash := sha256.Sum256([]byte(postmarkServerPassword))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

const (
	postmarkAcceptStatus = http.StatusOK
	postmarkRejectStatus = http.StatusForbidden
	postmarkRetryStatus  = http.StatusInternalServerError
)

func printRequestOutput(r *http.Request, status int) {
	fmt.Printf("[%d] %s\n", status, r.URL.Path)
	if status > 200 {
		jdb.RecordRejection()
	}
}

// okResponse sends a 200 response to postmark
func okResponse(w *http.ResponseWriter, r *http.Request) {
	(*w).WriteHeader(postmarkAcceptStatus)
	fmt.Fprintf(*w, "ok")
	printRequestOutput(r, postmarkAcceptStatus)
}

// errorResponse sends a 403 response to postmark
// according to postmark docs 200 and 403 are the correct response codes
// to prevent retries
func errorResponse(w *http.ResponseWriter, r *http.Request, err error) {
	(*w).WriteHeader(postmarkRejectStatus)
	fmt.Fprintf(*w, "error")
	log.Printf("an error occured: %v\n", err)
	printRequestOutput(r, postmarkRejectStatus)
}

func retryResponse(w *http.ResponseWriter, r *http.Request, err error) {
	(*w).WriteHeader(postmarkRetryStatus)
	fmt.Fprintf(*w, "error-retry")
	log.Printf("an error occured: %v\n", err)
	printRequestOutput(r, postmarkRetryStatus)
	jdb.RecordRetry()
}

// © Arthur Gladfield
