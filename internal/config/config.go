package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	DirName      = ".dots"
	ManifestName = "dots.yaml"
)

type Manifest struct {
	Files []FileEntry `yaml:"files"`
}

type FileEntry struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

func DotsDir(home string) string {
	return filepath.Join(home, DirName)
}

func ManifestPath(home string) string {
	return filepath.Join(DotsDir(home), ManifestName)
}

func EnsureDotsDir(home string) (string, error) {
	path := DotsDir(home)
	if err := os.MkdirAll(path, 0o755); err != nil {
		return "", fmt.Errorf("create dots directory: %w", err)
	}
	return path, nil
}

func Load(home string) (*Manifest, error) {
	manifestPath := ManifestPath(home)
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return &Manifest{Files: []FileEntry{}}, nil
		}
		return nil, fmt.Errorf("read manifest: %w", err)
	}
	var manifest Manifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("parse manifest: %w", err)
	}
	if manifest.Files == nil {
		manifest.Files = []FileEntry{}
	}
	return &manifest, nil
}

func Save(home string, manifest *Manifest) error {
	data, err := yaml.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("encode manifest: %w", err)
	}
	manifestPath := ManifestPath(home)
	if err := os.WriteFile(manifestPath, data, 0o644); err != nil {
		return fmt.Errorf("write manifest: %w", err)
	}
	return nil
}

func FindEntry(manifest *Manifest, target string) (FileEntry, bool) {
	for _, entry := range manifest.Files {
		if entry.Target == target {
			return entry, true
		}
	}
	return FileEntry{}, false
}

func UpsertEntry(manifest *Manifest, entry FileEntry) {
	for i, existing := range manifest.Files {
		if existing.Target == entry.Target {
			manifest.Files[i] = entry
			return
		}
	}
	manifest.Files = append(manifest.Files, entry)
}
