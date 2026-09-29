package auth

import (
	"context"
	"log"
	"net/http"
	"strings"

	firebaseAuth "firebase.google.com/go/v4/auth"
)

type contextKey string

const firebaseUIDKey contextKey = "firebaseUID"

type Middleware struct {
	authClient *firebaseAuth.Client
}

func NewMiddleware(authClient *firebaseAuth.Client) *Middleware {
	return &Middleware{
		authClient: authClient,
	}
}

func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "authorization header is required", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
			return
		}

		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		if idToken == "" {
			http.Error(w, "token is required", http.StatusUnauthorized)
			return
		}

		token, err := m.authClient.VerifyIDToken(r.Context(), idToken)
		if err != nil {
			log.Printf("VerifyIDToken failed: %v", err)

			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(
			r.Context(),
			firebaseUIDKey,
			token.UID,
		)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func FirebaseUIDFromContext(ctx context.Context) (string, bool) {
	uid, ok := ctx.Value(firebaseUIDKey).(string)
	return uid, ok
}
