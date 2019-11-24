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
	//id := uuid.New().String()
	namespace := "default"
	resourceLockName := "test"

	// We only care for in ClusterConfig atm
	config, err := rest.InClusterConfig()

	if err != nil {
		klog.Fatalf("Error getting cluster config: %s", err.Error())
	}

	client := kubernetes.NewForConfigOrDie(config)

	// Try config map lock for now
	lock, err := resourcelock.New("configmaps", namespace, resourceLockName, client.CoreV1(), client.CoordinationV1(), resourcelock.ResourceLockConfig{Identity: id})

	if err != nil {
		// could not create resourcelock
		klog.Fatalf("Could not create resource lock: %s", err.Error())
	}

	callbacks := leaderelection.LeaderCallbacks{
		OnStartedLeading: func(ctx context.Context) {
			// Whaat to do with ctx here??
		},
		OnStoppedLeading: func() {
			klog.Infof("Leader lost: %s", id)
		},
		OnNewLeader: func(identity string) {
			callback(identity)
		},
	}

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
