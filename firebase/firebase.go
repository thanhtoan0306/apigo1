package firebase

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var (
	FirestoreClient *firestore.Client
	AuthClient      *auth.Client
)

// InitializeFirebase initializes Firebase Admin SDK
func InitializeFirebase(ctx context.Context) error {
	// Get Firebase credentials from environment variable or file
	credentialsPath := os.Getenv("FIREBASE_CREDENTIALS")
	
	var app *firebase.App
	var err error
	var config *firebase.Config

	if credentialsPath != "" {
		// Use credentials file if provided
		opt := option.WithCredentialsFile(credentialsPath)
		app, err = firebase.NewApp(ctx, nil, opt)
	} else if jsonCreds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_JSON"); jsonCreds != "" {
		// Use JSON credentials from environment variable
		// Parse JSON to get project_id
		var creds map[string]interface{}
		if err := json.Unmarshal([]byte(jsonCreds), &creds); err != nil {
			log.Printf("Error parsing JSON credentials: %v", err)
			return err
		}

		// Extract project_id from credentials
		projectID, ok := creds["project_id"].(string)
		if !ok || projectID == "" {
			log.Printf("Error: project_id not found in credentials JSON")
			return errors.New("project_id is required in credentials JSON")
		}

		log.Printf("Initializing Firebase with project_id: %s", projectID)

		// Set project ID in config
		config = &firebase.Config{
			ProjectID: projectID,
		}

		opt := option.WithCredentialsJSON([]byte(jsonCreds))
		app, err = firebase.NewApp(ctx, config, opt)
		if err != nil {
			log.Printf("Error initializing Firebase with JSON credentials: %v", err)
			return err
		}
	} else {
		// Use Application Default Credentials (ADC)
		// This works if running on GCP or if GOOGLE_APPLICATION_CREDENTIALS is set
		app, err = firebase.NewApp(ctx, nil)
	}

	if err != nil {
		log.Printf("Error initializing Firebase: %v", err)
		return err
	}

	// Initialize Firestore
	FirestoreClient, err = app.Firestore(ctx)
	if err != nil {
		return err
	}

	// Initialize Auth (optional, for future use)
	AuthClient, err = app.Auth(ctx)
	if err != nil {
		log.Printf("Warning: Failed to initialize Auth client: %v", err)
	}

	log.Println("Firebase initialized successfully")
	return nil
}

// Close closes Firebase connections
func Close() error {
	if FirestoreClient != nil {
		return FirestoreClient.Close()
	}
	return nil
}

