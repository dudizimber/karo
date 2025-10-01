package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	alertreactionv1alpha1 "github.com/dudizimber/k8s-alert-reaction-operator/api/v1alpha1"
	"github.com/dudizimber/k8s-alert-reaction-operator/controllers"
	"github.com/dudizimber/k8s-alert-reaction-operator/webhook"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(alertreactionv1alpha1.AddToScheme(scheme))
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var webhookPort string

	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.StringVar(&webhookPort, "webhook-port", "9090", "The port for the webhook server.")

	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                  scheme,
		Metrics:                 metricsserver.Options{BindAddress: metricsAddr},
		HealthProbeBindAddress:  probeAddr,
		LeaderElection:          enableLeaderElection,
		LeaderElectionID:        "f1c5ece8.alertreaction.io",
		LeaderElectionNamespace: "default",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.AlertReactionReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AlertReaction")
		os.Exit(1)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	// Create webhook server
	alertReactionController := &controllers.AlertReactionReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}
	webhookServer := webhook.NewWebhookServer(alertReactionController, webhookPort)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle termination signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start webhook server in a goroutine
	go func() {
		if err := webhookServer.Start(ctx); err != nil {
			setupLog.Error(err, "problem running webhook server")
			cancel()
		}
	}()

	// Start manager in a goroutine
	go func() {
		setupLog.Info("starting manager")
		if err := mgr.Start(ctx); err != nil {
			setupLog.Error(err, "problem running manager")
			cancel()
		}
	}()

	// Print webhook configuration
	setupLog.Info("Alert Reaction Operator started successfully")
	setupLog.Info("Webhook server configuration:")
	setupLog.Info(fmt.Sprintf("Webhook endpoint: http://localhost:%s/webhook", webhookPort))
	setupLog.Info("Add this to your AlertManager configuration:")
	setupLog.Info(webhookServer.GetWebhookConfig(""))

	// Wait for termination signal
	<-sigChan
	setupLog.Info("Received termination signal, shutting down...")
	cancel()

	setupLog.Info("Alert Reaction Operator stopped")
}
