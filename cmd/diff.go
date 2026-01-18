package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/subcode-labs/dots/internal/config"
	"github.com/subcode-labs/dots/internal/dotfile"
)

var diffCmd = &cobra.Command{
	Use:   "diff [file]",
	Short: "Show diffs between tracked files and originals",
	Args:  cobra.MaximumNArgs(1),
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
		if len(manifest.Files) == 0 {
			color.New(color.FgYellow).Println("No tracked dotfiles.")
			return nil
		}

		if len(args) == 1 {
			return diffSingle(manifest, args[0])
		}

		return diffAll(manifest)
	},
}

func diffSingle(manifest *config.Manifest, target string) error {
	resolvedTarget, err := filepath.Abs(target)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}
	entry, found := config.FindEntry(manifest, resolvedTarget)
	if !found {
		return fmt.Errorf("file not tracked: %s", resolvedTarget)
	}
	output, err := runDiff(entry)
	if err != nil {
		return err
	}
	if strings.TrimSpace(output) == "" {
		color.New(color.FgYellow).Printf("No differences for %s\n", entry.Target)
		return nil
	}
	printDiff(output)
	return nil
}

func diffAll(manifest *config.Manifest) error {
	var outputs []string
	for _, entry := range manifest.Files {
		status, err := dotfile.ContentStatus(entry)
		if err != nil {
			return err
		}
		if status.Status != dotfile.StatusDiverged {
			continue
		}
		output, err := runDiff(entry)
		if err != nil {
			return err
		}
		if strings.TrimSpace(output) == "" {
			continue
		}
		outputs = append(outputs, output)
	}

	if len(outputs) == 0 {
		color.New(color.FgYellow).Println("No diverged files.")
		return nil
	}

	for i, output := range outputs {
		if i > 0 {
			fmt.Println()
		}
		printDiff(output)
	}
	return nil
}

func runDiff(entry config.FileEntry) (string, error) {
	cmd := exec.Command("diff", "-u", entry.Target, entry.Source)
	output, err := cmd.CombinedOutput()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			if exitErr.ExitCode() == 1 {
				return string(output), nil
			}
		}
		return "", fmt.Errorf("diff failed: %w (%s)", err, strings.TrimSpace(string(output)))
	}
	return string(output), nil
}

func printDiff(output string) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++"):
			color.New(color.FgGreen).Println(line)
		case strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---"):
			color.New(color.FgRed).Println(line)
		default:
			fmt.Println(line)
		}
	}
}
