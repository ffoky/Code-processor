package middleware

import (
	"context"
	"http_server/usecases/service"
	"net/http"
)

func AuthMiddleware(sessionService *service.SessionService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			sid := authHeader[7:]
			session, err := sessionService.Get(sid)
			if err != nil || session == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "session", session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
