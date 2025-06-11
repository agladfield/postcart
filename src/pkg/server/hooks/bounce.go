package hooks

import (
	"net/http"

	"github.com/agladfield/postcart/pkg/cards"
	"github.com/agladfield/postcart/pkg/jdb"
	"github.com/agladfield/postcart/pkg/postmark"
)

func bounceHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var bounceData postmark.BounceData

	decodeBodyErr := postmark.DecodeToStruct(r.Body, &bounceData)
	if decodeBodyErr != nil {
		errorResponse(&w, r, decodeBodyErr)
		return
	}
	// record as bounce
	jdb.RecordBounce()
	// block email from receiving emails
	jdb.BlockRecipient(bounceData.Email)
	jdb.RecordBlockedRecipient()
	// block sender from sending emails
	if bounceData.Metadata != nil {
		senderEmail, senderEmailExists := bounceData.Metadata["sender_email"]
		if senderEmailExists {
			cards.BlockSender(senderEmail)
		}
	}

	okResponse(&w, r)
}

// © Arthur Gladfield
