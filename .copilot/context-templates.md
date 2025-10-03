# Project Context Templates

## Full Project Context
```
I'm working on the k8s-alert-reaction-operator, a Kubernetes operator that:

**Purpose**: Creates Kubernetes Jobs in response to Prometheus alerts received via AlertManager webhooks

**Architecture**:
- Built with Go 1.24+ and controller-runtime
- Uses AlertReaction CRD (v1alpha1) with Prometheus-style matchers
- Webhook endpoint receives alerts from AlertManager
- Controller processes alerts and creates Jobs based on matching conditions

**Key Components**:
- AlertReaction CRD with matchers (=, !=, =~, !~), actions, volumes
- Controller that supports multiple AlertReactions per alert
- Webhook handler for AlertManager integration
- Comprehensive matcher evaluation with regex support

**Current Features**:
- Prometheus-style alert matching with 4 operators
- Multiple actions per AlertReaction
- Volume mounting and service account support
- Environment variable substitution from alert data
- Git hooks for automated code quality
- Multi-architecture container builds
- Comprehensive test suite

**Development Practices**:
- Conventional commits with git hooks
- Table-driven tests with comprehensive coverage
- GitHub Actions CI/CD with security scanning
- Kubebuilder patterns and controller-runtime best practices

{your_specific_question_or_request}
```

## Kubernetes Operator Context
```
This is a Kubernetes operator built with:
- **Framework**: Kubebuilder + controller-runtime
- **API Version**: v1alpha1 (storage version)
- **CRD**: AlertReaction with spec containing alertName, matchers, actions, volumes
- **Controller**: Reconciles AlertReaction resources and creates Jobs
- **Webhook**: HTTP endpoint for AlertManager webhook integration
- **RBAC**: Minimal permissions for AlertReactions, Jobs, ConfigMaps, Secrets

**Controller Pattern**:
- Watches AlertReaction resources
- Processes alerts via webhook (not traditional reconciliation)
- Creates Jobs with owner references for garbage collection
- Updates status with processing results and error conditions

{your_operator_specific_question}
```

## Alert Processing Context
```
The k8s-alert-reaction-operator processes alerts with this workflow:

1. **Alert Reception**: AlertManager sends webhook to /webhook endpoint
2. **Alert Parsing**: Extract labels, annotations, status from Prometheus alert
3. **Matcher Evaluation**: Find AlertReactions with matching conditions using:
   - Exact match (=): label.value = "expected"
   - Not equal (!=): label.value != "unwanted"  
   - Regex match (=~): label.instance =~ "web-.*"
   - Regex not match (!~): label.env !~ "(dev|test)"
4. **Job Creation**: Create Kubernetes Jobs for each matching action
5. **Status Update**: Record processing results in AlertReaction status

**Matcher Logic**:
- All matchers in an AlertReaction must match (AND logic)
- Multiple AlertReactions can match the same alert
- Empty matchers list matches all alerts
- Supports alert label and annotation matching

{your_alert_processing_question}
```

## Testing Context  
```
The k8s-alert-reaction-operator uses comprehensive testing:

**Test Structure**:
- Table-driven tests for multiple scenarios
- Fake Kubernetes client for controller testing
- Mock HTTP servers for webhook testing
- Test fixtures for AlertReaction resources and alert payloads

**Coverage Areas**:
- Unit tests: Individual functions and methods
- Integration tests: Controller with fake client
- Matcher tests: All operator combinations with various inputs
- Webhook tests: HTTP endpoint behavior
- Error scenario tests: Failure conditions and recovery

**Test Patterns**:
```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected OutputType
        wantErr  bool
    }{
        // Test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

{your_testing_question}
```