package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/kkosmrli/leader-elector/pkg/election"
	"k8s.io/klog"
)

var leader Leader

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	callback := func(name string) {
		klog.Infof("Currently leading: %s", name)
		leader = Leader{name}
	}

	go election.NewElection(ctx, callback)

	http.HandleFunc("/", leaderHandler)
	klog.Fatal(http.ListenAndServe(":4040", nil))
}
