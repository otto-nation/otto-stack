class OttoStack < Formula
  desc "Development stack management tool for streamlined local development automation"
  homepage "https://github.com/otto-nation/otto-stack"
  version "1.1.0"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/otto-nation/otto-stack/releases/download/otto-stack-v1.1.0/otto-stack-darwin-amd64"
      sha256 "17730a7fe5798bdcd4375310d499b4ae7538229ec5b003c93ed61f0d3b808d14"
    end

    on_arm do
      url "https://github.com/otto-nation/otto-stack/releases/download/otto-stack-v1.1.0/otto-stack-darwin-arm64"
      sha256 "636b0aa865fbd179e816f2f31f886b52cb9200ceef4c2d85c27413d223fe42c4"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/otto-nation/otto-stack/releases/download/otto-stack-v1.1.0/otto-stack-linux-amd64"
      sha256 "d11f901c9e2af3678d0686f86ceece403b529a69cca2e70e85ec2adcb53f98f8"
    end

    on_arm do
      url "https://github.com/otto-nation/otto-stack/releases/download/otto-stack-v1.1.0/otto-stack-linux-arm64"
      sha256 "2afc1c7808b9696e69c84e978deb53ae2852fed12adafae9884b6387edd82525"
    end
  end

  def install
    case Hardware::CPU.arch
    when :x86_64
      arch_suffix = "amd64"
    when :arm64
      arch_suffix = "arm64"
    else
      raise "Unsupported architecture: #{Hardware::CPU.arch}"
    end

    os_name = OS.mac? ? "darwin" : "linux"
    binary_name = "otto-stack-#{os_name}-#{arch_suffix}"

    bin.install binary_name => "otto-stack"
  end

  def caveats
    <<~EOS
      To get started with otto-stack:
        otto-stack init

      For more information:
        otto-stack --help

      Documentation: https://otto-nation.github.io/otto-stack/
    EOS
  end

  test do
    assert_match(/\d+\.\d+\.\d+/, shell_output("#{bin}/otto-stack --version"))

    # Test basic functionality
    system bin/"otto-stack", "--help"

    # Test that the binary is properly linked
    assert_path_exists bin/"otto-stack"
    assert_predicate bin/"otto-stack", :executable?
  end
end
