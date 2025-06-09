// Package env wraps consistent access to environment variables for configuring
// the program
package env

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	googleCredsPath           = "GCP_CRED_PATH" // #nosec G101
	useAIKey                  = "USE_AI"
	postmarkServerUsernameKey = "POSTMARK_SERVER_USER"
	postmarkServerPasswordKey = "POSTMARK_SERVER_PASS"
	postmarkServerTokenKey    = "POSTMARK_SERVER_TOKEN"
	postmarkInboundEmailKey   = "POSTMARK_INBOUND_EMAIL"
	useDBKey                  = "USE_DB"
	installFontsKey           = "INSTALL_FONTS"
)

var (
	errNoGoogleCredsPath        = errors.New("no gcp creds path provided")
	errNoPostmarkServerUsername = errors.New("no postmark server username provided")
	errNoPostmarkServerPassword = errors.New("no postmark server password provided")
	errNoPostmarkServerToken    = errors.New("no postmark server token provided")
	errNoPostmarkInboundEmail   = errors.New("no postmark inbound email provided")
)

var (
	gcpCredsPath           string
	useAI                  bool
	postmarkServerUsername string
	postmarkServerPassword string
	postmarkServerToken    string
	postmarkInboundEmail   string
	installFonts           bool
)

func Configure() error {
	var errs []error

	// essential for running the postmark integration
	postmarkServerUsername = os.Getenv(postmarkServerUsernameKey)
	if postmarkServerUsername == "" {
		errs = append(errs, errNoPostmarkServerUsername)
	}
	postmarkServerPassword = os.Getenv(postmarkServerPasswordKey)
	if postmarkServerPassword == "" {
		errs = append(errs, errNoPostmarkServerPassword)
	}
	postmarkServerToken = os.Getenv(postmarkServerTokenKey)
	if postmarkServerToken == "" {
		errs = append(errs, errNoPostmarkServerToken)
	}
	postmarkInboundEmail = os.Getenv(postmarkInboundEmailKey)
	if postmarkInboundEmail == "" {
		errs = append(errs, errNoPostmarkInboundEmail)
	}

	// optional for improved operation
	useAI = strings.ToLower(os.Getenv(useAIKey)) != "false"
	if useAI {
		gcpCredsPath = os.Getenv(gcpCredsPath)
		if gcpCredsPath == "" {
			errs = append(errs, errNoGoogleCredsPath)
		}
	}

	installFonts = strings.ToLower(os.Getenv(installFontsKey)) != "false"

	joinedErrs := errors.Join(errs...)
	if joinedErrs != nil {
		return fmt.Errorf("environment configuration err: %w", joinedErrs)
	}

	return nil
}

func GCPCredsPath() string {
	return gcpCredsPath
}

func UseAI() bool {
	return useAI
}

func PostmarkServerUsername() string {
	return postmarkServerUsername
}

func PostmarkServerPassword() string {
	return postmarkServerPassword
}

func PostmarkServerToken() string {
	return postmarkServerToken
}

func PostmarkInboundEmail() string {
	return postmarkInboundEmail
}

func InstallFonts() bool {
	return installFonts
}

// We need to have an email domain for deliveries and no-reply
func EmailDomain() string {
	return ""
}
