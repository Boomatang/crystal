package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/boomatang/crystal/internal/workflow"
)

var worldMain *workflow.World

func echoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Method: %s\n", r.Method)
	fmt.Fprintf(w, "URL: %s\n", r.URL.String())

	fmt.Fprint(w, "Header:")
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Fprintf(w, "  %s: %s\n", name, value)
		}
	}

	fmt.Fprintln(w, "Body:")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()
	fmt.Fprintln(w, string(body))

}

func workflowGraphHandler(world *workflow.World) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		fmt.Print(world.GetName())
		render := workflow.WorldGraph{World: world}
		fmt.Fprintln(w, render.Render())
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

	// err := worldMain.RunAction()
	// if err != nil {
	// 	fmt.Printf("error was raised, %s", err)
	// }

	// render := workflow.WorldGraph{World: worldMain}
	// fmt.Println(render.Render())

	http.HandleFunc("/echo", echoHandler)
	http.HandleFunc("/graph", workflowGraphHandler(worldMain))

	port := ":8000"
	fmt.Println("Echo server running on http://localhost" + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Server Failed:", err)
	}

}
