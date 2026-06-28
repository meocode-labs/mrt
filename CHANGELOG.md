# Changelog

All notable changes to MRT (Meo Reduce Token) are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.3.0] — 2026-06-28

### ⚠ BREAKING CHANGES

- **CLI command renamed**: the binary is now `mrt` instead of `meo`.
  - Old: `meo --help`, `meo init`, `meo target set cursor`, etc.
  - New: `mrt --help`, `mrt init`, `mrt target set cursor`, etc.
  - **Migration**: delete the old binary manually:
    - npm global: `rm "$(npm prefix -g)/bin/meo"` and `npm uninstall -g meo-reduce-token`
    - curl install: `sudo rm /usr/local/bin/meo`
    - brew: `brew uninstall mrt && brew install mrt`
    - Or set `MRT_AUTO_MIGRATE=1` before install to auto-remove.
  - **NPM package renamed**: `meo-reduce-token` → `mrt`. `npm install -g mrt`.
  - **Aliases**: any shell aliases or scripts calling `meo` need to be updated.
  - **Shell hook**: the `source "$HOME/mrt-hook.sh"` line written by `mrt init`
    (previously `meo init`) is unchanged — it still works.

### Changed

- Homebrew Formula now installs binary directly as `mrt` (no rename). Added
  `post_install` hook that auto-removes any leftover `meo` from v1.2.x.
- `install.sh` prints a migration notice when old `meo` binary detected.
- NPM postinstall (`install.js`) prints migration notice + cleans up `bin/meo`
  leftover inside the package itself.
- Package metadata: name `meo-reduce-token` → `mrt`, version `1.2.0` → `1.3.0`.
- Added `LICENSE` file (was missing — release artifacts were incomplete).
- Added `CHANGELOG.md` (this file).
- Added `bin/mrt` lazy-fallback shim so `npx mrt` works before postinstall runs.
- Auto-release workflow release-name uses new package name.

## [1.2.0] — 2026-06-24

### Added

- `install.sh` — universal curl installer (auto-detects OS/arch, supports
  `MRT_INSTALL_DIR` and `MRT_VERSION` env vars, falls back to `sudo` if
  install dir is not writable).
- Homebrew tap support: `brew tap meocode-labs/mrt https://github.com/meocode-labs/mrt`
  + `brew install mrt` (formula at `Formula/mrt.rb`).
- New `opencode` target profile (for opencode.io engine).
- `docs(install): warn about PATH conflicts between install methods` —
  README now explains how to detect multiple installations.

### Fixed

- Logo replaced from placeholder "DUNGEON" banner to proper `MRT` opening logo
  (was rendering as 7 random letters in earlier versions).
- Homebrew formula: `-o mrt` placed after `std_go_args` so it wins over the
  cellar-path `-o` flag injected by Homebrew itself.

## [1.1.0] — 2026-06-23

### Fixed

- `package.json` `bugs` object syntax (`]` → `}`).
- GitHub org name corrected to `meocode-labs/mrt` (was `meocodelabs/mrt`).

## [1.0.2] — 2026-06-23

### Fixed

- NPM installer URL corrected to use `meocode-labs/mrt` repo.
- Various CI workflow syntax fixes.

## [1.0.1] — 2026-06-22

### Fixed

- `package.json` `bugs` object syntax (`]` → `}`).
- CI workflow `strategy` moved to job level.

## [1.0.0] — 2026-06-22

### Added

- Initial public release.
- Core compression: ANSI stripping, duplicate removal, path collapsing.
- AI tool target profiles: cursor, claude-code, copilot, default.
- Live dashboard via `mrt dashboard` (was `meo dashboard` pre-1.3.0).
- Token savings dashboard via `mrt gain` (was `meo gain`).
- Shell integration via `mrt init` (was `meo init`).
- Distribution via NPM (`meo-reduce-token`, now `mrt`), Homebrew (`meocode-labs/tap/mrt`),
  cURL (`install.sh`), and `go install`.

[1.3.0]: #130--2026-06-28
[1.2.0]: #120--2026-06-24
[1.1.0]: #110--2026-06-23
[1.0.2]: #102--2026-06-23
[1.0.1]: #101--2026-06-22
[1.0.0]: #100--2026-06-22