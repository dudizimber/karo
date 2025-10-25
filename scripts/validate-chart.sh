#!/bin/bash

# validate-chart.sh - Validation script for Helm charts
# This script performs comprehensive validation of the Helm chart

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
    echo -e "${BLUE}[CHART VALIDATION]${NC} $1"
}

# Check if required tools are available
check_tools() {
    local missing_tools=()
    
    if ! command -v helm >/dev/null 2>&1; then
        missing_tools+=("helm")
    fi
    
    if ! command -v yq >/dev/null 2>&1; then
        print_warning "yq not found - YAML validation will be limited"
    fi
    
    if [[ ${#missing_tools[@]} -gt 0 ]]; then
        print_error "Missing required tools: ${missing_tools[*]}"
        echo "Please install the missing tools:"
        for tool in "${missing_tools[@]}"; do
            case $tool in
                helm)
                    echo "  - Helm: https://helm.sh/docs/intro/install/"
                    ;;
            esac
        done
        return 1
    fi
    
    return 0
}

# Validate Chart.yaml structure
validate_chart_yaml() {
    local chart_file="$CHART_DIR/Chart.yaml"
    
    print_status "Validating Chart.yaml structure..."
    
    if [[ ! -f "$chart_file" ]]; then
        print_error "Chart.yaml not found at $chart_file"
        return 1
    fi
    
    # Check required fields
    local required_fields=("name" "version" "description" "type" "appVersion")
    local missing_fields=()
    
    for field in "${required_fields[@]}"; do
        if ! grep -q "^${field}:" "$chart_file"; then
            missing_fields+=("$field")
        fi
    done
    
    if [[ ${#missing_fields[@]} -gt 0 ]]; then
        print_error "Missing required fields in Chart.yaml: ${missing_fields[*]}"
        return 1
    fi
    
    # Validate version format
    local version
    version=$(grep "^version:" "$chart_file" | awk '{print $2}' | tr -d '"')
    if [[ ! $version =~ ^[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.-]+)?$ ]]; then
        print_error "Invalid version format in Chart.yaml: $version"
        return 1
    fi
    
    # Validate app version format
    local app_version
    app_version=$(grep "^appVersion:" "$chart_file" | awk '{print $2}' | tr -d '"')
    if [[ -z "$app_version" ]]; then
        print_error "appVersion is empty in Chart.yaml"
        return 1
    fi
    
    print_status "✅ Chart.yaml structure is valid"
    return 0
}

# Validate values.yaml structure
validate_values_yaml() {
    local values_file="$CHART_DIR/values.yaml"
    
    print_status "Validating values.yaml structure..."
    
    if [[ ! -f "$values_file" ]]; then
        print_warning "values.yaml not found - this is optional but recommended"
        return 0
    fi
    
    # Check if it's valid YAML
    if command -v yq >/dev/null 2>&1; then
        if ! yq eval '.' "$values_file" > /dev/null 2>&1; then
            print_error "values.yaml is not valid YAML"
            return 1
        fi
    fi
    
    # Check for common configuration sections
    local common_sections=("image" "service" "resources")
    local found_sections=()
    
    for section in "${common_sections[@]}"; do
        if grep -q "^${section}:" "$values_file"; then
            found_sections+=("$section")
        fi
    done
    
    if [[ ${#found_sections[@]} -gt 0 ]]; then
        print_status "Found configuration sections: ${found_sections[*]}"
    fi
    
    print_status "✅ values.yaml structure is valid"
    return 0
}

# Validate template directory
validate_templates() {
    local templates_dir="$CHART_DIR/templates"
    
    print_status "Validating templates directory..."
    
    if [[ ! -d "$templates_dir" ]]; then
        print_error "Templates directory not found at $templates_dir"
        return 1
    fi
    
    # Check for common template files
    local template_files
    template_files=$(find "$templates_dir" -name "*.yaml" -o -name "*.yml" | wc -l)
    
    if [[ $template_files -eq 0 ]]; then
        print_warning "No YAML template files found in templates directory"
    else
        print_status "Found $template_files template files"
    fi
    
    # Check for NOTES.txt
    if [[ -f "$templates_dir/NOTES.txt" ]]; then
        print_status "✅ NOTES.txt found"
    else
        print_warning "NOTES.txt not found - consider adding installation notes"
    fi
    
    print_status "✅ Templates directory is valid"
    return 0
}

# Run Helm lint
run_helm_lint() {
    print_status "Running Helm lint..."
    
    if helm lint "$CHART_DIR" --quiet; then
        print_status "✅ Helm lint passed"
        return 0
    else
        print_error "❌ Helm lint failed"
        print_status "Running lint again with verbose output:"
        helm lint "$CHART_DIR"
        return 1
    fi
}

# Test template rendering
test_template_rendering() {
    print_status "Testing template rendering..."
    
    local temp_dir
    temp_dir=$(mktemp -d)
    local output_file="$temp_dir/rendered.yaml"
    
    # Test with default values
    if helm template test-release "$CHART_DIR" > "$output_file" 2>&1; then
        print_status "✅ Template rendering with default values successful"
        
        # Validate rendered YAML
        if command -v yq >/dev/null 2>&1; then
            if yq eval '.' "$output_file" > /dev/null 2>&1; then
                print_status "✅ Rendered templates are valid YAML"
            else
                print_error "❌ Rendered templates are not valid YAML"
                rm -rf "$temp_dir"
                return 1
            fi
        fi
    else
        print_error "❌ Template rendering failed"
        cat "$output_file"
        rm -rf "$temp_dir"
        return 1
    fi
    
    # Test with custom values if available
    if [[ -f "$CHART_DIR/values.yaml" ]]; then
        local custom_output="$temp_dir/custom-rendered.yaml"
        if helm template test-release "$CHART_DIR" -f "$CHART_DIR/values.yaml" > "$custom_output" 2>&1; then
            print_status "✅ Template rendering with custom values successful"
        else
            print_warning "⚠️  Template rendering with custom values failed"
            cat "$custom_output"
        fi
    fi
    
    rm -rf "$temp_dir"
    return 0
}

# Check for security best practices
check_security() {
    print_status "Checking security best practices..."
    
    local issues=()
    
    # Check for hard-coded secrets in templates
    if grep -r "password\|secret\|key" "$CHART_DIR/templates" --include="*.yaml" --include="*.yml" 2>/dev/null | grep -v "secretName\|secretKeyRef\|valueFrom" | grep -q "="; then
        issues+=("Potential hard-coded secrets found in templates")
    fi
    
    # Check if security context is configured
    if ! grep -r "securityContext" "$CHART_DIR/templates" --include="*.yaml" --include="*.yml" >/dev/null 2>&1; then
        issues+=("No securityContext found in templates - consider adding security configurations")
    fi
    
    # Check for resource limits
    if ! grep -r "resources:" "$CHART_DIR/templates" --include="*.yaml" --include="*.yml" >/dev/null 2>&1; then
        issues+=("No resource limits found - consider adding resource configurations")
    fi
    
    if [[ ${#issues[@]} -gt 0 ]]; then
        print_warning "Security recommendations:"
        for issue in "${issues[@]}"; do
            echo "  - $issue"
        done
    else
        print_status "✅ No obvious security issues found"
    fi
    
    return 0
}

# Generate chart summary
generate_summary() {
    print_header "Chart Summary"
    
    local chart_file="$CHART_DIR/Chart.yaml"
    local values_file="$CHART_DIR/values.yaml"
    local templates_dir="$CHART_DIR/templates"
    
    if [[ -f "$chart_file" ]]; then
        local name version app_version description
        name=$(grep "^name:" "$chart_file" | awk '{print $2}' | tr -d '"')
        version=$(grep "^version:" "$chart_file" | awk '{print $2}' | tr -d '"')
        app_version=$(grep "^appVersion:" "$chart_file" | awk '{print $2}' | tr -d '"')
        description=$(grep "^description:" "$chart_file" | cut -d':' -f2- | sed 's/^ *//' | tr -d '"')
        
        echo "Name: $name"
        echo "Chart Version: $version"
        echo "App Version: $app_version"
        echo "Description: $description"
    fi
    
    if [[ -d "$templates_dir" ]]; then
        local template_count
        template_count=$(find "$templates_dir" -name "*.yaml" -o -name "*.yml" | wc -l)
        echo "Template Files: $template_count"
    fi
    
    if [[ -f "$values_file" ]]; then
        echo "Default Values: Available"
    else
        echo "Default Values: Not found"
    fi
    
    echo
}

# Main validation function
main() {
    local exit_code=0
    
    print_header "Helm Chart Validation"
    
    # Check if chart directory exists
    if [[ ! -d "$CHART_DIR" ]]; then
        print_error "Chart directory not found: $CHART_DIR"
        exit 1
    fi
    
    print_status "Validating chart at: $CHART_DIR"
    echo
    
    # Check required tools
    if ! check_tools; then
        exit 1
    fi
    
    # Run all validations
    validate_chart_yaml || exit_code=1
    validate_values_yaml || exit_code=1
    validate_templates || exit_code=1
    
    if [[ $exit_code -eq 0 ]]; then
        run_helm_lint || exit_code=1
        test_template_rendering || exit_code=1
        check_security
    fi
    
    echo
    generate_summary
    
    if [[ $exit_code -eq 0 ]]; then
        print_status "🎉 Chart validation completed successfully!"
    else
        print_error "❌ Chart validation failed"
    fi
    
    exit $exit_code
}

# Check if script is run directly (not sourced)
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi