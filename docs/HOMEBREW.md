# Homebrew Tap Setup

This guide explains how to set up and maintain the Homebrew tap for GoKanon.

## Overview

Homebrew "taps" are third-party repositories that allow users to install packages not in the main Homebrew repository. For GoKanon, we'll create a tap at `alenon/homebrew-tap`.

## Setting Up the Tap Repository

### 1. Create the Tap Repository

Create a new GitHub repository named `homebrew-tap` under your account (e.g., `alenon/homebrew-tap`).

**Important naming convention:** Homebrew requires the repository to be named `homebrew-<tap-name>`. For a tap named "tap", the repository must be `homebrew-tap`.

### 2. Create the Formula

In the `homebrew-tap` repository, create a file at `Formula/gokanon.rb` with the following content:

```ruby
class GoKanon < Formula
  desc "Powerful CLI tool for running and comparing Go benchmark tests"
  homepage "https://github.com/alenon/gokanon"
  version "1.0.0"  # Update this with each release
  license "MIT"    # Update if different

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/alenon/gokanon/releases/download/v#{version}/gokanon-darwin-arm64.tar.gz"
      sha256 "REPLACE_WITH_ARM64_SHA256"  # Calculate after release
    else
      url "https://github.com/alenon/gokanon/releases/download/v#{version}/gokanon-darwin-amd64.tar.gz"
      sha256 "REPLACE_WITH_AMD64_SHA256"  # Calculate after release
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
    bin.install "gokanon-darwin-arm64" => "gokanon" if Hardware::CPU.arm? && OS.mac?
    bin.install "gokanon-darwin-amd64" => "gokanon" if Hardware::CPU.intel? && OS.mac?
    bin.install "gokanon-linux-arm64" => "gokanon" if Hardware::CPU.arm? && OS.linux?
    bin.install "gokanon-linux-amd64" => "gokanon" if Hardware::CPU.intel? && OS.linux?
  end

  test do
    system "#{bin}/gokanon", "--version"
  end
end
```

### 3. Initial Setup

```bash
# Create the homebrew-tap repository
mkdir homebrew-tap
cd homebrew-tap

# Create the Formula directory
mkdir -p Formula

# Add the formula file (gokanon.rb)
# (Create the file with the template above)

# Initialize git
git init
git add .
git commit -m "Initial commit: Add gokanon formula"

# Push to GitHub
git remote add origin https://github.com/alenon/homebrew-tap.git
git branch -M main
git push -u origin main
```

## Updating the Formula for New Releases

When you release a new version of GoKanon, you need to update the Homebrew formula:

### 1. Calculate SHA256 Checksums

After creating a release, calculate the SHA256 for each platform:

```bash
# macOS Intel
curl -sL https://github.com/alenon/gokanon/releases/download/v1.0.0/gokanon-darwin-amd64.tar.gz | shasum -a 256

# macOS Apple Silicon
curl -sL https://github.com/alenon/gokanon/releases/download/v1.0.0/gokanon-darwin-arm64.tar.gz | shasum -a 256

# Linux x86_64
curl -sL https://github.com/alenon/gokanon/releases/download/v1.0.0/gokanon-linux-amd64.tar.gz | shasum -a 256

# Linux ARM64
curl -sL https://github.com/alenon/gokanon/releases/download/v1.0.0/gokanon-linux-arm64.tar.gz | shasum -a 256
```

### 2. Update the Formula

```bash
cd homebrew-tap
# Edit Formula/gokanon.rb
# - Update version number
# - Update SHA256 checksums
git add Formula/gokanon.rb
git commit -m "Update gokanon to v1.0.0"
git push
```

### 3. Automated Updates (Optional)

You can create a script to automate formula updates:

```bash
#!/bin/bash
# update-formula.sh

VERSION=$1
if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

echo "Calculating checksums for version ${VERSION}..."

DARWIN_AMD64_SHA=$(curl -sL "https://github.com/alenon/gokanon/releases/download/${VERSION}/gokanon-darwin-amd64.tar.gz" | shasum -a 256 | awk '{print $1}')
DARWIN_ARM64_SHA=$(curl -sL "https://github.com/alenon/gokanon/releases/download/${VERSION}/gokanon-darwin-arm64.tar.gz" | shasum -a 256 | awk '{print $1}')
LINUX_AMD64_SHA=$(curl -sL "https://github.com/alenon/gokanon/releases/download/${VERSION}/gokanon-linux-amd64.tar.gz" | shasum -a 256 | awk '{print $1}')
LINUX_ARM64_SHA=$(curl -sL "https://github.com/alenon/gokanon/releases/download/${VERSION}/gokanon-linux-arm64.tar.gz" | shasum -a 256 | awk '{print $1}')

echo "Updating formula..."
# Use sed or your preferred method to update the formula
# This is a simplified example - adjust paths as needed
sed -i '' "s/version \".*\"/version \"${VERSION#v}\"/" Formula/gokanon.rb
sed -i '' "s/REPLACE_WITH_AMD64_SHA256/${DARWIN_AMD64_SHA}/" Formula/gokanon.rb
sed -i '' "s/REPLACE_WITH_ARM64_SHA256/${DARWIN_ARM64_SHA}/" Formula/gokanon.rb
# ... update other SHAs

echo "Formula updated successfully!"
```

## Usage

Once the tap is set up, users can install GoKanon with:

```bash
# Add the tap
brew tap alenon/tap

# Install gokanon
brew install gokanon

# Or in one command
brew install alenon/tap/gokanon
```

## Testing the Formula

Before pushing updates, test the formula locally:

```bash
# Test installation
brew install --build-from-source Formula/gokanon.rb

# Run audit
brew audit --strict Formula/gokanon.rb

# Test uninstall
brew uninstall gokanon
```

## Troubleshooting

### Formula Not Found

Make sure:
- Repository is named `homebrew-tap` (not just `tap`)
- Formula file is at `Formula/gokanon.rb`
- Repository is public

### SHA256 Mismatch

If users report SHA256 mismatches:
1. Verify the release files haven't been modified
2. Recalculate the SHA256 checksums
3. Update the formula with correct values

### Installation Fails

Check:
- URLs are correct and accessible
- Binary names match what's in the tarball
- The `install` block correctly handles all platforms

## Resources

- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [Homebrew Taps Documentation](https://docs.brew.sh/Taps)
- [How to Create and Maintain a Tap](https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap)
