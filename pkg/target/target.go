package target

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
)

type Profile struct {
	Name              string
	StripANSI         bool
	RemoveDupLines    bool
	CollapsePaths     bool
	Skeletonize      bool
	CompressionLevel  float64
	KeepReasoning     bool
	KeepSuggestions   bool
	KeepSyntax        bool
}

var Profiles = map[string]Profile{
	"cursor": {
		Name:             "Cursor",
		StripANSI:        true,
		RemoveDupLines:   true,
		CollapsePaths:    true,
		Skeletonize:      false,
		CompressionLevel: 0.80,
		KeepReasoning:    false,
		KeepSuggestions:  true,
		KeepSyntax:       true,
	},
	"claude-code": {
		Name:             "Claude Code",
		StripANSI:        true,
		RemoveDupLines:   true,
		CollapsePaths:    true,
		Skeletonize:      false,
		CompressionLevel: 0.60,
		KeepReasoning:    true,
		KeepSuggestions:  false,
		KeepSyntax:       true,
	},
	"copilot": {
		Name:             "GitHub Copilot",
		StripANSI:        true,
		RemoveDupLines:   false,
		CollapsePaths:    false,
		Skeletonize:      false,
		CompressionLevel: 0.40,
		KeepReasoning:    false,
		KeepSuggestions:  true,
		KeepSyntax:       true,
	},
	"default": {
		Name:             "Default",
		StripANSI:        true,
		RemoveDupLines:   true,
		CollapsePaths:    true,
		Skeletonize:      false,
		CompressionLevel: 0.70,
		KeepReasoning:    false,
		KeepSuggestions:  false,
		KeepSyntax:       true,
	},
}

type Manager struct {
	configDir string
	profiles  map[string]Profile
}

func NewManager() (*Manager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(home, ".mrt", "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	return &Manager{
		configDir: configDir,
		profiles:  Profiles,
	}, nil
}

func (m *Manager) GetProfile(name string) Profile {
	if p, ok := m.profiles[name]; ok {
		return p
	}
	return Profiles["default"]
}

func (m *Manager) GetCurrentTarget() string {
	targetPath := filepath.Join(m.configDir, "target")
	data, err := os.ReadFile(targetPath)
	if err != nil {
		return "default"
	}
	return string(data)
}

func (m *Manager) SetTarget(name string) error {
	if _, ok := m.profiles[name]; !ok {
		return fmt.Errorf("unknown target: %s", name)
	}

	targetPath := filepath.Join(m.configDir, "target")
	return os.WriteFile(targetPath, []byte(name), 0644)
}

func (m *Manager) ListTargets() []string {
	var names []string
	for name := range m.profiles {
		names = append(names, name)
	}
	return names
}

func (m *Manager) SaveProfile(name string, profile Profile) error {
	data := fmt.Sprintf("%s:%v:%v:%v:%v:%.2f:%v:%v:%v",
		profile.Name,
		profile.StripANSI,
		profile.RemoveDupLines,
		profile.CollapsePaths,
		profile.Skeletonize,
		profile.CompressionLevel,
		profile.KeepReasoning,
		profile.KeepSuggestions,
		profile.KeepSyntax,
	)

	profilePath := filepath.Join(m.configDir, "profiles", name)
	return os.WriteFile(profilePath, []byte(data), 0644)
}

func RenderTargetInfo(target string, profile Profile) string {
	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		Bold(true).Render

	metric := lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")).Render

	value := lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")).
		Bold(true).Render

	box := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("99")).
		Padding(1, 2)

	info := fmt.Sprintf(`Target Profile: %s
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Compression Level:  %s (%.0f%%)
Strip ANSI Colors:  %s
Remove Dup Lines:    %s
Collapse Paths:     %s
Skeletonize Code:   %s
Keep Reasoning:     %s
Keep Suggestions:   %s
Keep Syntax:        %s`,
		header(target),
		value(fmt.Sprintf("%.0f%%", profile.CompressionLevel*100)),
		profile.CompressionLevel,
		boolIcon(profile.StripANSI),
		boolIcon(profile.RemoveDupLines),
		boolIcon(profile.CollapsePaths),
		boolIcon(profile.Skeletonize),
		boolIcon(profile.KeepReasoning),
		boolIcon(profile.KeepSuggestions),
		boolIcon(profile.KeepSyntax),
	)

	return box.Render(info)
}

func boolIcon(b bool) string {
	if b {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("✓")
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("✗")
}
