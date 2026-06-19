package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	boldStyle    = lipgloss.NewStyle().Bold(true)
)

var (
	Version   = "1.0.0"
	Commit    = "dev"
	Year      = "2026"
	MeoAuthor = "Meo Code Labs"
)

var rootCmd = &cobra.Command{
	Use:   "meo",
	Short: fmt.Sprintf("%s ✦ Meo Reduce Token (MRT)", boldStyle.Render("meo")),
	Long: `
╔══════════════════════════════════════════════════════════════╗
║  ██████╗ ██╗   ██╗███╗   ██╗ ██████╗ ███████╗ ██████╗ ███╗   ██╗║
║  ██╔══██╗██║   ██║████╗  ██║██╔═══██╗██╔════╝██╔═══██╗████╗  ██║║
║  ██║  ██║██║   ██║██╔██╗ ██║██║   ██║███████╗██║   ██║██╔██╗ ██║║
║  ██║  ██║██║   ██║██║╚██╗██║██║   ██║╚════██║██║   ██║██║╚██╗██║║
║  ██████╔╝╚██████╔╝██║ ╚████║╚██████╔╝███████║╚██████╔╝██║ ╚████║║
║  ╚═════╝  ╚═════╝ ╚═╝  ╚═══╝ ╚═════╝ ╚══════╝ ╚═════╝ ╚═╝  ╚═══╝║
║                                                              ║
║  %s  │  %s  │  %s           ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝

A powerful token reduction tool for AI-assisted development.
Reduce terminal output noise and save tokens while coding.

Type 'meo --help' for available commands.`,
	DisableAutoGenTag: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\n%s Version %s (commit: %s)\n", infoStyle.Render("MRT"), successStyle.Render(Version), Commit)
		fmt.Printf("Developed by %s  |  meocode.com\n\n", infoStyle.Render(MeoAuthor))
		fmt.Printf("Usage %s for available commands.\n", infoStyle.Render("meo --help"))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s %v\n", errorStyle.Render("Error:"), err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Suppress all output except errors")
}
