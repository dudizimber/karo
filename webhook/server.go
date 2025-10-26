package webhook

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/dudizimber/karo/controllers"
)

// AlertManagerWebhook represents an AlertManager webhook payload
type AlertManagerWebhook struct {
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
	TruncatedAlerts   int               `json:"truncatedAlerts"`
	Status            string            `json:"status"`
	Receiver          string            `json:"receiver"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Alerts            []Alert           `json:"alerts"`
}

// Alert represents a single alert in the webhook payload
type Alert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt"`
	EndsAt       time.Time         `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Fingerprint  string            `json:"fingerprint"`
}

// WebhookServer handles incoming webhook requests from AlertManager
type WebhookServer struct {
	controller *controllers.AlertReactionReconciler
	port       string
}

// NewWebhookServer creates a new webhook server
func NewWebhookServer(controller *controllers.AlertReactionReconciler, port string) *WebhookServer {
	return &WebhookServer{
		controller: controller,
		port:       port,
	}
}

// Start starts the webhook server
func (ws *WebhookServer) Start(ctx context.Context) error {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Webhook endpoint for AlertManager
	router.POST("/webhook", ws.handleWebhook)

	// Webhook endpoint with receiver name (for multiple receivers)
	router.POST("/webhook/:receiver", ws.handleWebhook)

	server := &http.Server{
		Addr:              ":" + ws.port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	logger := log.FromContext(ctx)
	logger.Info("Starting webhook server", "port", ws.port)

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(err, "Failed to start webhook server")
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	logger.Info("Shutting down webhook server")

	// Shutdown server gracefully
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return server.Shutdown(shutdownCtx)
}

func (ws *WebhookServer) handleWebhook(c *gin.Context) {
	var webhook AlertManagerWebhook

	if err := c.ShouldBindJSON(&webhook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid JSON: %v", err)})
		return
	}

	logger := log.Log.WithValues("receiver", webhook.Receiver, "alertsCount", len(webhook.Alerts))
	logger.Info("Received webhook from AlertManager")

	// Process each alert
	for _, alert := range webhook.Alerts {
		if alert.Status != "firing" {
			logger.Info("Skipping non-firing alert", "alertName", alert.Labels["alertname"], "status", alert.Status)
			continue
		}

		alertName := alert.Labels["alertname"]
		if alertName == "" {
			logger.Info("Skipping alert without alertname label")
			continue
		}

		// Convert alert to map for processing
		alertData := ws.alertToMap(alert)

		// Process the alert
		ctx := context.Background()
		if err := ws.controller.ProcessAlert(ctx, alertName, alertData); err != nil {
			logger.Error(err, "Failed to process alert", "alertName", alertName)
			// Continue processing other alerts even if one fails
		} else {
			logger.Info("Successfully processed alert", "alertName", alertName)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}

func (ws *WebhookServer) alertToMap(alert Alert) map[string]interface{} {
	// Convert Alert struct to a map for easier field access
	alertMap := make(map[string]interface{})

	// Top-level fields
	alertMap["status"] = alert.Status
	alertMap["startsAt"] = alert.StartsAt.Format(time.RFC3339)
	alertMap["endsAt"] = alert.EndsAt.Format(time.RFC3339)
	alertMap["generatorURL"] = alert.GeneratorURL
	alertMap["fingerprint"] = alert.Fingerprint

	// Labels
	if alert.Labels != nil {
		labels := make(map[string]interface{})
		for k, v := range alert.Labels {
			labels[k] = v
			// Also add direct access for common patterns
			alertMap["labels."+k] = v
		}
		alertMap["labels"] = labels
	}

	// Annotations
	if alert.Annotations != nil {
		annotations := make(map[string]interface{})
		for k, v := range alert.Annotations {
			annotations[k] = v
			// Also add direct access for common patterns
			alertMap["annotations."+k] = v
		}
		alertMap["annotations"] = annotations
	}

	return alertMap
}

// GetWebhookURL returns the webhook URL that should be configured in AlertManager
func (ws *WebhookServer) GetWebhookURL(baseURL string) string {
	if baseURL == "" {
		baseURL = "http://localhost:" + ws.port
	}
	return baseURL + "/webhook"
}

// GetWebhookConfig returns an example AlertManager configuration
func (ws *WebhookServer) GetWebhookConfig(baseURL string) string {
	webhookURL := ws.GetWebhookURL(baseURL)

	config := fmt.Sprintf(`# AlertManager configuration example
# Add this to your AlertManager configuration file

route:
  group_by: ['alertname']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 1h
  receiver: 'karo'

receivers:
- name: 'karo'
  webhook_configs:
  - url: '%s'
    send_resolved: false
    http_config:
      timeout: 10s
    max_alerts: 0  # Send all alerts, no limit
`, webhookURL)

	return config
}
