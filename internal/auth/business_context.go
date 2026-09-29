package auth

import (
	"context"
	"net/http"
)

const businessIDKey contextKey = "businessID"

type BusinessIDProvider interface {
	GetBusinessIDByFirebaseUID(
		ctx context.Context,
		firebaseUID string,
	) (int64, error)
}

type BusinessContext struct {
	businessIDProvider BusinessIDProvider
}

func NewBusinessContext(
	businessIDProvider BusinessIDProvider,
) *BusinessContext {
	return &BusinessContext{
		businessIDProvider: businessIDProvider,
	}
}

func (bc *BusinessContext) Resolve(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		firebaseUID, ok := FirebaseUIDFromContext(r.Context())
		if !ok || firebaseUID == "" {
			http.Error(w, "firebase uid not found", http.StatusUnauthorized)
			return
		}

		businessID, err := bc.businessIDProvider.GetBusinessIDByFirebaseUID(
			r.Context(),
			firebaseUID,
		)
		if err != nil {
			http.Error(w, "business not found", http.StatusForbidden)
			return
		}

		// Store the business ID in the request context
		ctx := context.WithValue(
			r.Context(),
			businessIDKey,
			businessID,
		)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func BusinessIDFromContext(ctx context.Context) (int64, bool) {
	businessID, ok := ctx.Value(businessIDKey).(int64)
	return businessID, ok
}
