class OttoStack < Formula
  desc "Development stack management tool for streamlined local development automation"
  homepage "https://github.com/otto-nation/otto-stack"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/otto-nation/otto-stack/releases/latest/download/otto-stack-darwin-amd64"
      # sha256 will be updated automatically by release process
    end

    on_arm do
      url "https://github.com/otto-nation/otto-stack/releases/latest/download/otto-stack-darwin-arm64"
      # sha256 will be updated automatically by release process
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/otto-nation/otto-stack/releases/latest/download/otto-stack-linux-amd64"
      # sha256 will be updated automatically by release process
    end

    on_arm do
      url "https://github.com/otto-nation/otto-stack/releases/latest/download/otto-stack-linux-arm64"
      # sha256 will be updated automatically by release process
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

      Documentation: https://github.com/otto-nation/otto-stack/tree/main/docs-site
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
