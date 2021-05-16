package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/kkosmrli/leader-elector/pkg/election"
	"k8s.io/klog"
)

var (
	args struct {
		LockName      string `arg:"--election,env:ELECTION_NAME" default:"default" help:"Name of this election"`
		Namespace     string `arg:"env:ELECTION_NAMESPACE" default:"default" help:"Namespace of this election"`
		LockType      string `arg:"env:ELECTION_TYPE" default:"configmaps" help:"Resource lock type, must be one of the following: configmaps, endpoints, leases"`
		RenewDeadline time.Duration `arg:"--renew-deadline,env:ELECTION_RENEW_DEADLINE" default:"10s" help:"Duration that the acting leader will retry refreshing leadership before giving up"`
		RetryPeriod   time.Duration `arg:"--retry-period,env:ELECTION_RETRY_PERIOD" default:"2s" help:"Duration between each action retry"`
		LeaseDuration time.Duration `arg:"--lease-duration,env:ELECTION_LEASE_DURATION" default:"15s" help:"Duration that non-leader candidates will wait after observing a leadership renewal until attempting to acquire leadership of a led but unrenewed leader slot"`
		Port          string `arg:"env:ELECTION_PORT" default:"4040" help:"Port on which to query the leader"`
	}
	leader Leader
)

// Leader contains the name of the current leader of this election
type Leader struct {
	Name string `json:"name"`
}

func leaderHandler(res http.ResponseWriter, req *http.Request) {
	data, err := json.Marshal(leader)
	if err != nil {
		klog.Errorf("Error while marshaling leader response: %s", err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Write(data)
}

func main() {
	arg.MustParse(&args)

	// configuring context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// configuring signal handling
	terminationSignal := make(chan os.Signal, 1)
	signal.Notify(terminationSignal, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-terminationSignal
		klog.Infoln("Received termination signal, shutting down")
		cancel()
	}()

	// configuring HTTP server
	http.HandleFunc("/", leaderHandler)
	server := &http.Server{Addr: ":" + args.Port, Handler: nil}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			klog.Fatal(err)
		}
	}()

	// configuring Leader Election loop
	callback := func(name string) {
		klog.Infof("Currently leading: %s", name)
		leader = Leader{name}
	}

	electionConfig := election.Config{
		LockName:      args.LockName,
		LockNamespace: args.Namespace,
		LockType:      args.LockType,
		RenewDeadline: args.RenewDeadline,
		RetryPeriod:   args.RetryPeriod,
		LeaseDuration: args.LeaseDuration,
		Callback:      callback,
	}
	election.Run(ctx, electionConfig)

	// gracefully stop HTTP server
	srvCtx, srvCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer srvCancel()
	if err := server.Shutdown(srvCtx); err != nil {
		klog.Fatal(err)
	}
}
