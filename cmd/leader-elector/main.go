package main

import (
	"context"
	"net/http"

	"github.com/kkosmrli/leader-elector/pkg/election"
	"k8s.io/klog"
)

var leader string

func leaderHandler(res http.ResponseWriter, req *http.Request) {

}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	callback := func(name string) {
		klog.Infof("Currently leading: %s", leader)
		leader = name
	}

	go election.NewElection(ctx, callback)

	http.HandleFunc("/", leaderHandler)
	klog.Fatal(http.ListenAndServe(":4040", nil))
}
