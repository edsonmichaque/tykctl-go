// Package keyring provides a simple wrapper around the zalando/go-keyring library
// for storing and retrieving secrets from the system keyring on different platforms.
//
// The keyring package supports:
//   - macOS: Uses the macOS Keychain via the security command-line tool
//   - Linux/Unix: Uses the Secret Service API via D-Bus
//   - Windows: Uses the Windows Credential Manager
//   - Fallback: Returns errors for unsupported platforms
//
// Basic usage:
//
//	import (
//		"context"
//		"github.com/edsonmichaque/tykctl-go/keyring"
//	)
//
//	ctx := context.Background()
//
//	// Store a secret
//	err := keyring.Set(ctx, "my-service", "my-user", "my-secret")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Retrieve a secret
//	secret, err := keyring.Get(ctx, "my-service", "my-user")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(secret)
//
//	// Delete a secret
//	err = keyring.Delete(ctx, "my-service", "my-user")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Purge all secrets for a service
//	err = keyring.Purge(ctx, "my-service")
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Context support allows for cancellation and timeouts:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	
//	secret, err := keyring.Get(ctx, "service", "user")
//	if err != nil {
//		// Handle timeout or cancellation
//	}
//
// The package automatically detects the platform and uses the appropriate
// keyring implementation. If the platform is not supported, operations will
// return ErrUnsupportedPlatform.
//
// This package is a simple wrapper around github.com/zalando/go-keyring.
package keyring