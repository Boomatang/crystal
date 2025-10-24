package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/boomatang/crystal/internal/applicaton"
	"github.com/boomatang/crystal/internal/logger"
	"github.com/boomatang/crystal/internal/workflow"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	events    []workflow.AdmissionReview
	mu        sync.Mutex
	eventChan chan bool
)

func init() {
	workflow.MustRegister()
}

func workflowGraphHandler(world *workflow.World) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		render := workflow.WorldGraph{World: world}
		fmt.Fprintln(w, render.Render())
	}

}

func nodelistGraphHandler(nodes *workflow.NodeList) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		kind := r.URL.Query().Get("kind")
		name := r.URL.Query().Get("name")
		directionStr := r.URL.Query().Get("direction")

		// If no parameters provided, render full node list
		if kind == "" || name == "" {
			fmt.Fprintln(w, nodes.Render())
			return
		}

		// Find the specific node
		node := nodes.Get(kind, name)
		if node == nil {
			http.Error(w, fmt.Sprintf("Node not found: %s/%s", kind, name), http.StatusNotFound)
			return
		}

		// Parse direction parameter
		var direction workflow.Direction
		switch directionStr {
		case "up":
			direction = workflow.Up
		case "down":
			direction = workflow.Down
		case "both", "":
			direction = workflow.Both
		default:
			http.Error(w, "Invalid direction. Must be 'up', 'down', or 'both'", http.StatusBadRequest)
			return
		}

		// Create subgraph options
		opts := workflow.NodeOptList{
			Direction: direction,
		}

		// Get the subgraph
		subgraph, err := nodes.GetSubGraph(node, opts)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating subgraph: %v", err), http.StatusInternalServerError)
			return
		}

		// Render the subgraph
		fmt.Fprintln(w, subgraph.Render())
	}
}

func eventHandler(c chan bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Log.Info("event handler started")

		event := workflow.AdmissionReview{}
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			http.Error(w, "Invaild JSON", http.StatusBadRequest)
			return
		}

		mu.Lock()
		events = append(events, event)
		mu.Unlock()
		c <- true

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(event)
		logger.Log.Info("event handler completed")
	}

}

func listEventsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(events)
	}
}

func eventProcessor(world *workflow.World, nodes *workflow.NodeList) {
	for range eventChan {
		mu.Lock()
		load := events
		events = make([]workflow.AdmissionReview, 0)
		mu.Unlock()
		for _, l := range load {
			existing := nodes.Get(l.Request.Kind.Kind, l.Request.Name)

			logger.Log.Debug("node exists check", "existing", existing != nil)
			if existing == nil {
				node := workflow.NewNode(*l.Request)
				logger.Log.Info("adding new node",
					"node_kind", node.Kind,
					"node_name", node.Name)
				nodes.Add(node)
				existing = node
			}
			logger.Log.Debug("node count updated", "count", nodes.Len())
			nodes.Link(existing)
		}
		logger.Log.Info("workflow triggered by event")

		workflow.NodeCount.Set(float64(nodes.Len()))
		err := world.RunAction(nodes)
		if err != nil {
			logger.Log.Error("workflow execution error", "error", err)
		}
	}
}

func main() {
	// Initialize logger from environment variables
	logger.InitFromEnv()

	worldMain := applicaton.NewApplictaion()
	nodes := workflow.NewNodeList()
	nodes.SetLinker(applicaton.DeploymentConfigMapLinker)
	nodes.SetLinker(applicaton.ConfigSecretMapLinker)

	eventChan = make(chan bool, 100)
	go eventProcessor(worldMain, nodes)
	mux := http.NewServeMux()

	// Expose metrics at /metrics
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/graph", workflowGraphHandler(worldMain))
	mux.Handle("/nodelist", nodelistGraphHandler(nodes))
	mux.Handle("/event", eventHandler(eventChan))
	mux.Handle("/list", listEventsHandler())

	port := ":8000"
	logger.Log.Info("server started", "port", port)
	err := http.ListenAndServe(port, mux)
	if err != nil {
		logger.Log.Error("server failed to start", "error", err)
	}

}
