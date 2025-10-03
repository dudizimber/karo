# Quick Reference Prompts

## Feature Development
```
I need to add {feature_name} to the k8s-alert-reaction-operator.

Requirements:
- {requirement_1}
- {requirement_2}

Files likely affected:
- api/v1alpha1/alertreaction_types.go
- controllers/alertreaction_controller.go
- Tests and documentation

Please help me implement this feature following our existing patterns.
```

## Bug Fix
```
I'm experiencing {issue_description} in the k8s-alert-reaction-operator.

Expected: {expected_behavior}
Actual: {actual_behavior}

Error logs:
{paste_error_logs}

Please help me debug and fix this issue.
```

## Commit Message
```
I made changes to {files_modified} in the k8s-alert-reaction-operator.

Changes:
- {change_1}
- {change_2}

Type: feat/fix/docs/refactor/test/chore
Scope: controller/webhook/api/crd/ci

Please write a conventional commit message for these changes.
```

## Add Tests
```
I need comprehensive tests for {function_name} in {file_path}.

The function should be tested for:
- Happy path: {scenarios}
- Error cases: {error_scenarios}
- Edge cases: {edge_cases}

Please create table-driven tests following our project patterns.
```

## Code Review
```
Please review this code for the k8s-alert-reaction-operator:

{paste_code_to_review}

Focus on:
- Correctness and bug identification
- Go and Kubernetes best practices
- Performance and security
- Testing adequacy

Are there any issues or improvements you'd suggest?
```

## Performance Optimization
```
The k8s-alert-reaction-operator has performance issues:

Current metrics:
- Alert processing: {latency}
- Memory usage: {memory}
- Throughput: {throughput}

Target: {performance_goals}

Please analyze and suggest optimizations for:
{paste_performance_critical_code}
```

## Documentation
```
I need to document {feature_or_topic} for the k8s-alert-reaction-operator.

Target audience: {users/developers/operators}
Content needed:
- {topic_1}
- {topic_2}

Please create comprehensive documentation with examples.
```

## Refactoring
```
This code in the k8s-alert-reaction-operator needs refactoring:

{paste_code_to_refactor}

Issues:
- {problem_1}
- {problem_2}

Goals:
- Improve readability
- Better testability
- Follow Go best practices

Please suggest a refactored structure.
```

## CI/CD Updates
```
I want to improve the k8s-alert-reaction-operator CI/CD pipeline.

Current pipeline: {describe_current_state}

Desired improvements:
- {improvement_1}
- {improvement_2}

Please help update the GitHub Actions workflows.
```

## CRD Changes
```
I need to modify the AlertReaction CRD in the k8s-alert-reaction-operator.

Current structure: {describe_current_crd}

Proposed changes:
- Add field: {field_name} ({type}) - {purpose}
- Modify: {existing_field} - {how_it_changes}

Requirements:
- Backward compatibility
- Proper validation
- Controller integration

Please help implement these CRD changes.
```