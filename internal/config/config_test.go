package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDotsDir(t *testing.T) {
	home := "/home/testuser"
	want := "/home/testuser/.dots"
	got := DotsDir(home)
	if got != want {
		t.Errorf("DotsDir(%q) = %q, want %q", home, got, want)
	}
}

func TestManifestPath(t *testing.T) {
	home := "/home/testuser"
	want := "/home/testuser/.dots/dots.yaml"
	got := ManifestPath(home)
	if got != want {
		t.Errorf("ManifestPath(%q) = %q, want %q", home, got, want)
	}
}

func TestEnsureDotsDir(t *testing.T) {
	tmpDir := t.TempDir()

	path, err := EnsureDotsDir(tmpDir)
	if err != nil {
		t.Fatalf("EnsureDotsDir failed: %v", err)
	}

	expectedPath := filepath.Join(tmpDir, DirName)
	if path != expectedPath {
		t.Errorf("EnsureDotsDir returned %q, want %q", path, expectedPath)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("dots directory not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("dots path is not a directory")
	}
}

func TestEnsureDotsDirIdempotent(t *testing.T) {
	tmpDir := t.TempDir()

	// Create twice - should not fail
	_, err := EnsureDotsDir(tmpDir)
	if err != nil {
		t.Fatalf("first EnsureDotsDir failed: %v", err)
	}

	_, err = EnsureDotsDir(tmpDir)
	if err != nil {
		t.Fatalf("second EnsureDotsDir failed: %v", err)
	}
}

func TestLoadNonExistent(t *testing.T) {
	tmpDir := t.TempDir()

	manifest, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("Load failed on non-existent manifest: %v", err)
	}

	if manifest == nil {
		t.Fatal("Load returned nil manifest")
	}
	if manifest.Files == nil {
		t.Error("Load returned manifest with nil Files slice")
	}
	if len(manifest.Files) != 0 {
		t.Errorf("Load returned manifest with %d files, want 0", len(manifest.Files))
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	_, err := EnsureDotsDir(tmpDir)
	if err != nil {
		t.Fatalf("EnsureDotsDir failed: %v", err)
	}

	original := &Manifest{
		Files: []FileEntry{
			{Source: "/home/user/.dots/.bashrc", Target: "/home/user/.bashrc"},
			{Source: "/home/user/.dots/.vimrc", Target: "/home/user/.vimrc"},
		},
	}

	if err := Save(tmpDir, original); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(loaded.Files) != len(original.Files) {
		t.Fatalf("loaded %d files, want %d", len(loaded.Files), len(original.Files))
	}

	for i, entry := range loaded.Files {
		if entry.Source != original.Files[i].Source {
			t.Errorf("Files[%d].Source = %q, want %q", i, entry.Source, original.Files[i].Source)
		}
		if entry.Target != original.Files[i].Target {
			t.Errorf("Files[%d].Target = %q, want %q", i, entry.Target, original.Files[i].Target)
		}
	}
}

func TestLoadEmptyManifest(t *testing.T) {
	tmpDir := t.TempDir()
	_, err := EnsureDotsDir(tmpDir)
	if err != nil {
		t.Fatalf("EnsureDotsDir failed: %v", err)
	}

	// Write an empty YAML file
	manifestPath := ManifestPath(tmpDir)
	if err := os.WriteFile(manifestPath, []byte(""), 0o644); err != nil {
		t.Fatalf("failed to write empty manifest: %v", err)
	}

	manifest, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("Load failed on empty manifest: %v", err)
	}

	if manifest.Files == nil {
		t.Error("Files should not be nil for empty manifest")
	}
}

func TestLoadMalformedYAML(t *testing.T) {
	tmpDir := t.TempDir()
	_, err := EnsureDotsDir(tmpDir)
	if err != nil {
		t.Fatalf("EnsureDotsDir failed: %v", err)
	}

	manifestPath := ManifestPath(tmpDir)
	if err := os.WriteFile(manifestPath, []byte("not: valid: yaml: [[["), 0o644); err != nil {
		t.Fatalf("failed to write malformed manifest: %v", err)
	}

	_, err = Load(tmpDir)
	if err == nil {
		t.Error("Load should fail on malformed YAML")
	}
}

func TestFindEntry(t *testing.T) {
	manifest := &Manifest{
		Files: []FileEntry{
			{Source: "/dots/.bashrc", Target: "/home/user/.bashrc"},
			{Source: "/dots/.vimrc", Target: "/home/user/.vimrc"},
		},
	}

	tests := []struct {
		name      string
		target    string
		wantFound bool
		wantSrc   string
	}{
		{"existing entry", "/home/user/.bashrc", true, "/dots/.bashrc"},
		{"another existing", "/home/user/.vimrc", true, "/dots/.vimrc"},
		{"non-existent", "/home/user/.zshrc", false, ""},
		{"empty target", "", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry, found := FindEntry(manifest, tt.target)
			if found != tt.wantFound {
				t.Errorf("FindEntry found = %v, want %v", found, tt.wantFound)
			}
			if found && entry.Source != tt.wantSrc {
				t.Errorf("FindEntry Source = %q, want %q", entry.Source, tt.wantSrc)
			}
		})
	}
}

func TestFindEntryEmptyManifest(t *testing.T) {
	manifest := &Manifest{Files: []FileEntry{}}
	_, found := FindEntry(manifest, "/any/path")
	if found {
		t.Error("FindEntry should return false for empty manifest")
	}
}

func TestUpsertEntryNew(t *testing.T) {
	manifest := &Manifest{Files: []FileEntry{}}

	entry := FileEntry{Source: "/dots/.bashrc", Target: "/home/.bashrc"}
	UpsertEntry(manifest, entry)

	if len(manifest.Files) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(manifest.Files))
	}
	if manifest.Files[0] != entry {
		t.Errorf("entry not added correctly")
	}
}

func TestUpsertEntryUpdate(t *testing.T) {
	manifest := &Manifest{
		Files: []FileEntry{
			{Source: "/old/source", Target: "/home/.bashrc"},
		},
	}

	newEntry := FileEntry{Source: "/new/source", Target: "/home/.bashrc"}
	UpsertEntry(manifest, newEntry)

	if len(manifest.Files) != 1 {
		t.Fatalf("expected 1 entry after update, got %d", len(manifest.Files))
	}
	if manifest.Files[0].Source != "/new/source" {
		t.Errorf("Source = %q, want %q", manifest.Files[0].Source, "/new/source")
	}
}

func TestUpsertEntryMultiple(t *testing.T) {
	manifest := &Manifest{Files: []FileEntry{}}

	entries := []FileEntry{
		{Source: "/dots/.bashrc", Target: "/home/.bashrc"},
		{Source: "/dots/.vimrc", Target: "/home/.vimrc"},
		{Source: "/dots/.zshrc", Target: "/home/.zshrc"},
	}

	for _, e := range entries {
		UpsertEntry(manifest, e)
	}

	if len(manifest.Files) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(manifest.Files))
	}
}

func TestRemoveEntry(t *testing.T) {
	manifest := &Manifest{
		Files: []FileEntry{
			{Source: "/dots/.bashrc", Target: "/home/.bashrc"},
			{Source: "/dots/.vimrc", Target: "/home/.vimrc"},
		},
	}

	removed := RemoveEntry(manifest, "/home/.bashrc")
	if !removed {
		t.Error("RemoveEntry should return true for existing entry")
	}
	if len(manifest.Files) != 1 {
		t.Fatalf("expected 1 entry after removal, got %d", len(manifest.Files))
	}
	if manifest.Files[0].Target != "/home/.vimrc" {
		t.Error("wrong entry removed")
	}
}

func TestRemoveEntryNonExistent(t *testing.T) {
	manifest := &Manifest{
		Files: []FileEntry{
			{Source: "/dots/.bashrc", Target: "/home/.bashrc"},
		},
	}

	removed := RemoveEntry(manifest, "/home/.nonexistent")
	if removed {
		t.Error("RemoveEntry should return false for non-existent entry")
	}
	if len(manifest.Files) != 1 {
		t.Error("manifest should be unchanged")
	}
}

func TestRemoveEntryEmptyManifest(t *testing.T) {
	manifest := &Manifest{Files: []FileEntry{}}

	removed := RemoveEntry(manifest, "/any/path")
	if removed {
		t.Error("RemoveEntry should return false for empty manifest")
	}
}

func TestRemoveEntryLastItem(t *testing.T) {
	manifest := &Manifest{
		Files: []FileEntry{
			{Source: "/dots/.bashrc", Target: "/home/.bashrc"},
		},
	}

	removed := RemoveEntry(manifest, "/home/.bashrc")
	if !removed {
		t.Error("RemoveEntry should return true")
	}
	if len(manifest.Files) != 0 {
		t.Errorf("expected empty manifest, got %d entries", len(manifest.Files))
	}
}
