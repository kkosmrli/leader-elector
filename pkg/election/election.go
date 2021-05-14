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

type Config struct {
	LockType      string
	LockName      string
	LockNamespace string
	RetryPeriod   time.Duration
	LeaseDuration time.Duration
	RenewDeadline time.Duration
	Callback      func(leader string)
}

// Run creates and runs a new leader election
func Run(ctx context.Context, cfg Config) {
	id := os.Getenv("HOSTNAME")

	// We only care for inClusterConfig
	config, err := rest.InClusterConfig()

	if err != nil {
		klog.Fatalf("Error getting cluster config: %s", err.Error())
	}

	client := kubernetes.NewForConfigOrDie(config)

	// Create the lock resource
	lock, err := resourcelock.New(
		cfg.LockType,
		cfg.LockNamespace,
		cfg.LockName,
		client.CoreV1(),
		client.CoordinationV1(),
		resourcelock.ResourceLockConfig{Identity: id})

	if err != nil {
		klog.Fatalf("Could not create resource lock: %s", err.Error())
	}

	callbacks := leaderelection.LeaderCallbacks{
		OnStartedLeading: func(ctx context.Context) {
			cfg.Callback(id)
		},
		OnStoppedLeading: func() {
			klog.Infof("Leader lost: %s", id)
		},
		OnNewLeader: func(identity string) {
			cfg.Callback(identity)
		},
	}

	leaderElectionConf := leaderelection.LeaderElectionConfig{
		Lock:          lock,
		LeaseDuration: cfg.LeaseDuration,
		RenewDeadline: cfg.RenewDeadline,
		RetryPeriod:   cfg.RetryPeriod,
		Callbacks:     callbacks,
	}

	leaderelection.RunOrDie(ctx, leaderElectionConf)
	klog.Info("Exiting election loop.")
}
