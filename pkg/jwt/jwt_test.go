package jwt

import (
	"testing"
	"time"

	"github.com/koer/koer-module/pkg/config"
)

func TestGenerateAndValidateAccessToken(t *testing.T) {
	mgr := NewManager(config.JWTConfig{
		SecretKey:     "supersecret",
		AccessExpiry:  time.Hour,
		RefreshExpiry: 24 * time.Hour,
		Issuer:        "test",
	})

	token, err := mgr.GenerateAccessToken("user123", "admin")
	if err != nil {
		t.Fatalf("generate access token: %v", err)
	}

	claims, err := mgr.ValidateToken(token)
	if err != nil {
		t.Fatalf("validate token: %v", err)
	}

	if claims.UserID != "user123" {
		t.Errorf("expected user123, got %s", claims.UserID)
	}
	if claims.Role != "admin" {
		t.Errorf("expected admin, got %s", claims.Role)
	}
}

func TestGenerateAndValidateRefreshToken(t *testing.T) {
	mgr := NewManager(config.JWTConfig{
		SecretKey:     "supersecret",
		AccessExpiry:  time.Hour,
		RefreshExpiry: 24 * time.Hour,
		Issuer:        "test",
	})

	token, err := mgr.GenerateRefreshToken("user123")
	if err != nil {
		t.Fatalf("generate refresh token: %v", err)
	}

	claims, err := mgr.ValidateToken(token)
	if err != nil {
		t.Fatalf("validate token: %v", err)
	}

	if claims.UserID != "user123" {
		t.Errorf("expected user123, got %s", claims.UserID)
	}
}
