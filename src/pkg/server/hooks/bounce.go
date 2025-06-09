package hooks

import (
	"net/http"

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
	// block email from receiving emails
	// block sender from sending emails
	// update job status

	okResponse(&w, r)
}
