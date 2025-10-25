# Karo Helm Chart

This Helm chart deploys Karo (Kubernetes Alert Reaction Operator), a Kubernetes operator that creates Jobs in response to Prometheus alerts.

## Overview

Karo watches for `AlertReaction` custom resources and creates Kubernetes Jobs when matching Prometheus/Alertmanager alerts are received via webhook. This enables automated incident response and remediation workflows.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.8+
- Prometheus/Alertmanager (for sending alerts)

## Installation

### Add Helm Repository

```bash
# Add the Helm repository
helm repo add karo https://dudizimber.github.io/karo/
helm repo update
```

### Install from Repository

```bash
# Install with default values
helm install my-operator karo/karo

# Install with custom values
helm install my-operator karo/karo -f values.yaml

# Install in specific namespace
helm install my-operator karo/karo -n monitoring --create-namespace
```

### Install from OCI Registry

```bash
# Install directly from GitHub Container Registry
helm install my-operator oci://ghcr.io/dudizimber/charts/karo --version 0.1.0

# Install with custom values
helm install my-operator oci://ghcr.io/dudizimber/charts/karo --version 0.1.0 -f values.yaml
```

### Install from Release Assets

```bash
# Download chart from GitHub releases
curl -L https://github.com/dudizimber/karo/releases/download/v0.1.0/karo-0.1.0.tgz -o karo-0.1.0.tgz

# Install from local file
helm install my-operator ./karo-0.1.0.tgz
```

### Custom Resource Definitions (CRDs)

The chart automatically installs the required CRDs (`AlertReaction`) as part of the installation process. The CRDs are bundled with the chart and will be installed before the operator deployment.

**Note**: When upgrading the chart, CRDs are not automatically updated by Helm. If you need to update CRDs to a newer version, you can:

```bash
# Update CRDs manually (if needed during upgrades)
kubectl apply -f https://raw.githubusercontent.com/dudizimber/karo/main/config/crd/karo.io_alertreactions.yaml
```

## Configuration

### Basic Configuration

```yaml
# values.yaml
image:
  repository: docker.io/dudizimber/karo
  tag: "0.1.0"
  pullPolicy: IfNotPresent

replicaCount: 1

service:
  type: ClusterIP
  port: 8080
  targetPort: 8080

webhook:
  enabled: true
  port: 9443
  path: /webhook
```

### Security Configuration

```yaml
# values.yaml
podSecurityContext:
  runAsNonRoot: true
  runAsUser: 65532
  fsGroup: 65532

securityContext:
  allowPrivilegeEscalation: false
  capabilities:
    drop:
    - ALL
  readOnlyRootFilesystem: true
  runAsNonRoot: true
  runAsUser: 65532

resources:
  limits:
    cpu: 500m
    memory: 512Mi
  requests:
    cpu: 100m
    memory: 128Mi
```

### Advanced Configuration

```yaml
# values.yaml
# Custom environment variables
env:
  - name: LOG_LEVEL
    value: "info"
  - name: WEBHOOK_PORT
    value: "9443"

# Node selector
nodeSelector:
  kubernetes.io/os: linux

# Tolerations
tolerations:
  - effect: NoSchedule
    key: node-role.kubernetes.io/master
    operator: Exists

# Affinity
affinity:
  podAntiAffinity:
    preferredDuringSchedulingIgnoredDuringExecution:
    - weight: 100
      podAffinityTerm:
        labelSelector:
          matchExpressions:
          - key: app.kubernetes.io/name
            operator: In
            values:
            - karo
        topologyKey: kubernetes.io/hostname
```

## Configuration Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `image.repository` | Container image repository | `docker.io/dudizimber/karo` |
| `image.tag` | Container image tag | `"latest"` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `imagePullSecrets` | Image pull secrets | `[]` |
| `nameOverride` | Override chart name | `""` |
| `fullnameOverride` | Override full resource names | `""` |
| `replicaCount` | Number of operator replicas | `1` |
| `podAnnotations` | Pod annotations | `{}` |
| `podLabels` | Pod labels | `{}` |
| `podSecurityContext` | Pod security context | See values.yaml |
| `securityContext` | Container security context | See values.yaml |
| `service.type` | Service type | `ClusterIP` |
| `service.port` | Service port | `8080` |
| `service.targetPort` | Target port | `8080` |
| `webhook.enabled` | Enable webhook service | `true` |
| `webhook.port` | Webhook port | `9443` |
| `webhook.path` | Webhook path | `/webhook` |
| `webhook.tls.enabled` | Enable TLS for webhook | `true` |
| `webhook.tls.secretName` | TLS secret name | `""` |
| `ingress.enabled` | Enable ingress | `false` |
| `resources` | Resource limits and requests | `{}` |
| `nodeSelector` | Node selector | `{}` |
| `tolerations` | Tolerations | `[]` |
| `affinity` | Affinity | `{}` |
| `env` | Environment variables | `[]` |
| `serviceAccount.create` | Create service account | `true` |
| `serviceAccount.name` | Service account name | `""` |
| `serviceAccount.annotations` | Service account annotations | `{}` |
| `rbac.create` | Create RBAC resources | `true` |
| `monitoring.enabled` | Enable monitoring resources | `false` |
| `monitoring.serviceMonitor.enabled` | Create ServiceMonitor | `false` |

## Usage

### Creating AlertReaction Resources

After installing the operator, create `AlertReaction` resources to define automated responses:

```yaml
apiVersion: karo.io/v1alpha1
kind: AlertReaction
metadata:
  name: high-cpu-reaction
  namespace: default
spec:
  alertName: "HighCPUUsage"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: cpu-investigation
            image: busybox:1.35
            command: ["sh", "-c"]
            args:
            - |
              echo "Investigating high CPU usage..."
              # Add your remediation logic here
              kubectl top nodes
              kubectl top pods --all-namespaces --sort-by=cpu
          restartPolicy: Never
```

### Configuring Alertmanager

Configure Alertmanager to send webhooks to the operator:

```yaml
# alertmanager.yml
global:
  resolve_timeout: 5m

route:
  group_by: ['alertname']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 1h
  receiver: 'webhook'

receivers:
- name: 'webhook'
  webhook_configs:
  - url: 'http://karo.monitoring.svc.cluster.local:8080/webhook'
    http_config:
      basic_auth:
        username: 'alertmanager'
        password: 'secret'
```

## Monitoring

### Prometheus ServiceMonitor

Enable monitoring to scrape metrics from the operator:

```yaml
# values.yaml
monitoring:
  enabled: true
  serviceMonitor:
    enabled: true
    interval: 30s
    path: /metrics
    labels:
      prometheus: kube-prometheus
```

### Custom Metrics

The operator exposes the following metrics:

- `alert_reactions_total` - Total number of alert reactions processed
- `jobs_created_total` - Total number of jobs created
- `webhook_requests_total` - Total number of webhook requests received
- `errors_total` - Total number of errors encountered

## Troubleshooting

### Check Operator Logs

```bash
# Get operator pod name
kubectl get pods -l app.kubernetes.io/name=karo

# View logs
kubectl logs -f <pod-name>
```

### Verify Custom Resources

```bash
# Check AlertReaction resources
kubectl get alertreactions

# Describe specific AlertReaction
kubectl describe alertreaction <name>
```

### Check Webhook Connectivity

```bash
# Test webhook endpoint
kubectl run curl --image=curlimages/curl --rm -it --restart=Never -- \
  curl -X POST \
  http://karo.default.svc.cluster.local:8080/webhook \
  -H "Content-Type: application/json" \
  -d '{"alerts": [{"labels": {"alertname": "TestAlert"}}]}'
```

### Common Issues

1. **Webhook not receiving alerts**: Check Alertmanager configuration and network connectivity
2. **Jobs not being created**: Verify AlertReaction resource matches alert labels
3. **Permission errors**: Ensure RBAC permissions are correctly configured
4. **Image pull errors**: Check image repository and pull secrets

## Security Considerations

- The operator runs with minimal privileges using a non-root user
- RBAC is configured with least-privilege access
- Webhook endpoint can be secured with TLS certificates
- Consider network policies to restrict traffic

## Contributing

For development and contribution guidelines, see the [main repository](https://github.com/dudizimber/karo).

## License

This chart is licensed under the Apache License 2.0. See [LICENSE](https://github.com/dudizimber/karo/blob/main/LICENSE) for details.

## Support

- **Issues**: [GitHub Issues](https://github.com/dudizimber/karo/issues)
- **Discussions**: [GitHub Discussions](https://github.com/dudizimber/karo/discussions)
- **Official Actions**: [karo-reactions repository](https://github.com/dudizimber/karo-reactions)
