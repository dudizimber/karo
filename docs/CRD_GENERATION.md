# CRD Generation Guide

This document explains how to automatically generate CustomResourceDefinitions (CRDs) from Go types in this project.

## Overview

The AlertReaction operator uses Kubebuilder markers in Go code to automatically generate Kubernetes manifests, including CRDs. The process uses `controller-gen` from the `sigs.k8s.io/controller-tools` project.

## Available Make Targets

### `make manifests`
Generates CRDs, RBAC, and webhook configurations from Go code annotations.

```bash
make manifests
```

This command:
- Scans Go files for kubebuilder markers
- Generates CRD YAML files in `config/crd/`
- Generates RBAC YAML files in `config/rbac/`

### `make generate` 
Generates Go code for DeepCopy methods required by Kubernetes APIs.

```bash
make generate
```

### `make controller-gen`
Downloads and installs `controller-gen` tool locally if not present.

```bash
make controller-gen
```

## Kubebuilder Markers

The Go structs use special comments (markers) that `controller-gen` reads to generate CRDs:

### Package Level Markers
```go
// +kubebuilder:object:generate=true
// +groupName=karo.io
package v1
```

### Struct Level Markers
```go
// AlertReaction is the Schema for the alertreactions API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type AlertReaction struct {
    // ...
}
```

### Field Level Markers
```go
type AlertReactionSpec struct {
    // AlertName specifies the Prometheus alert name to react to
    // +kubebuilder:validation:Required
    AlertName string `json:"alertName"`
    
    // Actions defines the list of actions to perform
    // +kubebuilder:validation:Required
    // +kubebuilder:validation:MinItems=1
    Actions []Action `json:"actions"`
}
```

## Common Validation Markers

- `+kubebuilder:validation:Required` - Field is required
- `+kubebuilder:validation:MinItems=1` - Array must have at least 1 item
- `+kubebuilder:validation:MaxItems=10` - Array can have at most 10 items
- `+kubebuilder:validation:MinLength=1` - String minimum length
- `+kubebuilder:validation:MaxLength=100` - String maximum length
- `+kubebuilder:validation:Pattern="^[a-z]+$"` - Regex validation
- `+kubebuilder:validation:Enum=value1;value2` - Enum validation

## Printer Columns

Add custom columns to `kubectl get` output:

```go
// +kubebuilder:printcolumn:name="Alert Name",type=string,JSONPath=`.spec.alertName`
// +kubebuilder:printcolumn:name="Actions",type=integer,JSONPath=`.spec.actions[*].name | length`
type AlertReaction struct {
    // ...
}
```

## Manual CRD Updates

If `controller-gen` fails (due to parsing issues), you can manually update the CRD:

1. Edit `config/crd/karo.io_alertreactions.yaml`
2. Add missing fields to the appropriate `properties` section
3. Follow the OpenAPI v3 schema format
4. Test with `kubectl apply --dry-run=client -f config/crd/`

## Troubleshooting

### controller-gen Panics
This usually indicates a parsing issue in the Go code or version compatibility:
1. Run `go build ./api/...` to check for compilation errors
2. Check for syntax errors in struct tags
3. Ensure all imports are correct
4. **Verify kubebuilder markers are present on structs**:
   ```go
   // +kubebuilder:object:root=true
   // +kubebuilder:subresource:status
   type AlertReaction struct { ... }
   ```
5. Try updating to a newer version of controller-gen (we use v0.14.0)
6. Remove old controller-gen and reinstall: `rm bin/controller-gen && make controller-gen`

### Missing Fields in CRD
1. Verify the field has appropriate JSON tags: `json:"fieldName"`
2. Check that the struct is exported (starts with capital letter)
3. Run `make manifests` after making changes
4. If automatic generation fails, manually add to the CRD file

### CRD Validation Errors
1. Use `kubectl apply --dry-run=client -f config/crd/` to validate
2. Check the OpenAPI schema format in the generated CRD
3. Ensure required fields are marked correctly

## Best Practices

1. **Always run `make manifests` after changing API types**
2. **Test CRD changes with dry-run before applying**
3. **Use appropriate validation markers for better UX**
4. **Keep backwards compatibility when modifying existing fields**
5. **Add printer columns for important fields**

## Example Workflow

```bash
# 1. Modify Go structs in api/v1/
vim api/v1/alertreaction_types.go

# 2. Generate updated manifests
make manifests

# 3. Test the CRD
kubectl apply --dry-run=client -f config/crd/

# 4. Test with example resources
kubectl apply --dry-run=client -f examples/

# 5. Run tests
make test

# 6. Apply to cluster
kubectl apply -f config/crd/
```