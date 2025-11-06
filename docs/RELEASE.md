# Release Process

This document describes how to create a new release of GoKanon.

## Overview

The release process is largely automated through GitHub Actions. When you push a version tag, the workflow:
1. Builds binaries for all supported platforms
2. Creates tarballs with checksums
3. Creates a GitHub Release with all artifacts
4. Provides instructions for updating Homebrew

## Prerequisites

- Push access to the repository
- Ability to create tags
- (For Homebrew) Access to the `homebrew-tap` repository

## Release Steps

### 1. Prepare the Release

Before creating a release, ensure:

```bash
# Make sure you're on the main branch
git checkout main
git pull origin main

# Run tests
make test

# Build locally to verify
make build

# Test the binary
./bin/gokanon --version
./bin/gokanon run -pkg=./examples
```

### 2. Update Version Information

If you have version information in the code (e.g., in `main.go`), update it:

```go
var Version = "1.0.0"  // Update this
```

Commit any version updates:

```bash
git add .
git commit -m "Bump version to v1.0.0"
git push origin main
```

### 3. Create and Push Tag

```bash
# Create a new tag (use semantic versioning)
VERSION="v1.0.0"
git tag -a $VERSION -m "Release $VERSION"

# Push the tag to trigger the release workflow
git push origin $VERSION
```

### 4. Monitor the Release Build

1. Go to GitHub Actions: `https://github.com/alenon/gokanon/actions`
2. Watch the "Release" workflow
3. Ensure all builds complete successfully
4. The workflow will create a draft release

### 5. Update the Release Notes

1. Go to Releases: `https://github.com/alenon/gokanon/releases`
2. Find your release
3. Click "Edit"
4. Add a detailed description of changes:

```markdown
## What's New

- Feature 1: Description
- Feature 2: Description
- Bug fix: Description

## Breaking Changes

- List any breaking changes

## Installation

See the installation instructions below.

## Checksums

SHA256 checksums are provided for all binaries.
```

5. Publish the release

### 6. Update Homebrew Formula

After the release is published, update the Homebrew formula:

```bash
# Clone the homebrew-tap repository
git clone https://github.com/alenon/homebrew-tap.git
cd homebrew-tap

# Calculate SHA256 checksums
VERSION="v1.0.0"

echo "macOS Intel:"
curl -sL "https://github.com/alenon/gokanon/releases/download/${VERSION}/gokanon-darwin-amd64.tar.gz" | shasum -a 256

echo "macOS Apple Silicon:"
curl -sL "https://github.com/alenon/gokanon/releases/download/${VERSION}/gokanon-darwin-arm64.tar.gz" | shasum -a 256

echo "Linux x86_64:"
curl -sL "https://github.com/alenon/gokanon/releases/download/${VERSION}/gokanon-linux-amd64.tar.gz" | shasum -a 256

echo "Linux ARM64:"
curl -sL "https://github.com/alenon/gokanon/releases/download/${VERSION}/gokanon-linux-arm64.tar.gz" | shasum -a 256

# Edit Formula/gokanon.rb with:
# - New version number
# - Updated SHA256 checksums

# Commit and push
git add Formula/gokanon.rb
git commit -m "Update gokanon to ${VERSION}"
git push origin main
```

### 7. Test the Installation

Test that users can install the new release:

```bash
# Test direct download
curl -sSL https://raw.githubusercontent.com/alenon/gokanon/main/install.sh | bash

# Test Homebrew (if tap is set up)
brew uninstall gokanon  # if previously installed
brew install alenon/tap/gokanon

# Verify version
gokanon --version
```

### 8. Announce the Release

Consider announcing the release:
- Update the README if needed
- Post on social media
- Notify users in relevant communities

## Manual Release (Alternative)

If you need to create a release manually or the automated workflow fails:

```bash
# Build all binaries
make build-all

# Create release on GitHub
gh release create v1.0.0 \
  --title "Release v1.0.0" \
  --notes "Release notes here" \
  ./bin/gokanon-*

# Or upload to existing release
gh release upload v1.0.0 ./bin/gokanon-*
```

## Troubleshooting

### Workflow Fails

Check the GitHub Actions logs for errors. Common issues:
- Missing secrets or permissions
- Build errors (fix and create a new tag)
- Network timeouts (re-run the workflow)

### Wrong Tag Pushed

If you pushed the wrong tag:

```bash
# Delete local tag
git tag -d v1.0.0

# Delete remote tag
git push origin :refs/tags/v1.0.0

# Delete the release on GitHub if created
gh release delete v1.0.0

# Create correct tag and push again
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### Binary Not Working

If a released binary doesn't work:
1. Test locally with the same build command
2. Check for platform-specific issues
3. Create a patch release with fixes

## Version Numbering

Follow [Semantic Versioning](https://semver.org/):

- **MAJOR** version (v1.0.0 → v2.0.0): Breaking changes
- **MINOR** version (v1.0.0 → v1.1.0): New features, backward compatible
- **PATCH** version (v1.0.0 → v1.0.1): Bug fixes, backward compatible

Examples:
- `v1.0.0` - Initial release
- `v1.1.0` - Added new feature
- `v1.1.1` - Fixed bug
- `v2.0.0` - Breaking API change

## Pre-releases

For testing before official release:

```bash
# Create pre-release tag
git tag -a v1.0.0-beta.1 -m "Beta release"
git push origin v1.0.0-beta.1

# Mark as pre-release in GitHub
gh release create v1.0.0-beta.1 --prerelease
```

## Checklist

Before releasing, verify:

- [ ] All tests pass
- [ ] Code is merged to main branch
- [ ] Version numbers are updated
- [ ] CHANGELOG is updated (if you maintain one)
- [ ] Documentation is up to date
- [ ] Local build works
- [ ] Tag follows semantic versioning
- [ ] Release notes are prepared

After releasing:

- [ ] GitHub release is published
- [ ] All binaries are attached
- [ ] Installation script works
- [ ] Homebrew formula is updated
- [ ] Installation methods are tested
- [ ] Release is announced

## Automation Ideas

Consider automating more of the release process:

1. **Automatic changelog generation** from commit messages
2. **Automated Homebrew formula updates** in the workflow
3. **Release drafter** for automatic release notes
4. **Version bumping** scripts

## Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [GoReleaser](https://goreleaser.com/) - Alternative release automation tool
- [Semantic Versioning](https://semver.org/)
- [Keep a Changelog](https://keepachangelog.com/)
