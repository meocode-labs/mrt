package cmd

import (
	"fmt"
	"os"
	"strings"

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
	Long: `Meo Reduce Token (MRT) is a CLI tool that strips ANSI noise, dedupes
repetitive log lines, and collapses verbose paths so terminal output
costs fewer tokens when pasted into AI coding agents.

Supported targets:
  opencode      opencode.io engine           (high compression, code syntax)
  claude-code   Claude Code (Anthropic)       (medium compression, reasoning)
  copilot       GitHub Copilot CLI            (low compression, suggestions)
  cursor        Cursor AI IDE                 (high compression, code syntax)
  default       Generic / balanced

Type 'meo --help' for available commands.`,
	DisableAutoGenTag: true,
	Run: func(cmd *cobra.Command, args []string) {
		printLogo(Version)
		fmt.Println()
		fmt.Printf("Developed by %s  |  meocode.com\n", infoStyle.Render(MeoAuthor))
		fmt.Printf("Usage %s for available commands.\n\n", infoStyle.Render("meo --help"))
	},
}

const (
	bannerInnerWidth = 48 // box interior, between the two ║
	bannerLeftPad    = 3  // consistent left padding inside the box
)

// printLogo renders the MRT opening banner. Each row is built explicitly
// so the right border always lines up with the left, even when the
// version string changes length.
func printLogo(version string) {
	bannerStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Bold(true)

	subtitleStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	emptyRow := "║" + strings.Repeat(" ", bannerInnerWidth) + "║"
	borderTop := "╔" + strings.Repeat("═", bannerInnerWidth) + "╗"
	borderBot := "╚" + strings.Repeat("═", bannerInnerWidth) + "╝"

	artRows := []string{
		"║   ███╗   ███╗██████╗ ████████╗                  ║",
		"║   ████╗ ████║██╔══██╗╚══██╔══╝                  ║",
		"║   ██╔████╔██║██████╔╝   ██║                     ║",
		"║   ██║╚██╔╝██║██╔══██╗   ██║                     ║",
		"║   ██║ ╚═╝ ██║██║  ██║   ██║                     ║",
		"║   ╚═╝     ╚═╝╚═╝  ╚═╝   ╚═╝                     ║",
	}

	title := "Meo Reduce Token · v" + version
	targets := "opencode · claude-code · copilot · cursor"
	titleRow := buildContentRow(title)
	targetsRow := buildContentRow(targets)

	rows := []string{borderTop, emptyRow}
	rows = append(rows, artRows...)
	rows = append(rows, emptyRow, titleRow, targetsRow, emptyRow, borderBot)

	for _, r := range rows {
		if r == titleRow || r == targetsRow {
			fmt.Println(subtitleStyle.Render(r))
		} else {
			fmt.Println(bannerStyle.Render(r))
		}
	}
}

// buildContentRow assembles a banner row with left-padded content and
// right-padded fill so the closing ║ always lands at the same column.
func buildContentRow(text string) string {
	left := strings.Repeat(" ", bannerLeftPad)
	contentWidth := bannerInnerWidth - bannerLeftPad
	w := lipgloss.Width(text)
	if w >= contentWidth {
		return "║" + left + text[:contentWidth] + "║"
	}
	right := strings.Repeat(" ", contentWidth-w)
	return "║" + left + text + right + "║"
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
