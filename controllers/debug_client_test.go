package controllers

import (
	"context"
	"testing"

	"k8s.io/apimachinery/pkg/types"

	alertreactionv1alpha1 "github.com/dudizimber/k8s-alert-reaction-operator/api/v1alpha1"
)

func TestReconcilerClientDebug(t *testing.T) {
	reconciler, fakeClient, _ := setupTestWithAlertReaction()

	// Verify the AlertReaction exists using fakeClient
	var ar1 alertreactionv1alpha1.AlertReaction
	err := fakeClient.Get(context.TODO(), types.NamespacedName{Name: "test-alert-reaction", Namespace: "default"}, &ar1)
	if err != nil {
		t.Fatalf("fakeClient can't find AlertReaction: %v", err)
	}
	t.Logf("fakeClient found AlertReaction: %s", ar1.Name)

	// Verify the AlertReaction exists using reconciler's client
	var ar2 alertreactionv1alpha1.AlertReaction
	err = reconciler.Get(context.TODO(), types.NamespacedName{Name: "test-alert-reaction", Namespace: "default"}, &ar2)
	if err != nil {
		t.Fatalf("reconciler.Client can't find AlertReaction: %v", err)
	}
	t.Logf("reconciler.Client found AlertReaction: %s", ar2.Name)

	// Check if they're the same client
	t.Logf("Clients are same instance: %v", reconciler.Client == fakeClient)
}
