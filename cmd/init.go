package cmd

import (
	"fmt"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/subcode-labs/dots/internal/dotfile"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a dots repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := dotfile.HomeDir()
		if err != nil {
			return err
		}
		dotsDir, err := dotfile.Init(home)
		if err != nil {
			return err
		}
		if err := initGitRepo(dotsDir); err != nil {
			return err
		}
		color.New(color.FgGreen).Printf("Initialized dots at %s\n", dotsDir)
		return nil
	},
}

func initGitRepo(path string) error {
	cmd := exec.Command("git", "-C", path, "init")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("init git repo: %w (%s)", err, string(output))
	}
	return nil
}
