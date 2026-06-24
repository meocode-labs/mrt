# ✦ Meo Reduce Token (MRT)

<div align="center">

![MRT Banner](https://img.shields.io/badge/Meo%20Reduce%20Token-MRT-212?style=for-the-badge&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIyNCIgaGVpZ2h0PSIyNCIgdmlld0JveD0iMCAwIDI0IDI0IiBmaWxsPSJub25lIiBzdHJva2U9IiMyMTIiIHN0cm9rZS13aWR0aD0iMiI+PHBhdGggZD0iTTEyIDJDNi40NzIgMiAyIDI2LjQ3MiAyIDEyUzYuNDcyIDIgMTIgMiAyMiA2LjQ3MiAyMiAxMiAyNy41MjMgMiIgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIiBzdHJva2UtbGluZWpvaW49InJvdW5kIj48L3BhdGg+PHBhdGggZD0iTTE0IDhWMTIiLz48cGF0aCBkPSJNMTAgMTJIMTAiLz48cGF0aCBkPSJNMTQgMTJWMTQiLz48L3N2Zz4K)

**A powerful token reduction CLI tool for AI-assisted development**

Developed by [Meo Code Labs](https://meocode.com) | Maintained by [penadidik](https://github.com/penadidik)

[![License: MIT](https://img.shields.io/badge/License-MIT-212.svg?style=flat-square)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![npm](https://img.shields.io/badge/npm-v1.0.0-CB3837?style=flat-square&logo=npm)](https://www.npmjs.com/package/meo-reduce-token)

*Reduce terminal output noise. Save tokens. Code smarter.*

</div>

---

## ✨ Features

### 🎯 Token Reduction
- **ANSI Stripping** - Remove color codes and escape sequences
- **Duplicate Removal** - Collapse repetitive log lines (progress bars, download stats)
- **Path Collapsing** - Simplify verbose file paths to readable aliases
- **Code Skeletonization** - Keep signatures, hide implementation details
- **Local Token Counting** - Compute savings directly on your machine

### 🤖 AI Tool Integration
- **Cursor** - High compression for Cursor AI IDE
- **Claude Code** - Medium compression, preserve reasoning chains
- **GitHub Copilot** - Low compression, keep suggestions intact
- **OpenCode** - High compression for [opencode.io](https://opencode.io) engine, preserve code syntax and reasoning
- **Generic** - Balanced profile for any AI tool

### 📊 Live Dashboard
```
╔══════════════════════════════════════════════════════════╗
║      ✦  MRT Live Monitor  ✦  cursor                     ║
╚══════════════════════════════════════════════════════════╝

  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
  │ Session Stats   │  │ Current Rate    │  │ Token Flow      │
  │                 │  │                 │  │                 │
  │ Duration: 2h15m │  │ Tokens/min: 847 │  │    █            │
  │ Tokens: 142.8K  │  │ Lines/min: 45   │  │   ███           │
  │ Lines:  4,285   │  │ Savings: $0.52/h│  │  █████          │
  │ Duplicates: 892 │  │                 │  │ ███████         │
  └─────────────────┘  └─────────────────┘  └─────────────────┘

  [██████████████████████░░░░░░░░░░░] 67% compression
```

### 🔗 Shell Integration
- Transparent hook into your shell environment
- Works with `bash`, `zsh`, and `fish`
- Zero configuration required for basic usage

---

## 🚀 Quick Start

### Installation via npm (Recommended)

```bash
npm install -g meo-reduce-token
meo --help
```

### Installation via Homebrew

```bash
brew install --formula https://raw.githubusercontent.com/meocode-labs/mrt/main/homebrew-tap/Formula/mrt.rb
```

The formula lives in this repo at `homebrew-tap/Formula/mrt.rb` (see
[`homebrew-tap/README.md`](homebrew-tap/README.md) for the update procedure).

### Installation via cURL

```bash
# macOS/Linux (auto-detects OS/arch, installs to /usr/local/bin)
curl -fsSL https://raw.githubusercontent.com/meocode-labs/mrt/main/install.sh | bash

# Override install location or pin a version:
MRT_INSTALL_DIR=$HOME/bin MRT_VERSION=v1.2.0 \
  curl -fsSL https://raw.githubusercontent.com/meocode-labs/mrt/main/install.sh | bash
```

The installer script lives at [`install.sh`](install.sh) in this repo.

### Installation via Go

```bash
go install github.com/meocode-labs/mrt@latest
```

---

## 📖 Usage

### Initialize Shell Hook

```bash
meo init
# Restart your shell or: source ~/.bashrc
```

### Set Your AI Target

```bash
meo target set cursor        # For Cursor AI
meo target set claude-code   # For Claude Code
meo target set copilot       # For GitHub Copilot
meo target set opencode      # For OpenCode (opencode.io)
meo target info              # Show current config
```

### View Token Savings

```bash
meo gain                     # Show session stats
meo gain --live             # Live updating dashboard
meo gain --since 24h        # Last 24 hours
```

### Launch Dashboard

```bash
meo dashboard                # Full TUI dashboard
meo dashboard --compact      # Compact view
```

### Help

```bash
meo --help                   # General help
meo init --help              # Init options
meo target --help            # Target management
meo gain --help              # Gain command options
```

---

## 🎨 Visual Preview

### Terminal Output Compression

**Before MRT:**
```
[2K blob]
Downloading... 45% [================>                        ] 45MB/s 00:01:23
Downloading... 45% [================>                        ] 45MB/s 00:01:23
Downloading... 45% [================>                        ] 45MB/s 00:01:23
Downloading... 45% [================>                        ] 45MB/s 00:01:23
Downloading... 46% [=================>                       ] 46MB/s 00:01:22
...
```

**After MRT:**
```
Downloading... 45% [================>                        ] 45MB/s 00:01:23
...
```

### Token Savings Dashboard

```
╔══════════════════════════════════════════════════════════╗
║            ✦ MRT Token Savings ✦                         ║
╠══════════════════════════════════════════════════════════╣
║  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐      ║
║  │Tokens Saved  │ │  Cost Saved  │ │Lines Reduced │      ║
║  │   142,857    │ │    $2.34     │ │    4,285     │      ║
║  └──────────────┘ └──────────────┘ └──────────────┘      ║
║                                                          ║
║  Period: Current Session  │  Mode: ● LIVE                ║
╚══════════════════════════════════════════════════════════╝
```

---

## 🏗️ Architecture

```
mrt/
├── cmd/                    # CLI entry points
│   ├── root.go            # Root command (Cobra)
│   └── commands.go        # Subcommands (init, target, gain, dashboard)
├── pkg/
│   ├── compressor/        # Token reduction logic
│   │   └── compressor.go  # ANSI strip, dedup, path collapse
│   ├── monitor/           # TUI display
│   │   └── display.go     # Lipgloss rendering
│   ├── target/            # Profile management
│   │   └── target.go      # AI tool profiles
│   └── tui/               # Interactive UI components
├── main.go                # Binary entry point
├── go.mod
├── package.json          # NPM wrapper
├── install.js             # NPM install script
└── .goreleaser.yaml      # Release configuration
```

---

## 🛠️ Development

### Prerequisites

- Go 1.21+
- Node.js 16+ (for NPM packaging)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/meocode-labs/mrt.git
cd mrt

# Install dependencies
go mod download

# Build binary
go build -o meo main.go

# Run tests
go test ./...

# Run linter
golangci-lint run
```

### Release Process

```bash
# Create a semver tag
git tag v1.0.0
git push origin v1.0.0

# GitHub Actions will automatically:
# - Build binaries for all platforms
# - Create GitHub release with notes
# - Publish to npm registry
```

---

## 📦 NPM Package Structure

```
meo-reduce-token/
├── package.json      # NPM metadata & install script
├── install.js        # Platform-specific binary installer
└── README.md
```

The NPM package is a thin wrapper that downloads the pre-compiled Go binary for your platform on `npm install`.

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a [Pull Request](https://github.com/meocode-labs/mrt/pulls).

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📝 License

Copyright © 2026 [Meo Code Labs](https://meocode.com)

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) CLI framework
- TUI powered by [Charmbracelet Lipgloss](https://github.com/charmbracelet/lipgloss)
- Inspired by [Rust Token Killer (RTK)](https://github.com/...) for the core concept

---

<div align="center">

**Made with ❤️ by [Meo Code Labs](https://meocode.com) | Maintained by [penadidik](https://github.com/penadidik)**

*Reduce tokens. Save costs. Code smarter.*

</div>
