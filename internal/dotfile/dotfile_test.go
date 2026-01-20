package dotfile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/subcode-labs/dots/internal/config"
)

func TestHomeDir(t *testing.T) {
	home, err := HomeDir()
	if err != nil {
		t.Fatalf("HomeDir failed: %v", err)
	}
	if home == "" {
		t.Error("HomeDir returned empty string")
	}
}

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

func TestInit(t *testing.T) {
	tmpDir := t.TempDir()

	dotsPath, err := Init(tmpDir)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	expectedPath := filepath.Join(tmpDir, ".dots")
	if dotsPath != expectedPath {
		t.Errorf("Init returned %q, want %q", dotsPath, expectedPath)
	}

	// Check directory was created
	info, err := os.Stat(dotsPath)
	if err != nil {
		t.Fatalf("dots directory not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("dots path is not a directory")
	}

	// Check manifest was created
	manifestPath := filepath.Join(dotsPath, "dots.yaml")
	if _, err := os.Stat(manifestPath); err != nil {
		t.Errorf("manifest not created: %v", err)
	}
}

func TestInitIdempotent(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := Init(tmpDir)
	if err != nil {
		t.Fatalf("first Init failed: %v", err)
	}

	_, err = Init(tmpDir)
	if err != nil {
		t.Fatalf("second Init should succeed: %v", err)
	}
}

func TestCopyFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create source file
	srcPath := filepath.Join(tmpDir, "source.txt")
	content := []byte("test content\nline two\n")
	if err := os.WriteFile(srcPath, content, 0o644); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	// Copy file
	dstPath := filepath.Join(tmpDir, "dest.txt")
	if err := CopyFile(srcPath, dstPath); err != nil {
		t.Fatalf("CopyFile failed: %v", err)
	}

	// Verify content
	got, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("failed to read destination: %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("content mismatch: got %q, want %q", got, content)
	}
}

func TestCopyFilePreservesPermissions(t *testing.T) {
	tmpDir := t.TempDir()

	srcPath := filepath.Join(tmpDir, "executable.sh")
	if err := os.WriteFile(srcPath, []byte("#!/bin/bash\necho hi"), 0o755); err != nil {
		t.Fatalf("failed to create source: %v", err)
	}

	dstPath := filepath.Join(tmpDir, "copy.sh")
	if err := CopyFile(srcPath, dstPath); err != nil {
		t.Fatalf("CopyFile failed: %v", err)
	}

	srcInfo, _ := os.Stat(srcPath)
	dstInfo, _ := os.Stat(dstPath)

	if srcInfo.Mode().Perm() != dstInfo.Mode().Perm() {
		t.Errorf("permissions not preserved: src %v, dst %v", srcInfo.Mode().Perm(), dstInfo.Mode().Perm())
	}
}

func TestCopyFileNonExistent(t *testing.T) {
	tmpDir := t.TempDir()

	err := CopyFile(filepath.Join(tmpDir, "nonexistent"), filepath.Join(tmpDir, "dest"))
	if err == nil {
		t.Error("CopyFile should fail for non-existent source")
	}
}

func TestCopyIntoDots(t *testing.T) {
	tmpDir := t.TempDir()

	// Initialize dots
	_, err := Init(tmpDir)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// Create a source file
	srcPath := filepath.Join(tmpDir, ".bashrc")
	if err := os.WriteFile(srcPath, []byte("export PATH=/usr/bin"), 0o644); err != nil {
		t.Fatalf("failed to create source: %v", err)
	}

	// Copy into dots
	destPath, err := CopyIntoDots(tmpDir, srcPath)
	if err != nil {
		t.Fatalf("CopyIntoDots failed: %v", err)
	}

	expectedDest := filepath.Join(tmpDir, ".dots", ".bashrc")
	if destPath != expectedDest {
		t.Errorf("CopyIntoDots returned %q, want %q", destPath, expectedDest)
	}

	// Verify file exists
	if _, err := os.Stat(destPath); err != nil {
		t.Errorf("copied file doesn't exist: %v", err)
	}
}

func TestCopyIntoDotsRejectsDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	_, _ = Init(tmpDir)

	dirPath := filepath.Join(tmpDir, "mydir")
	if err := os.Mkdir(dirPath, 0o755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	_, err := CopyIntoDots(tmpDir, dirPath)
	if err == nil {
		t.Error("CopyIntoDots should reject directories")
	}
}

func TestCopyIntoDotsRejectsSymlink(t *testing.T) {
	tmpDir := t.TempDir()
	_, _ = Init(tmpDir)

	// Create a regular file and a symlink to it
	realFile := filepath.Join(tmpDir, "realfile")
	if err := os.WriteFile(realFile, []byte("content"), 0o644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	symlink := filepath.Join(tmpDir, "symlink")
	if err := os.Symlink(realFile, symlink); err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}

	_, err := CopyIntoDots(tmpDir, symlink)
	if err == nil {
		t.Error("CopyIntoDots should reject symlinks")
	}
}

func TestEnsureSymlink(t *testing.T) {
	tmpDir := t.TempDir()

	// Create source file
	source := filepath.Join(tmpDir, "source.txt")
	if err := os.WriteFile(source, []byte("content"), 0o644); err != nil {
		t.Fatalf("failed to create source: %v", err)
	}

	// Create symlink
	target := filepath.Join(tmpDir, "link.txt")
	if err := EnsureSymlink(target, source); err != nil {
		t.Fatalf("EnsureSymlink failed: %v", err)
	}

	// Verify symlink
	info, err := os.Lstat(target)
	if err != nil {
		t.Fatalf("symlink not created: %v", err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Error("target is not a symlink")
	}

	// Verify it points to the right place
	linkDest, err := os.Readlink(target)
	if err != nil {
		t.Fatalf("failed to read symlink: %v", err)
	}
	if linkDest != source {
		t.Errorf("symlink points to %q, want %q", linkDest, source)
	}
}

func TestEnsureSymlinkCreatesParentDirs(t *testing.T) {
	tmpDir := t.TempDir()

	source := filepath.Join(tmpDir, "source")
	if err := os.WriteFile(source, []byte("content"), 0o644); err != nil {
		t.Fatalf("failed to create source: %v", err)
	}

	// Target in nested non-existent directory
	target := filepath.Join(tmpDir, "a", "b", "c", "link")
	if err := EnsureSymlink(target, source); err != nil {
		t.Fatalf("EnsureSymlink should create parent dirs: %v", err)
	}

	if _, err := os.Lstat(target); err != nil {
		t.Errorf("symlink not created in nested directory: %v", err)
	}
}

func TestEnsureSymlinkOverwritesExistingFile(t *testing.T) {
	tmpDir := t.TempDir()

	source := filepath.Join(tmpDir, "source")
	if err := os.WriteFile(source, []byte("source content"), 0o644); err != nil {
		t.Fatalf("failed to create source: %v", err)
	}

	target := filepath.Join(tmpDir, "target")
	if err := os.WriteFile(target, []byte("existing content"), 0o644); err != nil {
		t.Fatalf("failed to create target: %v", err)
	}

	if err := EnsureSymlink(target, source); err != nil {
		t.Fatalf("EnsureSymlink should overwrite existing file: %v", err)
	}

	info, _ := os.Lstat(target)
	if info.Mode()&os.ModeSymlink == 0 {
		t.Error("target should now be a symlink")
	}
}

func TestEnsureSymlinkRejectsDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	source := filepath.Join(tmpDir, "source")
	if err := os.WriteFile(source, []byte("content"), 0o644); err != nil {
		t.Fatalf("failed to create source: %v", err)
	}

	target := filepath.Join(tmpDir, "targetdir")
	if err := os.Mkdir(target, 0o755); err != nil {
		t.Fatalf("failed to create target directory: %v", err)
	}

	err := EnsureSymlink(target, source)
	if err == nil {
		t.Error("EnsureSymlink should reject directory as target")
	}
}

func TestLinkStatusLinked(t *testing.T) {
	tmpDir := t.TempDir()

	source := filepath.Join(tmpDir, "source")
	if err := os.WriteFile(source, []byte("content"), 0o644); err != nil {
		t.Fatalf("failed to create source: %v", err)
	}

	target := filepath.Join(tmpDir, "target")
	if err := os.Symlink(source, target); err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}

	entry := config.FileEntry{Source: source, Target: target}
	status, err := LinkStatus(entry)
	if err != nil {
		t.Fatalf("LinkStatus failed: %v", err)
	}

	if status.Status != StatusLinked {
		t.Errorf("status = %v, want %v", status.Status, StatusLinked)
	}
}

func TestLinkStatusMissing(t *testing.T) {
	tmpDir := t.TempDir()

	entry := config.FileEntry{
		Source: filepath.Join(tmpDir, "source"),
		Target: filepath.Join(tmpDir, "nonexistent"),
	}

	status, err := LinkStatus(entry)
	if err != nil {
		t.Fatalf("LinkStatus failed: %v", err)
	}

	if status.Status != StatusMissing {
		t.Errorf("status = %v, want %v", status.Status, StatusMissing)
	}
}

func TestLinkStatusConflictNotSymlink(t *testing.T) {
	tmpDir := t.TempDir()

	source := filepath.Join(tmpDir, "source")
	target := filepath.Join(tmpDir, "target")

	if err := os.WriteFile(source, []byte("source"), 0o644); err != nil {
		t.Fatalf("failed to create source: %v", err)
	}
	if err := os.WriteFile(target, []byte("regular file"), 0o644); err != nil {
		t.Fatalf("failed to create target: %v", err)
	}

	entry := config.FileEntry{Source: source, Target: target}
	status, err := LinkStatus(entry)
	if err != nil {
		t.Fatalf("LinkStatus failed: %v", err)
	}

	if status.Status != StatusConflicts {
		t.Errorf("status = %v, want %v", status.Status, StatusConflicts)
	}
	if status.Info != "not a symlink" {
		t.Errorf("Info = %q, want %q", status.Info, "not a symlink")
	}
}

func TestLinkStatusConflictWrongTarget(t *testing.T) {
	tmpDir := t.TempDir()

	source := filepath.Join(tmpDir, "source")
	otherSource := filepath.Join(tmpDir, "other")
	target := filepath.Join(tmpDir, "target")

	if err := os.WriteFile(source, []byte("source"), 0o644); err != nil {
		t.Fatalf("failed to create source: %v", err)
	}
	if err := os.WriteFile(otherSource, []byte("other"), 0o644); err != nil {
		t.Fatalf("failed to create other: %v", err)
	}
	// Symlink points to wrong file
	if err := os.Symlink(otherSource, target); err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}

	entry := config.FileEntry{Source: source, Target: target}
	status, err := LinkStatus(entry)
	if err != nil {
		t.Fatalf("LinkStatus failed: %v", err)
	}

	if status.Status != StatusConflicts {
		t.Errorf("status = %v, want %v", status.Status, StatusConflicts)
	}
}

func TestContentStatusLinked(t *testing.T) {
	tmpDir := t.TempDir()

	source := filepath.Join(tmpDir, "source")
	if err := os.WriteFile(source, []byte("content"), 0o644); err != nil {
		t.Fatalf("failed to create source: %v", err)
	}

	target := filepath.Join(tmpDir, "target")
	if err := os.Symlink(source, target); err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}

	entry := config.FileEntry{Source: source, Target: target}
	status, err := ContentStatus(entry)
	if err != nil {
		t.Fatalf("ContentStatus failed: %v", err)
	}

	if status.Status != StatusLinked {
		t.Errorf("status = %v, want %v", status.Status, StatusLinked)
	}
}

func TestContentStatusMissingSourceFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Source doesn't exist, target doesn't exist
	entry := config.FileEntry{
		Source: filepath.Join(tmpDir, "nonexistent_source"),
		Target: filepath.Join(tmpDir, "nonexistent_target"),
	}

	status, err := ContentStatus(entry)
	if err != nil {
		t.Fatalf("ContentStatus failed: %v", err)
	}

	if status.Status != StatusMissing {
		t.Errorf("status = %v, want %v", status.Status, StatusMissing)
	}
}

func TestRelativePath(t *testing.T) {
	tests := []struct {
		home string
		path string
		want string
	}{
		{"/home/user", "/home/user/.bashrc", "~/.bashrc"},
		{"/home/user", "/home/user/.config/nvim/init.vim", "~/.config/nvim/init.vim"},
		{"/home/user", "/etc/hosts", "../etc/hosts"}, // outside home, gets relative path with ..
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := RelativePath(tt.home, tt.path)
			if got != tt.want {
				t.Errorf("RelativePath(%q, %q) = %q, want %q", tt.home, tt.path, got, tt.want)
			}
		})
	}
}

func TestSyncStatusConstants(t *testing.T) {
	// Ensure status constants have expected values
	if StatusMissing != "missing" {
		t.Errorf("StatusMissing = %q, want %q", StatusMissing, "missing")
	}
	if StatusLinked != "linked" {
		t.Errorf("StatusLinked = %q, want %q", StatusLinked, "linked")
	}
	if StatusDiverged != "diverged" {
		t.Errorf("StatusDiverged = %q, want %q", StatusDiverged, "diverged")
	}
	if StatusConflicts != "conflict" {
		t.Errorf("StatusConflicts = %q, want %q", StatusConflicts, "conflict")
	}
}
