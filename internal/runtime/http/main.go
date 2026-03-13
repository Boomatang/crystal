package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boomatang/crystal/internal/logger"
	"github.com/boomatang/crystal/internal/workflow"
)

var (
	EventChan chan bool
)

func WorkflowGraphHandler(world *workflow.World) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		render := workflow.WorldGraph{World: world}
		fmt.Fprintln(w, render.Render())
	}

}

func NodelistGraphHandler(nodes *workflow.NodeList) http.HandlerFunc {
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

func EventHandler(c chan bool, queue *workflow.EventQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		event := workflow.AdmissionReview{}
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		added, _ := queue.Add(event)
		if added {
			c <- true
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(event)
	}

}

func ListEventsHandler(queue *workflow.EventQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(queue.Snapshot())
	}
}

func EventProcessor(world *workflow.World, nodes *workflow.NodeList, queue *workflow.EventQueue) {
	for range EventChan {
		load := queue.DrainAndCopy()
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
			} else {
				existing.Data = l.Request.Object
				logger.Log.Info("updating existing node",
					"node_kind", existing.Kind,
					"node_name", existing.Name)
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
