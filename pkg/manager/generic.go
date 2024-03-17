package manager

import (
	networkv1alpha2 "github.com/openelb/openelb/api/v1alpha2"
	"github.com/openelb/openelb/pkg/client"
	"github.com/openelb/openelb/pkg/constant"
	"github.com/spf13/pflag"
	admissionv1 "k8s.io/api/admission/v1"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	nc "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

type GenericOptions struct {
	WebhookPort      int
	WebhookCertDir   string
	MetricsAddr      string
	ReadinessAddr    string
	LeaderElector    bool
	LeaderElectionID string
}

func NewGenericOptions() *GenericOptions {
	return &GenericOptions{
		WebhookPort:      443,
		WebhookCertDir:   "/tmp/k8s-webhook-server/serving-certs",
		MetricsAddr:      ":50052",
		ReadinessAddr:    "0",
		LeaderElector:    true,
		LeaderElectionID: constant.OpenELBControllerLocker,
	}
}

func (options *GenericOptions) AddFlags(fs *pflag.FlagSet) {
	fs.IntVar(&options.WebhookPort, "webhook-port", options.WebhookPort, "The port that the webhook server serves at")
	fs.StringVar(&options.WebhookCertDir, "webhook-cert-dir", options.WebhookCertDir, "Certificate directory used to setup webhooks.")
	fs.StringVar(&options.MetricsAddr, "metrics-addr", options.MetricsAddr, "The address the metric endpoint binds to.")
	fs.StringVar(&options.ReadinessAddr, "readiness-addr", options.ReadinessAddr, "The address readinessProbe used")
	fs.BoolVar(&options.LeaderElector, "leader-elect", options.LeaderElector, "Enable leader election for controller manager")
}

func NewManager(cfg *rest.Config, options *GenericOptions) (ctrl.Manager, error) {
	opts := ctrl.Options{Scheme: scheme}
	if options != nil {
		opts.WebhookServer = webhook.NewServer(webhook.Options{
			CertDir: options.WebhookCertDir,
			Port:    options.WebhookPort,
		})
		opts.Metrics.BindAddress = options.MetricsAddr
		opts.LeaderElection = options.LeaderElector
		opts.LeaderElectionID = options.LeaderElectionID
	}
	result, err := ctrl.NewManager(cfg, opts)

	if err == nil {
		client.Client, err = nc.New(cfg, nc.Options{Scheme: scheme})
	}

	return result, err
}

var (
	scheme = runtime.NewScheme()
)

func init() {
	_ = corev1.AddToScheme(scheme)
	_ = admissionv1.AddToScheme(scheme)
	_ = admissionv1beta1.AddToScheme(scheme)
	_ = networkv1alpha2.AddToScheme(scheme)
	_ = appsv1.AddToScheme(scheme)
}
