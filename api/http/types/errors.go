package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"http_server/repository"
	"net/http"
)

func ProcessError(w http.ResponseWriter, err error, resp any) {
	if err != nil {
		if errors.Is(err, repository.NotFound) {
			http.Error(w, "Id not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		}
		return
	}

	if resp != nil {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		}
	} else {
		_, err := fmt.Fprintln(w, "Not found")
		if err != nil {
			return
		}
	}
}
