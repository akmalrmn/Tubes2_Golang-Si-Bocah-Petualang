package main

import (
	"be/pkg/algorithms/bfs"
	config2 "be/pkg/config"
	"be/pkg/server"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"runtime/debug"
	"time"
)

func main() {
	// Load the configuration
	config := config2.NewTurboConfig()

	// Setup the process
	runtime.GOMAXPROCS(config.MaxProcessor)
	debug.SetMaxThreads(config.MaxThreads)

	//Create a new HTTP server
	var handler = server.NewHttpHandler()
	handler = server.AllowCORS(handler)
	log.Println("Server is available at http://localhost:8000")

	// Testing
	start := time.Now()
	bfs.BFS("/wiki/Joko_Widodo", "/wiki/Sleman_Regency", config)
	elapsed := time.Since(start)
	fmt.Printf("BFS took %s\n", elapsed)

	log.Fatal(http.ListenAndServe(":8000", handler))
}
