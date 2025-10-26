#!/bin/bash

# Test script for the Karo webhook endpoint
# This script simulates AlertManager webhook calls for testing

WEBHOOK_URL=${1:-"http://localhost:9090/webhook"}

echo "Testing Karo webhook at: $WEBHOOK_URL"

# Test 1: High CPU Usage Alert
echo ""
echo "=== Test 1: High CPU Usage Alert ==="
curl -X POST "$WEBHOOK_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "version": "4",
    "groupKey": "{}:{alertname=\"HighCPUUsage\"}",
    "truncatedAlerts": 0,
    "status": "firing",
    "receiver": "karo",
    "groupLabels": {
      "alertname": "HighCPUUsage"
    },
    "commonLabels": {
      "alertname": "HighCPUUsage",
      "instance": "server1.example.com:9100",
      "job": "node",
      "severity": "warning"
    },
    "commonAnnotations": {
      "description": "CPU usage is above 80% for more than 5 minutes",
      "summary": "High CPU usage detected"
    },
    "externalURL": "http://alertmanager.example.com:9093",
    "alerts": [
      {
        "status": "firing",
        "labels": {
          "alertname": "HighCPUUsage",
          "instance": "server1.example.com:9100",
          "job": "node",
          "severity": "warning",
          "deployment": "my-app"
        },
        "annotations": {
          "description": "CPU usage is above 80% for more than 5 minutes",
          "summary": "High CPU usage detected on server1"
        },
        "startsAt": "2024-01-15T10:00:00.000Z",
        "endsAt": "0001-01-01T00:00:00Z",
        "generatorURL": "http://prometheus.example.com:9090/graph?g0.expr=...",
        "fingerprint": "abc123"
      }
    ]
  }'

echo ""

# Test 2: Low Disk Space Alert
echo ""
echo "=== Test 2: Low Disk Space Alert ==="
curl -X POST "$WEBHOOK_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "version": "4",
    "groupKey": "{}:{alertname=\"LowDiskSpace\"}",
    "truncatedAlerts": 0,
    "status": "firing",
    "receiver": "karo",
    "groupLabels": {
      "alertname": "LowDiskSpace"
    },
    "commonLabels": {
      "alertname": "LowDiskSpace",
      "instance": "server2.example.com:9100",
      "job": "node",
      "severity": "critical"
    },
    "commonAnnotations": {
      "description": "Disk usage is above 90%",
      "summary": "Low disk space detected"
    },
    "externalURL": "http://alertmanager.example.com:9093",
    "alerts": [
      {
        "status": "firing",
        "labels": {
          "alertname": "LowDiskSpace",
          "instance": "server2.example.com:9100",
          "job": "node",
          "severity": "critical",
          "value": "95"
        },
        "annotations": {
          "description": "Disk usage is above 90%",
          "summary": "Only 5% disk space remaining on server2"
        },
        "startsAt": "2024-01-15T10:30:00.000Z",
        "endsAt": "0001-01-01T00:00:00Z",
        "generatorURL": "http://prometheus.example.com:9090/graph?g0.expr=...",
        "fingerprint": "def456"
      }
    ]
  }'

echo ""

# Test 3: Pod Crash Looping Alert
echo ""
echo "=== Test 3: Pod Crash Looping Alert ==="
curl -X POST "$WEBHOOK_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "version": "4",
    "groupKey": "{}:{alertname=\"PodCrashLooping\"}",
    "truncatedAlerts": 0,
    "status": "firing",
    "receiver": "karo",
    "groupLabels": {
      "alertname": "PodCrashLooping"
    },
    "commonLabels": {
      "alertname": "PodCrashLooping",
      "namespace": "default",
      "pod": "my-app-deployment-abc123-xyz789",
      "deployment": "my-app-deployment",
      "severity": "critical"
    },
    "commonAnnotations": {
      "description": "Pod has been restarting multiple times",
      "summary": "Pod crash loop detected"
    },
    "externalURL": "http://alertmanager.example.com:9093",
    "alerts": [
      {
        "status": "firing",
        "labels": {
          "alertname": "PodCrashLooping",
          "namespace": "default",
          "pod": "my-app-deployment-abc123-xyz789",
          "deployment": "my-app-deployment",
          "severity": "critical"
        },
        "annotations": {
          "description": "Pod has been restarting multiple times",
          "summary": "Pod my-app-deployment-abc123-xyz789 is crash looping"
        },
        "startsAt": "2024-01-15T11:00:00.000Z",
        "endsAt": "0001-01-01T00:00:00Z",
        "generatorURL": "http://prometheus.example.com:9090/graph?g0.expr=...",
        "fingerprint": "ghi789"
      }
    ]
  }'

echo ""

# Test health endpoint
echo ""
echo "=== Health Check ==="
curl -X GET "${WEBHOOK_URL%/webhook}/health"

echo ""
echo ""
echo "Test completed! Check the operator logs to see if alerts were processed correctly."
echo "You can also check for created jobs with: kubectl get jobs -l app.kubernetes.io/name=karo-job"
