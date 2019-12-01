package main

import (
	"context"
	"encoding/json"
	"flag"
	"net/http"

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	callback := func(name string) {
		klog.Infof("Currently leading: %s", name)
		leader = Leader{name}
	}

	go election.NewElection(ctx, electionName, namespace, callback)

	http.HandleFunc("/", leaderHandler)
	klog.Fatal(http.ListenAndServe(":"+port, nil))
}
