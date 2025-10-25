# Helm Chart Release Automation

This document describes the automated Helm chart release system for Karo (Kubernetes Alert Reaction Operator).

## Overview

The project now includes comprehensive automation for Helm chart releases that triggers when new Git tags are created. This ensures that every release includes properly versioned and validated Helm charts.

## Components

### 1. GitHub Actions Workflow (`.github/workflows/helm-release.yml`)

**Triggers**: Push events to tags matching `v*` pattern (e.g., `v1.0.0`, `v2.1.3-alpha.1`)

**Key Features**:
- Automatic version extraction from Git tags
- Chart metadata updates (version, appVersion, URLs)
- Comprehensive validation (lint, template rendering, YAML validation)
- Multi-format publishing (OCI registry, GitHub releases, Helm repository)
- Automated release notes generation
- Pre-release detection and handling

**Outputs**:
- Packaged chart files (`.tgz`)
- OCI artifacts in GitHub Container Registry
- GitHub Releases with installation instructions
- Updated Helm repository on `gh-pages` branch

### 2. Chart Update Script (`scripts/update-chart.sh`)

**Purpose**: Manual chart version management and local testing

**Features**:
- Version validation and format checking
- Automatic next version generation (major/minor/patch)
- Dry-run mode for testing changes
- Current version display
- Backup and rollback capabilities
- Integration with validation pipeline

**Usage Examples**:
```bash
# Show current versions
./scripts/update-chart.sh --current

# Update to specific version
./scripts/update-chart.sh 1.2.3

# Generate next patch version
./scripts/update-chart.sh --next patch

# Test changes without applying
./scripts/update-chart.sh --dry-run 2.0.0
```

### 3. Chart Validation Script (`scripts/validate-chart.sh`)

**Purpose**: Comprehensive chart validation and quality checks

**Validations**:
- Chart.yaml structure and required fields
- Version format validation
- values.yaml syntax and structure
- Template directory completeness
- Helm lint checks
- Template rendering tests
- Security best practices assessment
- YAML syntax validation

**Integration**: Used by both manual workflows and CI/CD automation

### 4. Enhanced Chart Metadata (`charts/karo/Chart.yaml`)

**Improvements**:
- Comprehensive keyword tagging
- Artifact Hub annotations
- Security and license metadata
- Kubernetes version compatibility
- Detailed descriptions and links
- Maintainer information
- OpenShift compatibility annotations

## Release Process

### Draft-First Release Process (Recommended)

The release process is now split into two phases for better control:

1. **Draft Release Creation** - Prepares and validates releases
2. **Helm Chart Release** - Publishes charts when draft is published

#### Phase 1: Create Draft Release

**Automatic (Push to main):**
```bash
# Ensure all changes are committed and CHANGELOG is updated
./scripts/manage-changelog.sh add added "New webhook endpoint"
./scripts/manage-changelog.sh add fixed "Memory leak in controller"

git add -A
git commit -m "feat: add new webhook endpoint and fix memory leak"
git push origin main

# Automation will create draft release with auto-generated version
```

**Manual (Specific version):**
```bash
# Prepare release with specific version
./scripts/prepare-release.sh v1.0.0

# Or use the GitHub CLI directly
gh workflow run draft-release.yml -f version=v1.0.0 -f prerelease=false
```

**What happens during draft creation:**
- Version validation and conflict checking
- CHANGELOG processing and release notes generation
- Release artifact preparation (binaries, configs, charts)
- Draft release creation in GitHub
- CHANGELOG updates for next version

#### Phase 2: Publish Release (Triggers Helm Automation)

1. **Review Draft Release**:
   - Check generated release notes
   - Verify attached artifacts
   - Test release if needed

2. **Publish Release**:
   - Navigate to GitHub Releases
   - Edit the draft release
   - Click "Publish release"

3. **Helm Automation Triggers**:
   - Chart metadata updates
   - Comprehensive validation
   - Multi-format packaging and publishing
   - OCI registry upload
   - Helm repository updates

### Manual Release Management

**Release Preparation:**
```bash
# Check if ready for release
./scripts/prepare-release.sh --check

# Suggest next version based on commits
./scripts/prepare-release.sh --suggest

# Prepare specific version with validation
./scripts/prepare-release.sh v1.0.0

# Force preparation (skip checks)
./scripts/prepare-release.sh v1.0.0 --force

# Dry run to see what would happen
./scripts/prepare-release.sh v1.0.0 --dry-run
```

**CHANGELOG Management:**
```bash
# Add changelog entries
./scripts/manage-changelog.sh add added "New feature description"
./scripts/manage-changelog.sh add fixed "Bug fix description"

# Show unreleased changes
./scripts/manage-changelog.sh show

# Validate changelog format
./scripts/manage-changelog.sh validate

# Release version in changelog
./scripts/manage-changelog.sh release 1.0.0
```

**Manual Chart Testing:**
```bash
# Validate chart
./scripts/validate-chart.sh

# Update chart version manually
./scripts/update-chart.sh 1.0.0-dev

# Package and test
helm package charts/karo/
helm install test-release ./karo-1.0.0-dev.tgz
```

## Installation Methods

After a release is published, users have multiple installation options:

### 1. OCI Registry (Recommended)
```bash
helm install karo \
  oci://ghcr.io/dudizimber/charts/karo \
  --version 1.0.0
```

### 2. Helm Repository
```bash
helm repo add karo \
  https://dudizimber.github.io/karo/
helm repo update
helm install karo \
  karo/karo --version 1.0.0
```

### 3. GitHub Release Assets
```bash
curl -L https://github.com/dudizimber/karo/releases/download/v1.0.0/karo-1.0.0.tgz -o chart.tgz
helm install karo ./chart.tgz
```

## Version Management

### Semantic Versioning

The system follows [Semantic Versioning](https://semver.org/):
- **MAJOR** (X.0.0): Breaking changes to chart API or defaults
- **MINOR** (X.Y.0): New features, backwards compatible
- **PATCH** (X.Y.Z): Bug fixes, security updates

### Pre-release Versions

Supported pre-release formats:
- `v1.0.0-alpha.1` - Alpha releases
- `v1.0.0-beta.2` - Beta releases  
- `v1.0.0-rc.1` - Release candidates

Pre-releases are automatically marked as such in GitHub releases.

## Quality Assurance

### Automated Checks

Every release includes:
- **Syntax Validation**: Chart.yaml and values.yaml syntax
- **Lint Checks**: `helm lint` validation
- **Template Rendering**: Test rendering with various value combinations
- **YAML Validation**: Ensure generated manifests are valid YAML
- **Security Scanning**: Basic security best practices assessment

### Manual Validation

Developers can run the same checks locally:
```bash
# Full validation suite
./scripts/validate-chart.sh

# Specific validations
helm lint charts/karo/
helm template test charts/karo/ | yq eval '.'
```

## Security Considerations

### Chart Security

- **Non-root execution**: Default security contexts use non-root user
- **Resource limits**: Recommended resource constraints
- **RBAC**: Minimal required permissions
- **Secret management**: Secure handling of sensitive data

### Release Security

- **Signed releases**: GitHub releases are signed and verifiable
- **OCI artifacts**: Container registry provides artifact signing
- **Access control**: Release process requires write permissions
- **Audit trail**: All changes tracked in Git history

## Troubleshooting

### Common Issues

1. **Draft Release Creation Fails**:
   - Check working directory is clean
   - Verify version format (vX.Y.Z)
   - Ensure version doesn't already exist
   - Review draft-release workflow logs

2. **Helm Release Workflow Fails**:
   - Ensure draft release was published (not just created)
   - Verify Chart.yaml syntax
   - Check OCI registry permissions
   - Review helm-release workflow logs

3. **Chart Validation Errors**:
   - Run `./scripts/validate-chart.sh` locally
   - Check template syntax with `helm template`
   - Verify values.yaml structure

4. **CHANGELOG Issues**:
   - Validate format with `./scripts/manage-changelog.sh validate`
   - Ensure [Unreleased] section exists
   - Check for duplicate version entries

### Debug Commands

```bash
# Release preparation
./scripts/prepare-release.sh --check          # Check release readiness
./scripts/prepare-release.sh --suggest        # Suggest next version

# CHANGELOG management
./scripts/manage-changelog.sh show            # Show unreleased changes
./scripts/manage-changelog.sh validate       # Validate format

# Chart validation
./scripts/validate-chart.sh                  # Full chart validation
./scripts/update-chart.sh --current          # Check current versions

# Workflow status
gh run list --workflow=draft-release.yml     # Draft release status
gh run list --workflow=helm-release.yml      # Helm release status
gh release list                              # List all releases
```

## Maintenance

### Regular Tasks

- **Dependency Updates**: Keep base images and dependencies current
- **Security Patches**: Apply security updates promptly
- **Documentation**: Keep chart README and values documentation updated
- **Testing**: Validate chart functionality with new Kubernetes versions

### Monitoring

- **Release Success**: Monitor GitHub Actions for release failures
- **Installation Metrics**: Track download and installation patterns
- **User Feedback**: Address issues reported in GitHub Issues

## Future Enhancements

### Planned Improvements

- **Chart Testing**: Integration with chart-testing tools
- **Multi-cluster Validation**: Test against different Kubernetes versions
- **Automated Security Scanning**: Integration with container security tools
- **Performance Metrics**: Chart installation and runtime performance tracking
- **Documentation Generation**: Automated chart documentation updates

This automation system ensures reliable, consistent, and secure Helm chart releases while maintaining high quality standards and comprehensive validation.