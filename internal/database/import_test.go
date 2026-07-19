package database

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestImportSeedFileReturnsReadError(t *testing.T) {
	missingPath := filepath.Join(t.TempDir(), "missing.json")

	err := ImportSeedFile(nil, missingPath)

	if err == nil {
		t.Fatal("expected a missing seed file error")
	}
	if !strings.Contains(err.Error(), "failed to read seed file") {
		t.Fatalf("expected a contextual seed file error, got %q", err)
	}
}
