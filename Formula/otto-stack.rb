class DevStack < Formula
  desc "Development stack management tool for streamlined local development automation"
  homepage "https://github.com/otto-nation/otto-stack"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/otto-nation/otto-stack/releases/download/otto-stack-v1.0.0/otto-stack-darwin-amd64"
      sha256 "0a222bad25491a4dbd64f9546e0fffb24b9221372194a1afa76c7a8444c1d5d7"
    end

    on_arm do
      url "https://github.com/otto-nation/otto-stack/releases/download/otto-stack-v1.0.0/otto-stack-darwin-arm64"
      sha256 "6897a44cbed7d6c24723489312355735023e8e98705d37dc84631008448809c0"
    end
  end

  def install
    bin.install "otto-stack-darwin-#{Hardware::CPU.arch}" => "otto-stack"
  end

  test do
    assert_match "otto-stack", shell_output("#{bin}/otto-stack --version")
  end
end
