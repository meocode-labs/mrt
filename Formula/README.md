# Homebrew formula for `mrt`

This directory holds the Homebrew formula for [mrt](https://github.com/meocode-labs/mrt),
exposed as the `mrt` command (renamed from `meo` in v1.3.0). The formula lives
in the mrt source repo (not a separate `homebrew-tap` repo) so the formula and
the source it builds from share a single tag.

## Install

```bash
brew install --formula https://raw.githubusercontent.com/meocode-labs/mrt/main/Formula/mrt.rb
```

Homebrew downloads the source tarball for the release pinned in `url`, builds
with `go build`, and installs the resulting `mrt` binary as `mrt` into
Homebrew's `bin/`.

## Updating to a new release

After tagging a new version (e.g. `v1.3.0`) and waiting for the auto-release
workflow to finish:

1. Compute the SHA256 of the new source tarball:
   ```bash
   curl -fsSL -o /tmp/mrt-vX.Y.Z.tar.gz \
     https://github.com/meocode-labs/mrt/archive/refs/tags/vX.Y.Z.tar.gz
   shasum -a 256 /tmp/mrt-vX.Y.Z.tar.gz
   ```

2. Edit `Formula/mrt.rb`:
   - Bump the `url` to the new tag tarball
   - Replace `sha256` with the value from step 1

3. Commit and push:
   ```bash
   git add Formula/mrt.rb
   git commit -m "mrt: bump formula to vX.Y.Z"
   git push origin main
   ```

## Why not a separate `homebrew-tap` repo?

The user prefers to keep the formula in the same repo as the source so:

- One tag = one source + one formula, no chance of drift
- The formula's `url` tag always matches the source it builds from
- Reviewers see formula bumps alongside the source changes they describe

If this convention changes later, move `homebrew-tap/Formula/mrt.rb` to a
dedicated `meocode-labs/homebrew-tap` repo and update this README.
