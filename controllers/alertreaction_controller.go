package controllers

import (
	"context"
	"fmt"
	"strings"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	alertreactionv1alpha1 "github.com/dudizimber/k8s-alert-reaction-operator/api/v1alpha1"
)

// AlertReactionReconciler reconciles an AlertReaction object
type AlertReactionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=alertreaction.io,resources=alertreactions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=alertreaction.io,resources=alertreactions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=alertreaction.io,resources=alertreactions/finalizers,verbs=update
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch

// Reconcile handles AlertReaction resources
func (r *AlertReactionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the AlertReaction instance
	var alertReaction alertreactionv1alpha1.AlertReaction
	if err := r.Get(ctx, req.NamespacedName, &alertReaction); err != nil {
		logger.Error(err, "unable to fetch AlertReaction")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info("Reconciling AlertReaction", "alertName", alertReaction.Spec.AlertName)

	// Update status conditions
	now := metav1.NewTime(time.Now())
	condition := metav1.Condition{
		Type:               "Ready",
		Status:             metav1.ConditionTrue,
		Reason:             "AlertReactionReady",
		Message:            "AlertReaction is ready to process alerts",
		LastTransitionTime: now,
	}

	// Update conditions if changed
	updated := false
	if len(alertReaction.Status.Conditions) == 0 {
		alertReaction.Status.Conditions = []metav1.Condition{condition}
		updated = true
	} else {
		lastCondition := alertReaction.Status.Conditions[len(alertReaction.Status.Conditions)-1]
		if lastCondition.Type != condition.Type || lastCondition.Status != condition.Status {
			alertReaction.Status.Conditions = append(alertReaction.Status.Conditions, condition)
			updated = true
		}
	}

	if updated {
		if err := r.Status().Update(ctx, &alertReaction); err != nil {
			logger.Error(err, "unable to update AlertReaction status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// ProcessAlert creates jobs for the given alert
func (r *AlertReactionReconciler) ProcessAlert(ctx context.Context, alertName string, alertData map[string]interface{}) error {
	logger := log.FromContext(ctx)

	// Find AlertReaction for this alert
	var alertReactionList alertreactionv1alpha1.AlertReactionList
	if err := r.List(ctx, &alertReactionList); err != nil {
		return fmt.Errorf("failed to list AlertReactions: %w", err)
	}

	var targetAlertReaction *alertreactionv1alpha1.AlertReaction
	for i := range alertReactionList.Items {
		if alertReactionList.Items[i].Spec.AlertName == alertName {
			targetAlertReaction = &alertReactionList.Items[i]
			break
		}
	}

	if targetAlertReaction == nil {
		logger.Info("No AlertReaction found for alert", "alertName", alertName)
		return nil
	}

	logger.Info("Processing alert", "alertName", alertName, "actionsCount", len(targetAlertReaction.Spec.Actions))

	var jobRefs []alertreactionv1alpha1.JobReference
	now := metav1.NewTime(time.Now())

	// Create a job for each action
	for _, action := range targetAlertReaction.Spec.Actions {
		job, err := r.createJobFromAction(ctx, targetAlertReaction, action, alertData)
		if err != nil {
			logger.Error(err, "failed to create job for action", "actionName", action.Name)
			continue
		}

		if err := r.Create(ctx, job); err != nil {
			logger.Error(err, "failed to create job", "jobName", job.Name)
			continue
		}

		logger.Info("Created job for action", "jobName", job.Name, "actionName", action.Name)

		jobRefs = append(jobRefs, alertreactionv1alpha1.JobReference{
			Name:       job.Name,
			Namespace:  job.Namespace,
			ActionName: action.Name,
			CreatedAt:  now,
		})
	}

	// Update AlertReaction status
	targetAlertReaction.Status.LastTriggered = &now
	targetAlertReaction.Status.TriggerCount++
	targetAlertReaction.Status.LastJobsCreated = jobRefs

	if err := r.Status().Update(ctx, targetAlertReaction); err != nil {
		logger.Error(err, "failed to update AlertReaction status")
		return err
	}

	return nil
}

func (r *AlertReactionReconciler) createJobFromAction(ctx context.Context, alertReaction *alertreactionv1alpha1.AlertReaction, action alertreactionv1alpha1.Action, alertData map[string]interface{}) (*batchv1.Job, error) {
	// Generate job name, limited to 63 chars (Kubernetes DNS label limit)
	var jobName string
	baseName := fmt.Sprintf("%s-%s-%d", alertReaction.Name, action.Name, time.Now().Unix())
	if len(baseName) > 63 {
		jobName = baseName[:63]
	} else {
		jobName = baseName
	}

	// Process environment variables
	env, err := r.processEnvVars(ctx, alertReaction.Namespace, action.Env, alertData)
	if err != nil {
		return nil, fmt.Errorf("failed to process environment variables: %w", err)
	}

	// Convert resource requirements
	var resources corev1.ResourceRequirements
	if action.Resources != nil {
		resources = corev1.ResourceRequirements{
			Limits:   make(corev1.ResourceList),
			Requests: make(corev1.ResourceList),
		}
		for k, v := range action.Resources.Limits {
			resources.Limits[corev1.ResourceName(k)] = parseQuantity(v)
		}
		for k, v := range action.Resources.Requests {
			resources.Requests[corev1.ResourceName(k)] = parseQuantity(v)
		}
	}

	// Convert volumes
	volumes, err := r.convertVolumes(alertReaction.Spec.Volumes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert volumes: %w", err)
	}

	// Convert volume mounts
	volumeMounts, err := r.convertVolumeMounts(action.VolumeMounts)
	if err != nil {
		return nil, fmt.Errorf("failed to convert volume mounts: %w", err)
	}

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: alertReaction.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":      "alert-reaction-job",
				"app.kubernetes.io/component": "job",
				"alert-reaction/alert-name":   alertReaction.Spec.AlertName,
				"alert-reaction/action-name":  action.Name,
				"alert-reaction/owner":        alertReaction.Name,
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: alertReaction.APIVersion,
					Kind:       alertReaction.Kind,
					Name:       alertReaction.Name,
					UID:        alertReaction.UID,
				},
			},
		},
		Spec: batchv1.JobSpec{
			TTLSecondsAfterFinished: int32Ptr(300), // Clean up jobs after 5 minutes
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy:      corev1.RestartPolicyNever,
					ServiceAccountName: action.ServiceAccount,
					Volumes:            volumes,
					Containers: []corev1.Container{
						{
							Name:         "action",
							Image:        action.Image,
							Command:      action.Command,
							Args:         action.Args,
							Env:          env,
							Resources:    resources,
							VolumeMounts: volumeMounts,
						},
					},
				},
			},
		},
	}

	return job, nil
}

func (r *AlertReactionReconciler) processEnvVars(ctx context.Context, namespace string, envVars []alertreactionv1alpha1.EnvVar, alertData map[string]interface{}) ([]corev1.EnvVar, error) {
	var result []corev1.EnvVar

	for _, envVar := range envVars {
		var value string
		var err error

		if envVar.Value != "" {
			value = envVar.Value
		} else if envVar.ValueFrom != nil {
			value, err = r.resolveEnvVarSource(ctx, namespace, envVar.ValueFrom, alertData)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve env var %s: %w", envVar.Name, err)
			}
		}

		result = append(result, corev1.EnvVar{
			Name:  envVar.Name,
			Value: value,
		})
	}

	return result, nil
}

func (r *AlertReactionReconciler) resolveEnvVarSource(ctx context.Context, namespace string, source *alertreactionv1alpha1.EnvVarSource, alertData map[string]interface{}) (string, error) {
	if source.AlertRef != nil {
		return r.getAlertFieldValue(alertData, source.AlertRef.FieldPath)
	}

	if source.ConfigMapKeyRef != nil {
		var cm corev1.ConfigMap
		key := types.NamespacedName{
			Name:      source.ConfigMapKeyRef.Name,
			Namespace: namespace,
		}
		if err := r.Get(ctx, key, &cm); err != nil {
			if source.ConfigMapKeyRef.Optional != nil && *source.ConfigMapKeyRef.Optional {
				return "", nil
			}
			return "", err
		}
		value, exists := cm.Data[source.ConfigMapKeyRef.Key]
		if !exists && (source.ConfigMapKeyRef.Optional == nil || !*source.ConfigMapKeyRef.Optional) {
			return "", fmt.Errorf("key %s not found in ConfigMap %s", source.ConfigMapKeyRef.Key, source.ConfigMapKeyRef.Name)
		}
		return value, nil
	}

	if source.SecretKeyRef != nil {
		var secret corev1.Secret
		key := types.NamespacedName{
			Name:      source.SecretKeyRef.Name,
			Namespace: namespace,
		}
		if err := r.Get(ctx, key, &secret); err != nil {
			if source.SecretKeyRef.Optional != nil && *source.SecretKeyRef.Optional {
				return "", nil
			}
			return "", err
		}
		value, exists := secret.Data[source.SecretKeyRef.Key]
		if !exists && (source.SecretKeyRef.Optional == nil || !*source.SecretKeyRef.Optional) {
			return "", fmt.Errorf("key %s not found in Secret %s", source.SecretKeyRef.Key, source.SecretKeyRef.Name)
		}
		return string(value), nil
	}

	return "", fmt.Errorf("no valid source specified")
}

func (r *AlertReactionReconciler) getAlertFieldValue(alertData map[string]interface{}, fieldPath string) (string, error) {
	// Simple field path resolution (can be enhanced for nested paths)
	// Supports: "labels.instance", "annotations.summary", "status", etc.

	value, exists := alertData[fieldPath]
	if !exists {
		// Try nested path resolution
		return r.getNestedField(alertData, fieldPath)
	}

	if str, ok := value.(string); ok {
		return str, nil
	}

	return fmt.Sprintf("%v", value), nil
}

func (r *AlertReactionReconciler) getNestedField(data map[string]interface{}, path string) (string, error) {
	// Basic implementation for dot-separated paths like "labels.instance"
	// This can be enhanced with more sophisticated path parsing

	parts := strings.Split(path, ".")

	current_data := data
	for i, part := range parts {
		if i == len(parts)-1 {
			// Last part - get the value
			if value, exists := current_data[part]; exists {
				if str, ok := value.(string); ok {
					return str, nil
				}
				return fmt.Sprintf("%v", value), nil
			}
			return "", fmt.Errorf("field %s not found", path)
		}

		// Navigate deeper
		if next, exists := current_data[part]; exists {
			if nextMap, ok := next.(map[string]interface{}); ok {
				current_data = nextMap
			} else {
				return "", fmt.Errorf("field %s is not a map", part)
			}
		} else {
			return "", fmt.Errorf("field %s not found", part)
		}
	}

	return "", fmt.Errorf("field %s not found", path)
}

// convertVolumes converts AlertReaction volumes to Kubernetes volumes
func (r *AlertReactionReconciler) convertVolumes(volumes []alertreactionv1alpha1.Volume) ([]corev1.Volume, error) {
	var result []corev1.Volume

	for _, vol := range volumes {
		k8sVol := corev1.Volume{
			Name: vol.Name,
		}

		// Convert volume source
		if vol.ConfigMap != nil {
			k8sVol.ConfigMap = &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: vol.ConfigMap.Name,
				},
				DefaultMode: vol.ConfigMap.DefaultMode,
				Optional:    vol.ConfigMap.Optional,
			}
		} else if vol.Secret != nil {
			k8sVol.Secret = &corev1.SecretVolumeSource{
				SecretName:  vol.Secret.SecretName,
				DefaultMode: vol.Secret.DefaultMode,
				Optional:    vol.Secret.Optional,
			}
		} else if vol.EmptyDir != nil {
			emptyDir := &corev1.EmptyDirVolumeSource{}
			if vol.EmptyDir.Medium != "" {
				emptyDir.Medium = corev1.StorageMedium(vol.EmptyDir.Medium)
			}
			if vol.EmptyDir.SizeLimit != "" {
				if sizeLimit, err := resource.ParseQuantity(vol.EmptyDir.SizeLimit); err == nil {
					emptyDir.SizeLimit = &sizeLimit
				}
			}
			k8sVol.EmptyDir = emptyDir
		} else if vol.PersistentVolumeClaim != nil {
			k8sVol.PersistentVolumeClaim = &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: vol.PersistentVolumeClaim.ClaimName,
				ReadOnly:  vol.PersistentVolumeClaim.ReadOnly,
			}
		} else if vol.HostPath != nil {
			hostPath := &corev1.HostPathVolumeSource{
				Path: vol.HostPath.Path,
			}
			if vol.HostPath.Type != "" {
				hostPathType := corev1.HostPathType(vol.HostPath.Type)
				hostPath.Type = &hostPathType
			}
			k8sVol.HostPath = hostPath
		} else {
			return nil, fmt.Errorf("volume %s has no valid source defined", vol.Name)
		}

		result = append(result, k8sVol)
	}

	return result, nil
}

// convertVolumeMounts converts AlertReaction volume mounts to Kubernetes volume mounts
func (r *AlertReactionReconciler) convertVolumeMounts(volumeMounts []alertreactionv1alpha1.VolumeMount) ([]corev1.VolumeMount, error) {
	var result []corev1.VolumeMount

	for _, vm := range volumeMounts {
		k8sVM := corev1.VolumeMount{
			Name:      vm.Name,
			MountPath: vm.MountPath,
			SubPath:   vm.SubPath,
			ReadOnly:  vm.ReadOnly,
		}
		result = append(result, k8sVM)
	}

	return result, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AlertReactionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&alertreactionv1alpha1.AlertReaction{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}

// Helper functions
func int32Ptr(i int32) *int32 {
	return &i
}

func parseQuantity(s string) resource.Quantity {
	// Simple implementation - in production, use resource.ParseQuantity
	return resource.MustParse(s)
}
