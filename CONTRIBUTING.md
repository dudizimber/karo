# Contributing to Karo

Thank you for your interest in contributing to Karo (Kubernetes Alert Reaction Operator)! This document provides guidelines and information for contributors.

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/2/1/code_of_conduct/). By participating, you agree to uphold this code.

## Getting Started

### Prerequisites

- Go 1.24 or later
- Docker
- Kubernetes cluster (for testing)
- kubectl
- Helm 3.0+

### Development Setup

1. **Fork and Clone**
   ```bash
   # Fork the repository on GitHub
   git clone https://github.com/dudizimber/karo.git
   cd karo
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   ```

3. **Install Development Tools**
   ```bash
   # Install controller-gen, kustomize, etc.
   make install-tools
   ```

4. **Run Tests**
   ```bash
   make test
   make test-integration
   ```

5. **Start Local Development**
   ```bash
   # Install CRDs
   make install

   # Run operator locally
   make run
   ```

## Development Workflow

### Creating a Feature Branch

```bash
git checkout -b feature/my-new-feature
# or
git checkout -b fix/issue-123
```

### Making Changes

1. **Code Changes**: Make your changes following the guidelines below
2. **Tests**: Add or update tests for your changes
3. **Documentation**: Update documentation if needed
4. **Commit**: Use clear, descriptive commit messages

### Commit Message Guidelines

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
type(scope): description

[optional body]

[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `refactor`: Code refactoring
- `chore`: Maintenance tasks
- `ci`: CI/CD changes

**Examples:**
```
feat(controller): add retry logic for failed jobs
fix(webhook): handle malformed alertmanager payload
docs: update installation instructions
test(controller): add test for job creation failure
```

### Testing

#### Unit Tests

```bash
# Run all unit tests
make test

# Run specific package tests
go test ./controllers/...

# Run tests with coverage
make test-coverage
```

#### Integration Tests

```bash
# Run integration tests (requires cluster)
make test-integration

# Run with specific kubeconfig
KUBECONFIG=/path/to/config make test-integration
```

#### Manual Testing

```bash
# Build and deploy locally
make docker-build
make deploy

# Create test AlertReaction
kubectl apply -f config/samples/
```

### Code Guidelines

#### Go Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` and `golint`
- Maximum line length: 120 characters
- Use meaningful variable and function names
- Add comments for exported functions and types

#### Controller Guidelines

- Follow controller-runtime best practices
- Use structured logging with `ctrl.Log`
- Handle errors gracefully with proper error wrapping
- Implement proper reconciliation logic
- Use controller-runtime client for Kubernetes API calls

#### Example Controller Function

```go
func (r *AlertReactionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    log := r.Log.WithValues("alertreaction", req.NamespacedName)
    
    // Fetch the AlertReaction instance
    var alertReaction alertreactionv1.AlertReaction
    if err := r.Get(ctx, req.NamespacedName, &alertReaction); err != nil {
        if apierrors.IsNotFound(err) {
            log.Info("AlertReaction resource not found. Ignoring since object must be deleted")
            return ctrl.Result{}, nil
        }
        log.Error(err, "Failed to get AlertReaction")
        return ctrl.Result{}, err
    }
    
    // Your reconciliation logic here
    
    return ctrl.Result{}, nil
}
```

#### API Design Guidelines

- Use clear, descriptive field names
- Add proper validation tags
- Include comprehensive documentation
- Follow Kubernetes API conventions
- Use appropriate field types (pointers for optional fields)

#### Example API Type

```go
// AlertReactionSpec defines the desired state of AlertReaction
type AlertReactionSpec struct {
    // Selector specifies which alerts trigger this reaction
    // +kubebuilder:validation:Required
    Selector AlertSelector `json:"selector"`
    
    // Actions define what to do when alerts match
    // +kubebuilder:validation:Required
    // +kubebuilder:validation:MinItems=1
    Actions []Action `json:"actions"`
    
    // Cooldown prevents rapid re-execution (optional)
    // +kubebuilder:validation:Optional
    // +kubebuilder:default="5m"
    Cooldown *metav1.Duration `json:"cooldown,omitempty"`
}
```

### Documentation

#### Code Documentation

- Document all exported functions and types
- Use godoc-style comments
- Include examples in documentation

```go
// CreateJob creates a Kubernetes Job based on the AlertReaction action.
// It returns the created job and any error encountered.
//
// Example:
//   job, err := r.CreateJob(ctx, alertReaction, action, alertData)
//   if err != nil {
//       return ctrl.Result{}, err
//   }
func (r *AlertReactionReconciler) CreateJob(ctx context.Context, ar *alertreactionv1.AlertReaction, action alertreactionv1.Action, alertData map[string]string) (*batchv1.Job, error) {
    // Implementation
}
```

#### README Updates

- Update feature documentation
- Add new configuration options
- Update examples if APIs change
- Keep troubleshooting section current

### Pull Request Process

1. **Create Pull Request**
   - Use descriptive title
   - Fill out PR template
   - Link related issues

2. **PR Requirements**
   - All tests pass
   - Code coverage maintained
   - Documentation updated
   - No linting errors

3. **Review Process**
   - Address reviewer feedback
   - Keep PR focused and small
   - Rebase if requested

4. **Merge Criteria**
   - All checks pass
   - At least one approval
   - Up-to-date with main branch

### Release Process

Releases are managed by maintainers:

1. **Version Bump**: Update version in relevant files
2. **Changelog**: Update CHANGELOG.md
3. **Tag**: Create Git tag with semantic version
4. **Release**: GitHub Actions builds and publishes artifacts

## Project Structure

```
karo/
├── api/v1/                     # API definitions
├── controllers/                # Controller implementations
├── webhook/                    # Webhook server
├── config/                     # Kubernetes manifests
├── charts/                     # Helm chart
├── tests/                      # Test files
├── scripts/                    # Utility scripts
├── docs/                       # Documentation
└── hack/                       # Build scripts
```

## Issue Reporting

### Bug Reports

Include:
- Kubernetes version
- Operator version
- Clear reproduction steps
- Expected vs actual behavior
- Relevant logs

### Feature Requests

Include:
- Use case description
- Proposed solution
- Alternative solutions considered
- Additional context

### Security Issues

Report security vulnerabilities privately to:
- Email: security@example.com (replace with actual email)
- GitHub Security Advisories

## Getting Help

- **Documentation**: Check [README](README.md) and [Wiki](https://github.com/dudizimber/karo/wiki)
- **Discussions**: [GitHub Discussions](https://github.com/dudizimber/karo/discussions)
- **Issues**: [GitHub Issues](https://github.com/dudizimber/karo/issues)

## Maintainers

- [@dudizimber](https://github.com/dudizimber)

## License

By contributing, you agree that your contributions will be licensed under the Apache License 2.0.

---

Thank you for contributing to Karo! 🎉

````
