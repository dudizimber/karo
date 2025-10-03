#!/bin/bash

# prepare-release.sh - Helper script for preparing releases
# This script helps with manual release preparation and validation

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${BLUE}[RELEASE PREP]${NC} $1"
}

# Function to validate version format
validate_version() {
    local version=$1
    if [[ ! $version =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.-]+)?$ ]]; then
        print_error "Invalid version format: $version"
        print_error "Expected format: vX.Y.Z or vX.Y.Z-suffix (e.g., v1.0.0, v1.0.0-alpha.1)"
        return 1
    fi
    return 0
}

# Function to check if version already exists
check_version_exists() {
    local version=$1
    
    # Check local tags
    if git tag -l | grep -q "^${version}$"; then
        print_error "Tag $version already exists locally"
        return 1
    fi
    
    # Check remote tags
    git fetch --tags >/dev/null 2>&1 || true
    if git ls-remote --tags origin | grep -q "refs/tags/${version}$"; then
        print_error "Tag $version already exists on remote"
        return 1
    fi
    
    # Check if release exists on GitHub
    if command -v gh >/dev/null 2>&1; then
        if gh release view "$version" >/dev/null 2>&1; then
            print_error "Release $version already exists on GitHub"
            return 1
        fi
    fi
    
    return 0
}

# Function to get latest version
get_latest_version() {
    git fetch --tags >/dev/null 2>&1 || true
    git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0"
}

# Function to generate next version
generate_next_version() {
    local current_version=$1
    local increment_type=${2:-patch}
    
    # Remove 'v' prefix and any pre-release suffix for calculation
    local base_version
    base_version=$(echo "$current_version" | sed 's/^v//' | cut -d'-' -f1)
    
    IFS='.' read -ra VERSION_PARTS <<< "$base_version"
    local major=${VERSION_PARTS[0]}
    local minor=${VERSION_PARTS[1]}
    local patch=${VERSION_PARTS[2]}
    
    case $increment_type in
        major)
            echo "v$((major + 1)).0.0"
            ;;
        minor)
            echo "v${major}.$((minor + 1)).0"
            ;;
        patch|*)
            echo "v${major}.${minor}.$((patch + 1))"
            ;;
    esac
}

# Function to analyze commit types
analyze_commits() {
    local from_tag=$1
    local breaking_changes feat_changes fix_changes
    
    if [[ "$from_tag" != "HEAD" ]]; then
        breaking_changes=$(git log --pretty=format:"%s" ${from_tag}..HEAD | grep -c -E "^feat!|^fix!|BREAKING CHANGE" || true)
        feat_changes=$(git log --pretty=format:"%s" ${from_tag}..HEAD | grep -c "^feat" || true)
        fix_changes=$(git log --pretty=format:"%s" ${from_tag}..HEAD | grep -c "^fix" || true)
    else
        breaking_changes=$(git log --pretty=format:"%s" | grep -c -E "^feat!|^fix!|BREAKING CHANGE" || true)
        feat_changes=$(git log --pretty=format:"%s" | grep -c "^feat" || true)
        fix_changes=$(git log --pretty=format:"%s" | grep -c "^fix" || true)
    fi
    
    if [[ $breaking_changes -gt 0 ]]; then
        echo "major"
    elif [[ $feat_changes -gt 0 ]]; then
        echo "minor"
    else
        echo "patch"
    fi
}

# Function to check working directory
check_working_directory() {
    if [[ -n "$(git status --porcelain)" ]]; then
        print_error "Working directory is not clean"
        echo "Please commit or stash your changes before preparing a release:"
        git status --short
        return 1
    fi
    return 0
}

# Function to validate CHANGELOG
validate_changelog() {
    local changelog_file="$PROJECT_ROOT/CHANGELOG.md"
    
    if [[ ! -f "$changelog_file" ]]; then
        print_warning "CHANGELOG.md not found"
        return 1
    fi
    
    if ! grep -q "## \[Unreleased\]" "$changelog_file"; then
        print_warning "No [Unreleased] section found in CHANGELOG.md"
        print_warning "Consider adding changes to the unreleased section before releasing"
        return 1
    fi
    
    # Check if unreleased section has meaningful content
    local unreleased_content
    # Extract content between [Unreleased] and the next ## section (or end of file)
    unreleased_content=$(awk '/## \[Unreleased\]/{flag=1; next} /^## \[/{flag=0} flag' "$changelog_file" | grep -v '^[[:space:]]*$' | head -10)
    
    if [[ -z "$unreleased_content" ]]; then
        print_warning "Unreleased section appears to be empty"
        print_warning "Consider adding meaningful changes before releasing"
        return 1
    fi
    
    # Check if it only contains placeholder content
    if echo "$unreleased_content" | grep -q -- "- Prepare for next release" && [[ $(echo "$unreleased_content" | wc -l) -eq 1 ]]; then
        print_warning "Unreleased section has only placeholder content"
        print_warning "Consider adding meaningful changes before releasing"
        return 1
    fi
    
    return 0
}

# Function to run pre-release checks
run_pre_release_checks() {
    print_status "Running pre-release checks..."
    
    local checks_passed=0
    local total_checks=0
    
    # Check 1: Working directory clean
    ((total_checks++))
    if check_working_directory; then
        print_status "✅ Working directory is clean"
        ((checks_passed++))
    else
        print_error "❌ Working directory check failed"
    fi
    
    # Check 2: On main branch
    ((total_checks++))
    local current_branch
    current_branch=$(git branch --show-current)
    if [[ "$current_branch" == "main" ]]; then
        print_status "✅ On main branch"
        ((checks_passed++))
    else
        print_warning "⚠️  Not on main branch (currently on: $current_branch)"
        print_warning "Releases should typically be made from main branch"
    fi
    
    # Check 3: Up to date with remote
    ((total_checks++))
    git fetch >/dev/null 2>&1 || true
    local local_commit remote_commit
    local_commit=$(git rev-parse HEAD)
    remote_commit=$(git rev-parse origin/main 2>/dev/null || echo "")
    
    if [[ "$local_commit" == "$remote_commit" ]]; then
        print_status "✅ Up to date with remote"
        ((checks_passed++))
    else
        print_warning "⚠️  Local branch may not be up to date with remote"
        print_warning "Consider pulling latest changes"
    fi
    
    # Check 4: Tests pass
    ((total_checks++))
    if command -v make >/dev/null 2>&1 && make test >/dev/null 2>&1; then
        print_status "✅ Tests pass"
        ((checks_passed++))
    else
        print_warning "⚠️  Could not run tests or tests failed"
    fi
    
    # Check 5: CHANGELOG validation
    ((total_checks++))
    if validate_changelog; then
        print_status "✅ CHANGELOG looks good"
        ((checks_passed++))
    else
        print_warning "⚠️  CHANGELOG validation failed"
    fi
    
    # Check 6: Chart validation
    ((total_checks++))
    if [[ -f "$PROJECT_ROOT/scripts/validate-chart.sh" ]] && "$PROJECT_ROOT/scripts/validate-chart.sh" >/dev/null 2>&1; then
        print_status "✅ Helm chart validation passed"
        ((checks_passed++))
    else
        print_warning "⚠️  Helm chart validation failed"
    fi
    
    echo
    print_status "Pre-release checks: $checks_passed/$total_checks passed"
    
    if [[ $checks_passed -eq $total_checks ]]; then
        return 0
    else
        return 1
    fi
}

# Function to show release summary
show_release_summary() {
    local version=$1
    local latest_version current_branch commit_count
    
    latest_version=$(get_latest_version)
    current_branch=$(git branch --show-current)
    
    if [[ "$latest_version" != "v0.0.0" ]]; then
        commit_count=$(git rev-list ${latest_version}..HEAD --count)
    else
        commit_count=$(git rev-list HEAD --count)
    fi
    
    print_header "Release Summary"
    echo "Target Version: $version"
    echo "Current Branch: $current_branch"
    echo "Latest Version: $latest_version"
    echo "Commits since last release: $commit_count"
    echo
    
    if [[ "$latest_version" != "v0.0.0" ]] && [[ $commit_count -gt 0 ]]; then
        echo "Recent changes:"
        git log --pretty=format:"  - %s (%h)" ${latest_version}..HEAD | head -10
        if [[ $commit_count -gt 10 ]]; then
            echo "  ... and $((commit_count - 10)) more commits"
        fi
    fi
    echo
}

# Function to create release branch
create_release_branch() {
    local version=$1
    local prerelease=${2:-false}
    
    local branch_name="release/${version}"
    
    print_status "Creating release branch: $branch_name"
    
    # Check if branch already exists
    if git branch -r | grep -q "origin/$branch_name"; then
        print_error "Release branch $branch_name already exists on remote"
        return 1
    fi
    
    if git branch | grep -q "$branch_name"; then
        print_error "Release branch $branch_name already exists locally"
        return 1
    fi
    
    # Create and switch to release branch
    if git checkout -b "$branch_name"; then
        print_status "✅ Created release branch: $branch_name"
        
        # Push the branch to trigger the workflow
        if git push -u origin "$branch_name"; then
            print_status "✅ Pushed release branch to remote"
            print_status "The draft release workflow will be triggered automatically"
            print_status "Monitor progress at: https://github.com/$(git remote get-url origin | sed 's/.*github.com[:/]\([^/]*\/[^/]*\)\.git/\1/' | sed 's/.*github.com[:/]\([^/]*\/[^/]*\)/\1/')/actions"
            return 0
        else
            print_error "❌ Failed to push release branch"
            git checkout main
            git branch -D "$branch_name"
            return 1
        fi
    else
        print_error "❌ Failed to create release branch"
        return 1
    fi
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS] [VERSION]"
    echo
    echo "Create release branches and trigger release workflows for the Alert Reaction Operator."
    echo
    echo "OPTIONS:"
    echo "  -h, --help          Show this help message"
    echo "  -c, --check         Run pre-release checks only"
    echo "  -s, --suggest       Suggest next version based on commits"
    echo "  -p, --prerelease    Mark as pre-release"
    echo "  -d, --dry-run       Show what would be done without creating release branch"
    echo "  -f, --force         Skip pre-release checks"
    echo
    echo "VERSION:"
    echo "  Semantic version with 'v' prefix (e.g., v1.0.0, v1.2.3-alpha.1)"
    echo "  If not provided, will suggest next version"
    echo
    echo "Examples:"
    echo "  $0 --check                    # Run pre-release checks"
    echo "  $0 --suggest                  # Suggest next version"
    echo "  $0 v1.0.0                     # Create release/v1.0.0 branch"
    echo "  $0 v1.0.0-alpha.1 --prerelease # Create pre-release branch"
    echo "  $0 --dry-run v1.0.0           # Show what would be done"
}

# Main function
main() {
    local version=""
    local check_only=false
    local suggest_only=false
    local prerelease=false
    local dry_run=false
    local force=false
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_usage
                exit 0
                ;;
            -c|--check)
                check_only=true
                shift
                ;;
            -s|--suggest)
                suggest_only=true
                shift
                ;;
            -p|--prerelease)
                prerelease=true
                shift
                ;;
            -d|--dry-run)
                dry_run=true
                shift
                ;;
            -f|--force)
                force=true
                shift
                ;;
            -*)
                print_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
            *)
                version=$1
                shift
                ;;
        esac
    done
    
    print_header "Release Preparation Tool"
    
    # Change to project root
    cd "$PROJECT_ROOT"
    
    if [[ "$check_only" == true ]]; then
        run_pre_release_checks
        exit $?
    fi
    
    if [[ "$suggest_only" == true ]]; then
        local latest_version suggested_type suggested_version
        latest_version=$(get_latest_version)
        suggested_type=$(analyze_commits "$latest_version")
        suggested_version=$(generate_next_version "$latest_version" "$suggested_type")
        
        print_status "Latest version: $latest_version"
        print_status "Suggested increment: $suggested_type"
        print_status "Suggested version: $suggested_version"
        exit 0
    fi
    
    # If no version provided, suggest one
    if [[ -z "$version" ]]; then
        local latest_version suggested_type suggested_version
        latest_version=$(get_latest_version)
        suggested_type=$(analyze_commits "$latest_version")
        suggested_version=$(generate_next_version "$latest_version" "$suggested_type")
        
        print_status "No version specified. Suggested version: $suggested_version"
        read -p "Use suggested version? (Y/n): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Nn]$ ]]; then
            read -p "Enter version: " version
        else
            version=$suggested_version
        fi
    fi
    
    # Validate version format
    if ! validate_version "$version"; then
        exit 1
    fi
    
    # Check if version already exists
    if ! check_version_exists "$version"; then
        exit 1
    fi
    
    # Show release summary
    show_release_summary "$version"
    
    # Run pre-release checks unless forced
    if [[ "$force" != true ]]; then
        if ! run_pre_release_checks; then
            print_error "Pre-release checks failed"
            print_error "Use --force to skip checks or fix the issues above"
            exit 1
        fi
    fi
    
    if [[ "$dry_run" == true ]]; then
        print_status "DRY RUN: Would create release branch with:"
        echo "  Version: $version"
        echo "  Branch: release/$version"
        echo "  Pre-release: $prerelease"
        exit 0
    fi
    
    # Confirm before creating release branch
    echo
    print_warning "This will create release branch 'release/$version' and trigger the draft release workflow"
    if [[ "$prerelease" == true ]]; then
        print_warning "This will be marked as a pre-release"
    fi
    read -p "Continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_status "Cancelled"
        exit 0
    fi
    
    # Create release branch
    if create_release_branch "$version" "$prerelease"; then
        print_status "🎉 Release branch created and workflow initiated!"
        echo
        print_status "Next steps:"
        echo "1. The draft release workflow is running automatically"
        echo "2. Monitor the workflow progress in GitHub Actions"
        echo "3. Make any final changes to the release branch if needed"
        echo "4. Review the draft release when ready"
        echo "5. Publish the release to trigger Helm chart automation"
        echo
        print_status "Current branch: $(git branch --show-current)"
    else
        print_error "Failed to create release branch"
        exit 1
    fi
}

# Check if script is run directly (not sourced)
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi