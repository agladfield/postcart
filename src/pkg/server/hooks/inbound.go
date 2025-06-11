package hooks

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/agladfield/postcart/pkg/cards"
	"github.com/agladfield/postcart/pkg/jdb"
	"github.com/agladfield/postcart/pkg/postmark"
	"github.com/agladfield/postcart/pkg/shared/env"
)

var (
	errInboundNoFulLRecipients   = errors.New("inbound email did not contain full recipients value")
	errInboundInvalidAddressee   = errors.New("inbound email not addressed to the correct recipient")
	errInboundMessageTooLong     = errors.New("inbound email message was too long")
	errInboundSenderIsBlocked    = errors.New("inbound email sender is blocked from sending more emails")
	errInboundRecipientIsBlocked = errors.New("inbound email intended recipient is blocked from receiving emails")
)

const validateInboundErrFmtStr = "inbound validation err: %w"

const maxEmailsPerSender = 3

func validateInboundRecipient(fullRecipients *[]postmark.EmailAddressFull) error {
	if fullRecipients == nil {
		return fmt.Errorf(validateInboundErrFmtStr, errInboundNoFulLRecipients)
	}
	for _, recipient := range *fullRecipients {
		if recipient.Email == env.PostmarkInboundEmail() {
			return nil
		}
	}

	return fmt.Errorf(validateInboundErrFmtStr, errInboundInvalidAddressee)
}

const maxMsgLen = 2048

func inboundHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var inboundRequest postmark.InboundData

	decodeBodyErr := postmark.DecodeToStruct(r.Body, &inboundRequest)
	if decodeBodyErr != nil {
		errorResponse(&w, r, decodeBodyErr)
		return
	}
	jdb.RecordInbound()

	if jdb.IsSenderBlocked(inboundRequest.FromFull.Email) {
		errorResponse(&w, r, errInboundSenderIsBlocked)
		// try blocking sender again as they must have gotten
		// past inbound rule
		cards.BlockSender(inboundRequest.FromFull.Email)
		return
	}

	senderCount := jdb.IncrementSender(inboundRequest.FromFull.Email)

	recipientErr := validateInboundRecipient(&inboundRequest.ToFull)
	if recipientErr != nil {
		// log error internally
		errorResponse(&w, r, recipientErr)
		return
	}

	if len(inboundRequest.TextBody) > maxMsgLen {
		errorResponse(&w, r, errInboundMessageTooLong)
		return
	}

	// parse whether valid postcard request:
	// if invalid, return 403
	job, parseErr := cards.Parse(&inboundRequest)
	if parseErr != nil {
		errorResponse(&w, r, parseErr)
		return
	}

	if jdb.IsRecipientBlocked(job.To.Email) {
		errorResponse(&w, r, errInboundRecipientIsBlocked)
		return
	}
	if senderCount > maxEmailsPerSender {
		cards.BlockSender(inboundRequest.FromFull.Email)
	}

	// save email/job
	addToQueueErr := cards.AddToQueue(job)
	if addToQueueErr != nil {
		retryResponse(&w, r, addToQueueErr)
		return
	}

	okResponse(&w, r)
}

// © Arthur Gladfield
