package hooks

import (
	"net/http"

	"github.com/agladfield/postcart/pkg/cards"
	"github.com/agladfield/postcart/pkg/jdb"
	"github.com/agladfield/postcart/pkg/postmark"
)

func spamComplaintHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var spamComplaintData postmark.SpamComplaintData

	decodeBodyErr := postmark.DecodeToStruct(r.Body, &spamComplaintData)
	if decodeBodyErr != nil {
		errorResponse(&w, r, decodeBodyErr)
		return
	}

	// record spam complaint
	jdb.RecordSpamComplaint()
	// block the recipient from recieving emails
	jdb.BlockRecipient(spamComplaintData.Email)
	jdb.RecordBlockedSender()
	// block the sender from sending emails
	if spamComplaintData.Metadata != nil {
		senderEmail, senderEmailExists := spamComplaintData.Metadata["sender_email"]
		if senderEmailExists {
			cards.BlockSender(senderEmail)
		}
	}

	okResponse(&w, r)
}

// © Arthur Gladfield
