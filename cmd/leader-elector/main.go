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
	electionName string
	namespace    string
	locktype     string
	port         string
	leader       Leader
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
	flag.StringVar(&electionName, "election", "default", "Name of this election")
	flag.StringVar(&namespace, "namespace", "default", "Namespace of this election")
	flag.StringVar(&locktype, "locktype", "configmaps", "Resource lock type, must be one of the following: configmaps, endpoints, leases")
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
	election.NewElection(ctx, namespace, electionName, locktype, callback)

	// gracefully stop HTTP server
	srvCtx, srvCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer srvCancel()
	if err := server.Shutdown(srvCtx); err != nil {
		klog.Fatal(err)
	}
}
