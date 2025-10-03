#!/bin/bash

# Test script for git hooks
# This script simulates the git hooks to test them without actually committing

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧪 Testing git hooks...${NC}"

# Function to test a hook
test_hook() {
    local hook_name="$1"
    local hook_path=".git/hooks/$hook_name"
    
    echo -e "${GREEN}📝 Testing $hook_name hook...${NC}"
    
    if [ -x "$hook_path" ]; then
        echo -e "${GREEN}✅ $hook_name hook found and executable${NC}"
        return 0
    else
        echo -e "${RED}❌ $hook_name hook not found or not executable${NC}"
        return 1
    fi
}

# Test individual hooks
echo -e "${BLUE}🔍 Checking hook installation...${NC}"

hooks_ok=true

for hook in pre-commit pre-push commit-msg; do
    if ! test_hook "$hook"; then
        hooks_ok=false
    fi
done

if [ "$hooks_ok" = false ]; then
    echo -e "${RED}❌ Some hooks are missing or not executable${NC}"
    echo -e "${YELLOW}💡 Run ./scripts/setup-hooks.sh to install hooks${NC}"
    exit 1
fi

echo
echo -e "${BLUE}🎯 Testing commit message validation...${NC}"

# Test commit message hook
if [ -x .git/hooks/commit-msg ]; then
    # Test valid commit messages
    valid_messages=(
        "feat: add new feature"
        "fix(controller): resolve memory leak"
        "docs: update README"
        "test: add unit tests"
        "refactor: simplify code"
        "chore: update dependencies"
    )
    
    for msg in "${valid_messages[@]}"; do
        echo "$msg" > /tmp/test_commit_msg
        if .git/hooks/commit-msg /tmp/test_commit_msg >/dev/null 2>&1; then
            echo -e "${GREEN}✅ Valid: $msg${NC}"
        else
            echo -e "${RED}❌ Should be valid: $msg${NC}"
        fi
    done
    
    # Test invalid commit messages
    invalid_messages=(
        "invalid message"
        "Fix bug"
        "add feature"
        "updated readme"
    )
    
    for msg in "${invalid_messages[@]}"; do
        echo "$msg" > /tmp/test_commit_msg
        if ! .git/hooks/commit-msg /tmp/test_commit_msg >/dev/null 2>&1; then
            echo -e "${GREEN}✅ Correctly rejected: $msg${NC}"
        else
            echo -e "${RED}❌ Should be invalid: $msg${NC}"
        fi
    done
    
    rm -f /tmp/test_commit_msg
    echo -e "${GREEN}✅ Commit message validation test passed${NC}"
else
    echo -e "${RED}❌ Commit message hook not found${NC}"
fi

echo
echo -e "${BLUE}🔧 Testing tool availability...${NC}"

# Check for required tools
tools_ok=true

check_tool() {
    local tool="$1"
    
    if command -v "$tool" >/dev/null 2>&1; then
        echo -e "${GREEN}✅ $tool found${NC}"
        return 0
    else
        echo -e "${RED}❌ $tool not found${NC}"
        return 1
    fi
}

required_tools=("go" "gofmt")
for tool in "${required_tools[@]}"; do
    if ! check_tool "$tool"; then
        tools_ok=false
    fi
done

optional_tools=("golangci-lint" "make")
for tool in "${optional_tools[@]}"; do
    if ! check_tool "$tool"; then
        echo -e "${YELLOW}💡 $tool is recommended but not required${NC}"
    fi
done

echo
echo -e "${BLUE}🏃 Running quick validation checks...${NC}"

# Test basic Go formatting
echo -e "${GREEN}🎨 Checking Go formatting...${NC}"
if unformatted=$(gofmt -l . | grep -v vendor | head -5); then
    if [ -n "$unformatted" ]; then
        echo -e "${YELLOW}⚠️  Some files need formatting:${NC}"
        echo "$unformatted"
        echo -e "${YELLOW}💡 Run 'gofmt -w .' to fix${NC}"
    else
        echo -e "${GREEN}✅ All Go files are properly formatted${NC}"
    fi
fi

# Test basic build
echo -e "${GREEN}🔨 Testing Go build...${NC}"
if go build ./... >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Go build successful${NC}"
else
    echo -e "${RED}❌ Go build failed${NC}"
    tools_ok=false
fi

# Test go vet
echo -e "${GREEN}🔍 Running go vet...${NC}"
if go vet ./... >/dev/null 2>&1; then
    echo -e "${GREEN}✅ go vet passed${NC}"
else
    echo -e "${YELLOW}⚠️  go vet found issues${NC}"
fi

echo
if [ "$hooks_ok" = true ] && [ "$tools_ok" = true ]; then
    echo -e "${GREEN}🎉 All hook tests passed!${NC}"
    echo -e "${BLUE}💡 Your git hooks are ready to use${NC}"
    echo
    echo -e "${BLUE}📚 What happens now:${NC}"
    echo -e "  • ${GREEN}pre-commit${NC} will run on every commit"
    echo -e "  • ${GREEN}pre-push${NC} will run on every push"
    echo -e "  • ${GREEN}commit-msg${NC} will validate commit message format"
    echo
    echo -e "${BLUE}🧪 To test manually:${NC}"
    echo -e "  • Make a test commit to see pre-commit in action"
    echo -e "  • Try an invalid commit message to test validation"
    exit 0
else
    echo -e "${RED}❌ Some tests failed${NC}"
    echo -e "${YELLOW}💡 Please address the issues above before proceeding${NC}"
    exit 1
fi