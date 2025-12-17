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
	jsonCreds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_JSON")
	
	log.Printf("Checking Firebase credentials...")
	log.Printf("FIREBASE_CREDENTIALS: %s", credentialsPath)
	if jsonCreds != "" {
		log.Printf("GOOGLE_APPLICATION_CREDENTIALS_JSON: found (%d chars)", len(jsonCreds))
	} else {
		log.Printf("GOOGLE_APPLICATION_CREDENTIALS_JSON: not set")
	}
	
	var app *firebase.App
	var err error
	var config *firebase.Config

	if credentialsPath != "" {
		// Use credentials file if provided
		log.Printf("Using credentials file: %s", credentialsPath)
		opt := option.WithCredentialsFile(credentialsPath)
		app, err = firebase.NewApp(ctx, nil, opt)
	} else if jsonCreds != "" {
		// Use JSON credentials from environment variable
		log.Printf("Using JSON credentials from environment variable")
		// Parse JSON to get project_id
		var creds map[string]interface{}
		if err := json.Unmarshal([]byte(jsonCreds), &creds); err != nil {
			log.Printf("Error parsing JSON credentials: %v", err)
			log.Printf("JSON preview (first 200 chars): %s", jsonCreds[:min(200, len(jsonCreds))])
			return err
		}

		log.Printf("JSON parsed successfully, found %d keys", len(creds))

		// Extract project_id from credentials
		projectID, ok := creds["project_id"].(string)
		if !ok || projectID == "" {
			log.Printf("Error: project_id not found in credentials JSON")
			log.Printf("Available keys: %v", getKeys(creds))
			return errors.New("project_id is required in credentials JSON")
		}

		log.Printf("Found project_id: %s", projectID)
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
		log.Printf("Using Application Default Credentials")
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

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

