package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/agladfield/postcart/pkg/jdb"
)

func statsHandler(w http.ResponseWriter, r *http.Request) {
	stats := jdb.GetStats()
	statsBytes, bytesErr := json.Marshal(stats)
	w.Header().Set("Content-Type", "application/json")
	if bytesErr != nil {
		w.WriteHeader(200)
		fmt.Fprintf(w, "{\"error\":%q}", bytesErr)
		return
	}

	w.WriteHeader(200)
	fmt.Fprint(w, string(statsBytes))
}
