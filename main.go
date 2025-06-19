package main

import (
	"fmt"
	"io"
	"net/http"
)

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

func main() {
	http.HandleFunc("/", echoHandler)

	port := ":8000"
	fmt.Println("Echo server running on http://localhost" + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Server Failed:", err)
	}
}
