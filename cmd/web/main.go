package main

import (
	"fmt"
	"net/http"
)

func main() {
	outputDir := "output"

	fs := http.FileServer(http.Dir(outputDir))

	http.Handle("/", fs)

	port := 8082

	fmt.Printf("server running at http://localhost:%d/\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n\n", err)
	}
}
