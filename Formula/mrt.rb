class Mrt < Formula
  desc "Meo Reduce Token - terminal output compression for AI coding agents"
  homepage "https://github.com/meocode-labs/mrt"
  url "https://github.com/meocode-labs/mrt/archive/refs/tags/v1.2.0.tar.gz"
  # SHA256 of the v1.2.0 source tarball. Update whenever `url` is bumped;
  # see homebrew-tap/README.md for the procedure.
  sha256 "3cbf5b9dce158a5e4865db39495b0290a1dd8707f54fd68c8c4e3b16a973385f"
  license "MIT"

  depends_on "go" => :build

  def install
    ldflags = %W[
      -s -w
      -X github.com/meocode-labs/mrt/cmd.Version=#{version}
      -X github.com/meocode-labs/mrt/cmd.Commit=#{tap.user}
    ]
    # -o mrt is required: `go build ./main.go` defaults to a binary
    # named `main`, but bin.install below expects a file called `mrt`.
    system "go", "build", "-o", "mrt", *std_go_args(ldflags: ldflags), "./main.go"
    bin.install "mrt" => "meo"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/meo --version 2>&1")
  end
end
