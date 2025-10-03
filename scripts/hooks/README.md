# Git Hooks

This directory contains git hooks to ensure code quality and consistency in the project.

## Available Hooks

### pre-commit
Runs before each commit to ensure:
- Go code is properly formatted (`gofmt`)
- Basic linting with `go vet`
- No trailing whitespace in files
- Files end with newlines
- YAML files are valid

### pre-push  
Runs before each push to ensure:
- All tests pass (`go test`)
- Build is successful (`go build`)
- CRD manifests are up-to-date (`make manifests`)
- Generated code is current (`make generate`)
- Full linting passes (if `golangci-lint` is available)

### commit-msg
Validates commit messages follow conventional commit format:
```
<type>[(scope)]: <description>

[optional body]

[optional footer]
```

Valid types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `ci`

Examples:
- ✅ `feat: add matcher field to AlertReaction CRD`
- ✅ `fix(controller): resolve memory leak in alert processing`
- ✅ `docs: update installation instructions`
- ❌ `updated readme` (missing type)
- ❌ `Fix bug` (type should be lowercase)

## Installation

### Automatic Setup
```bash
./scripts/setup-hooks.sh
```

This script will:
- Copy all hooks to `.git/hooks/`
- Make them executable
- Verify installation

### Manual Setup
```bash
# Copy hooks
cp scripts/hooks/* .git/hooks/

# Make executable
chmod +x .git/hooks/pre-commit .git/hooks/pre-push .git/hooks/commit-msg
```

## Testing

Test your hooks without committing:
```bash
./scripts/test-hooks.sh
```

This will verify:
- All hooks are installed and executable
- Commit message validation works correctly
- Required tools are available
- Basic code quality checks pass

## Skipping Hooks

In exceptional cases, you can skip hooks:

```bash
# Skip pre-commit hook
git commit --no-verify -m "emergency fix"

# Skip pre-push hook  
git push --no-verify
```

**Note:** Use `--no-verify` sparingly and only for genuine emergencies.

## Required Tools

### Essential
- `go` - Go compiler
- `gofmt` - Go code formatter

### Recommended
- `golangci-lint` - Advanced Go linter
- `make` - Build automation

Install golangci-lint:
```bash
# macOS
brew install golangci-lint

# Linux/Windows
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

## Troubleshooting

### Hook fails on commit
1. Check formatting: `gofmt -w .`
2. Run tests: `go test ./...`
3. Check for issues: `go vet ./...`

### Hook fails on push
1. Update manifests: `make manifests`
2. Update generated code: `make generate`
3. Run full test suite: `make test`

### Commit message rejected
Ensure your commit message follows the format:
```
type(scope): description
```

### Disable temporarily
If you need to bypass hooks temporarily:
```bash
mv .git/hooks .git/hooks.disabled
# ... do your work ...
mv .git/hooks.disabled .git/hooks
```

## Custom Configuration

### Skip specific files
Edit the hooks to add file patterns to skip:
```bash
# In pre-commit hook
if [[ "$file" =~ \.(pb|generated)\.go$ ]]; then
    continue  # Skip generated files
fi
```

### Modify commit message rules
Edit `scripts/hooks/commit-msg` to change the validation regex or add new rules.

### Add new checks
Edit the hook files to add new validation steps. Remember to:
1. Keep checks fast for pre-commit
2. Save expensive operations for pre-push
3. Provide clear error messages

## Team Setup

When onboarding new team members:
1. Have them run `./scripts/setup-hooks.sh`
2. Verify with `./scripts/test-hooks.sh`
3. Walk through the conventional commit format
4. Ensure they have required tools installed

The hooks will ensure consistent code quality across the entire team automatically.