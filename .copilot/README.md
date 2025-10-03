# Copilot Prompts for k8s-alert-reaction-operator

This directory contains prompt templates to streamline development with GitHub Copilot for the k8s-alert-reaction-operator project.

## How to Use

1. **Copy the relevant prompt** from the templates below
2. **Replace placeholders** (marked with `{placeholder}`) with your specific details
3. **Paste into Copilot chat** or use as inline comments
4. **Iterate and refine** based on Copilot's responses

## Available Prompts

### Development Tasks
- `add-feature.md` - Adding new features to the operator
- `fix-bug.md` - Debugging and fixing issues
- `add-tests.md` - Writing comprehensive tests
- `refactor-code.md` - Code refactoring and cleanup
- `optimize-performance.md` - Performance improvements

### Documentation
- `write-docs.md` - Creating and updating documentation
- `commit-messages.md` - Writing conventional commit messages
- `code-review.md` - Code review assistance

### Kubernetes Specific
- `crd-changes.md` - Modifying Custom Resource Definitions
- `controller-logic.md` - Controller development and updates
- `rbac-permissions.md` - RBAC and security configurations

### CI/CD and DevOps
- `ci-cd-updates.md` - GitHub Actions and pipeline improvements
- `deployment.md` - Deployment and installation scripts
- `monitoring.md` - Metrics and observability

## Project Context

When using these prompts, Copilot will have context about:

- **Project Type**: Kubernetes operator built with Kubebuilder
- **Language**: Go 1.24+
- **Framework**: controller-runtime
- **CRD**: AlertReaction with Prometheus-style matchers
- **Architecture**: Webhook-driven alert processing with Job creation
- **Testing**: Comprehensive test suite with table-driven tests
- **CI/CD**: GitHub Actions with security scanning and multi-arch builds

## Best Practices

1. **Be Specific**: Include file paths, function names, and exact requirements
2. **Provide Context**: Mention related code, dependencies, or constraints
3. **Include Examples**: Show expected input/output or similar patterns
4. **Iterate**: Refine prompts based on initial results
5. **Test Results**: Always validate generated code with tests and reviews

## Quick Reference

### Common Placeholders
- `{feature_name}` - Name of the feature being added
- `{issue_description}` - Description of the bug or issue
- `{file_path}` - Path to the file being modified
- `{function_name}` - Name of the function or method
- `{test_scenario}` - Description of what should be tested
- `{scope}` - Scope for commit messages (controller, webhook, crd, etc.)

### Project-Specific Terms
- **AlertReaction** - Main CRD that defines alert-to-job mappings
- **Matcher** - Prometheus-style alert filtering conditions
- **AlertManager** - Source of webhook alerts
- **Controller** - Reconciles AlertReaction resources and creates Jobs
- **Webhook** - HTTP endpoint that receives alerts from AlertManager

## Examples

### Quick Feature Addition
```
I need to add {feature_name} to the AlertReaction CRD. Please help me:
1. Update the API types in api/v1alpha1/alertreaction_types.go
2. Update the controller logic in controllers/alertreaction_controller.go
3. Add comprehensive tests
4. Update the CRD manifests

The feature should {detailed_description}.
```

### Quick Bug Fix
```
I'm experiencing {issue_description} in {file_path}. 
The expected behavior is {expected_behavior}.
The current behavior is {actual_behavior}.
Please help me debug and fix this issue.
```

### Quick Test Addition
```
I need comprehensive tests for {function_name} in {file_path}.
Please create table-driven tests that cover:
- Happy path scenarios
- Error conditions
- Edge cases
- {specific_test_requirements}
```

## Tips for Better Results

1. **Start Broad, Then Narrow**: Begin with high-level requirements, then get specific
2. **Reference Existing Code**: Point to similar implementations in the codebase
3. **Specify Constraints**: Mention performance, security, or compatibility requirements
4. **Ask for Explanations**: Request comments and documentation alongside code
5. **Validate Assumptions**: Ask Copilot to explain its approach before implementation