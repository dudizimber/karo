# Alert Reaction Operator

[![CI/CD Pipeline](https://github.com/dudizimber/k8s-alert-reaction-operator/actions/workflows/ci-cd.yml/badge.svg)](https://github.com/dudizimber/k8s-alert-reaction-operator/actions/workflows/ci-cd.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dudizimber/k8s-alert-reaction-operator)](https://goreportcard.com/report/github.com/dudizimber/k8s-alert-reaction-operator)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

A Kubernetes operator that creates Jobs in response to Prometheus alerts received via AlertManager webhooks.

## Overview

The Alert Reaction Operator bridges the gap between monitoring and automated remediation by allowing you to define specific actions (Kubernetes Jobs) that should be executed when certain alerts are triggered. This enables automatic incident response, scaling actions, diagnostic data collection, and other reactive operations.

### Key Features

- 🚨 **Alert-Driven Automation**: Respond to Prometheus alerts with predefined actions
- 🔄 **Job Creation**: Automatically create Kubernetes Jobs based on alert data
- 🎯 **Flexible Mapping**: One alert name maps to one AlertReaction manifest with multiple possible actions
- 🌐 **Webhook Integration**: Seamless integration with AlertManager webhooks
- 📊 **Monitoring Ready**: Built-in metrics and observability
- 🛡️ **Security Focused**: Minimal RBAC permissions and secure defaults
- ⚡ **High Performance**: Efficient controller-runtime based implementation

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────────┐
│   Prometheus    │    │   AlertManager   │    │ Alert Reaction      │
│                 │───▶│                  │───▶│ Operator            │
│   (Monitoring)  │    │   (Webhook)      │    │                     │
└─────────────────┘    └──────────────────┘    └─────────────────────┘
                                                          │
                                                          ▼
                                                ┌─────────────────────┐
                                                │ AlertReaction CRD   │
                                                │ (Configuration)     │
                                                └─────────────────────┘
                                                          │
                                                          ▼
                                                ┌─────────────────────┐
                                                │ Kubernetes Jobs     │
                                                │ (Actions)           │
                                                └─────────────────────┘
```

## Quick Start

### Prerequisites

- Kubernetes 1.19+
- Prometheus and AlertManager configured
- `kubectl` configured to access your cluster
  
### Installation

#### Option 1: Using Helm (Recommended)

```bash
# Install with default configuration
helm install alert-reaction-operator ./charts/alert-reaction-operator

# Install for development
helm install alert-reaction-operator ./charts/alert-reaction-operator \\\\
  -f ./charts/alert-reaction-operator/values-dev.yaml

# Install for production
helm install alert-reaction-operator ./charts/alert-reaction-operator \\\\
  -f ./charts/alert-reaction-operator/values-prod.yaml \\\\
  --namespace monitoring --create-namespace
```

#### Option 2: Using kubectl

```bash
# Install CRDs
kubectl apply -f config/crd/

# Install RBAC
kubectl apply -f config/rbac/

# Install the operator
kubectl apply -f config/manager/
```

#### Option 3: Using the installation script

```bash
# Make the script executable
chmod +x scripts/install-helm.sh

# Install for development
./scripts/install-helm.sh -e dev

# Install for production
./scripts/install-helm.sh -e prod -n monitoring
```

### Basic Usage

1. **Configure AlertManager** to send webhooks to the operator:

```yaml
# alertmanager.yml
route:
  group_by: ['alertname']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 1h
  receiver: 'alert-reaction-webhook'

receivers:
- name: 'alert-reaction-webhook'
  webhook_configs:
  - url: 'http://alert-reaction-operator-webhook.default.svc.cluster.local:9090/webhook'
    send_resolved: true
```

2. **Create an AlertReaction** resource to define responses:

```yaml
apiVersion: alertreaction.io/v1
kind: AlertReaction
metadata:
  name: high-cpu-reaction
  namespace: default
spec:
  alertName: "HighCPUUsage"
  actions:
  - name: "collect-diagnostics"
    image: "busybox:latest"
    command: ["sh", "-c"]
    args: ["echo 'Collecting diagnostics for high CPU on instance: $INSTANCE'"]
    env:
    - name: "INSTANCE"
      value: "labels.instance"
    - name: "ALERT_SEVERITY"
      value: "labels.severity"
  - name: "restart-service"
    image: "kubectl:latest"
    command: ["kubectl", "rollout", "restart", "deployment/my-app"]
```

3. **Test the setup** by triggering an alert or using the test script:

```bash
# Test webhook directly
scripts/test-webhook.sh http://localhost:9090/webhook

# Check created jobs
kubectl get jobs -l alert-reaction/alert-name=HighCPUUsage
```

## Configuration

### AlertReaction Custom Resource

The `AlertReaction` CRD defines how the operator should respond to specific alerts:

```yaml
apiVersion: alertreaction.io/v1
kind: AlertReaction
metadata:
  name: example-reaction
  namespace: default
spec:
  alertName: "AlertName"        # Must match the alertname label from Prometheus
  volumes:                      # Optional: Volumes to attach to jobs
  - name: "config-volume"
    configMap:
      name: "my-config"
  - name: "storage-volume"
    persistentVolumeClaim:
      claimName: "my-pvc"
  actions:                      # List of actions to execute
  - name: "action-name"         # Unique name for this action
    image: "image:tag"          # Container image to run
    command: ["cmd"]            # Command to execute
    args: ["arg1", "arg2"]      # Arguments (optional)
    serviceAccount: "action-service-account"  # Optional: Service account for this action
    env:                        # Environment variables (optional)
    - name: "VAR_NAME"
      value: "field.path"       # Dynamic value from alert data
    volumeMounts:               # Optional: Mount volumes in container
    - name: "config-volume"
      mountPath: "/config"
      readOnly: true
    - name: "storage-volume"
      mountPath: "/data"
```

### Environment Variable Substitution

Environment variables support dynamic values from alert data:

| Value Pattern | Description | Example |
|---------------|-------------|---------|
| `status` | Alert status (firing/resolved) | `firing` |
| `labels.labelname` | Alert label value | `labels.instance` → `server1.example.com` |
| `annotations.annotationname` | Alert annotation value | `annotations.summary` → `"High CPU usage detected"` |
| `static-value` | Literal string | `"production"` |

### Examples

#### Example 1: Database Backup on Critical Alert

```yaml
apiVersion: alertreaction.io/v1
kind: AlertReaction
metadata:
  name: database-backup-reaction
  namespace: production
spec:
  alertName: "DatabaseConnectionLoss"
  actions:
  - name: "emergency-backup"
    image: "postgres:13"
    command: ["pg_dump"]
    args: ["-h", "backup-server", "-U", "backup-user", "production_db"]
    env:
    - name: "ALERT_TIME"
      value: "annotations.timestamp"
    - name: "AFFECTED_INSTANCE"
      value: "labels.instance"
```

#### Example 2: Auto-scaling Response

```yaml
apiVersion: alertreaction.io/v1
kind: AlertReaction
metadata:
  name: scale-up-reaction
  namespace: default
spec:
  alertName: "HighMemoryUsage"
  actions:
  - name: "scale-deployment"
    image: "bitnami/kubectl:latest"
    command: ["kubectl"]
    args: ["scale", "deployment/web-app", "--replicas=5"]
  - name: "notify-team"
    image: "curlimages/curl:latest"
    command: ["curl"]
    args: ["-X", "POST", "https://hooks.slack.com/...", "-d", "Auto-scaled due to high memory"]
```

#### Example 3: Diagnostic Collection

```yaml
apiVersion: alertreaction.io/v1
kind: AlertReaction
metadata:
  name: diagnostics-reaction
  namespace: monitoring
spec:
  alertName: "PodCrashLooping"
  actions:
  - name: "collect-logs"
    image: "busybox:latest"
    command: ["sh", "-c"]
    args: ["kubectl logs $POD_NAME -n $NAMESPACE > /tmp/crash-logs-$(date +%s).log"]
    env:
    - name: "POD_NAME"
      value: "labels.pod"
    - name: "NAMESPACE"
      value: "labels.namespace"
  - name: "describe-pod"
    image: "bitnami/kubectl:latest"
    command: ["kubectl", "describe", "pod"]
    args: ["$POD_NAME", "-n", "$NAMESPACE"]
    env:
    - name: "POD_NAME"
      value: "labels.pod"
    - name: "NAMESPACE"
      value: "labels.namespace"
```

#### Example 4: Volume Mounting and Service Accounts

```yaml
apiVersion: alertreaction.io/v1
kind: AlertReaction
metadata:
  name: volume-example-reaction
  namespace: default
spec:
  alertName: "DiskSpaceLow"
  volumes:
  - name: "config-volume"
    configMap:
      name: "cleanup-config"
  - name: "temp-storage"
    emptyDir:
      medium: "Memory"
      sizeLimit: "1Gi"
  - name: "persistent-logs"
    persistentVolumeClaim:
      claimName: "log-storage-pvc"
  - name: "host-logs"
    hostPath:
      path: "/var/log"
      type: "Directory"
  - name: "secret-volume"
    secret:
      secretName: "cleanup-credentials"
      defaultMode: 0600
  actions:
  - name: "cleanup-disk"
    image: "alpine:latest"
    command: ["sh", "-c"]
    args: ["source /config/cleanup.sh && cleanup_old_logs /host-logs /persistent-logs"]
    serviceAccount: "cleanup-service-account"  # Service account with cleanup permissions
    volumeMounts:
    - name: "config-volume"
      mountPath: "/config"
      readOnly: true
    - name: "temp-storage"
      mountPath: "/tmp/work"
    - name: "persistent-logs"
      mountPath: "/persistent-logs"
    - name: "host-logs"
      mountPath: "/host-logs"
      readOnly: true
    - name: "secret-volume"
      mountPath: "/secrets"
      readOnly: true
    env:
    - name: "AFFECTED_NODE"
      value: "labels.instance"
    - name: "THRESHOLD"
      value: "annotations.threshold"
```

### Volume Types

The operator supports all Kubernetes volume types:

| Volume Type | Description | Use Case |
|-------------|-------------|----------|
| `configMap` | Mount ConfigMap as files | Configuration files, scripts |
| `secret` | Mount Secret as files | Credentials, certificates |
| `emptyDir` | Temporary storage | Scratch space, shared data |
| `persistentVolumeClaim` | Persistent storage | Databases, logs, artifacts |
| `hostPath` | Host filesystem access | System logs, device files |
| `downwardAPI` | Pod/container metadata | Runtime information |
| `projected` | Combine multiple sources | Complex configurations |

### Service Account Configuration

Service accounts provide identity and permissions for individual actions. Each action can specify its own service account, allowing for fine-grained security control:

```yaml
# Service account for cleanup operations
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cleanup-service-account
  namespace: default
---
# Service account for system monitoring
apiVersion: v1
kind: ServiceAccount  
metadata:
  name: monitoring-service-account
  namespace: default
---
# Role for cleanup operations
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: default
  name: cleanup-role
rules:
- apiGroups: [""]
  resources: ["pods", "configmaps"]
  verbs: ["get", "list"]
- apiGroups: [""]
  resources: ["persistentvolumeclaims"]
  verbs: ["get", "list"]
---
# Role for monitoring operations
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: default
  name: monitoring-role
rules:
- apiGroups: [""]
  resources: ["nodes", "pods"]
  verbs: ["get", "list"]
---
# RoleBinding for cleanup service account
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: cleanup-role-binding
  namespace: default
subjects:
- kind: ServiceAccount
  name: cleanup-service-account
  namespace: default
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: cleanup-role
---
# RoleBinding for monitoring service account
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: monitoring-role-binding
  namespace: default
subjects:
- kind: ServiceAccount
  name: monitoring-service-account
  namespace: default
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: cleanup-role
subjects:
- kind: ServiceAccount
  name: cleanup-service-account
  namespace: default
```

## Official Actions Repository

### Pre-built Alert Reactions

The [**dudizimber/alert-reactions**](https://github.com/dudizimber/alert-reactions) repository provides a curated collection of production-ready AlertReaction manifests for common operational scenarios.

#### Using Official Actions

1. **Browse Available Actions**: Visit [dudizimber/alert-reactions](https://github.com/dudizimber/alert-reactions) to explore all available reactions
2. **Review Configuration**: Each action includes detailed configuration examples and prerequisites
3. **Customize for Your Environment**: Modify resource limits, image versions, and environment-specific settings
4. **Deploy**: Apply the manifests to your cluster using `kubectl` or your GitOps workflow

#### Contributing Actions

Help grow the official actions library:

```bash
# Fork the repository
gh repo fork dudizimber/alert-reactions

# Create a new action category
mkdir -p my-category/my-action

# Add your AlertReaction manifest and documentation
# Submit a pull request with your contribution
```

## Operations

### Monitoring the Operator

#### Check Operator Status
```bash
# Check deployment status
kubectl get deployment alert-reaction-operator -n default

# View operator logs
kubectl logs -l app.kubernetes.io/name=alert-reaction-operator -n default -f

# Check webhook service
kubectl get svc alert-reaction-operator-webhook -n default
```

#### View AlertReaction Resources
```bash
# List all AlertReactions
kubectl get alertreactions

# Describe a specific AlertReaction
kubectl describe alertreaction high-cpu-reaction

# View AlertReaction with custom columns
kubectl get alertreactions -o custom-columns="NAME:.metadata.name,ALERT:.spec.alertName,ACTIONS:.spec.actions[*].name,TRIGGERED:.status.lastTriggered"
```

#### Monitor Created Jobs
```bash
# List jobs created by the operator
kubectl get jobs -l alert-reaction/alert-name

# View jobs for a specific alert
kubectl get jobs -l alert-reaction/alert-name=HighCPUUsage

# Check job status with details
kubectl get jobs -o wide
```

### Metrics and Observability

The operator exposes Prometheus metrics on port 8080:

```bash
# Port forward to access metrics
kubectl port-forward svc/alert-reaction-operator-metrics 8080:8080

# View metrics
curl http://localhost:8080/metrics
```

#### Key Metrics

- `alertreaction_alerts_received_total` - Total number of alerts received
- `alertreaction_jobs_created_total` - Total number of jobs created
- `alertreaction_reconcile_duration_seconds` - Time taken for reconciliation
- `controller_runtime_*` - Standard controller-runtime metrics

### Troubleshooting

#### Common Issues

**1. Webhook not receiving alerts**
```bash
# Check service and endpoint
kubectl get svc alert-reaction-operator-webhook
kubectl get endpoints alert-reaction-operator-webhook

# Test webhook manually
kubectl port-forward svc/alert-reaction-operator-webhook 9090:9090
curl -X POST http://localhost:9090/webhook \\\\
  -H "Content-Type: application/json" \\\\
  -d '{"alerts":[{"labels":{"alertname":"TestAlert"}}]}'
```

**2. Jobs not being created**
```bash
# Check if AlertReaction exists
kubectl get alertreactions

# Verify alert name matches
kubectl get alertreaction <name> -o jsonpath='{.spec.alertName}'

# Check operator logs for errors
kubectl logs -l app.kubernetes.io/name=alert-reaction-operator
```

**3. Permission issues**
```bash
# Check RBAC
kubectl get clusterrole alert-reaction-operator
kubectl get clusterrolebinding alert-reaction-operator

# Verify service account
kubectl get serviceaccount alert-reaction-operator
```

#### Debug Commands

```bash
# Get all operator-related resources
kubectl get all -l app.kubernetes.io/name=alert-reaction-operator

# Check events for issues
kubectl get events --field-selector involvedObject.name=alert-reaction-operator

# Describe operator deployment
kubectl describe deployment alert-reaction-operator

# Test webhook health
curl http://alert-reaction-operator-webhook.default.svc.cluster.local:9090/health
```

## Development

### Prerequisites

- Go 1.24+
- Docker
- kubectl
- kind (for local testing)
- make

### Building from Source

```bash
# Clone the repository
git clone https://github.com/dudizimber/k8s-alert-reaction-operator.git
cd k8s-alert-reaction-operator

# Set up development environment (including git hooks)
./scripts/setup-hooks.sh

# Build the operator
make build

# Run tests
make test

# Build Docker image
make docker-build IMG=dudizimber/alert-reaction-operator:latest
```

### Git Hooks

This project uses git hooks to ensure code quality:

```bash
# Install git hooks (automatic formatting, linting, testing)
./scripts/setup-hooks.sh

# Test hooks installation
./scripts/test-hooks.sh
```

The hooks will automatically:
- **pre-commit**: Format code, basic linting, check for common issues
- **pre-push**: Run full tests, generate manifests, comprehensive linting  
- **commit-msg**: Validate conventional commit message format

See [scripts/hooks/README.md](scripts/hooks/README.md) for detailed information.

### Copilot Prompt System

This project includes a comprehensive prompt library for GitHub Copilot to streamline development:

```bash
# Quick reference for common tasks
cat .copilot/quick-prompts.md

# Detailed templates for complex tasks
ls .copilot/prompts/

# Project context templates
cat .copilot/context-templates.md
```

**Available prompt templates:**
- **Feature Development** - `.copilot/prompts/add-feature.md`
- **Bug Fixing** - `.copilot/prompts/fix-bug.md` 
- **Testing** - `.copilot/prompts/add-tests.md`
- **Code Review** - `.copilot/prompts/code-review.md`
- **Documentation** - `.copilot/prompts/write-docs.md`
- **Performance** - `.copilot/prompts/optimize-performance.md`
- **Refactoring** - `.copilot/prompts/refactor-code.md`
- **CI/CD** - `.copilot/prompts/ci-cd-updates.md`
- **CRD Changes** - `.copilot/prompts/crd-changes.md`
- **Controller Logic** - `.copilot/prompts/controller-logic.md`

See [.copilot/README.md](.copilot/README.md) for complete usage guide.

### Running Locally

```bash
# Install CRDs
make install

# Run the operator locally
make run

# In another terminal, create a test AlertReaction
kubectl apply -f examples/alertreactions.yaml
```

### Testing

```bash
# Run unit tests
make test

# Run integration tests
make test-integration

# Lint code
make lint

# Test webhook manually
./scripts/test-webhook.sh http://localhost:9090/webhook
```

## Security

### RBAC Permissions

The operator requires minimal permissions:

```yaml
# AlertReaction CRD management
- alertreaction.io: alertreactions (all verbs)
- alertreaction.io: alertreactions/status (get, update, patch)
- alertreaction.io: alertreactions/finalizers (update)

# Job management
- batch: jobs (all verbs)

# Configuration access
- "": configmaps, secrets (get, list, watch)

# Leader election
- "": configmaps (all verbs for leader election)
- coordination.k8s.io: leases (all verbs for leader election)
```

### Security Best Practices

1. **Run as non-root user** (UID 65532)
2. **Read-only root filesystem**
3. **Dropped capabilities** (ALL)
4. **Network policies** to restrict traffic
5. **Resource limits** to prevent resource exhaustion
6. **Secure image scanning** in CI/CD pipeline

### Network Policies

Example network policy to restrict webhook access:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: alert-reaction-operator
spec:
  podSelector:
    matchLabels:
      app.kubernetes.io/name: alert-reaction-operator
  policyTypes:
  - Ingress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: monitoring  # AlertManager namespace
    ports:
    - protocol: TCP
      port: 9090
```

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup

1. **Fork the repository**
2. **Clone your fork**
   ```bash
   git clone https://github.com/dudizimber/k8s-alert-reaction-operator.git
   ```
3. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```
4. **Make your changes**
5. **Run tests**
   ```bash
   make test
   make lint
   ```
6. **Commit and push**
   ```bash
   git commit -m "Add your feature"
   git push origin feature/your-feature-name
   ```
7. **Create a Pull Request**

### Code Guidelines

- Follow Go best practices and formatting (`gofmt`, `golint`)
- Add tests for new functionality  
- Update documentation for API changes
- Ensure CI/CD pipeline passes
- Use [Conventional Commits](https://www.conventionalcommits.org/) format
- Git hooks will automatically enforce formatting and linting

## Versioning

This project follows [Semantic Versioning](https://semver.org/):

- **MAJOR**: Incompatible API changes
- **MINOR**: Backwards-compatible functionality additions
- **PATCH**: Backwards-compatible bug fixes

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a detailed history of changes.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Support

- **Documentation**: [GitHub Wiki](https://github.com/dudizimber/k8s-alert-reaction-operator/wiki)
- **Issues**: [GitHub Issues](https://github.com/dudizimber/k8s-alert-reaction-operator/issues)
- **Discussions**: [GitHub Discussions](https://github.com/dudizimber/k8s-alert-reaction-operator/discussions)

## Acknowledgments

- Built with [Kubebuilder](https://kubebuilder.io/)
- Powered by [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime)
- HTTP server using [Gin](https://github.com/gin-gonic/gin)
- Inspired by the Kubernetes and Prometheus communities

---

**Made with ❤️ for the Kubernetes community**