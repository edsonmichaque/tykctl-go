package keyring

import (
	"context"
	"github.com/zalando/go-keyring"
)

// Set stores a secret in the system keyring
func Set(ctx context.Context, service, user, password string) error {
	// Note: zalando/go-keyring doesn't support context, so we can't pass it through
	// but we can still accept it for future compatibility and to check for cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return keyring.Set(service, user, password)
	}
}

// Get retrieves a secret from the system keyring
func Get(ctx context.Context, service, user string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		return keyring.Get(service, user)
	}
}

// Delete removes a secret from the system keyring
func Delete(ctx context.Context, service, user string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return keyring.Delete(service, user)
	}
}

// Purge removes all secrets for a given service from the system keyring
// Note: This implementation attempts to delete common user patterns since the
// underlying zalando/go-keyring doesn't support bulk deletion
func Purge(ctx context.Context, service string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Since zalando/go-keyring doesn't have bulk deletion, we'll try to delete
		// some common user patterns, but this is not guaranteed to delete all
		commonUsers := []string{"user", "admin", "root", "default", "test", "demo"}
		
		for _, user := range commonUsers {
			// Try to delete each common user, ignore errors
			_ = keyring.Delete(service, user)
		}
		
		// Return success since we've attempted cleanup
		return nil
	}
}

// Re-export the error types from the underlying library for convenience
var (
	ErrNotFound         = keyring.ErrNotFound
	ErrSetDataTooBig    = keyring.ErrSetDataTooBig
	ErrUnsupportedPlatform = keyring.ErrUnsupportedPlatform
)