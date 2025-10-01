package controllers

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	alertreactionv1alpha1 "github.com/dudizimber/k8s-alert-reaction-operator/api/v1alpha1"
)

func TestFakeClientBasic(t *testing.T) {
	reconciler, fakeClient, _ := setupTestEmpty()

	// Create a simple AlertReaction
	alertReaction := &alertreactionv1alpha1.AlertReaction{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-basic",
			Namespace: "default",
		},
		Spec: alertreactionv1alpha1.AlertReactionSpec{
			AlertName: "TestAlert",
			Actions: []alertreactionv1alpha1.Action{
				{
					Name:    "test-action",
					Image:   "busybox:latest",
					Command: []string{"echo", "test"},
				},
			},
		},
	}

	// Create the AlertReaction
	err := fakeClient.Create(context.TODO(), alertReaction)
	if err != nil {
		t.Fatalf("Failed to create AlertReaction: %v", err)
	}

	// Try to get it back using the same client
	var retrieved alertreactionv1alpha1.AlertReaction
	err = fakeClient.Get(context.TODO(), types.NamespacedName{Name: "test-basic", Namespace: "default"}, &retrieved)
	if err != nil {
		t.Fatalf("Failed to retrieve AlertReaction with fakeClient: %v", err)
	}

	// Try to get it using the reconciler's client
	var retrievedByReconciler alertreactionv1alpha1.AlertReaction
	err = reconciler.Get(context.TODO(), types.NamespacedName{Name: "test-basic", Namespace: "default"}, &retrievedByReconciler)
	if err != nil {
		t.Fatalf("Failed to retrieve AlertReaction with reconciler: %v", err)
	}

	t.Logf("Successfully created and retrieved AlertReaction: %s", retrieved.Name)
}
