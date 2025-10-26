# Release Branch Workflow

This document describes the release branch workflow for Karo (Kubernetes Alert Reaction Operator).

## Overview

The release process uses dedicated release branches to manage releases, providing better control and isolation for release preparations.

## Workflow

### 1. Create Release Branch

Use the prepare-release script to create a release branch:

```bash
# Run pre-release checks
./scripts/prepare-release.sh --check

# Create release branch for v1.0.0
./scripts/prepare-release.sh v1.0.0

# Create pre-release branch
./scripts/prepare-release.sh v1.0.0-alpha.1 --prerelease
```

### 2. Automatic Draft Release

When a release branch is created and pushed:

1. **Branch Creation**: Creates `release/v1.0.0` branch
2. **Workflow Trigger**: Push to release branch triggers draft-release workflow
3. **Version Extraction**: Version is extracted from branch name
4. **CHANGELOG Processing**: Unreleased changes are moved to version section
5. **Draft Creation**: GitHub draft release is created with artifacts

### 3. Release Branch Benefits

- **Isolation**: Release preparation doesn't interfere with main development
- **Collaboration**: Team can collaborate on release-specific changes
- **Control**: Manual approval required before publishing release
- **Traceability**: Clear history of what went into each release

### 4. Branch Naming Convention

- **Standard Release**: `release/v1.0.0`
- **Pre-release**: `release/v1.0.0-alpha.1`
- **Release Candidate**: `release/v1.0.0-rc.1`

### 5. Release Process Steps

1. **Preparation**: 
   ```bash
   ./scripts/prepare-release.sh v1.0.0
   ```

2. **Automatic Workflow**: 
   - Draft release created
   - CHANGELOG updated
   - Artifacts built

3. **Review**: 
   - Check draft release
   - Test artifacts if needed
   - Make final adjustments to release branch

4. **Publish**: 
   - Publish draft release on GitHub
   - Triggers Helm chart release pipeline

5. **Cleanup**: 
   - Merge release branch back to main (optional)
   - Delete release branch (optional)

## Manual Release Branch Creation

You can also create release branches manually:

```bash
# Create and switch to release branch
git checkout -b release/v1.0.0

# Make any release-specific changes
git add .
git commit -m "chore: prepare release v1.0.0"

# Push to trigger workflow
git push -u origin release/v1.0.0
```

## Troubleshooting

### Branch Already Exists
If a release branch already exists, either:
- Use a different version number
- Delete the existing branch first
- Use the existing branch and push new commits

### Workflow Not Triggered
Ensure:
- Branch name follows `release/v*` pattern
- Branch is pushed to remote
- GitHub Actions are enabled

### Invalid Version Format
Version must follow semantic versioning:
- `v1.0.0` (standard release)
- `v1.0.0-alpha.1` (pre-release)
- `v1.0.0-rc.1` (release candidate)

## Configuration

The workflow is configured in `.github/workflows/draft-release.yml` and triggers on:
- Push to `release/**` branches
- Manual workflow dispatch

## Next Steps

After the draft release is created:
1. Review release notes and artifacts
2. Test the release if needed
3. Publish the release to trigger Helm chart automation
4. The published release triggers the Helm chart pipeline