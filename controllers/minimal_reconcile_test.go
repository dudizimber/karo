package controllers

import (
	"context"
	"testing"

	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	alertreactionv1alpha1 "github.com/dudizimber/karo/api/v1alpha1"
)

func TestMinimalReconcile(t *testing.T) {
	reconciler, _ := setupTestWithAlertReaction()

	req := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      "test-alert-reaction",
			Namespace: "default",
		},
	}

	ctx := context.TODO()

	// Minimal version of what Reconcile does
	var alertReaction alertreactionv1alpha1.AlertReaction
	if err := reconciler.Get(ctx, req.NamespacedName, &alertReaction); err != nil {
		t.Fatalf("Minimal reconcile Get failed: %v", err)
	}

	t.Logf("Minimal reconcile Get succeeded: %s", alertReaction.Name)
}
