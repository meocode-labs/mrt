class Mrt < Formula
  desc "Meo Reduce Token - terminal output compression for AI coding agents"
  homepage "https://github.com/meocode-labs/mrt"
  url "https://github.com/meocode-labs/mrt/archive/refs/tags/v1.2.0.tar.gz"
  # SHA256 of the v1.2.0 source tarball. Must be updated whenever `url`
  # is bumped. See homebrew-tap/README.md for the update procedure.
  # The placeholder below is invalid on purpose so `brew install` fails
  # loudly until the maintainer replaces it.
  sha256 "0000000000000000000000000000000000000000000000000000000000000000"
  license "MIT"

  depends_on "go" => :build

  def install
    ldflags = %W[
      -s -w
      -X github.com/meocode-labs/mrt/cmd.Version=#{version}
      -X github.com/meocode-labs/mrt/cmd.Commit=#{tap.user}
    ]
    system "go", "build", *std_go_args(ldflags: ldflags), "./main.go"
    bin.install "mrt" => "meo"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/meo --version 2>&1")
  end
end
