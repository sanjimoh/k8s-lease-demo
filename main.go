package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/klog/v2"
)

var (
	leaseName      = flag.String("lease-name", "leader-election-demo", "Name of the lease object")
	leaseNamespace = flag.String("lease-namespace", "default", "Namespace of the lease object")
	podName        = flag.String("pod-name", "", "Name of the pod")
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	if *podName == "" {
		*podName = os.Getenv("POD_NAME")
		if *podName == "" {
			klog.Fatal("Pod name must be provided via --pod-name or POD_NAME environment variable")
		}
	}

	// Create Kubernetes client
	config, err := rest.InClusterConfig()
	if err != nil {
		klog.Fatalf("Failed to create in-cluster config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Failed to create clientset: %v", err)
	}

	// Create a new lease lock
	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      *leaseName,
			Namespace: *leaseNamespace,
		},
		Client: clientset.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: *podName,
		},
	}

	// Create leader election config
	leaderelectionConfig := leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   15 * time.Second,
		RenewDeadline:   10 * time.Second,
		RetryPeriod:     2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				klog.Infof("Pod %s is now the leader", *podName)
				// Keep the leader running
				<-ctx.Done()
			},
			OnStoppedLeading: func() {
				klog.Infof("Pod %s lost leadership", *podName)
			},
		},
	}

	// Create a context that will be canceled on SIGTERM or SIGINT
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigChan
		klog.Info("Received shutdown signal, canceling context")
		cancel()
	}()

	// Start leader election
	leaderelection.RunOrDie(ctx, leaderelectionConfig)
}
