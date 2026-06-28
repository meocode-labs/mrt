class Mrt < Formula
  desc "Meo Reduce Token - terminal output compression for AI coding agents"
  homepage "https://github.com/meocode-labs/mrt"
  url "https://github.com/meocode-labs/mrt/archive/refs/tags/v1.3.0.tar.gz"
  # SHA256 of the v1.3.0 source tarball. Update whenever `url` is bumped;
  # see Formula/README.md for the procedure.
  sha256 "REPLACE_AT_RELEASE_TIME"
  license "MIT"

  depends_on "go" => :build

  def install
    ldflags = %W[
      -s -w
      -X github.com/meocode-labs/mrt/cmd.Version=#{version}
      -X github.com/meocode-labs/mrt/cmd.Commit=#{tap.user}
    ]
    # The Go binary builds with name `mrt`. v1.3.0 installs it as `mrt`
    # directly (no rename). If a previous install left `meo` around, bin.install
    # will overwrite via the keg link, but the old `meo` symlink at the
    # user's PATH location must be removed manually (see post_install below).
    system "go", "build", *std_go_args(ldflags: ldflags), "./main.go"
  end

  def post_install
    # Migration: v1.2.0 installed the binary as `meo`. v1.3.0 installs as `mrt`.
    # If an old `meo` binary exists alongside the new `mrt`, drop it so the user
    # doesn't get confused by two commands.
    old = bin/"meo"
    if old.exist?
      opoo "v1.3.0 renamed \`meo\` to \`mrt\`. Removing old #{old}."
      old.delete
    end
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/mrt --version 2>&1")
  end
end
