package cmd

import (
	"fmt"
	"sort"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/subcode-labs/dots/internal/config"
	"github.com/subcode-labs/dots/internal/dotfile"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show dotfile sync status",
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
		statuses := make([]dotfile.StatusEntry, 0, len(manifest.Files))
		for _, entry := range manifest.Files {
			status, err := dotfile.ContentStatus(entry)
			if err != nil {
				return err
			}
			statuses = append(statuses, status)
		}
		sort.Slice(statuses, func(i, j int) bool {
			return statuses[i].Entry.Target < statuses[j].Entry.Target
		})
		for _, status := range statuses {
			printStatus(status)
		}
		return nil
	},
}

func printStatus(status dotfile.StatusEntry) {
	label := string(status.Status)
	var painter *color.Color
	info := status.Info

	switch status.Status {
	case dotfile.StatusLinked:
		painter = color.New(color.FgGreen)
		label = "synced"
	case dotfile.StatusMissing:
		painter = color.New(color.FgYellow)
	case dotfile.StatusDiverged:
		painter = color.New(color.FgRed)
		label = "diverged"
	case dotfile.StatusConflicts:
		painter = color.New(color.FgMagenta)
		label = "conflict"
	default:
		painter = color.New(color.FgWhite)
	}

	statusLabel := painter.Sprintf("%-9s", label)
	if info != "" {
		fmt.Printf("%s %s (%s)\n", statusLabel, status.Entry.Target, info)
		return
	}
	fmt.Printf("%s %s\n", statusLabel, status.Entry.Target)
}
