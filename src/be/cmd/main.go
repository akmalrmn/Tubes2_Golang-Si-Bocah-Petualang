package main

import (
	"be/pkg/algorithms/bfs"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {

	// ! Profiler
	go func() {
		// Start the HTTP server on port 8080
		err := http.ListenAndServe("localhost:8080", nil)
		if err != nil {
			panic(err)
		}
	}()

	// Start the timer
	start := time.Now()

	// Call the BreadthFirstSearch function
	result := bfs.BreadthFirstSearch("/wiki/Joko_Widodo", "/wiki/Sleman_Regency")

	if result != nil {
		fmt.Println("Path found!")
	} else {
		fmt.Println("Path not found!")
	}

	fmt.Println("Total article checked: ", bfs.HowManyArticelChecked)

	elapsed := time.Since(start)
	fmt.Printf("Execution time: %d minute %d seconds %d miliseconds \n", int(elapsed.Minutes()), int(elapsed.Seconds())%60, (elapsed.Milliseconds())%1000)
	// Print the result

	return
}
