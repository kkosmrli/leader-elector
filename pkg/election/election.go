package election

import (
	"context"
	"time"

	"github.com/google/uuid"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/klog"
)

// NewElection creates and runs a new leader election
func NewElection(ctx context.Context, callback func(leader string)) {

	id := uuid.New().String()
	namespace := "default"
	resourceLockName := "test"

	// We only care for in ClusterConfig
	config, err := rest.InClusterConfig()

	if err != nil {
		panic(err.Error())
	}

	client := kubernetes.NewForConfigOrDie(config)

	// Try config map lock for now
	lock, err := resourcelock.New("ConfigMapsResourceLock", namespace, resourceLockName, client.CoreV1(), client.CoordinationV1(), resourcelock.ResourceLockConfig{Identity: id})

	if err != nil {
		// could not create resourcelock
		panic(err.Error())
	}

	callbacks := leaderelection.LeaderCallbacks{
		OnStartedLeading: func(ctx context.Context) {
			// Whaat to do with ctx here??
		},
		OnStoppedLeading: func() {
			klog.Infof("leader lost: %s", id)
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
}
