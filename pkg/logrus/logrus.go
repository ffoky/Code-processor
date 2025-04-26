package logrus

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)
}

func LogFullRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)
		}

		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		headers := ""
		for k, v := range r.Header {
			headers += k + ": " + sliceToString(v) + "\n"
		}

		logrus.WithFields(logrus.Fields{
			"method":  r.Method,
			"url":     r.URL.String(),
			"headers": headers,
			"body":    string(bodyBytes),
		}).Info("Incoming HTTP request")

		next.ServeHTTP(w, r)
	})
}

func sliceToString(s []string) string {
	var out string
	for i, v := range s {
		if i > 0 {
			out += ", "
		}
		out += v
	}
	return out
}
