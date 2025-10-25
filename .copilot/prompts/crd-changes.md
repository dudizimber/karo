# CRD Changes Prompt Template

## Context
# CRD Changes Prompt Template

## Prompt
```
I'm working on Karo (Kubernetes Alert Reaction Operator), a Kubernetes operator that uses a custom AlertReaction CRD. The CRD currently supports Prometheus-style matchers and defines how alerts should trigger Kubernetes Jobs.

## Current CRD Structure
**Current AlertReaction spec:**
```go
type AlertReactionSpec struct {
    AlertName    string         `json:"alertName"`
    Matchers     []AlertMatcher `json:"matchers,omitempty"`
    Actions      []Action       `json:"actions"`
    Volumes      []Volume       `json:"volumes,omitempty"`
}

type AlertMatcher struct {
    Name     string        `json:"name"`
    Operator MatchOperator `json:"operator"`
    Value    string        `json:"value"`
}

type Action struct {
    Name           string            `json:"name"`
    Image          string            `json:"image"`
    Command        []string          `json:"command,omitempty"`
    Args           []string          `json:"args,omitempty"`
    Env            []EnvVar          `json:"env,omitempty"`
    ServiceAccount string            `json:"serviceAccount,omitempty"`
    VolumeMounts   []VolumeMount     `json:"volumeMounts,omitempty"`
}
```

## Proposed Changes
**What I want to add/modify:**
{detailed_description_of_crd_changes}

**New fields needed:**
- `{field_name}` ({type}) - {description_and_purpose}
- `{field_name}` ({type}) - {description_and_purpose}

**Modified fields:**
- `{existing_field}` - {how_it_should_change}

**Removed fields (if any):**
- `{field_to_remove}` - {reason_for_removal}

## Requirements
**Functional requirements:**
- [ ] {requirement_1}
- [ ] {requirement_2}
- [ ] {requirement_3}

**Technical constraints:**
- [ ] Must maintain backward compatibility with existing AlertReaction resources
- [ ] Should follow Kubernetes API conventions
- [ ] Must include proper JSON/YAML tags
- [ ] Should have appropriate validation rules
- [ ] Must support OpenAPI schema generation

**Validation needs:**
- [ ] Required field validation
- [ ] Format validation (e.g., regex patterns)
- [ ] Range validation (e.g., min/max values)
- [ ] Enum validation for restricted values
- [ ] Cross-field validation

## Use Cases
**Primary use case:**
{describe_the_main_use_case_for_these_changes}

**Example configuration:**
```yaml
apiVersion: karo.io/v1alpha1
kind: AlertReaction
metadata:
  name: example-with-new-fields
spec:
  alertName: "HighCPUUsage"
  {new_field_1}: {example_value_1}
  {new_field_2}: {example_value_2}
  # ... rest of spec
```

**Expected behavior:**
{describe_how_the_controller_should_handle_these_new_fields}

## Backward Compatibility
**Migration strategy:**
- [ ] New fields are optional (have default values)
- [ ] Existing resources continue to work unchanged
- [ ] Controller handles both old and new field combinations
- [ ] Clear migration path for users who want to adopt new fields

**Default values:**
- `{new_field_1}`: {default_value_and_reasoning}
- `{new_field_2}`: {default_value_and_reasoning}

## Files to Update
**API types:**
- [ ] `api/v1alpha1/alertreaction_types.go` - Main CRD structure
- [ ] `api/v1alpha1/zz_generated.deepcopy.go` - Generated after running `make generate`

**Controller logic:**
- [ ] `controllers/alertreaction_controller.go` - Handle new fields in reconciliation
- [ ] Related controller helper functions

**Validation:**
- [ ] Webhook validation (if complex validation needed)
- [ ] OpenAPI schema markers for basic validation

**Generated manifests:**
- [ ] `config/crd/bases/karo.io_alertreactions.yaml` - Generated after `make manifests`

**Tests:**
- [ ] Unit tests for new API types
- [ ] Controller tests with new field scenarios
- [ ] Integration tests for backward compatibility

**Documentation:**
- [ ] README examples
- [ ] API documentation
- [ ] Migration guide (if needed)

## Request for Copilot
Please help me implement these CRD changes by:

1. **Updating the Go struct definitions** with proper tags and validation
2. **Adding appropriate OpenAPI schema markers** for validation
3. **Ensuring backward compatibility** with existing resources
4. **Providing example YAML** showing the new fields in use
5. **Updating controller logic** to handle the new fields
6. **Creating comprehensive tests** for the new functionality
7. **Documenting the changes** clearly

**Specific questions:**
1. What's the best way to structure these new fields in the Go types?
2. What validation markers should I use for the requirements above?
3. How should the controller logic change to handle these new fields?
4. Are there any security implications I should consider?
5. What's the migration impact on existing users?

---

## Usage Instructions
1. Replace `{placeholders}` with your specific CRD change requirements
2. Include current CRD structure for context
3. Describe the use cases and expected behavior clearly
4. Specify validation and compatibility requirements
5. Copy and paste into Copilot chat

## Example Usage
```
**What I want to add/modify:**
Add support for alert throttling and rate limiting to prevent spam from noisy alerts.

**New fields needed:**
- `throttle` (ThrottleConfig) - Configuration for limiting how often actions are executed
- `cooldownPeriod` (string) - Minimum time between action executions (e.g., "5m", "1h")

**Requirements:**
- Should prevent duplicate job creation for the same alert within the cooldown period
- Must track last execution time per AlertReaction
- Should be optional with sensible defaults (no throttling by default)
```