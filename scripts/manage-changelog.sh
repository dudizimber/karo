#!/bin/bash

# manage-changelog.sh - Helper script for CHANGELOG management
# This script helps maintain CHANGELOG.md in Keep a Changelog format

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
CHANGELOG_FILE="$PROJECT_ROOT/CHANGELOG.md"

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
    echo -e "${BLUE}[CHANGELOG]${NC} $1"
}

# Function to create initial CHANGELOG.md
create_initial_changelog() {
    local project_name=${1:-"Alert Reaction Operator"}
    local repo_url=${2:-"https://github.com/dudizimber/k8s-alert-reaction-operator"}
    
    cat > "$CHANGELOG_FILE" << EOF
# Changelog

All notable changes to the $project_name will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial project setup and documentation

EOF
    
    print_status "Created initial CHANGELOG.md"
}

# Function to add entry to unreleased section
add_entry() {
    local type=$1
    local description=$2
    
    # Validate type
    case $type in
        added|changed|deprecated|removed|fixed|security)
            ;;
        *)
            print_error "Invalid entry type: $type"
            print_error "Valid types: added, changed, deprecated, removed, fixed, security"
            return 1
            ;;
    esac
    
    if [[ ! -f "$CHANGELOG_FILE" ]]; then
        print_error "CHANGELOG.md not found. Create it first with --init"
        return 1
    fi
    
    # Check if unreleased section exists
    if ! grep "\#\# \[Unreleased\]" "$CHANGELOG_FILE"; then
        print_error "No [Unreleased] section found in CHANGELOG.md"
        return 1
    fi
    
    # Create backup
    cp "$CHANGELOG_FILE" "${CHANGELOG_FILE}.backup"
    
    # Add entry to appropriate section
    local section_header="### ${type^}"
    local temp_file=$(mktemp)
    
    # Check if the section already exists under [Unreleased]
    if sed -n '/## \[Unreleased\]/,/## \[/p' "$CHANGELOG_FILE" | grep -q "^${section_header}$"; then
        # Section exists, add to it
        awk -v section="$section_header" -v entry="- $description" '
        /^## \[Unreleased\]/ { unreleased=1 }
        /^## \[/ && !/^## \[Unreleased\]/ { unreleased=0 }
        unreleased && $0 ~ "^" section "$" {
            print $0
            getline
            print "- " entry
            print $0
            next
        }
        { print }
        ' "$CHANGELOG_FILE" > "$temp_file"
    else
        # Section doesn't exist, create it
        awk -v section="$section_header" -v entry="- $description" '
        /^## \[Unreleased\]/ {
            print $0
            print ""
            print section
            print entry
            next
        }
        { print }
        ' "$CHANGELOG_FILE" > "$temp_file"
    fi
    
    mv "$temp_file" "$CHANGELOG_FILE"
    print_status "Added $type entry: $description"
}

# Function to release current unreleased version
release_version() {
    local version=$1
    local date=${2:-$(date +%Y-%m-%d)}
    
    if [[ ! -f "$CHANGELOG_FILE" ]]; then
        print_error "CHANGELOG.md not found"
        return 1
    fi
    
    # Remove 'v' prefix if present
    version=${version#v}
    
    # Validate version format
    if [[ ! $version =~ ^[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.-]+)?$ ]]; then
        print_error "Invalid version format: $version"
        return 1
    fi
    
    # Check if version already exists
    if grep "\#\# \[${version}\]" "$CHANGELOG_FILE"; then
        print_error "Version $version already exists in CHANGELOG"
        return 1
    fi
    
    # Create backup
    cp "$CHANGELOG_FILE" "${CHANGELOG_FILE}.backup"
    
    # Replace [Unreleased] with version and add new [Unreleased] section
    local temp_file=$(mktemp)
    
    awk -v version="$version" -v date="$date" '
    /^## \[Unreleased\]/ {
        print "## [Unreleased]"
        print ""
        print "### Added"
        print "- Prepare for next release"
        print ""
        print "## [" version "] - " date
        next
    }
    { print }
    ' "$CHANGELOG_FILE" > "$temp_file"
    
    mv "$temp_file" "$CHANGELOG_FILE"
    
    # Update or add comparison links at the end
    update_comparison_links "$version"
    
    print_status "Released version $version in CHANGELOG"
}

# Function to update comparison links
update_comparison_links() {
    local version=$1
    local repo_url
    
    # Try to determine repository URL
    if git remote get-url origin >/dev/null 2>&1; then
        repo_url=$(git remote get-url origin | sed 's/\.git$//' | sed 's/git@github.com:/https:\/\/github.com\//')
    else
        repo_url="https://github.com/dudizimber/k8s-alert-reaction-operator"
    fi
    
    # Check if comparison links section exists
    if grep -q "\[Unreleased\]:" "$CHANGELOG_FILE"; then
        # Update existing links
        sed -i.tmp "s|\[Unreleased\]:.*|[Unreleased]: ${repo_url}/compare/v${version}...HEAD|" "$CHANGELOG_FILE"
        sed -i.tmp "/\[Unreleased\]:/a\\
[${version}]: ${repo_url}/releases/tag/v${version}" "$CHANGELOG_FILE"
    else
        # Add new links section
        echo "" >> "$CHANGELOG_FILE"
        echo "[Unreleased]: ${repo_url}/compare/v${version}...HEAD" >> "$CHANGELOG_FILE"
        echo "[${version}]: ${repo_url}/releases/tag/v${version}" >> "$CHANGELOG_FILE"
    fi
    
    # Clean up temp file
    rm -f "${CHANGELOG_FILE}.tmp"
}

# Function to show unreleased changes
show_unreleased() {
    if [[ ! -f "$CHANGELOG_FILE" ]]; then
        print_error "CHANGELOG.md not found"
        return 1
    fi
    
    if ! grep "\#\# \[Unreleased\]" "$CHANGELOG_FILE"; then
        print_error "No [Unreleased] section found"
        return 1
    fi
    
    print_header "Unreleased Changes"
    sed -n '/## \[Unreleased\]/,/## \[/p' "$CHANGELOG_FILE" | head -n -1 | tail -n +2
}

# Function to show version changes
show_version() {
    local version=$1
    
    if [[ ! -f "$CHANGELOG_FILE" ]]; then
        print_error "CHANGELOG.md not found"
        return 1
    fi
    
    # Remove 'v' prefix if present
    version=${version#v}
    
    if ! grep "\#\# \[${version}\]" "$CHANGELOG_FILE"; then
        print_error "Version $version not found in CHANGELOG"
        return 1
    fi
    
    print_header "Changes in version $version"
    sed -n "/## \[${version}\]/,/## \[/p" "$CHANGELOG_FILE" | head -n -1
}

# Function to validate CHANGELOG format
validate_changelog() {
    if [[ ! -f "$CHANGELOG_FILE" ]]; then
        print_error "CHANGELOG.md not found"
        return 1
    fi
    
    local issues=0
    
    # Check for required header
    if ! head -5 "$CHANGELOG_FILE" | grep -q "# Changelog"; then
        print_error "Missing main 'Changelog' header"
        ((issues++))
    fi
    
    # Check for Keep a Changelog reference
    if ! grep -q "Keep a Changelog" "$CHANGELOG_FILE"; then
        print_warning "Missing reference to Keep a Changelog format"
        ((issues++))
    fi
    
    # Check for Semantic Versioning reference
    if ! grep -q "Semantic Versioning" "$CHANGELOG_FILE"; then
        print_warning "Missing reference to Semantic Versioning"
        ((issues++))
    fi
    
    # Check for [Unreleased] section
    if ! grep "\#\# \[Unreleased\]" "$CHANGELOG_FILE"; then
        print_error "Missing [Unreleased] section"
        ((issues++))
    fi
    
    # Check version format
    local invalid_versions
    invalid_versions=$(grep "^## \[" "$CHANGELOG_FILE" | grep -v "Unreleased" | grep -v -E "^## \[[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.-]+)?\]" || true)
    if [[ -n "$invalid_versions" ]]; then
        print_error "Invalid version formats found:"
        echo "$invalid_versions"
        ((issues++))
    fi
    
    if [[ $issues -eq 0 ]]; then
        print_status "✅ CHANGELOG format is valid"
        return 0
    else
        print_error "❌ Found $issues issues in CHANGELOG format"
        return 1
    fi
}

# Function to generate release notes
generate_release_notes() {
    local version=$1
    local output_file=${2:-"release-notes.md"}
    
    if [[ ! -f "$CHANGELOG_FILE" ]]; then
        print_error "CHANGELOG.md not found"
        return 1
    fi
    
    # Remove 'v' prefix if present
    version=${version#v}
    
    if ! grep "\#\# \[${version}\]" "$CHANGELOG_FILE"; then
        print_error "Version $version not found in CHANGELOG"
        return 1
    fi
    
    # Extract version section
    sed -n "/## \[${version}\]/,/## \[/p" "$CHANGELOG_FILE" | head -n -1 > "$output_file"
    
    print_status "Generated release notes for $version in $output_file"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo
    echo "Manage CHANGELOG.md in Keep a Changelog format."
    echo
    echo "COMMANDS:"
    echo "  init                    Create initial CHANGELOG.md"
    echo "  add TYPE DESCRIPTION    Add entry to unreleased section"
    echo "  release VERSION [DATE]  Release current unreleased changes"
    echo "  show [VERSION]          Show unreleased or specific version changes"
    echo "  validate                Validate CHANGELOG format"
    echo "  notes VERSION [FILE]    Generate release notes for version"
    echo
    echo "ENTRY TYPES:"
    echo "  added       New features"
    echo "  changed     Changes in existing functionality"
    echo "  deprecated  Soon-to-be removed features"
    echo "  removed     Removed features"
    echo "  fixed       Bug fixes"
    echo "  security    Vulnerability fixes"
    echo
    echo "OPTIONS:"
    echo "  -h, --help             Show this help message"
    echo
    echo "Examples:"
    echo "  $0 init                                    # Create initial CHANGELOG"
    echo "  $0 add added 'New webhook endpoint'        # Add new feature"
    echo "  $0 add fixed 'Fix memory leak in operator' # Add bug fix"
    echo "  $0 release 1.0.0                          # Release version 1.0.0"
    echo "  $0 show                                    # Show unreleased changes"
    echo "  $0 show 1.0.0                             # Show changes in v1.0.0"
    echo "  $0 validate                                # Validate format"
    echo "  $0 notes 1.0.0                            # Generate release notes"
}

# Main function
main() {
    if [[ $# -eq 0 ]]; then
        show_usage
        exit 1
    fi
    
    case $1 in
        init)
            create_initial_changelog "${2:-}" "${3:-}"
            ;;
        add)
            if [[ $# -lt 3 ]]; then
                print_error "Usage: $0 add TYPE DESCRIPTION"
                exit 1
            fi
            add_entry "$2" "$3"
            ;;
        release)
            if [[ $# -lt 2 ]]; then
                print_error "Usage: $0 release VERSION [DATE]"
                exit 1
            fi
            release_version "$2" "${3:-}"
            ;;
        show)
            if [[ $# -eq 1 ]]; then
                show_unreleased
            else
                show_version "$2"
            fi
            ;;
        validate)
            validate_changelog
            ;;
        notes)
            if [[ $# -lt 2 ]]; then
                print_error "Usage: $0 notes VERSION [FILE]"
                exit 1
            fi
            generate_release_notes "$2" "${3:-}"
            ;;
        -h|--help)
            show_usage
            ;;
        *)
            print_error "Unknown command: $1"
            show_usage
            exit 1
            ;;
    esac
}

# Check if script is run directly (not sourced)
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi