# AlertReaction Matchers

AlertReaction supports advanced matching conditions using Prometheus-style operators through the `matchers` field. This allows you to specify fine-grained conditions for triggering reactions based on alert labels and annotations, just like Prometheus silences.

## Matcher Configuration

Matchers are defined as an array of conditions that must **all** be satisfied for the reaction to trigger. If no matchers are specified, only the `alertName` is used for matching.

```yaml
spec:
  alertName: "ServiceDown"
  matchers:
    - name: severity
      operator: "="
      value: critical
    - name: service
      operator: "=~"
      value: "web-server|api-gateway"
```

## Matcher Fields

### `name` (required)
The name of the alert label or annotation to match against:
- For labels: Use the label name directly (e.g., `severity`, `instance`, `service`)
- For annotations: Prefix with `annotations.` (e.g., `annotations.runbook`, `annotations.summary`)

### `operator` (required)
Prometheus-style matching operators:

| Operator | Description | Example Usage |
|----------|-------------|---------------|
| `=` | Exact equality match | `severity = "critical"` |
| `!=` | Not equal match | `environment != "test"` |
| `=~` | Regular expression match | `instance =~ "prod-.*"` |
| `!~` | Negative regular expression match | `service !~ ".*test.*"` |

### `value` (required)
The value to match against:
- For `=` and `!=`: Exact string value
- For `=~` and `!~`: Valid regular expression pattern

## Examples

### Basic Label Matching
```yaml
matchers:
  - name: severity
    operator: "="
    value: critical
```

### Environment Exclusion
```yaml
matchers:
  - name: environment
    operator: "!="
    value: test
```

### Multiple Services with Regex
```yaml
matchers:
  - name: service
    operator: "=~"
    value: "web-server|api-gateway|database"
```

### Instance Pattern Matching
```yaml
matchers:
  - name: instance
    operator: "=~"
    value: "prod-.*-[0-9]+"
  - name: environment
    operator: "!~"
    value: ".*test.*|.*staging.*"
```

### Annotation Matching
```yaml
matchers:
  - name: annotations.runbook
    operator: "=~"
    value: ".*emergency.*"
  - name: annotations.severity
    operator: "="
    value: critical
```

### Complex Real-World Example
```yaml
apiVersion: alertreaction.io/v1alpha1
kind: AlertReaction
metadata:
  name: production-critical-alerts
spec:
  alertName: "ServiceDown"
  matchers:
    # Must be critical severity
    - name: severity
      operator: "="
      value: critical
    
    # Must be production environment
    - name: environment
      operator: "="
      value: production
    
    # Must be one of our core services
    - name: service
      operator: "=~"
      value: "web-server|api-gateway|database|cache"
    
    # Must have emergency runbook
    - name: annotations.runbook
      operator: "=~"
      value: ".*emergency.*"
    
    # Exclude test instances
    - name: instance
      operator: "!~"
      value: ".*test.*|.*dev.*"
  
  actions:
    - name: "emergency-response"
      image: "my-registry/emergency-responder:latest"
      # ... action configuration
```

## Prometheus Compatibility

The matcher syntax is designed to be familiar to Prometheus users:

| Prometheus Silence | AlertReaction Matcher |
|-------------------|----------------------|
| `severity="critical"` | `name: severity, operator: "=", value: critical` |
| `instance!="test"` | `name: instance, operator: "!=", value: test` |
| `service=~"web-.*"` | `name: service, operator: "=~", value: "web-.*"` |
| `env!~".*test.*"` | `name: env, operator: "!~", value: ".*test.*"` |

## Backwards Compatibility

The `matchers` field is optional. Existing AlertReaction resources without matchers will continue to work exactly as before, matching only on the `alertName` field.

## Matcher Evaluation Logic

1. All matchers must evaluate to `true` for the reaction to trigger (AND logic)
2. If any matcher evaluates to `false`, the reaction is skipped
3. If no matchers are specified, only the `alertName` is checked
4. Label matching is done against alert labels directly
5. Annotation matching requires the `annotations.` prefix

## Performance Considerations

- Matchers are evaluated for every incoming alert
- Regular expression matchers (`=~`, `!~`) are more expensive than equality checks (`=`, `!=`)
- Use regex patterns efficiently - consider anchoring patterns when appropriate
- Place most selective matchers first to fail fast when possible