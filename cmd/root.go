package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dots",
	Short: "Dots manages dotfiles with a YAML manifest",
	Long:  "dots manages your dotfiles by tracking them in a ~/.dots directory and a YAML manifest.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		color.New(color.FgRed).Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.SetVersionTemplate("dots {{.Version}}\n")
	rootCmd.Version = "0.1.0"
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(applyCmd)
}

