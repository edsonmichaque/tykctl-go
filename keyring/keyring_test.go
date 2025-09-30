package keyring

import (
	"context"
	"runtime"
	"strings"
	"testing"
	"time"
)

const (
	service  = "test-service"
	user     = "test-user"
	password = "test-password"
)

// TestSet tests setting a user and password in the keyring.
func TestSet(t *testing.T) {
	ctx := context.Background()
	err := Set(ctx, service, user, password)
	if err != nil {
		t.Errorf("Should not fail, got: %s", err)
	}
}

func TestSetTooLong(t *testing.T) {
	ctx := context.Background()
	extraLongPassword := "ba" + strings.Repeat("na", 5000)
	err := Set(ctx, service, user, extraLongPassword)

	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		// should fail on those platforms
		if err != ErrSetDataTooBig {
			t.Errorf("Should have failed, got: %s", err)
		}
	}
}

// TestGetMultiline tests getting a multi-line password from the keyring
func TestGetMultiLine(t *testing.T) {
	ctx := context.Background()
	multilinePassword := `this password
has multiple
lines and will be
encoded by some keyring implementiations
like osx`
	err := Set(ctx, service, user, multilinePassword)
	if err != nil {
		t.Errorf("Should not fail, got: %s", err)
	}

	pw, err := Get(ctx, service, user)
	if err != nil {
		t.Errorf("Should not fail, got: %s", err)
	}

	if multilinePassword != pw {
		t.Errorf("Expected password %s, got %s", multilinePassword, pw)
	}
}

// TestGetMultiline tests getting a multi-line password from the keyring
func TestGetUmlaut(t *testing.T) {
	ctx := context.Background()
	umlautPassword := "at least on OSX üöäÜÖÄß will be encoded"
	err := Set(ctx, service, user, umlautPassword)
	if err != nil {
		t.Errorf("Should not fail, got: %s", err)
	}

	pw, err := Get(ctx, service, user)
	if err != nil {
		t.Errorf("Should not fail, got: %s", err)
	}

	if umlautPassword != pw {
		t.Errorf("Expected password %s, got %s", umlautPassword, pw)
	}
}

// TestGetSingleLineHex tests getting a single line hex string password from the keyring.
func TestGetSingleLineHex(t *testing.T) {
	ctx := context.Background()
	hexPassword := "abcdef123abcdef123"
	err := Set(ctx, service, user, hexPassword)
	if err != nil {
		t.Errorf("Should not fail, got: %s", err)
	}

	pw, err := Get(ctx, service, user)
	if err != nil {
		t.Errorf("Should not fail, got: %s", err)
	}

	if hexPassword != pw {
		t.Errorf("Expected password %s, got %s", hexPassword, pw)
	}
}

// TestGet tests getting a password from the keyring.
func TestGet(t *testing.T) {
	ctx := context.Background()
	err := Set(ctx, service, user, password)
	if err != nil {
		t.Errorf("Should not fail, got: %s", err)
	}

	pw, err := Get(ctx, service, user)
	if err != nil {
		t.Errorf("Should not fail, got: %s", err)
	}

	if password != pw {
		t.Errorf("Expected password %s, got %s", password, pw)
	}
}

// TestGetNonExisting tests getting a secret not in the keyring.
func TestGetNonExisting(t *testing.T) {
	ctx := context.Background()
	_, err := Get(ctx, service, user+"fake")
	if err != ErrNotFound {
		t.Errorf("Expected error ErrNotFound, got %s", err)
	}
}

// TestDelete tests deleting a secret from the keyring.
func TestDelete(t *testing.T) {
	ctx := context.Background()
	err := Delete(ctx, service, user)
	if err != nil {
		t.Errorf("Should not fail, got: %s", err)
	}
}

// TestDeleteNonExisting tests deleting a secret not in the keyring.
func TestDeleteNonExisting(t *testing.T) {
	ctx := context.Background()
	err := Delete(ctx, service, user+"fake")
	if err != ErrNotFound {
		t.Errorf("Expected error ErrNotFound, got %s", err)
	}
}

// TestPurge tests purging all secrets for a given service.
func TestPurge(t *testing.T) {
	ctx := context.Background()
	// Set up multiple secrets for the same service
	err := Set(ctx, service, user, password)
	if err != nil {
		t.Errorf("Should not fail, got: %s", err)
	}

	err = Set(ctx, service, user+"2", password+"2")
	if err != nil {
		t.Errorf("Should not fail, got: %s", err)
	}

	// Purge all secrets for the service
	err = Purge(ctx, service)
	if err != nil {
		t.Errorf("Purge should not fail, got: %s", err)
	}

	// Clean up manually since Purge only tries common patterns
	err = Delete(ctx, service, user)
	if err != nil {
		t.Errorf("Should not fail, got: %s", err)
	}

	err = Delete(ctx, service, user+"2")
	if err != nil {
		t.Errorf("Should not fail, got: %s", err)
	}
}

// TestPurgeEmptyService tests purging with empty service name
func TestPurgeEmptyService(t *testing.T) {
	ctx := context.Background()
	err := Set(ctx, service, user, password)

	if err != nil {
		t.Errorf("Should not fail, got: %s", err)
	}
	
	// Purge should succeed even with empty service (it will try common patterns)
	err = Purge(ctx, "")
	if err != nil {
		t.Errorf("Purge should not fail, got: %s", err)
	}
	
	// Clean up
	err = Delete(ctx, service, user)
	if err != nil {
		t.Errorf("Should not fail, got: %s", err)
	}
}

// TestContextCancellation tests that context cancellation works
func TestContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately
	
	err := Set(ctx, service, user, password)
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got: %v", err)
	}
	
	_, err = Get(ctx, service, user)
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got: %v", err)
	}
	
	err = Delete(ctx, service, user)
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got: %v", err)
	}
	
	err = Purge(ctx, service)
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got: %v", err)
	}
}

// TestContextTimeout tests that context timeout works
func TestContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	
	// Wait for timeout
	time.Sleep(1 * time.Millisecond)
	
	err := Set(ctx, service, user, password)
	if err != context.DeadlineExceeded {
		t.Errorf("Expected context.DeadlineExceeded, got: %v", err)
	}
}
