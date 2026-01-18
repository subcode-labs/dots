package dotfile

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/subcode-labs/dots/internal/config"
)

type SyncStatus string

const (
	StatusMissing   SyncStatus = "missing"
	StatusLinked    SyncStatus = "linked"
	StatusDiverged  SyncStatus = "diverged"
	StatusConflicts SyncStatus = "conflict"
)

type StatusEntry struct {
	Entry  config.FileEntry
	Status SyncStatus
	Info   string
}

func HomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory: %w", err)
	}
	return home, nil
}

func DotsDir(home string) string {
	return config.DotsDir(home)
}

func ManifestPath(home string) string {
	return config.ManifestPath(home)
}

func Init(home string) (string, error) {
	if _, err := config.EnsureDotsDir(home); err != nil {
		return "", err
	}
	manifest, err := config.Load(home)
	if err != nil {
		return "", err
	}
	if err := config.Save(home, manifest); err != nil {
		return "", err
	}
	return config.DotsDir(home), nil
}

func CopyIntoDots(home, sourcePath string) (string, error) {
	info, err := os.Lstat(sourcePath)
	if err != nil {
		return "", fmt.Errorf("inspect source: %w", err)
	}
	if info.IsDir() {
		return "", fmt.Errorf("source must be a file, got directory")
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return "", fmt.Errorf("source must be a regular file, got symlink")
	}
	base := filepath.Base(sourcePath)
	destination := filepath.Join(config.DotsDir(home), base)
	if err := CopyFile(sourcePath, destination); err != nil {
		return "", err
	}
	return destination, nil
}

func CopyFile(src, dst string) error {
	return copyFile(src, dst)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer in.Close()

	info, err := in.Stat()
	if err != nil {
		return fmt.Errorf("stat source: %w", err)
	}

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode().Perm())
	if err != nil {
		return fmt.Errorf("create destination: %w", err)
	}
	defer func() {
		_ = out.Close()
	}()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("copy data: %w", err)
	}
	if err := out.Sync(); err != nil {
		return fmt.Errorf("sync destination: %w", err)
	}
	return nil
}

func EnsureSymlink(target, source string) error {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("ensure parent dir: %w", err)
	}
	if info, err := os.Lstat(target); err == nil {
		if info.IsDir() {
			return fmt.Errorf("target %s is a directory", target)
		}
		if err := os.Remove(target); err != nil {
			return fmt.Errorf("remove existing target: %w", err)
		}
	} else if !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("stat target: %w", err)
	}
	if err := os.Symlink(source, target); err != nil {
		return fmt.Errorf("create symlink: %w", err)
	}
	return nil
}

func LinkStatus(entry config.FileEntry) (StatusEntry, error) {
	info, err := os.Lstat(entry.Target)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return StatusEntry{Entry: entry, Status: StatusMissing}, nil
		}
		return StatusEntry{}, fmt.Errorf("stat target: %w", err)
	}

	if info.Mode()&os.ModeSymlink == 0 {
		return StatusEntry{Entry: entry, Status: StatusConflicts, Info: "not a symlink"}, nil
	}

	linkPath, err := os.Readlink(entry.Target)
	if err != nil {
		return StatusEntry{}, fmt.Errorf("read symlink: %w", err)
	}

	if linkPath == entry.Source {
		return StatusEntry{Entry: entry, Status: StatusLinked}, nil
	}

	return StatusEntry{Entry: entry, Status: StatusConflicts, Info: fmt.Sprintf("links to %s", linkPath)}, nil
}

func ContentStatus(entry config.FileEntry) (StatusEntry, error) {
	status, err := LinkStatus(entry)
	if err != nil {
		return StatusEntry{}, err
	}
	if status.Status == StatusLinked {
		return status, nil
	}

	if status.Status == StatusMissing || status.Status == StatusConflicts {
		if _, err := os.Stat(entry.Source); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				status.Status = StatusMissing
				status.Info = "stored file missing"
				return status, nil
			}
			return StatusEntry{}, fmt.Errorf("stat stored file: %w", err)
		}
		if status.Status == StatusConflicts {
			return status, nil
		}
	}

	match, err := sameContent(entry.Source, entry.Target)
	if err != nil {
		return StatusEntry{}, err
	}
	if match {
		return StatusEntry{Entry: entry, Status: StatusMissing, Info: "target missing"}, nil
	}
	return StatusEntry{Entry: entry, Status: StatusDiverged}, nil
}

func sameContent(source, target string) (bool, error) {
	sourceHash, err := fileHash(source)
	if err != nil {
		return false, err
	}
	targetHash, err := fileHash(target)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return sourceHash == targetHash, nil
}

func fileHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("hash file: %w", err)
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func RelativePath(home, path string) string {
	if rel, err := filepath.Rel(home, path); err == nil {
		return filepath.Join("~", rel)
	}
	return path
}
