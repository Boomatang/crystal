package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/boomatang/crystal/internal/applicaton"
	"github.com/boomatang/crystal/internal/workflow"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	events    []workflow.Event
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
	return func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, nodes.Render())
	}

}

func eventHandler(c chan bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("running the event handler")

		var event workflow.Event
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
		events = make([]workflow.Event, 0)
		mu.Unlock()
		for _, l := range load {
			existing := nodes.Get(l.Kind, l.Name)
			if existing == nil {
				node := workflow.NewNode(l)
				nodes.Add(node)
				existing = node
			}
			nodes.Link(existing)
		}
		fmt.Println("triggered by event")

		workflow.NodeCount.Set(float64(nodes.Len()))
		err := world.RunAction()
		if err != nil {
			fmt.Printf("error was raised, %s", err)
		}
	}
}

func main() {

	worldMain := applicaton.NewApplictaion()
	eventChan = make(chan bool, 100)
	link := workflow.Link{Parent: "crd", Child: "cr", LinkFunc: func(p, c *workflow.Node) bool {
		fmt.Println("Happy holidays")
		return true
	}}
	nodes := workflow.NewNodeList()
	nodes.SetLinker(link)
	go eventProcessor(worldMain, nodes)

	mux := http.NewServeMux()

	// Expose metrics at /metrics
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/graph", workflowGraphHandler(worldMain))
	mux.Handle("/nodelist", nodelistGraphHandler(nodes))
	mux.Handle("/event", eventHandler(eventChan))
	mux.Handle("/list", listEventsHandler())

	port := ":8000"
	fmt.Println("Echo server running on http://localhost" + port)
	err := http.ListenAndServe(port, mux)
	if err != nil {
		fmt.Println("Server Failed:", err)
	}

}
