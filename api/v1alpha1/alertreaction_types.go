// Package v1alpha1 contains API Schema definitions for the alertreaction v1alpha1 API group
// +kubebuilder:object:generate=true
// +groupName=alertreaction.io
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "alertreaction.io", Version: "v1alpha1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

// addKnownTypes adds the set of types defined in this package to the supplied scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(GroupVersion,
		&AlertReaction{},
		&AlertReactionList{},
	)
	metav1.AddToGroupVersion(scheme, GroupVersion)
	return nil
}

// AlertMatcher defines conditions for matching alerts based on their attributes
type AlertMatcher struct {
	// Name of the alert attribute to match against (e.g., "labels.instance", "annotations.severity", "status")
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Operator defines how to compare the alert attribute value with the expected value
	// Supported operators: Equal, NotEqual, In, NotIn, Exists, DoesNotExist, GreaterThan, LessThan, Regex, NotRegex
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=Equal;NotEqual;In;NotIn;Exists;DoesNotExist;GreaterThan;LessThan;Regex;NotRegex
	Operator MatchOperator `json:"operator"`

	// Values are the expected values to match against (not required for Exists/DoesNotExist operators)
	// For "In" and "NotIn" operators, multiple values can be specified
	// For other operators, only the first value is used
	Values []string `json:"values,omitempty"`
}

// MatchOperator defines the type of matching operation to perform
// +kubebuilder:validation:Enum=Equal;NotEqual;In;NotIn;Exists;DoesNotExist;GreaterThan;LessThan;Regex;NotRegex
type MatchOperator string

const (
	// MatchOperatorEqual checks if the alert attribute equals the specified value
	MatchOperatorEqual MatchOperator = "Equal"
	// MatchOperatorNotEqual checks if the alert attribute does not equal the specified value
	MatchOperatorNotEqual MatchOperator = "NotEqual"
	// MatchOperatorIn checks if the alert attribute value is in the list of specified values
	MatchOperatorIn MatchOperator = "In"
	// MatchOperatorNotIn checks if the alert attribute value is not in the list of specified values
	MatchOperatorNotIn MatchOperator = "NotIn"
	// MatchOperatorExists checks if the alert attribute exists (regardless of value)
	MatchOperatorExists MatchOperator = "Exists"
	// MatchOperatorDoesNotExist checks if the alert attribute does not exist
	MatchOperatorDoesNotExist MatchOperator = "DoesNotExist"
	// MatchOperatorGreaterThan checks if the alert attribute value is greater than the specified value (numeric comparison)
	MatchOperatorGreaterThan MatchOperator = "GreaterThan"
	// MatchOperatorLessThan checks if the alert attribute value is less than the specified value (numeric comparison)
	MatchOperatorLessThan MatchOperator = "LessThan"
	// MatchOperatorRegex checks if the alert attribute value matches the specified regular expression
	MatchOperatorRegex MatchOperator = "Regex"
	// MatchOperatorNotRegex checks if the alert attribute value does not match the specified regular expression
	MatchOperatorNotRegex MatchOperator = "NotRegex"
)

// AlertReactionSpec defines the desired state of AlertReaction
type AlertReactionSpec struct {
	// AlertName specifies the Prometheus alert name to react to
	// +kubebuilder:validation:Required
	AlertName string `json:"alertName"`

	// Matchers defines additional conditions that must be met for the alert to trigger this reaction
	// All matchers must match for the reaction to be triggered
	// If no matchers are specified, only the AlertName is used for matching
	Matchers []AlertMatcher `json:"matchers,omitempty"`

	// Actions defines the list of actions to perform when the alert is received
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Actions []Action `json:"actions"`

	// Volumes defines volumes that can be mounted by actions in this AlertReaction
	// These volumes will be available to all jobs created by this AlertReaction
	Volumes []Volume `json:"volumes,omitempty"`
}

// Action defines a single action to perform when an alert is received
type Action struct {
	// Name of the action
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Image to use for the job
	// +kubebuilder:validation:Required
	Image string `json:"image"`

	// Command to execute in the container
	// +kubebuilder:validation:Required
	Command []string `json:"command"`

	// Args for the command (optional)
	Args []string `json:"args,omitempty"`

	// Environment variables for the job (optional)
	Env []EnvVar `json:"env,omitempty"`

	// Resources for the job (optional)
	Resources *ResourceRequirements `json:"resources,omitempty"`

	// VolumeMounts specifies the volumes to mount into this action's container
	// The volumes must be defined in the AlertReaction's spec.volumes
	VolumeMounts []VolumeMount `json:"volumeMounts,omitempty"`

	// ServiceAccount specifies the service account to use for the job created by this action
	// If not specified, the default service account will be used
	ServiceAccount string `json:"serviceAccount,omitempty"`
}

// EnvVar represents an environment variable present in a Container.
type EnvVar struct {
	// Name of the environment variable
	Name string `json:"name"`

	// Value of the environment variable
	Value string `json:"value,omitempty"`

	// Source for the environment variable's value
	ValueFrom *EnvVarSource `json:"valueFrom,omitempty"`
}

// EnvVarSource represents a source for the value of an EnvVar.
type EnvVarSource struct {
	// Selects a field of the alert
	AlertRef *AlertFieldSelector `json:"alertRef,omitempty"`

	// Selects a key of a ConfigMap
	ConfigMapKeyRef *ConfigMapKeySelector `json:"configMapKeyRef,omitempty"`

	// Selects a key of a secret in the pod's namespace
	SecretKeyRef *SecretKeySelector `json:"secretKeyRef,omitempty"`
}

// AlertFieldSelector selects a field from the alert
type AlertFieldSelector struct {
	// Path to the field in the alert (e.g., "labels.instance", "annotations.summary")
	FieldPath string `json:"fieldPath"`
}

// ConfigMapKeySelector selects a key from a ConfigMap
type ConfigMapKeySelector struct {
	// Name of the ConfigMap
	Name string `json:"name"`

	// Key to select from the ConfigMap
	Key string `json:"key"`

	// Specify whether the ConfigMap or its key must be defined
	Optional *bool `json:"optional,omitempty"`
}

// SecretKeySelector selects a key from a Secret
type SecretKeySelector struct {
	// Name of the Secret
	Name string `json:"name"`

	// Key to select from the Secret
	Key string `json:"key"`

	// Specify whether the Secret or its key must be defined
	Optional *bool `json:"optional,omitempty"`
}

// ResourceRequirements describes the compute resource requirements.
type ResourceRequirements struct {
	// Limits describes the maximum amount of compute resources allowed
	Limits map[string]string `json:"limits,omitempty"`

	// Requests describes the minimum amount of compute resources required
	Requests map[string]string `json:"requests,omitempty"`
}

// Volume represents a named volume in a job
type Volume struct {
	// Name of the volume
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// VolumeSource represents the location and type of the mounted volume
	VolumeSource `json:",inline"`
}

// VolumeSource represents the source of a volume to mount
type VolumeSource struct {
	// ConfigMap represents a configMap that should populate this volume
	ConfigMap *ConfigMapVolumeSource `json:"configMap,omitempty"`

	// Secret represents a secret that should populate this volume
	Secret *SecretVolumeSource `json:"secret,omitempty"`

	// EmptyDir represents a temporary directory that shares a job's lifetime
	EmptyDir *EmptyDirVolumeSource `json:"emptyDir,omitempty"`

	// PersistentVolumeClaim represents a PVC that should be mounted
	PersistentVolumeClaim *PersistentVolumeClaimVolumeSource `json:"persistentVolumeClaim,omitempty"`

	// HostPath represents a pre-existing file or directory on the host machine
	HostPath *HostPathVolumeSource `json:"hostPath,omitempty"`
}

// ConfigMapVolumeSource adapts a ConfigMap into a volume
type ConfigMapVolumeSource struct {
	// Name of the ConfigMap
	Name string `json:"name"`

	// Optional: mode bits to use on created files by default
	DefaultMode *int32 `json:"defaultMode,omitempty"`

	// Optional: specify whether the ConfigMap or its keys must be defined
	Optional *bool `json:"optional,omitempty"`
}

// SecretVolumeSource adapts a Secret into a volume
type SecretVolumeSource struct {
	// Name of the Secret
	SecretName string `json:"secretName"`

	// Optional: mode bits to use on created files by default
	DefaultMode *int32 `json:"defaultMode,omitempty"`

	// Optional: specify whether the Secret or its keys must be defined
	Optional *bool `json:"optional,omitempty"`
}

// EmptyDirVolumeSource represents a temporary directory
type EmptyDirVolumeSource struct {
	// What type of storage medium should back this directory
	Medium string `json:"medium,omitempty"`

	// Total amount of local storage required for this EmptyDir volume
	SizeLimit string `json:"sizeLimit,omitempty"`
}

// PersistentVolumeClaimVolumeSource references a PVC in the same namespace
type PersistentVolumeClaimVolumeSource struct {
	// ClaimName is the name of a PersistentVolumeClaim in the same namespace
	ClaimName string `json:"claimName"`

	// ReadOnly will force the ReadOnly setting in VolumeMounts
	ReadOnly bool `json:"readOnly,omitempty"`
}

// HostPathVolumeSource represents a host path mapped into a job
type HostPathVolumeSource struct {
	// Path of the directory on the host
	Path string `json:"path"`

	// Type for HostPath Volume
	Type string `json:"type,omitempty"`
}

// VolumeMount describes a mounting of a Volume within a container
type VolumeMount struct {
	// Name must match the name of a volume defined in spec.volumes
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Path within the container at which the volume should be mounted
	// +kubebuilder:validation:Required
	MountPath string `json:"mountPath"`

	// SubPath within the volume from which the container's volume should be mounted
	SubPath string `json:"subPath,omitempty"`

	// Mounted read-only if true, read-write otherwise (false or unspecified)
	ReadOnly bool `json:"readOnly,omitempty"`
}

// AlertReactionStatus defines the observed state of AlertReaction
type AlertReactionStatus struct {
	// LastTriggered indicates when this AlertReaction was last triggered
	LastTriggered *metav1.Time `json:"lastTriggered,omitempty"`

	// TriggerCount indicates how many times this AlertReaction has been triggered
	TriggerCount int64 `json:"triggerCount,omitempty"`

	// LastJobsCreated contains references to the last batch of jobs created
	LastJobsCreated []JobReference `json:"lastJobsCreated,omitempty"`

	// Conditions represent the latest available observations of the AlertReaction's state
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// JobReference contains a reference to a created job
type JobReference struct {
	// Name of the job
	Name string `json:"name"`

	// Namespace of the job
	Namespace string `json:"namespace"`

	// ActionName that created this job
	ActionName string `json:"actionName"`

	// CreatedAt timestamp
	CreatedAt metav1.Time `json:"createdAt"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Alert Name",type=string,JSONPath=`.spec.alertName`
// +kubebuilder:printcolumn:name="Actions",type=integer,JSONPath=`.spec.actions[*].name | length`
// +kubebuilder:printcolumn:name="Last Triggered",type=date,JSONPath=`.status.lastTriggered`
// +kubebuilder:printcolumn:name="Trigger Count",type=integer,JSONPath=`.status.triggerCount`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// AlertReaction is the Schema for the alertreactions API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type AlertReaction struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AlertReactionSpec   `json:"spec,omitempty"`
	Status AlertReactionStatus `json:"status,omitempty"`
}

// AlertReactionList contains a list of AlertReaction
// +kubebuilder:object:root=true
type AlertReactionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AlertReaction `json:"items"`
}
