# Bug Fix Prompt Template

## Context
# Bug Fix Prompt Template

## Prompt
```
I'm working on Karo (Kubernetes Alert Reaction Operator), a Kubernetes operator built with Go 1.24+ and controller-runtime. The operator processes Prometheus alerts via webhooks and creates Kubernetes Jobs based on AlertReaction CRDs with Prometheus-style matchers.

## Bug Report
**Issue:** {bug_title}

### Problem Description
**What's happening:** {detailed_description_of_the_issue}

**Expected behavior:** {what_should_happen}

**Actual behavior:** {what_is_actually_happening}

### Environment
- **Go version:** {go_version}
- **Kubernetes version:** {k8s_version}
- **Operator version/commit:** {version_or_commit}
- **Alert source:** {prometheus_alertmanager_version}

### Reproduction Steps
1. {step_1}
2. {step_2}
3. {step_3}
4. Observe: {what_you_observe}

### Error Messages/Logs
```
{paste_relevant_error_messages_or_logs}
```

### Affected Files
Based on the symptoms, these files might be involved:
- [ ] `{primary_suspected_file}`
- [ ] `{secondary_suspected_file}`
- [ ] `{other_potential_files}`

### Code Context
**Relevant code snippet:**
```go
{paste_relevant_code_snippet_if_known}
```

**Related functions/methods:**
- `{function_name_1}()` in `{file_path_1}`
- `{function_name_2}()` in `{file_path_2}`

## Investigation Areas
Please help me investigate:

1. **Root Cause Analysis**: What could be causing this issue?
2. **Code Review**: Are there bugs in the suspected code areas?
3. **Logic Flaws**: Are there logical errors in the implementation?
4. **Edge Cases**: Are there unhandled edge cases?
5. **Race Conditions**: Could this be a concurrency issue?
6. **Resource Issues**: Could this be related to resource limits or permissions?

## Fix Requirements
- [ ] Must maintain backward compatibility
- [ ] Should include comprehensive tests for the fix
- [ ] Must not introduce new bugs or regressions
- [ ] Should include logging for better debugging
- [ ] Must follow existing code patterns and conventions

## Testing Strategy
How should I test the fix?
- [ ] Unit tests for the specific bug scenario
- [ ] Integration tests with real AlertManager webhooks
- [ ] Edge case testing
- [ ] Regression testing for related functionality

## Questions for Copilot
1. What's the most likely root cause of this issue?
2. What debugging steps should I take to confirm the cause?
3. What's the safest way to fix this without breaking existing functionality?
4. What tests should I add to prevent this regression?
5. Are there similar issues elsewhere in the codebase?

---

## Usage Instructions
1. Replace all `{placeholders}` with specific details about your bug
2. Include actual error messages, logs, and code snippets
3. Be as specific as possible about reproduction steps
4. Copy and paste into Copilot chat for debugging assistance

## Example Usage
```
**Issue:** Controller creates duplicate Jobs for the same alert

### Problem Description
**What's happening:** When AlertManager sends the same alert multiple times (which is normal), the controller creates multiple Jobs instead of recognizing it's the same alert.

**Expected behavior:** Controller should create only one Job per unique alert occurrence and handle duplicate webhook calls gracefully.

**Actual behavior:** Multiple Jobs are created for alerts with the same fingerprint, leading to resource waste and duplicate actions.

### Reproduction Steps
1. Configure AlertManager to send webhooks to the operator
2. Create an AlertReaction for "HighCPUUsage" alerts
3. Trigger a HighCPUUsage alert that persists for >1 minute
4. Observe: Multiple Jobs created with labels showing the same alert fingerprint
```