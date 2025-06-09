package hooks

import (
	"net/http"

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

	// prevent sender from sending emails again for at least 1 week
	// db.
	// record spam complaint
	// block the recipient from recieving emails
	// block the sender from sending emails

	okResponse(&w, r)
}
