package auth

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	firebaseAuth "firebase.google.com/go/v4/auth"
)

func NewFirebaseAuthClient(ctx context.Context) (*firebaseAuth.Client, error) {
	config := &firebase.Config{
		ProjectID:        "quote-project-a2672-dev",
		ServiceAccountID: "firebase-adminsdk-fbsvc@quote-project-a2672-dev.iam.gserviceaccount.com",
	}

	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize firebase app: %w", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize firebase auth client: %w", err)
	}

	return client, nil
}
