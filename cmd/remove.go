package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/subcode-labs/dots/internal/config"
	"github.com/subcode-labs/dots/internal/dotfile"
)

var removeCmd = &cobra.Command{
	Use:   "remove <file>",
	Short: "Remove a file from dots management",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := dotfile.HomeDir()
		if err != nil {
			return err
		}
		if err := ensureManifestExists(home); err != nil {
			return err
		}
		manifest, err := config.Load(home)
		if err != nil {
			return err
		}
		resolvedTarget, err := filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("resolve path: %w", err)
		}
		entry, found := config.FindEntry(manifest, resolvedTarget)
		if !found {
			return fmt.Errorf("file not tracked: %s", resolvedTarget)
		}

		if err := removeSymlink(entry.Target); err != nil {
			return err
		}
		if err := restoreFile(entry); err != nil {
			return err
		}
		if removed := config.RemoveEntry(manifest, resolvedTarget); !removed {
			return fmt.Errorf("failed to remove manifest entry for %s", resolvedTarget)
		}
		if err := config.Save(home, manifest); err != nil {
			return err
		}
		if err := os.Remove(entry.Source); err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return fmt.Errorf("remove stored file: %w", err)
			}
		}
		color.New(color.FgGreen).Printf("Removed %s from dots\n", entry.Target)
		return nil
	},
}

func removeSymlink(target string) error {
	info, err := os.Lstat(target)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("stat target: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("target %s is a directory", target)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("target %s is not a symlink", target)
	}
	if err := os.Remove(target); err != nil {
		return fmt.Errorf("remove symlink: %w", err)
	}
	return nil
}

func restoreFile(entry config.FileEntry) error {
	if err := os.MkdirAll(filepath.Dir(entry.Target), 0o755); err != nil {
		return fmt.Errorf("ensure target dir: %w", err)
	}
	if err := dotfile.CopyFile(entry.Source, entry.Target); err != nil {
		return err
	}
	return nil
}
