//go:build integration
// +build integration

package integration

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	alertreactionv1alpha1 "github.com/dudizimber/karo/api/v1alpha1"
	"github.com/dudizimber/karo/controllers"
)

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment
var ctx context.Context
var cancel context.CancelFunc

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	ctx, cancel = context.WithCancel(context.TODO())

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "config", "crd")},
		ErrorIfCRDPathMissing: true,
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = alertreactionv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = batchv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	err = (&controllers.AlertReactionReconciler{
		Client: k8sManager.GetClient(),
		Scheme: k8sManager.GetScheme(),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err = k8sManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred(), "failed to run manager")
	}()
})

var _ = AfterSuite(func() {
	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("AlertReaction Controller", func() {
	const (
		AlertReactionName      = "test-alertreaction"
		AlertReactionNamespace = "default"
		AlertName              = "TestAlert"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When creating an AlertReaction", func() {
		It("Should create the AlertReaction successfully", func() {
			By("Creating a new AlertReaction")
			ctx := context.Background()
			alertReaction := &alertreactionv1alpha1.AlertReaction{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "karo.io/v1",
					Kind:       "AlertReaction",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      AlertReactionName,
					Namespace: AlertReactionNamespace,
				},
				Spec: alertreactionv1alpha1.AlertReactionSpec{
					AlertName: AlertName,
					Actions: []alertreactionv1alpha1.Action{
						{
							Name:    "test-action",
							Image:   "busybox:latest",
							Command: []string{"echo", "test"},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, alertReaction)).Should(Succeed())

			alertReactionLookupKey := types.NamespacedName{Name: AlertReactionName, Namespace: AlertReactionNamespace}
			createdAlertReaction := &alertreactionv1alpha1.AlertReaction{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, alertReactionLookupKey, createdAlertReaction)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			Expect(createdAlertReaction.Spec.AlertName).Should(Equal(AlertName))
			Expect(createdAlertReaction.Spec.Actions).Should(HaveLen(1))
			Expect(createdAlertReaction.Spec.Actions[0].Name).Should(Equal("test-action"))
		})

		It("Should update the status with conditions", func() {
			alertReactionLookupKey := types.NamespacedName{Name: AlertReactionName, Namespace: AlertReactionNamespace}
			createdAlertReaction := &alertreactionv1alpha1.AlertReaction{}

			Eventually(func() int {
				err := k8sClient.Get(context.Background(), alertReactionLookupKey, createdAlertReaction)
				if err != nil {
					return 0
				}
				return len(createdAlertReaction.Status.Conditions)
			}, timeout, interval).Should(BeNumerically(">", 0))

			Expect(createdAlertReaction.Status.Conditions[0].Type).Should(Equal("Ready"))
			Expect(createdAlertReaction.Status.Conditions[0].Status).Should(Equal(metav1.ConditionTrue))
		})
	})

	Context("When processing alerts", func() {
		It("Should create jobs for matching alerts", func() {
			By("Processing an alert that matches the AlertReaction")
			ctx := context.Background()

			// Get the controller
			reconciler := &controllers.AlertReactionReconciler{
				Client: k8sClient,
				Scheme: scheme.Scheme,
			}

			alertData := map[string]interface{}{
				"status": "firing",
				"labels": map[string]interface{}{
					"alertname": AlertName,
					"instance":  "test-instance",
				},
				"annotations": map[string]interface{}{
					"summary": "Test alert",
				},
				"labels.alertname": AlertName,
				"labels.instance":  "test-instance",
			}

			err := reconciler.ProcessAlert(ctx, AlertName, alertData)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() int {
				jobList := &batchv1.JobList{}
				err := k8sClient.List(ctx, jobList, client.InNamespace(AlertReactionNamespace))
				if err != nil {
					return 0
				}
				return len(jobList.Items)
			}, timeout, interval).Should(BeNumerically(">", 0))

			// Verify job properties
			jobList := &batchv1.JobList{}
			err = k8sClient.List(ctx, jobList, client.InNamespace(AlertReactionNamespace))
			Expect(err).ToNot(HaveOccurred())
			Expect(jobList.Items).To(HaveLen(1))

			job := jobList.Items[0]
			Expect(job.Labels).To(HaveKey("alert-reaction/alert-name"))
			Expect(job.Labels["alert-reaction/alert-name"]).To(Equal(AlertName))
			Expect(job.Labels).To(HaveKey("alert-reaction/action-name"))
			Expect(job.Labels["alert-reaction/action-name"]).To(Equal("test-action"))
		})

		It("Should update AlertReaction status after processing alerts", func() {
			alertReactionLookupKey := types.NamespacedName{Name: AlertReactionName, Namespace: AlertReactionNamespace}
			updatedAlertReaction := &alertreactionv1alpha1.AlertReaction{}

			Eventually(func() int64 {
				err := k8sClient.Get(context.Background(), alertReactionLookupKey, updatedAlertReaction)
				if err != nil {
					return 0
				}
				return updatedAlertReaction.Status.TriggerCount
			}, timeout, interval).Should(BeNumerically(">", 0))

			Expect(updatedAlertReaction.Status.LastTriggered).ToNot(BeNil())
			Expect(updatedAlertReaction.Status.LastJobsCreated).To(HaveLen(1))
			Expect(updatedAlertReaction.Status.LastJobsCreated[0].ActionName).To(Equal("test-action"))
		})
	})

	Context("When cleaning up", func() {
		It("Should delete the AlertReaction", func() {
			By("Deleting the AlertReaction")
			ctx := context.Background()
			alertReactionLookupKey := types.NamespacedName{Name: AlertReactionName, Namespace: AlertReactionNamespace}
			alertReaction := &alertreactionv1alpha1.AlertReaction{}

			err := k8sClient.Get(ctx, alertReactionLookupKey, alertReaction)
			Expect(err).ToNot(HaveOccurred())

			err = k8sClient.Delete(ctx, alertReaction)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				err := k8sClient.Get(ctx, alertReactionLookupKey, alertReaction)
				return err != nil
			}, timeout, interval).Should(BeTrue())
		})
	})
})
