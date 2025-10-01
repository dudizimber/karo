#!/bin/bash

set -e

echo "Installing Kubernetes Alert Reaction Operator..."

# Create namespace if it doesn't exist
kubectl create namespace alert-reaction-system --dry-run=client -o yaml | kubectl apply -f -

# Install CRD
echo "Installing Custom Resource Definition..."
kubectl apply -f config/crd/alertreaction.io_alertreactions.yaml

# Install RBAC
echo "Installing RBAC..."
kubectl apply -f config/rbac/rbac.yaml

# Install Deployment
echo "Installing Operator..."
kubectl apply -f config/manager/deployment.yaml

# Wait for deployment to be ready
echo "Waiting for operator to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/alert-reaction-operator

# Get webhook URL
echo ""
echo "Installation completed successfully!"
echo ""
echo "Webhook endpoint: http://alert-reaction-operator-webhook.default.svc.cluster.local:9090/webhook"
echo ""
echo "To configure AlertManager, add this receiver to your configuration:"
echo ""
cat << 'EOF'
receivers:
- name: 'k8s-alert-reaction-operator'
  webhook_configs:
  - url: 'http://alert-reaction-operator-webhook.default.svc.cluster.local:9090/webhook'
    send_resolved: false
    http_config:
      timeout: 10s
EOF

echo ""
echo "You can now create AlertReaction resources. See examples/ directory for samples."
