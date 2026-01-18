package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/subcode-labs/dots/internal/config"
	"github.com/subcode-labs/dots/internal/dotfile"
)

var addCmd = &cobra.Command{
	Use:   "add <file>",
	Short: "Add a file to dots management",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := dotfile.HomeDir()
		if err != nil {
			return err
		}
		sourcePath, err := filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("resolve path: %w", err)
		}
		if err := ensureManifestExists(home); err != nil {
			return err
		}
		manifest, err := config.Load(home)
		if err != nil {
			return err
		}

		destination, err := dotfile.CopyIntoDots(home, sourcePath)
		if err != nil {
			return err
		}
		entry := config.FileEntry{
			Source: destination,
			Target: sourcePath,
		}
		config.UpsertEntry(manifest, entry)
		if err := config.Save(home, manifest); err != nil {
			return err
		}
		if err := dotfile.EnsureSymlink(sourcePath, destination); err != nil {
			return err
		}
		color.New(color.FgGreen).Printf("Tracked %s -> %s\n", sourcePath, destination)
		return nil
	},
}

func ensureManifestExists(home string) error {
	manifestPath := config.ManifestPath(home)
	if _, err := os.Stat(manifestPath); err != nil {
		return fmt.Errorf("manifest not found, run 'dots init' first")
	}
	return nil
}
