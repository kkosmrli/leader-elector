package main

import (
	"context"
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kkosmrli/leader-elector/pkg/election"
	"k8s.io/klog"
)

var (
	electionName      string
	electionNamespace string
	lockType          string
	renewDeadline     time.Duration
	retryPeriod       time.Duration
	leaseDuration     time.Duration
	port              string
	leader            Leader
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

func parseFlags() {
	flag.StringVar(&electionName, "election", "default", "Name of the resource used for this election")
	flag.StringVar(&electionNamespace, "namespace", "default", "Namespace of the resource used for this election")
	flag.StringVar(&lockType, "locktype", "configmaps",
		"Resource lock type, must be one of the following: configmaps, endpoints, leases")
	flag.DurationVar(&renewDeadline, "renew-deadline", 10*time.Second,
		"Duration that the acting leader will retry refreshing leadership before giving up")
	flag.DurationVar(&leaseDuration, "lease-duration", 15*time.Second,
		`Duration that non-leader candidates will wait after observing a leadership
		renewal until attempting to acquire leadership of a led but unrenewed leader slot`)
	flag.DurationVar(&retryPeriod, "retry-period", 2*time.Second, "Duration between each action retry")
	flag.StringVar(&port, "port", "4040", "Port on which to query the leader")
	flag.Parse()
}

func main() {
	parseFlags()

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
	server := &http.Server{Addr: ":" + port, Handler: nil}
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
		LockName:      electionName,
		LockNamespace: electionNamespace,
		LockType:      lockType,
		RenewDeadline: renewDeadline,
		RetryPeriod:   retryPeriod,
		LeaseDuration: leaseDuration,
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
