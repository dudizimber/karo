#!/bin/bash

# Git hooks setup script for k8s-alert-reaction-operator
# This script installs git hooks for automated linting and formatting

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔧 Setting up Git hooks for k8s-alert-reaction-operator...${NC}"

# Check if we're in a git repository
if [ ! -d ".git" ]; then
    echo -e "${RED}❌ This script must be run from the root of the git repository${NC}"
    exit 1
fi

# Function to copy and make executable
install_hook() {
    local hook_name="$1"
    local source_path="scripts/hooks/$hook_name"
    local target_path=".git/hooks/$hook_name"
    
    if [ ! -f "$source_path" ]; then
        echo -e "${RED}❌ Hook source file not found: $source_path${NC}"
        return 1
    fi
    
    echo -e "${GREEN}📝 Installing $hook_name hook...${NC}"
    cp "$source_path" "$target_path"
    chmod +x "$target_path"
    echo -e "${GREEN}✅ $hook_name hook installed${NC}"
}

# Create hooks directory in scripts if it doesn't exist
mkdir -p scripts/hooks

# Check if hooks already exist in .git/hooks and back them up
backup_existing_hook() {
    local hook_name="$1"
    local hook_path=".git/hooks/$hook_name"
    
    if [ -f "$hook_path" ] && [ ! -L "$hook_path" ]; then
        echo -e "${YELLOW}⚠️  Existing $hook_name hook found, backing up...${NC}"
        mv "$hook_path" "$hook_path.backup.$(date +%Y%m%d_%H%M%S)"
        echo -e "${GREEN}✅ Backed up existing $hook_name hook${NC}"
    fi
}

# Copy our hooks to the scripts directory for version control
echo -e "${GREEN}📁 Setting up hooks in scripts directory...${NC}"

# If hooks exist in .git/hooks, copy them to scripts/hooks for version control
for hook in pre-commit pre-push; do
    if [ -f ".git/hooks/$hook" ]; then
        echo -e "${GREEN}📋 Copying $hook to scripts/hooks/ for version control...${NC}"
        cp ".git/hooks/$hook" "scripts/hooks/$hook"
    fi
done

# Create a commit-msg hook for conventional commits
cat > scripts/hooks/commit-msg << 'EOF'
#!/bin/bash

# Commit message hook to enforce conventional commit format
# See: https://www.conventionalcommits.org/

commit_regex='^(feat|fix|docs|style|refactor|test|chore|perf|ci|build|revert)(\(.+\))?: .{1,50}'

error_msg="❌ Invalid commit message format!

Please use the conventional commit format:
  <type>[optional scope]: <description>

Examples:
  feat: add webhook server
  fix(controller): resolve memory leak
  docs: update README with examples
  test: add unit tests for matcher
  refactor: simplify alert processing logic
  chore: update dependencies

Valid types: feat, fix, docs, style, refactor, test, chore, perf, ci, build, revert
"

if ! grep -qE "$commit_regex" "$1"; then
    echo "$error_msg" >&2
    exit 1
fi
EOF

chmod +x scripts/hooks/commit-msg

# Install hooks
echo -e "${GREEN}🔧 Installing git hooks...${NC}"

for hook in pre-commit pre-push commit-msg; do
    backup_existing_hook "$hook"
    if [ -f "scripts/hooks/$hook" ]; then
        echo -e "${GREEN}📝 Installing $hook hook...${NC}"
        cp "scripts/hooks/$hook" ".git/hooks/$hook"
        chmod +x ".git/hooks/$hook"
        echo -e "${GREEN}✅ $hook hook installed${NC}"
    fi
done

# Check for required tools
echo -e "${GREEN}🔍 Checking for required tools...${NC}"

check_tool() {
    local tool="$1"
    local install_cmd="$2"
    
    if command -v "$tool" >/dev/null 2>&1; then
        echo -e "${GREEN}✅ $tool found${NC}"
        return 0
    else
        echo -e "${YELLOW}⚠️  $tool not found${NC}"
        if [ -n "$install_cmd" ]; then
            echo -e "${BLUE}💡 Install with: $install_cmd${NC}"
        fi
        return 1
    fi
}

# Check required tools
all_tools_available=true

if ! check_tool "go" "https://golang.org/doc/install"; then
    all_tools_available=false
fi

if ! check_tool "gofmt" "Included with Go installation"; then
    all_tools_available=false
fi

# Check optional but recommended tools
if ! check_tool "golangci-lint" "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; then
    echo -e "${YELLOW}💡 golangci-lint is recommended for better linting${NC}"
fi

if ! check_tool "make" "Install build-essential (Ubuntu) or Xcode command line tools (macOS)"; then
    echo -e "${YELLOW}💡 make is recommended for running project targets${NC}"
fi

echo

if [ "$all_tools_available" = false ]; then
    echo -e "${YELLOW}⚠️  Some required tools are missing. Git hooks may not work correctly until they are installed.${NC}"
else
    echo -e "${GREEN}✅ All required tools are available${NC}"
fi

# Create a test script
echo -e "${GREEN}📝 Creating test script for git hooks...${NC}"

cat > scripts/test-hooks.sh << 'EOF'
#!/bin/bash

# Test script for git hooks
# This script simulates the git hooks to test them without actually committing

set -e

echo "🧪 Testing git hooks..."

echo "📝 Testing pre-commit hook..."
if [ -x .git/hooks/pre-commit ]; then
    # Create a temporary commit to test pre-commit
    echo "# Test file" > test_hook_file.go
    git add test_hook_file.go
    
    # Run pre-commit hook
    .git/hooks/pre-commit
    
    # Clean up
    git reset HEAD test_hook_file.go
    rm -f test_hook_file.go
    
    echo "✅ Pre-commit hook test passed"
else
    echo "❌ Pre-commit hook not found or not executable"
fi

echo
echo "📤 Testing pre-push hook..."
if [ -x .git/hooks/pre-push ]; then
    # Note: This will run the actual pre-push checks
    echo "⚠️  Pre-push hook will run actual checks..."
    # .git/hooks/pre-push
    echo "ℹ️  Skipping pre-push test to avoid side effects"
    echo "💡 Run manually with: .git/hooks/pre-push"
else
    echo "❌ Pre-push hook not found or not executable"
fi

echo
echo "✅ Hook tests completed"
EOF

chmod +x scripts/test-hooks.sh

echo
echo -e "${GREEN}🎉 Git hooks setup completed!${NC}"
echo
echo -e "${BLUE}📋 What was installed:${NC}"
echo -e "  • ${GREEN}pre-commit${NC} - Runs formatting, linting, and basic checks before commits"
echo -e "  • ${GREEN}pre-push${NC} - Runs comprehensive checks before pushes"
echo -e "  • ${GREEN}commit-msg${NC} - Enforces conventional commit message format"
echo
echo -e "${BLUE}💡 Additional files created:${NC}"
echo -e "  • ${GREEN}scripts/hooks/${NC} - Version controlled hooks"
echo -e "  • ${GREEN}scripts/test-hooks.sh${NC} - Test script for hooks"
echo
echo -e "${BLUE}🔧 To test the setup:${NC}"
echo -e "  • Run: ${GREEN}./scripts/test-hooks.sh${NC}"
echo -e "  • Or make a test commit to see pre-commit hook in action"
echo
echo -e "${BLUE}📚 Documentation:${NC}"
echo -e "  • Pre-commit runs: formatting, go vet, golangci-lint, build, quick tests"
echo -e "  • Pre-push runs: all pre-commit checks + full tests + manifest generation"
echo -e "  • Commit messages must follow conventional commit format"
echo
echo -e "${GREEN}✨ Happy coding!${NC}"