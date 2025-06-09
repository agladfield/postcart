package hooks

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/agladfield/postcart/pkg/cards"
	"github.com/agladfield/postcart/pkg/postmark"
	"github.com/agladfield/postcart/pkg/shared/env"
)

var (
	errInboundNoFulLRecipients = errors.New("inbound email did not contain full recipients value")
	errInboundInvalidAddressee = errors.New("inbound email not addressed to the correct recipient")
)

const validateInboundErrFmtStr = "inbound validation err: %w"

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

func inboundHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var inboundRequest postmark.InboundData

	decodeBodyErr := postmark.DecodeToStruct(r.Body, &inboundRequest)
	if decodeBodyErr != nil {
		errorResponse(&w, r, decodeBodyErr)
		return
	}

	recipientErr := validateInboundRecipient(&inboundRequest.ToFull)
	if recipientErr != nil {
		// log error internally
		errorResponse(&w, r, recipientErr)
		// record sender bad record, 3 strikes
		return
	}

	// parse whether valid postcard request:
	// if invalid, return 403
	job, parseErr := cards.Parse(&inboundRequest)
	if parseErr != nil {
		// log error internally
		errorResponse(&w, r, parseErr)
		// record sender bad record, 3 strikes
		return
	}

	// start := time.Now()
	// // we get the database connection, record the postmark inbound email details and
	// // the parsed job parameter details
	// db := pdb.Obtain()
	// tx, txErr := db.DB.Begin()
	// if txErr != nil {
	// 	// should signal retry
	// 	retryResponse(&w, r, txErr)
	// 	return
	// }
	// defer tx.Rollback()

	// qtx := db.WithTx(tx)

	// // record inbound
	// setInboundErr := qtx.SetInboundEmail(context.Background(), pdb.SetInboundEmailParams{
	// 	ID:       inboundRequest.MessageID,
	// 	Received: time.Now().Unix(),
	// 	Email:    inboundRequest.FromFull.Email,
	// 	FromName: inboundRequest.FromFull.Name,
	// 	Subject:  inboundRequest.Subject,
	// 	Message:  inboundRequest.TextBody,
	// })

	// if setInboundErr != nil {
	// 	// signal retry
	// 	retryResponse(&w, r, setInboundErr)
	// 	return
	// }

	// // get the user (sender)
	// sender, senderErr := cards.GetSenderByEmail(inboundRequest.FromFull.Email)
	// if senderErr != nil {
	// 	// signal retry
	// 	retryResponse(&w, r, senderErr)
	// 	return
	// }

	// fmt.Println("sender:", sender)
	// fmt.Println("dur", time.Since(start))
	// retryResponse(&w, r, errors.New("retry"))
	// return

	// record as queued job
	// queuedRequest := job.ToQueueRequest(sender.ID)
	// setQueuedRequestErr := qtx.SetQueuedRequest(context.Background(), *queuedRequest)
	// if setQueuedRequestErr != nil {
	// 	// signal retry
	// 	retryResponse(&w, r, setQueuedRequestErr)
	// 	return
	// }

	// commitErr := tx.Commit()
	// if commitErr != nil {
	// 	retryResponse(&w, r, commitErr)
	// 	return
	// }

	// record the two as a transaction

	cards.AddToQueue(job)

	// pass postcardReq over to cards job channel

	// if parse is ok we add it to job queue/channel
	// and database
	// then respond with

	okResponse(&w, r)
}
