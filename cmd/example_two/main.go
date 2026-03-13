package main

import (
	"net/http"

	"github.com/boomatang/crystal/internal/applicaton"
	"github.com/boomatang/crystal/internal/logger"
	rhttp "github.com/boomatang/crystal/internal/runtime/http"
	"github.com/boomatang/crystal/internal/workflow"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Initialize logger from environment variables
	logger.InitFromEnv()
	logger.Log.Info("running example two")
	logger.Log.Info("adds sleep in workflow to show queuing")

	worldMain := applicaton.NewApplictaionTwo()
	nodes := workflow.NewNodeList()
	nodes.SetLinker(applicaton.DeploymentConfigMapLinker)
	nodes.SetLinker(applicaton.ConfigSecretMapLinker)

	queue := workflow.NewEventQueue()
	rhttp.EventChan = make(chan bool, 100)
	go rhttp.EventProcessor(worldMain, nodes, queue)
	mux := http.NewServeMux()

	// Expose metrics at /metrics
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/graph", rhttp.WorkflowGraphHandler(worldMain))
	mux.Handle("/nodelist", rhttp.NodelistGraphHandler(nodes))
	mux.Handle("/event", rhttp.EventHandler(rhttp.EventChan, queue))
	mux.Handle("/list", rhttp.ListEventsHandler(queue))

	port := ":8000"
	logger.Log.Info("server started", "port", port)
	err := http.ListenAndServe(port, mux)
	if err != nil {
		logger.Log.Error("server failed to start", "error", err)
	}

}
