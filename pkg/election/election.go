package election

import (
	"context"
	"os"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/klog"
)

// NewElection creates and runs a new leader election
func NewElection(ctx context.Context, callback func(leader string)) {
	id := os.Getenv("HOSTNAME")
	namespace := "default"
	resourceLockName := "test"

	// We only care for inClusterConfig
	config, err := rest.InClusterConfig()

	if err != nil {
		klog.Fatalf("Error getting cluster config: %s", err.Error())
	}

	client := kubernetes.NewForConfigOrDie(config)

	// Create the lock resource
	lock, err := resourcelock.New("configmaps", namespace, resourceLockName, client.CoreV1(), client.CoordinationV1(), resourcelock.ResourceLockConfig{Identity: id})

	if err != nil {
		klog.Fatalf("Could not create resource lock: %s", err.Error())
	}

	callbacks := leaderelection.LeaderCallbacks{
		OnStartedLeading: func(ctx context.Context) {
			// ToDo: return own id to callback?
		},
		OnStoppedLeading: func() {
			klog.Infof("Leader lost: %s", id)
		},
		OnNewLeader: func(identity string) {
			callback(identity)
		},
	}

	// ToDo: Revisit values
	conf := leaderelection.LeaderElectionConfig{
		Lock:          lock,
		LeaseDuration: time.Second * 15,
		RenewDeadline: time.Second * 10,
		RetryPeriod:   time.Second * 2,
		Callbacks:     callbacks,
	}

	leaderelection.RunOrDie(ctx, conf)
	klog.Info("Exiting election loop.")
}
