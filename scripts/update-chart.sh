#!/bin/bash

# update-chart.sh - Helper script for updating Helm chart versions
# Usage: ./scripts/update-chart.sh [version]

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
CHART_DIR="${PROJECT_ROOT}/charts/karo"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
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
    echo -e "${BLUE}[HELM CHART UPDATE]${NC} $1"
}

# Function to validate version format
validate_version() {
    local version=$1
    if [[ ! $version =~ ^[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.-]+)?$ ]]; then
        print_error "Invalid version format: $version"
        print_error "Expected format: X.Y.Z or X.Y.Z-suffix (e.g., 1.0.0, 1.0.0-alpha.1)"
        exit 1
    fi
}

# Function to get current version from Chart.yaml
get_current_version() {
    if [[ -f "$CHART_DIR/Chart.yaml" ]]; then
        grep "^version:" "$CHART_DIR/Chart.yaml" | awk '{print $2}' | tr -d '"'
    else
        echo "0.1.0"
    fi
}

# Function to get current app version from Chart.yaml
get_current_app_version() {
    if [[ -f "$CHART_DIR/Chart.yaml" ]]; then
        grep "^appVersion:" "$CHART_DIR/Chart.yaml" | awk '{print $2}' | tr -d '"'
    else
        echo "latest"
    fi
}

# Function to update Chart.yaml
update_chart_yaml() {
    local new_version=$1
    local chart_file="$CHART_DIR/Chart.yaml"
    
    print_status "Updating Chart.yaml with version $new_version"
    
    # Create backup
    cp "$chart_file" "$chart_file.backup"
    
    # Update version and appVersion
    sed -i.tmp "s/^version: .*/version: $new_version/" "$chart_file"
    sed -i.tmp "s/^appVersion: .*/appVersion: \"$new_version\"/" "$chart_file"
    
    # Clean up temp file
    rm -f "$chart_file.tmp"
    
    print_status "Chart.yaml updated successfully"
}

# Function to update values.yaml
update_values_yaml() {
    local new_version=$1
    local values_file="$CHART_DIR/values.yaml"
    
    if [[ -f "$values_file" ]]; then
        print_status "Updating values.yaml with version $new_version"
        
        # Create backup
        cp "$values_file" "$values_file.backup"
        
        # Update image tag
        sed -i.tmp "s/tag: .*/tag: \"$new_version\"/" "$values_file"
        
        # Clean up temp file
        rm -f "$values_file.tmp"
        
        print_status "values.yaml updated successfully"
    else
        print_warning "values.yaml not found, skipping update"
    fi
}

# Function to lint the chart
lint_chart() {
    print_status "Linting Helm chart..."
    
    if command -v helm >/dev/null 2>&1; then
        if helm lint "$CHART_DIR"; then
            print_status "✅ Chart lint passed"
        else
            print_error "❌ Chart lint failed"
            exit 1
        fi
    else
        print_warning "Helm not found, skipping lint check"
    fi
}

# Function to test template rendering
test_template() {
    print_status "Testing chart template rendering..."
    
    if command -v helm >/dev/null 2>&1; then
        local temp_dir=$(mktemp -d)
        if helm template test-release "$CHART_DIR" > "$temp_dir/rendered.yaml"; then
            print_status "✅ Template rendering successful"
            
            # Check if output is valid YAML
            if command -v yq >/dev/null 2>&1; then
                if yq eval '.' "$temp_dir/rendered.yaml" > /dev/null; then
                    print_status "✅ Rendered templates are valid YAML"
                else
                    print_error "❌ Rendered templates are not valid YAML"
                    exit 1
                fi
            fi
        else
            print_error "❌ Template rendering failed"
            exit 1
        fi
        rm -rf "$temp_dir"
    else
        print_warning "Helm not found, skipping template test"
    fi
}

# Function to show changes
show_changes() {
    print_status "Changes made:"
    
    if [[ -f "$CHART_DIR/Chart.yaml.backup" ]]; then
        echo "Chart.yaml changes:"
        diff "$CHART_DIR/Chart.yaml.backup" "$CHART_DIR/Chart.yaml" || true
        echo
    fi
    
    if [[ -f "$CHART_DIR/values.yaml.backup" ]]; then
        echo "values.yaml changes:"
        diff "$CHART_DIR/values.yaml.backup" "$CHART_DIR/values.yaml" || true
        echo
    fi
}

# Function to clean up backups
cleanup_backups() {
    print_status "Cleaning up backup files..."
    rm -f "$CHART_DIR/Chart.yaml.backup"
    rm -f "$CHART_DIR/values.yaml.backup"
}

# Function to generate next version
generate_next_version() {
    local current_version=$1
    local version_type=${2:-patch}
    
    # Remove any pre-release suffix for calculation
    local base_version=$(echo "$current_version" | cut -d'-' -f1)
    
    IFS='.' read -ra VERSION_PARTS <<< "$base_version"
    local major=${VERSION_PARTS[0]}
    local minor=${VERSION_PARTS[1]}
    local patch=${VERSION_PARTS[2]}
    
    case $version_type in
        major)
            echo "$((major + 1)).0.0"
            ;;
        minor)
            echo "$major.$((minor + 1)).0"
            ;;
        patch|*)
            echo "$major.$minor.$((patch + 1))"
            ;;
    esac
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS] [VERSION]"
    echo
    echo "Update Helm chart version and related files."
    echo
    echo "OPTIONS:"
    echo "  -h, --help          Show this help message"
    echo "  -c, --current       Show current versions"
    echo "  -n, --next [TYPE]   Generate next version (TYPE: major, minor, patch)"
    echo "  -l, --lint-only     Only run lint and template tests"
    echo "  -d, --dry-run       Show what would be changed without making changes"
    echo
    echo "VERSION:"
    echo "  Semantic version (e.g., 1.0.0, 1.2.3-alpha.1)"
    echo "  If not provided, will prompt for input"
    echo
    echo "Examples:"
    echo "  $0 1.0.0              # Update to version 1.0.0"
    echo "  $0 --next patch       # Generate next patch version"
    echo "  $0 --current          # Show current versions"
    echo "  $0 --lint-only        # Only run validation"
}

# Main function
main() {
    local version=""
    local show_current=false
    local generate_next=false
    local next_type="patch"
    local lint_only=false
    local dry_run=false
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_usage
                exit 0
                ;;
            -c|--current)
                show_current=true
                shift
                ;;
            -n|--next)
                generate_next=true
                if [[ $# -gt 1 && $2 =~ ^(major|minor|patch)$ ]]; then
                    next_type=$2
                    shift
                fi
                shift
                ;;
            -l|--lint-only)
                lint_only=true
                shift
                ;;
            -d|--dry-run)
                dry_run=true
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
    
    print_header "Helm Chart Update Tool"
    
    # Check if chart directory exists
    if [[ ! -d "$CHART_DIR" ]]; then
        print_error "Chart directory not found: $CHART_DIR"
        exit 1
    fi
    
    local current_version=$(get_current_version)
    local current_app_version=$(get_current_app_version)
    
    if [[ "$show_current" == true ]]; then
        print_status "Current chart version: $current_version"
        print_status "Current app version: $current_app_version"
        exit 0
    fi
    
    if [[ "$lint_only" == true ]]; then
        lint_chart
        test_template
        print_status "✅ Chart validation completed successfully"
        exit 0
    fi
    
    if [[ "$generate_next" == true ]]; then
        version=$(generate_next_version "$current_version" "$next_type")
        print_status "Generated next $next_type version: $version"
    fi
    
    # If no version provided, prompt for input
    if [[ -z "$version" ]]; then
        echo "Current chart version: $current_version"
        echo "Current app version: $current_app_version"
        echo
        read -p "Enter new version (e.g., 1.0.0): " version
        
        if [[ -z "$version" ]]; then
            print_error "Version is required"
            exit 1
        fi
    fi
    
    # Validate version format
    validate_version "$version"
    
    if [[ "$version" == "$current_version" ]]; then
        print_warning "New version is the same as current version: $version"
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_status "Cancelled"
            exit 0
        fi
    fi
    
    if [[ "$dry_run" == true ]]; then
        print_status "DRY RUN: Would update from $current_version to $version"
        print_status "DRY RUN: Chart.yaml version field would be updated"
        print_status "DRY RUN: Chart.yaml appVersion field would be updated"
        if [[ -f "$CHART_DIR/values.yaml" ]]; then
            print_status "DRY RUN: values.yaml image tag would be updated"
        fi
        exit 0
    fi
    
    # Perform updates
    print_status "Updating chart from version $current_version to $version"
    
    update_chart_yaml "$version"
    update_values_yaml "$version"
    
    # Validate changes
    lint_chart
    test_template
    
    # Show what changed
    show_changes
    
    # Clean up backups
    cleanup_backups
    
    print_status "✅ Chart updated successfully to version $version"
    print_status "Next steps:"
    echo "  1. Review the changes"
    echo "  2. Commit the changes: git add charts/ && git commit -m 'chore: update chart to v$version'"
    echo "  3. Tag the release: git tag v$version && git push origin v$version"
}

# Check if script is run directly (not sourced)
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi