#!/bin/bash

set -e

echo "Uninstalling Kubernetes Alert Reaction Operator..."

# Delete examples if they exist
echo "Removing example AlertReaction resources..."
kubectl delete -f examples/ --ignore-not-found=true

# Delete deployment
echo "Removing operator deployment..."
kubectl delete -f config/manager/deployment.yaml --ignore-not-found=true

# Delete RBAC
echo "Removing RBAC..."
kubectl delete -f config/rbac/rbac.yaml --ignore-not-found=true

# Delete CRD (this will also delete all AlertReaction resources)
echo "Removing Custom Resource Definition..."
kubectl delete -f config/crd/alertreaction.io_alertreactions.yaml --ignore-not-found=true

echo ""
echo "Alert Reaction Operator has been uninstalled successfully!"
