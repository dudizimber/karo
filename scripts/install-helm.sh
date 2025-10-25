#!/bin/bash
set -e

# Alert Reaction Operator Helm Chart Installation Script

CHART_DIR="charts/karo"
RELEASE_NAME="karo"
NAMESPACE="default"
VALUES_FILE=""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Print colored output
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Help function
show_help() {
    cat << EOF
Alert Reaction Operator Helm Installation Script

Usage: $0 [OPTIONS]

Options:
    -n, --namespace NAMESPACE    Target namespace (default: default)
    -r, --release RELEASE        Release name (default: karo)
    -f, --values VALUES_FILE     Custom values file
    -e, --environment ENV        Use predefined environment (dev|prod)
    --dry-run                    Perform a dry run
    --upgrade                    Upgrade existing installation
    --uninstall                  Uninstall the release
    -h, --help                   Show this help message

Examples:
    $0                           # Install with default values
    $0 -e dev                    # Install for development
    $0 -e prod -n monitoring     # Install for production in monitoring namespace
    $0 --upgrade -f my-values.yaml  # Upgrade with custom values
    $0 --uninstall               # Uninstall the operator

EOF
}

# Check if helm is installed
check_helm() {
    if ! command -v helm &> /dev/null; then
        print_error "Helm is not installed. Please install Helm first."
        exit 1
    fi
}

# Check if chart directory exists
check_chart() {
    if [ ! -d "$CHART_DIR" ]; then
        print_error "Chart directory $CHART_DIR not found."
        print_info "Please run this script from the repository root directory."
        exit 1
    fi
}

# Install or upgrade the chart
install_chart() {
    local action="install"
    local extra_args=""
    
    if [ "$UPGRADE" = "true" ]; then
        action="upgrade"
    fi
    
    if [ "$DRY_RUN" = "true" ]; then
        extra_args="--dry-run"
    fi
    
    if [ -n "$VALUES_FILE" ]; then
        extra_args="$extra_args -f $VALUES_FILE"
    fi
    
    print_info "Running: helm $action $RELEASE_NAME $CHART_DIR --namespace $NAMESPACE --create-namespace $extra_args"
    
    if helm $action $RELEASE_NAME $CHART_DIR \
        --namespace $NAMESPACE \
        --create-namespace \
        $extra_args; then
        
        if [ "$DRY_RUN" != "true" ]; then
            print_info "Successfully ${action}ed Alert Reaction Operator!"
            print_info ""
            print_info "To check the status:"
            print_info "  kubectl get deployment $RELEASE_NAME -n $NAMESPACE"
            print_info ""
            print_info "To view the webhook service:"
            print_info "  kubectl get svc ${RELEASE_NAME}-webhook -n $NAMESPACE"
            print_info ""
            print_info "To see the operator logs:"
            print_info "  kubectl logs -l app.kubernetes.io/name=karo -n $NAMESPACE"
        fi
    else
        print_error "Failed to $action Alert Reaction Operator"
        exit 1
    fi
}

# Uninstall the chart
uninstall_chart() {
    print_warning "This will uninstall the Alert Reaction Operator and remove all related resources."
    print_warning "AlertReaction custom resources will remain unless manually deleted."
    read -p "Are you sure? (y/N): " -n 1 -r
    echo
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        if helm uninstall $RELEASE_NAME --namespace $NAMESPACE; then
            print_info "Successfully uninstalled Alert Reaction Operator!"
            print_warning "Note: CRDs and custom resources may still exist."
            print_info "To remove CRDs: kubectl delete crd alertreactions.karo.io"
        else
            print_error "Failed to uninstall Alert Reaction Operator"
            exit 1
        fi
    else
        print_info "Uninstall cancelled."
    fi
}

# Parse command line arguments
DRY_RUN=false
UPGRADE=false
UNINSTALL=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -n|--namespace)
            NAMESPACE="$2"
            shift 2
            ;;
        -r|--release)
            RELEASE_NAME="$2"
            shift 2
            ;;
        -f|--values)
            VALUES_FILE="$2"
            shift 2
            ;;
        -e|--environment)
            case $2 in
                dev)
                    VALUES_FILE="$CHART_DIR/values-dev.yaml"
                    ;;
                prod)
                    VALUES_FILE="$CHART_DIR/values-prod.yaml"
                    ;;
                *)
                    print_error "Unknown environment: $2"
                    print_info "Available environments: dev, prod"
                    exit 1
                    ;;
            esac
            shift 2
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --upgrade)
            UPGRADE=true
            shift
            ;;
        --uninstall)
            UNINSTALL=true
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Main execution
print_info "Alert Reaction Operator Helm Installation"
print_info "Chart: $CHART_DIR"
print_info "Release: $RELEASE_NAME"
print_info "Namespace: $NAMESPACE"

if [ -n "$VALUES_FILE" ]; then
    print_info "Values file: $VALUES_FILE"
fi

check_helm
check_chart

if [ "$UNINSTALL" = "true" ]; then
    uninstall_chart
else
    # Lint the chart first
    print_info "Linting chart..."
    if ! helm lint $CHART_DIR; then
        print_error "Chart linting failed"
        exit 1
    fi
    
    install_chart
fi
