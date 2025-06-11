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
	googleCredsPathKey        = "GCP_CRED_PATH"
	googleProjectKey          = "GCP_PROJECT"
	googleBucketKey           = "GCP_BUCKET"
	useAIKey                  = "USE_AI"
	postmarkServerUsernameKey = "POSTMARK_SERVER_USER"
	postmarkServerPasswordKey = "POSTMARK_SERVER_PASS"
	postmarkServerTokenKey    = "POSTMARK_SERVER_TOKEN"
	postmarkInboundEmailKey   = "POSTMARK_INBOUND_EMAIL"
	postmarkEmailDomainKey    = "POSTMARK_EMAIL_DOMAIN"
	installFontsKey           = "INSTALL_FONTS"
	allowAttachmentsKey       = "ALLOW_ATTACHMENTS"
)

var (
	errNoGoogleCredsPath        = errors.New("no gcp creds path provided")
	errNoGoogleProject          = errors.New("no gcp project provided")
	errNoPostmarkServerUsername = errors.New("no postmark server username provided")
	errNoPostmarkServerPassword = errors.New("no postmark server password provided")
	errNoPostmarkServerToken    = errors.New("no postmark server token provided")
	errNoPostmarkInboundEmail   = errors.New("no postmark inbound email provided")
	errNoPostmarkEmailDomain    = errors.New("no postmark email domain provided")
)

var (
	gcpCredsPath           string
	gcpProject             string
	gcpBucket              string
	useAI                  bool
	postmarkServerUsername string
	postmarkServerPassword string
	postmarkServerToken    string
	postmarkInboundEmail   string
	postmarkEmailDomain    string
	installFonts           bool
	allowAttachments       bool
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
	postmarkEmailDomain = os.Getenv(postmarkEmailDomainKey)
	if postmarkEmailDomain == "" {
		errs = append(errs, errNoPostmarkEmailDomain)
	}

	// optional for improved operation
	gcpCredsPath = os.Getenv(googleCredsPathKey)
	gcpProject = os.Getenv(googleProjectKey)
	useAI = strings.ToLower(os.Getenv(useAIKey)) != "false"
	if useAI {
		if gcpCredsPath == "" {
			errs = append(errs, errNoGoogleCredsPath)
		}
		if gcpProject == "" {
			errs = append(errs, errNoGoogleProject)
		}
	}

	gcpBucket = os.Getenv(googleBucketKey)

	installFonts = strings.ToLower(os.Getenv(installFontsKey)) == "true"
	allowAttachments = strings.ToLower(os.Getenv(installFontsKey)) != "false"

	joinedErrs := errors.Join(errs...)
	if joinedErrs != nil {
		return fmt.Errorf("environment configuration err: %w", joinedErrs)
	}

	return nil
}

func GCPCredsPath() string {
	return gcpCredsPath
}

func GCPProject() string {
	return gcpProject
}

func GCPBucket() string {
	return gcpBucket
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

func PostmarkEmailDomain() string {
	return postmarkEmailDomain
}

func InstallFonts() bool {
	return installFonts
}

func AllowAttachments() bool {
	return allowAttachments
}

// © Arthur Gladfield
