package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/boomatang/crystal/internal/workflow"
)

var (
	events    []workflow.Event
	mu        sync.Mutex
	eventChan chan bool
)

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
		fmt.Println("This should be printing")
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
		fmt.Println("this should happen")
		for _, l := range load {
			existing := nodes.Get(l.Kind, l.Name)
			if existing == nil {
				node := workflow.NewNode(l)
				nodes.Add(node)
				existing = node
			}
			nodes.Link(existing)
			fmt.Print("existing parents: ")
			fmt.Println(existing.Parents)
			fmt.Print("existing Childern: ")
			fmt.Println(existing.Childern)
		}
		fmt.Println("triggered by event")
		fmt.Println(nodes)

		// err := world.RunAction()
		// if err != nil {
		// 	fmt.Printf("error was raised, %s", err)
		// }
	}
}

func alan() error {
	fmt.Println("This was a function call")
	return nil
}
func tom() error {
	fmt.Println("This was a function call")
	return nil
}
func john() error {
	fmt.Println("This was a function call")
	return nil
}
func mark() error {
	fmt.Println("This was a function call")
	return nil
}

func main() {
	actionAlan := workflow.NewPoint(alan)
	actionTom := workflow.NewPoint(tom)
	actionJohn := workflow.NewPoint(john)
	actionMark := workflow.NewPoint(mark)
	worldOne := workflow.NewWorld("Temporay world")
	worldOne.PreCondition = actionTom
	worldOne.PostCondition = actionJohn
	worldOne.ErrorHandler = actionMark
	for i := 0; i < 10; i++ {
		worldOne.AddAction(actionAlan)
	}

	worldMain := workflow.NewWorld("Main World")
	worldMain.PreCondition = worldOne
	worldMain.AddAction(actionAlan)
	worldMain.AddAction(actionTom)
	worldMain.AddAction(actionMark)
	worldMain.AddAction(actionJohn)
	worldMain.AddAction(worldOne)
	worldMain.PostCondition = worldOne

	eventChan = make(chan bool, 100)
	link := workflow.Link{Parent: "crd", Child: "cr", LinkFunc: func(p, c *workflow.Node) bool {
		fmt.Println("Happy holidays")
		return true
	}}
	nodes := workflow.NewNodeList()
	nodes.SetLinker(link)
	go eventProcessor(worldMain, nodes)

	http.HandleFunc("/graph", workflowGraphHandler(worldMain))
	http.HandleFunc("/nodelist", nodelistGraphHandler(nodes))
	http.HandleFunc("/event", eventHandler(eventChan))
	http.HandleFunc("/list", listEventsHandler())

	port := ":8000"
	fmt.Println("Echo server running on http://localhost" + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Server Failed:", err)
	}

}
