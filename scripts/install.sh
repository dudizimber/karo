#!/bin/bash

set -e

echo "Installing Karo (Kubernetes Alert Reaction Operator)..."

# Create namespace
kubectl create namespace karo-system --dry-run=client -o yaml | kubectl apply -f -

# Install CRD
echo "Installing Custom Resource Definition..."
kubectl apply -f config/crd/karo.io_alertreactions.yaml

# Install RBAC
echo "Installing RBAC..."
kubectl apply -f config/rbac/rbac.yaml

# Install Deployment
echo "Installing Operator..."
kubectl apply -f config/manager/deployment.yaml

# Wait for deployment to be ready
echo "Waiting for operator to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/karo

# Get webhook URL
echo ""
echo "Installation completed successfully!"
echo ""
echo "Webhook endpoint: http://karo-webhook.default.svc.cluster.local:9090/webhook"
echo ""
echo "To configure AlertManager, add this receiver to your configuration:"
echo ""
cat << 'EOF'
receivers:
- name: 'karo'
  webhook_configs:
  - url: 'http://karo-webhook.default.svc.cluster.local:9090/webhook'
    send_resolved: false
    http_config:
      timeout: 10s
EOF

echo ""
echo "You can now create AlertReaction resources. See examples/ directory for samples."
