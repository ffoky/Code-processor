package types

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"http_server/repository"
	"net/http"
)

func ProcessError(w http.ResponseWriter, err error, resp any, statusCode int) {
	logrus.Debug("Debug", statusCode)

	if err != nil {
		if errors.Is(err, repository.NotFound) {
			http.Error(w, "Id not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		}
		logrus.Errorf("Error %v", err)
		return
	}

	if resp != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		logrus.WithFields(logrus.Fields{
			"status": statusCode,
		}).Debug("debug info")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Not found", http.StatusNotFound)
	}
}
