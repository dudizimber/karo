package controllers

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	alertreactionv1alpha1 "github.com/dudizimber/k8s-alert-reaction-operator/api/v1alpha1"
)

func TestAlertMatching(t *testing.T) {
	reconciler, fakeClient := setupTestEmpty()

	// Create AlertReactions with different matchers
	tests := []struct {
		name           string
		alertReaction  *alertreactionv1alpha1.AlertReaction
		alertData      map[string]interface{}
		shouldMatch    bool
		description    string
	}{
		{
			name: "no-matchers-should-match",
			alertReaction: &alertreactionv1alpha1.AlertReaction{
				ObjectMeta: metav1.ObjectMeta{Name: "no-matchers", Namespace: "default"},
				Spec: alertreactionv1alpha1.AlertReactionSpec{
					AlertName: "TestAlert",
					Actions: []alertreactionv1alpha1.Action{
						{Name: "test-action", Image: "busybox:latest", Command: []string{"echo", "hello"}},
					},
				},
			},
			alertData: map[string]interface{}{
				"labels": map[string]interface{}{"severity": "critical"},
			},
			shouldMatch: true,
			description: "AlertReaction with no matchers should match when alertName matches",
		},
		{
			name: "equal-matcher-should-match",
			alertReaction: &alertreactionv1alpha1.AlertReaction{
				ObjectMeta: metav1.ObjectMeta{Name: "equal-matcher", Namespace: "default"},
				Spec: alertreactionv1alpha1.AlertReactionSpec{
					AlertName: "TestAlert",
					Matchers: []alertreactionv1alpha1.AlertMatcher{
						{Name: "severity", Operator: "=", Value: "critical"},
					},
					Actions: []alertreactionv1alpha1.Action{
						{Name: "test-action", Image: "busybox:latest", Command: []string{"echo", "hello"}},
					},
				},
			},
			alertData: map[string]interface{}{
				"labels": map[string]interface{}{"severity": "critical"},
			},
			shouldMatch: true,
			description: "AlertReaction with equal matcher should match when label value matches",
		},
		{
			name: "equal-matcher-should-not-match",
			alertReaction: &alertreactionv1alpha1.AlertReaction{
				ObjectMeta: metav1.ObjectMeta{Name: "equal-no-match", Namespace: "default"},
				Spec: alertreactionv1alpha1.AlertReactionSpec{
					AlertName: "TestAlert",
					Matchers: []alertreactionv1alpha1.AlertMatcher{
						{Name: "severity", Operator: "=", Value: "critical"},
					},
					Actions: []alertreactionv1alpha1.Action{
						{Name: "test-action", Image: "busybox:latest", Command: []string{"echo", "hello"}},
					},
				},
			},
			alertData: map[string]interface{}{
				"labels": map[string]interface{}{"severity": "warning"},
			},
			shouldMatch: false,
			description: "AlertReaction with equal matcher should not match when label value differs",
		},
		{
			name: "not-equal-matcher-should-match",
			alertReaction: &alertreactionv1alpha1.AlertReaction{
				ObjectMeta: metav1.ObjectMeta{Name: "not-equal-matcher", Namespace: "default"},
				Spec: alertreactionv1alpha1.AlertReactionSpec{
					AlertName: "TestAlert",
					Matchers: []alertreactionv1alpha1.AlertMatcher{
						{Name: "environment", Operator: "!=", Value: "test"},
					},
					Actions: []alertreactionv1alpha1.Action{
						{Name: "test-action", Image: "busybox:latest", Command: []string{"echo", "hello"}},
					},
				},
			},
			alertData: map[string]interface{}{
				"labels": map[string]interface{}{"environment": "production"},
			},
			shouldMatch: true,
			description: "AlertReaction with not-equal matcher should match when label value is different",
		},
		{
			name: "regex-matcher-should-match",
			alertReaction: &alertreactionv1alpha1.AlertReaction{
				ObjectMeta: metav1.ObjectMeta{Name: "regex-matcher", Namespace: "default"},
				Spec: alertreactionv1alpha1.AlertReactionSpec{
					AlertName: "TestAlert",
					Matchers: []alertreactionv1alpha1.AlertMatcher{
						{Name: "instance", Operator: "=~", Value: "prod-.*"},
					},
					Actions: []alertreactionv1alpha1.Action{
						{Name: "test-action", Image: "busybox:latest", Command: []string{"echo", "hello"}},
					},
				},
			},
			alertData: map[string]interface{}{
				"labels": map[string]interface{}{"instance": "prod-server-01"},
			},
			shouldMatch: true,
			description: "AlertReaction with regex matcher should match when pattern matches",
		},
		{
			name: "negative-regex-matcher-should-match",
			alertReaction: &alertreactionv1alpha1.AlertReaction{
				ObjectMeta: metav1.ObjectMeta{Name: "negative-regex-matcher", Namespace: "default"},
				Spec: alertreactionv1alpha1.AlertReactionSpec{
					AlertName: "TestAlert",
					Matchers: []alertreactionv1alpha1.AlertMatcher{
						{Name: "instance", Operator: "!~", Value: "test-.*"},
					},
					Actions: []alertreactionv1alpha1.Action{
						{Name: "test-action", Image: "busybox:latest", Command: []string{"echo", "hello"}},
					},
				},
			},
			alertData: map[string]interface{}{
				"labels": map[string]interface{}{"instance": "prod-server-01"},
			},
			shouldMatch: true,
			description: "AlertReaction with negative regex matcher should match when pattern doesn't match",
		},
		{
			name: "annotation-matcher-should-match",
			alertReaction: &alertreactionv1alpha1.AlertReaction{
				ObjectMeta: metav1.ObjectMeta{Name: "annotation-matcher", Namespace: "default"},
				Spec: alertreactionv1alpha1.AlertReactionSpec{
					AlertName: "TestAlert",
					Matchers: []alertreactionv1alpha1.AlertMatcher{
						{Name: "annotations.runbook", Operator: "=", Value: "emergency-runbook"},
					},
					Actions: []alertreactionv1alpha1.Action{
						{Name: "test-action", Image: "busybox:latest", Command: []string{"echo", "hello"}},
					},
				},
			},
			alertData: map[string]interface{}{
				"annotations": map[string]interface{}{"runbook": "emergency-runbook"},
			},
			shouldMatch: true,
			description: "AlertReaction should match annotation values correctly",
		},
		{
			name: "multiple-matchers-all-match",
			alertReaction: &alertreactionv1alpha1.AlertReaction{
				ObjectMeta: metav1.ObjectMeta{Name: "multiple-matchers-match", Namespace: "default"},
				Spec: alertreactionv1alpha1.AlertReactionSpec{
					AlertName: "TestAlert",
					Matchers: []alertreactionv1alpha1.AlertMatcher{
						{Name: "severity", Operator: "=", Value: "critical"},
						{Name: "environment", Operator: "!=", Value: "test"},
					},
					Actions: []alertreactionv1alpha1.Action{
						{Name: "test-action", Image: "busybox:latest", Command: []string{"echo", "hello"}},
					},
				},
			},
			alertData: map[string]interface{}{
				"labels": map[string]interface{}{
					"severity":    "critical",
					"environment": "production",
				},
			},
			shouldMatch: true,
			description: "AlertReaction with multiple matchers should match when all matchers match",
		},
		{
			name: "multiple-matchers-one-fails",
			alertReaction: &alertreactionv1alpha1.AlertReaction{
				ObjectMeta: metav1.ObjectMeta{Name: "multiple-matchers-fail", Namespace: "default"},
				Spec: alertreactionv1alpha1.AlertReactionSpec{
					AlertName: "TestAlert",
					Matchers: []alertreactionv1alpha1.AlertMatcher{
						{Name: "severity", Operator: "=", Value: "critical"},
						{Name: "environment", Operator: "=", Value: "production"},
					},
					Actions: []alertreactionv1alpha1.Action{
						{Name: "test-action", Image: "busybox:latest", Command: []string{"echo", "hello"}},
					},
				},
			},
			alertData: map[string]interface{}{
				"labels": map[string]interface{}{
					"severity":    "critical",
					"environment": "test",
				},
			},
			shouldMatch: false,
			description: "AlertReaction with multiple matchers should not match when one matcher fails",
		},
	}

	// Create all AlertReactions
	for _, test := range tests {
		err := fakeClient.Create(context.TODO(), test.alertReaction)
		if err != nil {
			t.Fatalf("Failed to create AlertReaction %s: %v", test.name, err)
		}
	}

	// Test each scenario
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			matches := reconciler.alertMatches(test.alertReaction, "TestAlert", test.alertData)
			if matches != test.shouldMatch {
				t.Errorf("%s: expected match=%v, got match=%v", test.description, test.shouldMatch, matches)
			}
		})
	}
}

func TestProcessMultipleAlertReactions(t *testing.T) {
	reconciler, fakeClient := setupTestEmpty()

	// Create multiple AlertReactions that should all match the same alert
	alertReaction1 := &alertreactionv1alpha1.AlertReaction{
		ObjectMeta: metav1.ObjectMeta{Name: "reaction-1", Namespace: "default"},
		Spec: alertreactionv1alpha1.AlertReactionSpec{
			AlertName: "ServiceDown",
			Matchers: []alertreactionv1alpha1.AlertMatcher{
				{Name: "severity", Operator: "=", Value: "critical"},
			},
			Actions: []alertreactionv1alpha1.Action{
				{Name: "notify-oncall", Image: "notification:latest", Command: []string{"notify"}},
			},
		},
	}

	alertReaction2 := &alertreactionv1alpha1.AlertReaction{
		ObjectMeta: metav1.ObjectMeta{Name: "reaction-2", Namespace: "default"},
		Spec: alertreactionv1alpha1.AlertReactionSpec{
			AlertName: "ServiceDown",
			Matchers: []alertreactionv1alpha1.AlertMatcher{
				{Name: "environment", Operator: "=", Value: "production"},
			},
			Actions: []alertreactionv1alpha1.Action{
				{Name: "restart-service", Image: "kubectl:latest", Command: []string{"kubectl", "rollout", "restart"}},
			},
		},
	}

	alertReaction3 := &alertreactionv1alpha1.AlertReaction{
		ObjectMeta: metav1.ObjectMeta{Name: "reaction-3", Namespace: "default"},
		Spec: alertreactionv1alpha1.AlertReactionSpec{
			AlertName: "ServiceDown",
			// No matchers - should always match when alertName matches
			Actions: []alertreactionv1alpha1.Action{
				{Name: "log-event", Image: "logger:latest", Command: []string{"log"}},
			},
		},
	}

	// Create all AlertReactions
	for _, ar := range []*alertreactionv1alpha1.AlertReaction{alertReaction1, alertReaction2, alertReaction3} {
		err := fakeClient.Create(context.TODO(), ar)
		if err != nil {
			t.Fatalf("Failed to create AlertReaction %s: %v", ar.Name, err)
		}
	}

	// Process alert that should match alertReaction1 and alertReaction3
	alertData := map[string]interface{}{
		"labels": map[string]interface{}{
			"severity":    "critical",
			"environment": "staging", // This won't match alertReaction2
		},
	}

	err := reconciler.ProcessAlert(context.TODO(), "ServiceDown", alertData)
	if err != nil {
		t.Fatalf("ProcessAlert returned error: %v", err)
	}

	// Check that alertReaction1 was triggered
	var updatedReaction1 alertreactionv1alpha1.AlertReaction
	err = fakeClient.Get(context.TODO(), types.NamespacedName{Name: "reaction-1", Namespace: "default"}, &updatedReaction1)
	if err != nil {
		t.Fatalf("Failed to get updated AlertReaction 1: %v", err)
	}

	if updatedReaction1.Status.TriggerCount != 1 {
		t.Errorf("Expected AlertReaction 1 trigger count 1, got %d", updatedReaction1.Status.TriggerCount)
	}

	// Check that alertReaction2 was NOT triggered (environment doesn't match)
	var updatedReaction2 alertreactionv1alpha1.AlertReaction
	err = fakeClient.Get(context.TODO(), types.NamespacedName{Name: "reaction-2", Namespace: "default"}, &updatedReaction2)
	if err != nil {
		t.Fatalf("Failed to get updated AlertReaction 2: %v", err)
	}

	if updatedReaction2.Status.TriggerCount != 0 {
		t.Errorf("Expected AlertReaction 2 trigger count 0, got %d", updatedReaction2.Status.TriggerCount)
	}

	// Check that alertReaction3 was triggered (no matchers, so always matches)
	var updatedReaction3 alertreactionv1alpha1.AlertReaction
	err = fakeClient.Get(context.TODO(), types.NamespacedName{Name: "reaction-3", Namespace: "default"}, &updatedReaction3)
	if err != nil {
		t.Fatalf("Failed to get updated AlertReaction 3: %v", err)
	}

	if updatedReaction3.Status.TriggerCount != 1 {
		t.Errorf("Expected AlertReaction 3 trigger count 1, got %d", updatedReaction3.Status.TriggerCount)
	}
}