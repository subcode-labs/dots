package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/subcode-labs/dots/internal/config"
	"github.com/subcode-labs/dots/internal/dotfile"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Create symlinks for all tracked dotfiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := dotfile.HomeDir()
		if err != nil {
			return err
		}
		manifest, err := config.Load(home)
		if err != nil {
			return err
		}
		if len(manifest.Files) == 0 {
			color.New(color.FgYellow).Println("No tracked dotfiles.")
			return nil
		}
		for _, entry := range manifest.Files {
			if err := dotfile.EnsureSymlink(entry.Target, entry.Source); err != nil {
				return err
			}
			color.New(color.FgGreen).Printf("Linked %s -> %s\n", entry.Target, entry.Source)
		}
		return nil
	},
}
