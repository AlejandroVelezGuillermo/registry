package database

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
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
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("expected a wrapped missing-file error, got %v", err)
	}
}

func TestReadSeedFileReturnsParseError(t *testing.T) {
	invalidPath := filepath.Join(t.TempDir(), "invalid.json")
	if err := os.WriteFile(invalidPath, []byte(`[{"id":42}]`), 0o600); err != nil {
		t.Fatalf("write invalid seed file: %v", err)
	}

	_, err := readSeedFile(invalidPath)

	if err == nil {
		t.Fatal("expected an invalid seed file error")
	}
	if !strings.Contains(err.Error(), "failed to parse seed JSON") {
		t.Fatalf("expected a contextual parse error, got %q", err)
	}
	var typeError *json.UnmarshalTypeError
	if !errors.As(err, &typeError) {
		t.Fatalf("expected a wrapped JSON type error, got %v", err)
	}
}

func TestSeedFileIncludesCurrentXquikServer(t *testing.T) {
	seedPath := filepath.Join("..", "..", "data", "seed_2025_05_16.json")
	servers, err := readSeedFile(seedPath)
	if err != nil {
		t.Fatalf("read seed file: %v", err)
	}

	var matches []int
	for index, server := range servers {
		if server.Name == "com.xquik/mcp" {
			matches = append(matches, index)
		}
	}
	if len(matches) != 1 {
		t.Fatalf("expected 1 Xquik server, got %d", len(matches))
	}

	server := servers[matches[0]]
	if server.VersionDetail.Version != "2.6.0" {
		t.Fatalf("expected Xquik version 2.6.0, got %q", server.VersionDetail.Version)
	}
	if server.Repository.ID != "1168120544" {
		t.Fatalf("unexpected Xquik repository ID %q", server.Repository.ID)
	}
	if len(server.Packages) != 1 || server.Packages[0].Version != "2.6.0" {
		t.Fatalf("unexpected Xquik package metadata: %#v", server.Packages)
	}
	if len(server.Remotes) != 1 ||
		server.Remotes[0].TransportType != "streamable-http" ||
		server.Remotes[0].URL != "https://xquik.com/mcp" {
		t.Fatalf("unexpected Xquik remote metadata: %#v", server.Remotes)
	}
}
