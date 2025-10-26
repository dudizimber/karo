package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	alertreactionv1alpha1 "github.com/dudizimber/karo/api/v1alpha1"
	"github.com/dudizimber/karo/controllers"
)

func setupWebhookTest() (*WebhookServer, *controllers.AlertReactionReconciler) {
	s := runtime.NewScheme()
	_ = scheme.AddToScheme(s)
	_ = alertreactionv1alpha1.AddToScheme(s)

	fakeClient := fake.NewClientBuilder().WithScheme(s).Build()

	controller := &controllers.AlertReactionReconciler{
		Client: fakeClient,
		Scheme: s,
	}

	webhookServer := NewWebhookServer(controller, "9090")

	return webhookServer, controller
}

func TestWebhookServer_NewWebhookServer(t *testing.T) {
	controller := &controllers.AlertReactionReconciler{}
	server := NewWebhookServer(controller, "8080")

	if server.controller != controller {
		t.Error("Controller not set correctly")
	}

	if server.port != "8080" {
		t.Errorf("Expected port 8080, got %s", server.port)
	}
}

func TestWebhookServer_HandleWebhook_ValidPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)

	webhookServer, controller := setupWebhookTest()

	// Create an AlertReaction in the fake client
	alertReaction := &alertreactionv1alpha1.AlertReaction{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-alert-reaction",
			Namespace: "default",
		},
		Spec: alertreactionv1alpha1.AlertReactionSpec{
			AlertName: "HighCPUUsage",
			Actions: []alertreactionv1alpha1.Action{
				{
					Name:    "scale-up",
					Image:   "busybox:latest",
					Command: []string{"echo", "scaling up"},
				},
			},
		},
	}

	err := controller.Create(context.TODO(), alertReaction)
	if err != nil {
		t.Fatalf("Failed to create AlertReaction: %v", err)
	}

	// Create test webhook payload
	webhook := AlertManagerWebhook{
		Version:  "4",
		GroupKey: "{}:{alertname=\"HighCPUUsage\"}",
		Status:   "firing",
		Receiver: "karo",
		GroupLabels: map[string]string{
			"alertname": "HighCPUUsage",
		},
		CommonLabels: map[string]string{
			"alertname": "HighCPUUsage",
			"instance":  "server1.example.com:9100",
			"job":       "node",
			"severity":  "warning",
		},
		CommonAnnotations: map[string]string{
			"description": "CPU usage is above 80% for more than 5 minutes",
			"summary":     "High CPU usage detected",
		},
		ExternalURL: "http://alertmanager.example.com:9093",
		Alerts: []Alert{
			{
				Status: "firing",
				Labels: map[string]string{
					"alertname": "HighCPUUsage",
					"instance":  "server1.example.com:9100",
					"job":       "node",
					"severity":  "warning",
				},
				Annotations: map[string]string{
					"description": "CPU usage is above 80% for more than 5 minutes",
					"summary":     "High CPU usage detected on server1",
				},
				StartsAt:     time.Now(),
				EndsAt:       time.Time{},
				GeneratorURL: "http://prometheus.example.com:9090/graph?g0.expr=...",
				Fingerprint:  "abc123",
			},
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(webhook)
	if err != nil {
		t.Fatalf("Failed to marshal webhook: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "/webhook", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create recorder and router
	w := httptest.NewRecorder()
	router := gin.New()
	router.POST("/webhook", webhookServer.handleWebhook)

	// Perform request
	router.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "Webhook processed successfully" {
		t.Errorf("Expected success message, got %s", response["message"])
	}
}

func TestWebhookServer_HandleWebhook_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	webhookServer, _ := setupWebhookTest()

	// Create invalid JSON
	invalidJSON := `{"invalid": json}`

	// Create HTTP request
	req, err := http.NewRequest("POST", "/webhook", bytes.NewBufferString(invalidJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create recorder and router
	w := httptest.NewRecorder()
	router := gin.New()
	router.POST("/webhook", webhookServer.handleWebhook)

	// Perform request
	router.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestWebhookServer_HandleWebhook_NonFiringAlert(t *testing.T) {
	gin.SetMode(gin.TestMode)

	webhookServer, _ := setupWebhookTest()

	// Create test webhook payload with resolved alert
	webhook := AlertManagerWebhook{
		Version:  "4",
		Status:   "resolved",
		Receiver: "karo",
		Alerts: []Alert{
			{
				Status: "resolved",
				Labels: map[string]string{
					"alertname": "HighCPUUsage",
				},
				StartsAt: time.Now().Add(-time.Hour),
				EndsAt:   time.Now(),
			},
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(webhook)
	if err != nil {
		t.Fatalf("Failed to marshal webhook: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "/webhook", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create recorder and router
	w := httptest.NewRecorder()
	router := gin.New()
	router.POST("/webhook", webhookServer.handleWebhook)

	// Perform request
	router.ServeHTTP(w, req)

	// Check response - should still be OK even if no actions taken
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestWebhookServer_AlertToMap(t *testing.T) {
	webhookServer, _ := setupWebhookTest()

	alert := Alert{
		Status: "firing",
		Labels: map[string]string{
			"alertname": "TestAlert",
			"instance":  "server1.example.com",
			"severity":  "warning",
		},
		Annotations: map[string]string{
			"summary":     "Test summary",
			"description": "Test description",
		},
		StartsAt:     time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		EndsAt:       time.Time{},
		GeneratorURL: "http://prometheus.example.com/graph",
		Fingerprint:  "abc123",
	}

	alertMap := webhookServer.alertToMap(alert)

	// Test top-level fields
	if alertMap["status"] != "firing" {
		t.Errorf("Expected status 'firing', got %v", alertMap["status"])
	}

	if alertMap["generatorURL"] != "http://prometheus.example.com/graph" {
		t.Errorf("Expected generatorURL 'http://prometheus.example.com/graph', got %v", alertMap["generatorURL"])
	}

	if alertMap["fingerprint"] != "abc123" {
		t.Errorf("Expected fingerprint 'abc123', got %v", alertMap["fingerprint"])
	}

	// Test labels
	labels, ok := alertMap["labels"].(map[string]interface{})
	if !ok {
		t.Error("Expected labels to be a map")
	} else {
		if labels["alertname"] != "TestAlert" {
			t.Errorf("Expected labels.alertname 'TestAlert', got %v", labels["alertname"])
		}
	}

	// Test direct label access
	if alertMap["labels.alertname"] != "TestAlert" {
		t.Errorf("Expected labels.alertname 'TestAlert', got %v", alertMap["labels.alertname"])
	}

	if alertMap["labels.instance"] != "server1.example.com" {
		t.Errorf("Expected labels.instance 'server1.example.com', got %v", alertMap["labels.instance"])
	}

	// Test annotations
	annotations, ok := alertMap["annotations"].(map[string]interface{})
	if !ok {
		t.Error("Expected annotations to be a map")
	} else {
		if annotations["summary"] != "Test summary" {
			t.Errorf("Expected annotations.summary 'Test summary', got %v", annotations["summary"])
		}
	}

	// Test direct annotation access
	if alertMap["annotations.summary"] != "Test summary" {
		t.Errorf("Expected annotations.summary 'Test summary', got %v", alertMap["annotations.summary"])
	}

	// Test time formatting
	if alertMap["startsAt"] != "2024-01-15T10:00:00Z" {
		t.Errorf("Expected startsAt '2024-01-15T10:00:00Z', got %v", alertMap["startsAt"])
	}
}

func TestWebhookServer_GetWebhookURL(t *testing.T) {
	webhookServer, _ := setupWebhookTest()

	tests := []struct {
		baseURL  string
		expected string
	}{
		{"", "http://localhost:9090/webhook"},
		{"http://example.com:8080", "http://example.com:8080/webhook"},
		{"https://alertmanager.example.com", "https://alertmanager.example.com/webhook"},
	}

	for _, test := range tests {
		t.Run(test.baseURL, func(t *testing.T) {
			url := webhookServer.GetWebhookURL(test.baseURL)
			if url != test.expected {
				t.Errorf("Expected URL %s, got %s", test.expected, url)
			}
		})
	}
}

func TestWebhookServer_GetWebhookConfig(t *testing.T) {
	webhookServer, _ := setupWebhookTest()

	config := webhookServer.GetWebhookConfig("http://example.com:9090")

	// Check that the config contains expected elements
	expectedStrings := []string{
		"karo",
		"http://example.com:9090/webhook",
		"send_resolved: false",
		"timeout: 10s",
	}

	for _, expected := range expectedStrings {
		if !contains(config, expected) {
			t.Errorf("Expected config to contain '%s'", expected)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			bytes.Contains([]byte(s), []byte(substr)))
}
