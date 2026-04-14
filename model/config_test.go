package models

import (
	"os"
	"testing"
)

func TestLoadConfigUsesEnvironmentValues(t *testing.T) {
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalWD)
	})

	t.Setenv("PORT", "9090")
	t.Setenv("GIN_MODE", "debug")
	t.Setenv("MONGO_URI", "mongodb://example:27017")
	t.Setenv("MONGO_DATABASE", "movies_test")
	t.Setenv("MONGO_COLLECTION", "entries")

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if config.Port != "9090" {
		t.Fatalf("expected port 9090, got %s", config.Port)
	}

	if config.GinMode != "debug" {
		t.Fatalf("expected gin mode debug, got %s", config.GinMode)
	}

	if config.MongoURI != "mongodb://example:27017" {
		t.Fatalf("expected mongo uri mongodb://example:27017, got %s", config.MongoURI)
	}

	if config.MongoDatabase != "movies_test" {
		t.Fatalf("expected mongo database movies_test, got %s", config.MongoDatabase)
	}

	if config.MongoCollection != "entries" {
		t.Fatalf("expected mongo collection entries, got %s", config.MongoCollection)
	}
}

func TestLoadConfigRequiresMongoURI(t *testing.T) {
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalWD)
	})

	t.Setenv("PORT", "")
	t.Setenv("GIN_MODE", "")
	t.Setenv("MONGO_URI", "")
	t.Setenv("MONGO_DATABASE", "")
	t.Setenv("MONGO_COLLECTION", "")

	_, err = LoadConfig()
	if err == nil {
		t.Fatal("expected an error when MONGO_URI is missing")
	}
}
