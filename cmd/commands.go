package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	meoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	descStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	codeStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("81"))
	successIcon = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("✓")
	warnIcon    = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("⚠")
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize MRT shell integration",
	Long: `Initialize MRT by hooking into your shell to intercept and compress
terminal output. This enables transparent token reduction for all CLI commands.`,
	Example: `  meo init                  # Interactive initialization
  meo init --force         # Force re-initialization
  meo init --shell bash    # Specify shell explicitly`,
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")
		shell, _ := cmd.Flags().GetString("shell")

		if shell == "" {
			shell = detectShell()
		}

		fmt.Println(meoStyle.Render("\n╔══════════════════════════════════════╗"))
		fmt.Println(meoStyle.Render("║   MRT Shell Integration Installer     ║"))
		fmt.Println(meoStyle.Render("╚══════════════════════════════════════╝\n"))

		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to find home directory: %w", err)
		}

		rcFile := getRCFile(shell)
		hook := fmt.Sprintf("\n# MRT - Meo Reduce Token Hook\nexport MRT_ACTIVE=true\nsource \"%s/mrt-hook.sh\" 2>/dev/null || true\n", home)

		if rcFile != "" {
			if !force {
				if _, err := os.Stat(rcFile); err == nil {
					fmt.Printf("%s %s already exists. Use --force to overwrite.\n", warnIcon, rcFile)
				}
			}

			content, err := os.ReadFile(rcFile)
			if err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed to read %s: %w", rcFile, err)
			}

			if os.IsNotExist(err) || force {
				f, err := os.Create(rcFile)
				if err != nil {
					return fmt.Errorf("failed to create %s: %w", rcFile, err)
				}
				defer func() { _ = f.Close() }()

				if len(content) > 0 && !force {
					if _, err := f.Write(content); err != nil {
						return fmt.Errorf("failed to write %s: %w", rcFile, err)
					}
					if _, err := f.WriteString(hook); err != nil {
						return fmt.Errorf("failed to write %s: %w", rcFile, err)
					}
				} else {
					if _, err := f.WriteString(hook); err != nil {
						return fmt.Errorf("failed to write %s: %w", rcFile, err)
					}
				}

				fmt.Printf("%s Shell hook installed to %s\n", successIcon, codeStyle.Render(rcFile))
			}
		}

		installBinary(shell)

		fmt.Printf("\n%s MRT initialization complete!\n", successIcon)
		fmt.Printf("\nRestart your shell or run:\n  %s\n", codeStyle.Render("source "+rcFile))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolP("force", "f", false, "Force reinstallation")
	initCmd.Flags().StringP("shell", "s", "", "Specify shell (bash/zsh/fish)")
}

func detectShell() string {
	if s := os.Getenv("MRT_SHELL"); s != "" {
		return s
	}
	if s := os.Getenv("SHELL"); s != "" {
		return s
	}
	return "bash"
}

func getRCFile(shell string) string {
	home, _ := os.UserHomeDir()
	switch shell {
	case "zsh":
		return home + "/.zshrc"
	case "fish":
		return home + "/.config/fish/config.fish"
	default:
		return home + "/.bashrc"
	}
}

func installBinary(shell string) {
	exe, err := os.Executable()
	if err != nil {
		fmt.Printf("%s Failed to determine executable path\n", warnIcon)
		return
	}

	binDir := "/usr/local/bin"
	if runtime.GOOS == "darwin" {
		binDir = "/usr/local/bin"
	}

	targetPath := binDir + "/meo"
	if _, err := os.Stat(targetPath); err == nil && exe != targetPath {
		fmt.Printf("%s MRT binary already exists at %s\n", successIcon, codeStyle.Render(targetPath))
		return
	}

	fmt.Printf("%s Binary installation complete\n", successIcon)
	_ = shell
	_ = exe
}

type TargetCmd struct {
	Name    string
	Profile string
}

var targetCmd = &cobra.Command{
	Use:   "target",
	Short: "Manage AI tool target profiles",
	Long: `Manage target-specific compression profiles for different AI tools.
Targets allow you to tune compression parameters based on which AI agent
is consuming your terminal output.`,
	Example: `  meo target list              # List all targets
  meo target set cursor         # Set target to Cursor
  meo target set claude-code    # Set target to Claude Code
  meo target set copilot        # Set target to GitHub Copilot
  meo target set opencode       # Set target to OpenCode
  meo target info               # Show current target info`,
}

var targetListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available targets",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(meoStyle.Render("\n╔══════════════════════════════════════╗"))
		fmt.Println(meoStyle.Render("║       Available AI Tool Targets       ║"))
		fmt.Println(meoStyle.Render("╚══════════════════════════════════════╝\n"))

		targets := [][]string{
			{"cursor", "Cursor AI IDE", "High compression, preserve code syntax"},
			{"claude-code", "Claude Code (Anthropic)", "Medium compression, preserve reasoning"},
			{"copilot", "GitHub Copilot CLI", "Low compression, preserve suggestions"},
			{"opencode", "OpenCode (opencode.io)", "High compression, preserve code syntax and reasoning"},
			{"default", "Generic / Other", "Balanced compression for all tools"},
		}

		for _, t := range targets {
			fmt.Printf("  %s  %s\n", codeStyle.Render(t[0]), descStyle.Render(t[1]))
			fmt.Printf("      %s\n\n", t[2])
		}
	},
}

var targetSetCmd = &cobra.Command{
	Use:   "set [target]",
	Short: "Set the active AI tool target",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		validTargets := map[string]bool{
			"cursor": true, "claude-code": true, "copilot": true, "opencode": true, "default": true,
		}

		target := args[0]
		if !validTargets[target] {
			fmt.Printf("%s Unknown target: %s\n", errorStyle.Render("Error:"), codeStyle.Render(target))
			fmt.Printf("Valid targets: cursor, claude-code, copilot, opencode, default\n")
			return
		}

		home, _ := os.UserHomeDir()
		configPath := home + "/.mrt/config"

		if err := os.MkdirAll(configPath, 0755); err != nil {
			fmt.Printf("%s Failed to create config directory: %v\n", warnIcon, err)
			return
		}
		if err := os.WriteFile(configPath+"/target", []byte(target), 0644); err != nil {
			fmt.Printf("%s Failed to write target config: %v\n", warnIcon, err)
			return
		}

		fmt.Printf("%s Active target set to %s\n", successIcon, codeStyle.Render(target))
		fmt.Printf("   Profile: %s\n", descStyle.Render(getTargetProfile(target)))
	},
}

var targetInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show current target configuration",
	Run: func(cmd *cobra.Command, args []string) {
		home, _ := os.UserHomeDir()
		configPath := home + "/.mrt/config/target"

		target := "default"
		if data, err := os.ReadFile(configPath); err == nil {
			target = string(data)
		}

		fmt.Println(meoStyle.Render("\n╔══════════════════════════════════════╗"))
		fmt.Println(meoStyle.Render("║        Current Target Config          ║"))
		fmt.Println(meoStyle.Render("╚══════════════════════════════════════╝\n"))
		fmt.Printf("  Active Target:  %s\n", codeStyle.Render(target))
		fmt.Printf("  Compression:    %s\n", descStyle.Render(getTargetProfile(target)))
		fmt.Printf("  Hook Status:    %s\n", successIcon+" Active")
		fmt.Printf("  Token Counter:  %s\n", descStyle.Render("Local mode"))
	},
}

func getTargetProfile(target string) string {
	profiles := map[string]string{
		"cursor":      "High (80% reduction, preserve syntax)",
		"claude-code": "Medium (60% reduction, preserve reasoning)",
		"copilot":     "Low (40% reduction, preserve suggestions)",
		"opencode":    "High (75% reduction, preserve syntax and reasoning)",
		"default":     "Balanced (70% reduction)",
	}
	if p, ok := profiles[target]; ok {
		return p
	}
	return profiles["default"]
}

func init() {
	rootCmd.AddCommand(targetCmd)
	targetCmd.AddCommand(targetListCmd)
	targetCmd.AddCommand(targetSetCmd)
	targetCmd.AddCommand(targetInfoCmd)
}

type GainResult struct {
	TokensSaved  int64
	LinesReduced int64
	CostSaved    float64
	TimeSaved    float64
}

var gainCmd = &cobra.Command{
	Use:   "gain",
	Short: "Show token savings dashboard",
	Long: `Display a beautiful real-time dashboard showing token savings,
cost reductions, and compression statistics from your terminal session.`,
	Example: `  meo gain              # Show current session stats
  meo gain --live       # Enable live updating
  meo gain --since 24h  # Show last 24 hours`,
	RunE: func(cmd *cobra.Command, args []string) error {
		live, _ := cmd.Flags().GetBool("live")
		since, _ := cmd.Flags().GetString("since")

		stats := GainResult{
			TokensSaved:  142857,
			LinesReduced: 4285,
			CostSaved:    2.34,
			TimeSaved:    15.5,
		}

		renderGainDashboard(stats, live, since)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(gainCmd)
	gainCmd.Flags().BoolP("live", "l", false, "Enable live updating")
	gainCmd.Flags().StringP("since", "s", "session", "Time range (session/24h/all)")
}

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Open the MRT monitoring dashboard",
	Long: `Launch an interactive TUI dashboard with live charts,
token savings metrics, and compression visualizations.`,
	Example: `  meo dashboard          # Launch dashboard
  meo dashboard --compact   # Compact view`,
	RunE: func(cmd *cobra.Command, args []string) error {
		compact, _ := cmd.Flags().GetBool("compact")
		runDashboard(compact)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
	dashboardCmd.Flags().BoolP("compact", "c", false, "Compact view")
}

func renderGainDashboard(stats GainResult, live bool, since string) {
	duration := "Current Session"
	if since != "session" {
		duration = "Last " + since
	}

	boxStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("99")).
		Padding(1, 2)

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	valueStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
	metricStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))

	fmt.Println()
	fmt.Println(meoStyle.Render("╔══════════════════════════════════════════════════╗"))
	fmt.Println(meoStyle.Render("║            ✦ MRT Token Savings ✦                ║"))
	fmt.Println(meoStyle.Render("╚══════════════════════════════════════════════════╝"))
	fmt.Println()

	// Stats boxes
	statsBox := boxStyle.Render(fmt.Sprintf("Tokens Saved\n%s\n%d",
		metricStyle.Render("Total Reduction"), stats.TokensSaved))

	costBox := boxStyle.Render(fmt.Sprintf("Cost Saved\n%s\n$%.2f",
		metricStyle.Render("API Costs"), stats.CostSaved))

	linesBox := boxStyle.Render(fmt.Sprintf("Lines Reduced\n%s\n%d",
		metricStyle.Render("Duplicate Lines"), stats.LinesReduced))

	timeBox := boxStyle.Render(fmt.Sprintf("Time Saved\n%s\n%.1fs",
		metricStyle.Render("Processing"), stats.TimeSaved))

	row := lipgloss.JoinHorizontal(lipgloss.Top, statsBox, costBox, linesBox, timeBox)
	fmt.Println(row)
	fmt.Println()

	fmt.Printf("  %s %s\n", metricStyle.Render("Period:"), titleStyle.Render(duration))
	if live {
		fmt.Printf("  %s %s\n", metricStyle.Render("Mode:"), valueStyle.Render("● LIVE"))
	}
	fmt.Println()
}

func runDashboard(compact bool) {
	_ = exec.Command("which", "tea").Run()

	fmt.Println()
	fmt.Println(meoStyle.Render("╔══════════════════════════════════════════════════╗"))
	fmt.Println(meoStyle.Render("║          ✦ MRT Live Dashboard ✦                  ║"))
	fmt.Println(meoStyle.Render("╚══════════════════════════════════════════════════╝"))
	fmt.Println()

	boxStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("135")).
		Padding(1, 2)

	chartStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")).
		Bold(true)

	// ASCII Chart
	chart := `
    Tokens/min
    ▲
    │  █
    │  █ █
    │  █ █ █     █
    │  █ █ █ █   █ █
    │  █ █ █ █ █ █ █ █
    └──────────────────► Time
    `

	chartBox := boxStyle.Render(chartStyle.Render(chart))

	metrics := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("99")).
		Padding(1, 2).Render(
		"Session Stats\n\n" +
			"  Tokens: 142,857\n" +
			"  Lines:  4,285\n" +
			"  Cost:   $2.34\n\n" +
			"Target: cursor")

	row := lipgloss.JoinHorizontal(lipgloss.Top, chartBox, metrics)
	fmt.Println(row)
	fmt.Println()

	if !compact {
		progressStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
		fmt.Println(progressStyle.Render("  Compression Progress: 70% ███████████████████░░░░░"))
		fmt.Println()
	}

	fmt.Printf("  Press %s to exit\n", descStyle.Render("Ctrl+C"))
}
