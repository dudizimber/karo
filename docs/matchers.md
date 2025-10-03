# AlertReaction Matchers

AlertReaction now supports advanced matching conditions through the `matchers` field. This allows you to specify fine-grained conditions for triggering reactions based on alert attributes beyond just the alert name.

## Matcher Configuration

Matchers are defined as an array of conditions that must **all** be satisfied for the reaction to trigger. If no matchers are specified, only the `alertName` is used for matching.

```yaml
spec:
  alertName: "ServiceDown"
  matchers:
    - name: "labels.severity"
      operator: Equal
      values: ["critical"]
    - name: "labels.service"
      operator: In
      values: ["web-server", "api-gateway"]
```

## Matcher Fields

### `name` (required)
The path to the alert attribute to match against. Common paths include:
- `labels.<label-name>` - Match against alert labels
- `annotations.<annotation-name>` - Match against alert annotations  
- `status` - Match against alert status
- `startsAt` - Match against alert start time
- `endsAt` - Match against alert end time

### `operator` (required)
The comparison operator to use. Supported operators:

| Operator | Description | Example |
|----------|-------------|---------|
| `Equal` | Exact string match | `value == "critical"` |
| `NotEqual` | String not equal | `value != "warning"` |
| `In` | Value in list | `value in ["critical", "high"]` |
| `NotIn` | Value not in list | `value not in ["low", "info"]` |
| `Exists` | Attribute exists | Label/annotation is present |
| `DoesNotExist` | Attribute missing | Label/annotation is not present |
| `GreaterThan` | Numeric comparison | `value > "90"` |
| `LessThan` | Numeric comparison | `value < "10"` |
| `Regex` | Regular expression | `value matches ".*prod.*"` |
| `NotRegex` | Negative regex | `value does not match ".*test.*"` |

### `values` (optional)
Array of values to match against. Not required for `Exists` and `DoesNotExist` operators.
- For `In` and `NotIn`: Multiple values are supported
- For other operators: Only the first value is used
- For numeric operators: Values are converted to numbers for comparison

## Examples

### Basic Label Matching
```yaml
matchers:
  - name: "labels.severity"
    operator: Equal
    values: ["critical"]
```

### Multiple Service Matching
```yaml
matchers:
  - name: "labels.service"
    operator: In
    values: ["web-server", "api-gateway", "database"]
```

### Existence Check
```yaml
matchers:
  - name: "labels.instance"
    operator: Exists
  - name: "annotations.runbook"
    operator: DoesNotExist
```

### Numeric Thresholds
```yaml
matchers:
  - name: "annotations.value"
    operator: GreaterThan
    values: ["90"]
```

### Regular Expression Matching
```yaml
matchers:
  - name: "labels.instance"
    operator: Regex
    values: ["prod-.*-[0-9]+"]
  - name: "labels.environment"
    operator: NotRegex
    values: [".*test.*", ".*staging.*"]
```

### Complex Example
```yaml
apiVersion: alerts.davidzimberknopf.io/v1alpha1
kind: AlertReaction
metadata:
  name: production-critical-alerts
spec:
  alertName: "ServiceDown"
  matchers:
    # Must be critical severity
    - name: "labels.severity"
      operator: Equal
      values: ["critical"]
    
    # Must be production environment
    - name: "labels.environment"
      operator: Equal
      values: ["production"]
    
    # Must be one of our core services
    - name: "labels.service"
      operator: In
      values: ["web-server", "api-gateway", "database", "cache"]
    
    # Must have a runbook
    - name: "annotations.runbook"
      operator: Exists
    
    # Exclude test instances
    - name: "labels.instance"
      operator: NotRegex
      values: [".*test.*", ".*dev.*"]
  
  actions:
    - name: "emergency-response"
      image: "my-registry/emergency-responder:latest"
      # ... action configuration
```

## Backwards Compatibility

The `matchers` field is optional. Existing AlertReaction resources without matchers will continue to work exactly as before, matching only on the `alertName` field.

## Matcher Evaluation Logic

1. All matchers must evaluate to `true` for the reaction to trigger
2. If any matcher evaluates to `false`, the reaction is skipped
3. If no matchers are specified, only the `alertName` is checked
4. Matchers are evaluated in the order they are specified (though all must pass)

## Performance Considerations

- Matchers are evaluated for every incoming alert
- Regular expression matchers (`Regex`, `NotRegex`) are more expensive than simple equality checks
- Consider using `In`/`NotIn` operators instead of multiple `Equal`/`NotEqual` matchers when possible
- Place most selective matchers first to fail fast when possible