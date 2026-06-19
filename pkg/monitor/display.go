package monitor

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type DashboardTheme struct {
	Header lipgloss.Style
	Box    lipgloss.Style
	Metric lipgloss.Style
	Value  lipgloss.Style
	Chart  lipgloss.Style
	Footer lipgloss.Style
}

func DefaultTheme() DashboardTheme {
	return DashboardTheme{
		Header: lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Bold(true).
			Padding(0, 1),
		Box: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("99")).
			Padding(1, 2).
			Margin(1),
		Metric: lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")),
		Value: lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(true),
		Chart: lipgloss.NewStyle().
			Foreground(lipgloss.Color("81")),
		Footer: lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")),
	}
}

type DashboardData struct {
	SessionStart   time.Time
	TokensSaved    int64
	LinesProcessed int64
	Duplicates     int64
	CostSaved      float64
	CurrentRate    float64
	Target         string
}

func RenderDashboard(data DashboardData, theme DashboardTheme) {
	elapsed := time.Since(data.SessionStart).Round(time.Second)

	header := theme.Header.Render(fmt.Sprintf(`╔══════════════════════════════════════════════════════════╗
║      ✦  MRT Live Monitor  ✦  %s                  ║
╚══════════════════════════════════════════════════════════╝`, data.Target))

	stats := theme.Box.Render(fmt.Sprintf(`Session Stats
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Duration:    %s
  Tokens:      %s
  Lines:       %s
  Duplicates:  %s
  Cost:        $%.4f`, elapsed, formatInt(data.TokensSaved), formatInt(data.LinesProcessed), formatInt(data.Duplicates), data.CostSaved))

	rate := theme.Box.Render(fmt.Sprintf(`Current Rate
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Tokens/min:  %.1f
  Lines/min:   %.1f
  Savings:     $%.2f/hr`, data.CurrentRate, data.CurrentRate*0.03, data.CurrentRate*TOKEN_RATE*60))

	chart := renderMiniChart(data)
	chartBox := theme.Box.Render(chart)

	footer := theme.Footer.Render("Ctrl+C to exit  |  Press 'c' to clear  |  Press 't' to change target")

	fmt.Println()
	fmt.Println(header)
	fmt.Println()
	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, stats, rate, chartBox))
	fmt.Println()
	fmt.Println(footer)
}

func renderMiniChart(data DashboardData) string {
	rate := int(data.CurrentRate)
	if rate <= 0 {
		rate = 50
	}
	if rate > 100 {
		rate = 100
	}

	chart := "Token Flow\n"
	chart += "━━━━━━━━━━━━━━━━━━\n"

	bars := []string{"▏", "▎", "▍", "▌", "▋", "▊", "▉", "█"}

	for i := 0; i < 8; i++ {
		height := (rate * (i + 1) / 100) * 8
		if height > 8 {
			height = 8
		}
		bar := ""
		for j := 0; j < height; j++ {
			bar += bars[j%len(bars)]
		}
		chart += fmt.Sprintf("  %2d%% %s\n", (i+1)*12, bar)
	}

	return chart
}

func RenderProgressBar(current, total int, width int) string {
	ratio := float64(current) / float64(total)
	filled := int(ratio * float64(width))

	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else if i == filled {
			bar += "▓"
		} else {
			bar += "░"
		}
	}

	percent := int(ratio * 100)
	return fmt.Sprintf("[%s] %d%%", bar, percent)
}

func RenderTable(headers []string, rows [][]string) string {
	colWidths := make([]int, len(headers))
	for i, h := range headers {
		colWidths[i] = len(h)
	}

	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	var result string

	headerRow := ""
	for i, h := range headers {
		headerRow += fmt.Sprintf("%-*s", colWidths[i]+2, h)
	}
	result += lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).Render(headerRow) + "\n"

	separator := ""
	for _, w := range colWidths {
		separator += fmt.Sprintf("%s──", strings.Repeat("─", w))
	}
	result += lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Render(separator) + "\n"

	for _, row := range rows {
		rowStr := ""
		for i, cell := range row {
			rowStr += fmt.Sprintf("%-*s", colWidths[i]+2, cell)
		}
		result += rowStr + "\n"
	}

	return result
}

func formatInt(n int64) string {
	if n >= 1_000_000 {
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	}
	if n >= 1_000 {
		return fmt.Sprintf("%.1fK", float64(n)/1_000)
	}
	return fmt.Sprintf("%d", n)
}

const TOKEN_RATE = 0.00001

func stringsReplicate(s string, count int) string {
	var result string
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

type strings struct{}

func (strings) Repeat(s string, count int) string {
	return stringsReplicate(s, count)
}

var _strings strings
