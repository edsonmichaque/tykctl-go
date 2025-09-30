package keyring

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Example demonstrates basic usage of the keyring package with global functions.
func Example() {
	ctx := context.Background()
	service := "example-service"
	user := "example-user"
	password := "example-password"

	// Store a secret
	fmt.Println("Storing secret...")
	err := Set(ctx, service, user, password)
	if err != nil {
		log.Printf("Failed to store secret: %v", err)
		return
	}
	fmt.Println("Secret stored successfully")

	// Retrieve the secret
	fmt.Println("Retrieving secret...")
	retrievedPassword, err := Get(ctx, service, user)
	if err != nil {
		log.Printf("Failed to retrieve secret: %v", err)
		return
	}
	fmt.Printf("Retrieved password: %s\n", retrievedPassword)

	// Verify the password matches
	if retrievedPassword == password {
		fmt.Println("Password verification successful")
	} else {
		fmt.Println("Password verification failed")
	}

	// Delete the secret
	fmt.Println("Deleting secret...")
	err = Delete(ctx, service, user)
	if err != nil {
		log.Printf("Failed to delete secret: %v", err)
		return
	}
	fmt.Println("Secret deleted successfully")

	// Try to retrieve the deleted secret
	fmt.Println("Attempting to retrieve deleted secret...")
	_, err = Get(ctx, service, user)
	if err == ErrNotFound {
		fmt.Println("Secret not found (as expected)")
	} else if err != nil {
		log.Printf("Unexpected error: %v", err)
	}
}

// ExampleMultipleSecrets demonstrates storing multiple secrets.
func ExampleMultipleSecrets() {
	ctx := context.Background()
	
	// Store multiple secrets
	secrets := map[string]string{
		"user1": "password1",
		"user2": "password2",
		"user3": "password3",
	}

	service := "example-service"
	fmt.Println("Storing multiple secrets...")
	for user, pass := range secrets {
		err := Set(ctx, service, user, pass)
		if err != nil {
			log.Printf("Failed to store secret for %s: %v", user, err)
			return
		}
		fmt.Printf("Stored secret for %s\n", user)
	}

	// Retrieve and verify secrets
	fmt.Println("Retrieving secrets...")
	for user, expectedPass := range secrets {
		pass, err := Get(ctx, service, user)
		if err != nil {
			log.Printf("Failed to retrieve secret for %s: %v", user, err)
			continue
		}
		if pass == expectedPass {
			fmt.Printf("Secret for %s verified successfully\n", user)
		} else {
			fmt.Printf("Secret for %s verification failed\n", user)
		}
	}

	// Delete individual secrets
	fmt.Println("Deleting individual secrets...")
	for user := range secrets {
		err := Delete(ctx, service, user)
		if err != nil {
			log.Printf("Failed to delete secret for %s: %v", user, err)
		} else {
			fmt.Printf("Deleted secret for %s\n", user)
		}
	}
}

// ExamplePurge demonstrates the Purge functionality.
func ExamplePurge() {
	ctx := context.Background()
	service := "example-service"

	// Store multiple secrets for the same service
	fmt.Println("Storing multiple secrets...")
	secrets := map[string]string{
		"user1": "password1",
		"user2": "password2",
		"user3": "password3",
	}

	for user, pass := range secrets {
		err := Set(ctx, service, user, pass)
		if err != nil {
			log.Printf("Failed to store secret for %s: %v", user, err)
			return
		}
		fmt.Printf("Stored secret for %s\n", user)
	}

	// Try to purge all secrets for the service
	fmt.Println("Attempting to purge all secrets for service...")
	err := Purge(ctx, service)
	if err != nil {
		log.Printf("Failed to purge all secrets: %v", err)
		return
	}
	fmt.Println("Purge completed (attempted to clean common patterns)")

	// Since Purge only tries common patterns, clean up individually
	fmt.Println("Cleaning up remaining secrets individually...")
	for user := range secrets {
		err := Delete(ctx, service, user)
		if err != nil {
			log.Printf("Failed to delete secret for %s: %v", user, err)
		} else {
			fmt.Printf("Deleted secret for %s\n", user)
		}
	}

	// Verify all secrets are deleted
	fmt.Println("Verifying deletion...")
	for user := range secrets {
		_, err := Get(ctx, service, user)
		if err == ErrNotFound {
			fmt.Printf("Secret for %s not found (as expected)\n", user)
		} else if err != nil {
			log.Printf("Unexpected error retrieving secret for %s: %v", user, err)
		}
	}
}

// ExampleContext demonstrates context usage with timeout and cancellation.
func ExampleContext() {
	service := "example-service"
	user := "example-user"
	password := "example-password"

	// Example with timeout
	fmt.Println("Example with timeout...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := Set(ctx, service, user, password)
	if err != nil {
		log.Printf("Failed to store secret with timeout: %v", err)
		return
	}
	fmt.Println("Secret stored successfully with timeout")

	// Example with cancellation
	fmt.Println("Example with cancellation...")
	ctx2, cancel2 := context.WithCancel(context.Background())
	cancel2() // Cancel immediately

	err = Set(ctx2, service, user+"2", password+"2")
	if err == context.Canceled {
		fmt.Println("Operation was cancelled as expected")
	} else if err != nil {
		log.Printf("Unexpected error: %v", err)
	}

	// Clean up
	ctx3 := context.Background()
	Delete(ctx3, service, user)
}
