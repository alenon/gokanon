# Release Process

This document describes how to create a new release of gokanon.

## Overview

Gokanon uses semantic versioning (MAJOR.MINOR.PATCH) and automated GitHub Actions workflows for releases.

## Version Management

The current version is stored in the `VERSION` file at the project root. This file contains the version number in the format `MAJOR.MINOR.PATCH` (e.g., `0.1.0`).

### Updating the Version

Use the provided Makefile targets to bump the version:

```bash
# Bump patch version (0.0.X) - for bug fixes
make version-bump-patch

# Bump minor version (0.X.0) - for new features
make version-bump-minor

# Bump major version (X.0.0) - for breaking changes
make version-bump-major

# Set a specific version
make version-set VERSION=1.2.3

# Display current version
make version
```

## Release Workflow

### 1. Prepare the Release

1. Ensure all changes for the release are merged to the main branch
2. Update the version number:
   ```bash
   make version-bump-minor  # or patch/major as appropriate
   ```

3. Commit the VERSION file:
   ```bash
   git add VERSION
   git commit -m "Bump version to $(cat VERSION)"
   git push origin main
   ```

### 2. Create and Push the Release Tag

Create a git tag for the release:

```bash
# This will create a tag with the version from the VERSION file
make tag-release

# Push the tag to trigger the release workflow
git push origin v$(cat VERSION)
```

Or manually:

```bash
# Get current version
VERSION=$(cat VERSION)

# Create annotated tag
git tag -a "v${VERSION}" -m "Release v${VERSION}"

# Push the tag
git push origin "v${VERSION}"
```

### 3. Automated Release Process

When you push a tag starting with `v` (e.g., `v0.1.0`), the GitHub Actions release workflow will automatically:

1. **Build binaries** for multiple platforms:
   - Linux (amd64, arm64)
   - macOS (amd64 Intel, arm64 Apple Silicon)
   - Windows (amd64)

2. **Generate release notes** that include:
   - Changes since the previous version (commit history)
   - Comparison link between versions
   - Installation instructions for all platforms
   - Checksums for verification

3. **Create a GitHub Release** with:
   - All platform binaries (tar.gz for Unix, zip for Windows)
   - SHA256 checksums for each binary
   - Detailed release notes

4. **Provide Homebrew update instructions** (manual step required)

### 4. Manual Release Trigger

You can also trigger a release manually from the GitHub Actions UI:

1. Go to Actions → Release workflow
2. Click "Run workflow"
3. Enter the version (e.g., `v0.1.0`)
4. Click "Run workflow"

## Release Notes

Release notes are automatically generated and include:

- **What's Changed**: Commit history since the previous release
- **Installation Instructions**: Platform-specific installation commands
- **Checksums**: SHA256 verification information
- **Comparison Link**: GitHub comparison between versions

The release notes are generated from:
- Commit messages between the previous and current tag
- Static installation instructions for all platforms

## Version in Binaries

The version information is embedded in the binaries during build time using Go ldflags. This includes:

- `Version`: The semantic version (from the VERSION file or git tag)
- `GitCommit`: The git commit hash
- `BuildDate`: The timestamp of the build

Users can check the version with:

```bash
gokanon version
# or
gokanon --version
```

## Best Practices

1. **Semantic Versioning**: Follow semantic versioning principles:
   - PATCH: Bug fixes and minor changes
   - MINOR: New features, backward compatible
   - MAJOR: Breaking changes

2. **Commit Messages**: Write clear commit messages that will appear in release notes
   - Use conventional commit format when possible (feat:, fix:, docs:, etc.)
   - Keep messages concise but descriptive

3. **Testing**: Before creating a release tag:
   - Run `make test` to ensure all tests pass
   - Run `make build` to verify the build works
   - Run `make check` for quick validation

4. **Release Cadence**:
   - Create releases when there are meaningful changes
   - Consider creating patch releases for important bug fixes
   - Bundle related features into minor releases

## Troubleshooting

### Build Failures

If the release workflow fails during the build step:

1. Check the GitHub Actions logs for specific errors
2. Test the build locally: `make build-all`
3. Verify the VERSION file contains a valid version number
4. Ensure go.mod is up to date: `make mod-tidy`

### Tag Already Exists

If you need to recreate a tag:

```bash
# Delete local tag
git tag -d v0.1.0

# Delete remote tag (use with caution!)
git push origin :refs/tags/v0.1.0

# Create new tag
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

### Release Notes Missing Changes

Release notes are generated from git commits. If changes are missing:

1. Ensure commits are on the main branch
2. Check that commits are between the previous and current tag
3. Verify git history: `git log <previous-tag>..HEAD`

## Updating Homebrew Formula

After a successful release, update the Homebrew formula:

1. The release workflow will print instructions with checksums
2. Update the formula in the `homebrew-tap` repository
3. Update version and SHA256 checksums for both architectures
4. Submit a PR to the homebrew-tap repository

## Related Files

- `VERSION`: Current version number
- `.github/workflows/release.yml`: Release automation workflow
- `Makefile`: Build and version management targets
- `internal/cli/version.go`: Version variables for the application
