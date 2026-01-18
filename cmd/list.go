package cmd

import (
	"fmt"
	"sort"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/subcode-labs/dots/internal/config"
	"github.com/subcode-labs/dots/internal/dotfile"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tracked dotfiles",
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
		sort.Slice(manifest.Files, func(i, j int) bool {
			return manifest.Files[i].Target < manifest.Files[j].Target
		})
		for _, entry := range manifest.Files {
			fmt.Println(entry.Target)
		}
		return nil
	},
}
