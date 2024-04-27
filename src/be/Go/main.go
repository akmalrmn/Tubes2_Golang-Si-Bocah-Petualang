package main

import (

	config2 "be/pkg/config"
	"be/pkg/server"
	"log"
	"net/http"
	"runtime"
	"runtime/debug"
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

	log.Fatal(http.ListenAndServe(":8000", handler))
}
