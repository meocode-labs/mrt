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
    # Place `-o mrt` AFTER std_go_args so it wins over the `-o=...`
    # that std_go_args injects for the cellar path. Result: binary
    # named `mrt` is written to CWD, then bin.install renames it to
    # `meo` and moves it to the install path.
    system "go", "build", *std_go_args(ldflags: ldflags), "-o", "mrt", "./main.go"
    bin.install "mrt" => "meo"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/meo --version 2>&1")
  end
end
