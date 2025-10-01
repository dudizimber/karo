package controllers

import (
	"context"
	"testing"

	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func TestReconcileNonExistent(t *testing.T) {
	reconciler, _ := setupTestEmpty()

	// Try to reconcile a non-existent AlertReaction
	req := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      "non-existent",
			Namespace: "default",
		},
	}

	result, err := reconciler.Reconcile(context.TODO(), req)
	if err != nil {
		t.Errorf("Reconcile should not fail for non-existent resource, got error: %v", err)
	}

	if result.RequeueAfter != 0 {
		t.Error("Should not requeue for non-existent resource")
	}

	t.Log("Reconcile handled non-existent resource correctly")
}
