# This is a template for the Homebrew formula
# Copy this to your homebrew-tap repository at: Formula/gokanon.rb
# See docs/HOMEBREW.md for detailed setup instructions

class GoKanon < Formula
  desc "Powerful CLI tool for running and comparing Go benchmark tests"
  homepage "https://github.com/alenon/gokanon"
  version "1.0.0"  # Update this with each release
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/alenon/gokanon/releases/download/v#{version}/gokanon-darwin-arm64.tar.gz"
      sha256 "REPLACE_WITH_ARM64_SHA256"  # Calculate with: curl -sL <URL> | shasum -a 256
    else
      url "https://github.com/alenon/gokanon/releases/download/v#{version}/gokanon-darwin-amd64.tar.gz"
      sha256 "REPLACE_WITH_AMD64_SHA256"  # Calculate with: curl -sL <URL> | shasum -a 256
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/alenon/gokanon/releases/download/v#{version}/gokanon-linux-arm64.tar.gz"
      sha256 "REPLACE_WITH_LINUX_ARM64_SHA256"
    else
      url "https://github.com/alenon/gokanon/releases/download/v#{version}/gokanon-linux-amd64.tar.gz"
      sha256 "REPLACE_WITH_LINUX_AMD64_SHA256"
    end
  end

  def install
    # Install the binary with the correct name based on platform
    if OS.mac?
      if Hardware::CPU.arm?
        bin.install "gokanon-darwin-arm64" => "gokanon"
      else
        bin.install "gokanon-darwin-amd64" => "gokanon"
      end
    elsif OS.linux?
      if Hardware::CPU.arm?
        bin.install "gokanon-linux-arm64" => "gokanon"
      else
        bin.install "gokanon-linux-amd64" => "gokanon"
      end
    end
  end

  test do
    system "#{bin}/gokanon", "--version"
  end
end
