package hooks

import (
	"net/http"

	"github.com/agladfield/postcart/pkg/postmark"
)

func deliveredHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var deliveredData postmark.DeliveredData

	decodeBodyErr := postmark.DecodeToStruct(r.Body, &deliveredData)
	if decodeBodyErr != nil {
		errorResponse(&w, r, decodeBodyErr)
		return
	}

	// record as delivered
	// update job status

	okResponse(&w, r)
}
